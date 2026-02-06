package runtime

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

// Server runs the local gateway HTTP runtime.
type Server struct {
	cfg    *config.Config
	routes []config.Route
}

func NewServer(cfg *config.Config) *Server {
	routes := append([]config.Route(nil), cfg.Routes...)
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].Path) > len(routes[j].Path)
	})
	return &Server{cfg: cfg, routes: routes}
}

func (s *Server) ListenAndServe() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok\n"))
	})
	mux.HandleFunc("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ready\n"))
	})
	mux.HandleFunc("/", s.handleRequest)

	h := loggingMiddleware(mux)
	server := &http.Server{
		Addr:              s.cfg.Gateway.ListenAddr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
	}
	log.Printf("gateway %q listening on %s", s.cfg.Gateway.Name, s.cfg.Gateway.ListenAddr)
	return server.ListenAndServe()
}

func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
	route, ok := s.matchRoute(r.URL.Path)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if err := s.enforceAuth(route, r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	server, ok := s.lookupServer(route.Server)
	if !ok {
		http.Error(w, "route server not found", http.StatusBadGateway)
		return
	}

	switch server.Transport {
	case "http":
		s.proxyHTTP(server.URL, w, r)
	case "stdio":
		s.respondNotImplemented(route, server, w)
	default:
		http.Error(w, "unsupported server transport", http.StatusBadGateway)
	}
}

func (s *Server) matchRoute(path string) (config.Route, bool) {
	for _, r := range s.routes {
		if strings.HasPrefix(path, r.Path) {
			return r, true
		}
	}
	return config.Route{}, false
}

func (s *Server) lookupServer(name string) (config.Server, bool) {
	for _, server := range s.cfg.Servers {
		if server.Name == name {
			return server, true
		}
	}
	return config.Server{}, false
}

func (s *Server) enforceAuth(route config.Route, r *http.Request) error {
	authType := routeAuthType(s.cfg, route)
	switch authType {
	case "none":
		return nil
	case "apiKey":
		header := "X-API-Key"
		keys := []string{}
		if route.Auth != nil {
			if strings.TrimSpace(route.Auth.HeaderName) != "" {
				header = route.Auth.HeaderName
			}
			keys = route.Auth.APIKeys
		}
		v := r.Header.Get(header)
		if v == "" {
			return fmt.Errorf("missing API key")
		}
		if len(keys) == 0 {
			return nil
		}
		for _, k := range keys {
			if v == k {
				return nil
			}
		}
		return fmt.Errorf("invalid API key")
	case "jwt":
		authz := strings.TrimSpace(r.Header.Get("Authorization"))
		if !strings.HasPrefix(authz, "Bearer ") || strings.TrimSpace(strings.TrimPrefix(authz, "Bearer ")) == "" {
			return fmt.Errorf("missing bearer token")
		}
		return nil
	default:
		return fmt.Errorf("unsupported auth type %q", authType)
	}
}

func (s *Server) proxyHTTP(rawURL string, w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(rawURL)
	if err != nil {
		http.Error(w, "invalid upstream URL", http.StatusBadGateway)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	proxy.ErrorHandler = func(rw http.ResponseWriter, _ *http.Request, e error) {
		http.Error(rw, fmt.Sprintf("upstream error: %v", e), http.StatusBadGateway)
	}
	proxy.ServeHTTP(w, r)
}

func (s *Server) respondNotImplemented(route config.Route, server config.Server, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotImplemented)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"error":     "stdio transport bridge not yet wired",
		"route":     route.Name,
		"server":    server.Name,
		"transport": server.Transport,
	})
}

func routeAuthType(cfg *config.Config, route config.Route) string {
	if route.Auth == nil {
		if cfg.Auth.RequireAuth {
			return "apiKey"
		}
		return "none"
	}
	if strings.TrimSpace(route.Auth.Type) == "" {
		if cfg.Auth.RequireAuth {
			return "apiKey"
		}
		return "none"
	}
	return route.Auth.Type
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("method=%s path=%s duration=%s", r.Method, r.URL.Path, time.Since(start))
	})
}
