package apiserver

// ServerRunOptions provides various option which will be used
// in creating a rest server
type ServerRunOptions struct {
	Addr                  string
	EnableTLS             bool
	CertFile              string
	KeyFile               string
	CorsAllowedOriginList []string
	CorsAllowedMethodList []string
}
