package config

// Config is the root gateway configuration schema.
type Config struct {
	APIVersion string       `yaml:"apiVersion"`
	Kind       string       `yaml:"kind"`
	Gateway    Gateway      `yaml:"gateway"`
	Auth       AuthDefaults `yaml:"auth"`
	Servers    []Server     `yaml:"servers"`
	Routes     []Route      `yaml:"routes"`
}

// Gateway contains listener and runtime options.
type Gateway struct {
	Name       string `yaml:"name"`
	ListenAddr string `yaml:"listenAddr"`
	AdminAddr  string `yaml:"adminAddr"`
	LogLevel   string `yaml:"logLevel"`
}

// AuthDefaults sets secure-by-default behavior.
type AuthDefaults struct {
	RequireAuth bool `yaml:"requireAuth"`
}

// Server defines an MCP upstream.
type Server struct {
	Name      string   `yaml:"name"`
	Transport string   `yaml:"transport"` // http or stdio
	URL       string   `yaml:"url,omitempty"`
	Command   string   `yaml:"command,omitempty"`
	Args      []string `yaml:"args,omitempty"`
}

// Route maps a public path to an upstream server.
type Route struct {
	Name   string      `yaml:"name"`
	Path   string      `yaml:"path"`
	Server string      `yaml:"server"`
	Auth   *RouteAuth  `yaml:"auth,omitempty"`
	Policy RoutePolicy `yaml:"policy"`
}

// RouteAuth allows per-route auth overrides.
type RouteAuth struct {
	Type       string   `yaml:"type"` // apiKey, jwt, none
	Require    *bool    `yaml:"require,omitempty"`
	HeaderName string   `yaml:"headerName,omitempty"`
	APIKeys    []string `yaml:"apiKeys,omitempty"`
	Issuer     string   `yaml:"issuer,omitempty"`
	Audience   string   `yaml:"audience,omitempty"`
}

// RoutePolicy contains baseline traffic control settings.
type RoutePolicy struct {
	TimeoutMs    int `yaml:"timeoutMs"`
	RetryCount   int `yaml:"retryCount"`
	RateLimitRPS int `yaml:"rateLimitRps"`
}
