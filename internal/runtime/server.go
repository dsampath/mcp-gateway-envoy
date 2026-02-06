package runtime

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/djsam/mcp-gateway-envoy/internal/config"
)

// Server runs the local gateway HTTP runtime.
type Server struct {
	cfg       *config.Config
	routes    []config.Route
	logBodies bool
}

func NewServer(cfg *config.Config) *Server {
	routes := append([]config.Route(nil), cfg.Routes...)
	sort.Slice(routes, func(i, j int) bool {
		return len(routes[i].Path) > len(routes[j].Path)
	})
	return &Server{
		cfg:       cfg,
		routes:    routes,
		logBodies: strings.EqualFold(os.Getenv("GATEWAY_LOG_BODIES"), "true"),
	}
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
	normalizeRequestPath(r, route.Path)
	if err := s.enforceAuth(route, r); err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	reqBodyPreview, reqBodyTruncated := s.captureRequestBody(r)
	if s.logBodies && reqBodyPreview != "" {
		log.Printf("mcp_request route=%s path=%s body=%q truncated=%t", route.Name, r.URL.Path, reqBodyPreview, reqBodyTruncated)
	}
	server, ok := s.lookupServer(route.Server)
	if !ok {
		http.Error(w, "route server not found", http.StatusBadGateway)
		return
	}

	switch server.Transport {
	case "http":
		s.proxyHTTP(route, server.URL, w, r)
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

func (s *Server) proxyHTTP(route config.Route, rawURL string, w http.ResponseWriter, r *http.Request) {
	target, err := url.Parse(rawURL)
	if err != nil {
		http.Error(w, "invalid upstream URL", http.StatusBadGateway)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(target)
	if s.logBodies {
		proxy.ModifyResponse = func(resp *http.Response) error {
			resp.Body = newLoggingReadCloser(resp.Body, 16*1024, func(preview string, truncated bool) {
				log.Printf("mcp_response route=%s status=%d content_type=%q body=%q truncated=%t",
					route.Name, resp.StatusCode, resp.Header.Get("Content-Type"), preview, truncated)
			})
			return nil
		}
	}
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

func normalizeRequestPath(r *http.Request, routePath string) {
	if routePath == "/" {
		return
	}
	// Some clients probe with a trailing slash; upstream MCP endpoints are often strict.
	if r.URL.Path == routePath+"/" {
		r.URL.Path = routePath
	}
}

func (s *Server) captureRequestBody(r *http.Request) (string, bool) {
	if r.Body == nil {
		return "", false
	}
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		return "", false
	}
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewReader(b))
	preview, truncated := capBytes(b, 16*1024)
	return strings.TrimSpace(preview), truncated
}

func capBytes(b []byte, max int) (string, bool) {
	if len(b) <= max {
		return string(b), false
	}
	return string(b[:max]), true
}

type loggingReadCloser struct {
	rc         io.ReadCloser
	limit      int
	preview    bytes.Buffer
	overflow   bool
	emitted    bool
	onComplete func(preview string, truncated bool)
}

func newLoggingReadCloser(rc io.ReadCloser, limit int, onComplete func(preview string, truncated bool)) io.ReadCloser {
	return &loggingReadCloser{
		rc:         rc,
		limit:      limit,
		onComplete: onComplete,
	}
}

func (l *loggingReadCloser) Read(p []byte) (int, error) {
	n, err := l.rc.Read(p)
	if n > 0 && l.preview.Len() < l.limit {
		remaining := l.limit - l.preview.Len()
		if n > remaining {
			_, _ = l.preview.Write(p[:remaining])
			l.overflow = true
		} else {
			_, _ = l.preview.Write(p[:n])
		}
	} else if n > 0 {
		l.overflow = true
	}
	if err == io.EOF {
		l.emit()
	}
	return n, err
}

func (l *loggingReadCloser) Close() error {
	l.emit()
	return l.rc.Close()
}

func (l *loggingReadCloser) emit() {
	if l.emitted {
		return
	}
	l.emitted = true
	if l.onComplete != nil {
		l.onComplete(strings.TrimSpace(l.preview.String()), l.overflow)
	}
}
