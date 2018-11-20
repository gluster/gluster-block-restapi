package apiserver

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"net"
	"net/http"

	"github.com/gluster/gluster-block-restapi/glusterblockrestd/handlers"

	muxhandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// Server represents gluster-block rest server
type Server struct {
	isTLS                 bool
	certFile              string
	keyFile               string
	addr                  string
	srv                   *http.Server
	srvMux                http.Handler
	corsAllowedOriginList []string
	corsAllowedMethodList []string
	middlewares           []func(http.Handler) http.Handler
}

// NewServer returns a gluster-block rest server
func NewServer(conf *ServerRunOptions) *Server {
	server := &Server{
		srv:                   &http.Server{},
		srvMux:                mux.NewRouter(),
		isTLS:                 conf.EnableTLS,
		certFile:              conf.CertFile,
		keyFile:               conf.KeyFile,
		addr:                  conf.Addr,
		corsAllowedOriginList: conf.CorsAllowedOriginList,
		corsAllowedMethodList: conf.CorsAllowedMethodList,
	}
	server.AddMiddleware(muxhandlers.CORS(
		muxhandlers.AllowedOrigins(conf.CorsAllowedOriginList),
		muxhandlers.AllowedMethods(conf.CorsAllowedMethodList),
	))
	return server
}

// AddMiddleware adds http middleware to gluster-block rest server.
// eg.. we can add middleware like authentication, tracing ..etc
func (s *Server) AddMiddleware(middleware ...func(http.Handler) http.Handler) {
	s.middlewares = append(s.middlewares, middleware...)
}

// Run starts gluster-block rest server
func (s *Server) Run(errCh chan<- error) {
	s.registerRoutes()
	listener, err := s.listener()
	if err != nil {
		errCh <- err
		return
	}
	s.srv.Handler = s.srvMux
	log.WithField("address", listener.Addr().String()).Info("starting glusterblock api server")
	errCh <- s.srv.Serve(listener)
}

// Stop stops gluster-block rest server
func (s *Server) Stop() error {
	log.Info("stopping glusterblock api server")
	return s.srv.Shutdown(context.Background())
}

func (s *Server) registerRoutes() {
	s.srvMux = handlers.NewRoutes()

	for i := len(s.middlewares) - 1; i >= 0; i-- {
		s.srvMux = s.middlewares[i](s.srvMux)
	}
}

func (s *Server) listener() (net.Listener, error) {
	listener, err := net.Listen("tcp", s.addr)

	if err != nil {
		return nil, err
	}

	if !s.isTLS {
		return listener, nil
	}

	certificate, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		return nil, err
	}

	tlsConf := &tls.Config{
		MinVersion:   tls.VersionTLS12, // force TLS 1.2
		Certificates: []tls.Certificate{certificate},
		Rand:         rand.Reader,
	}

	return tls.NewListener(listener, tlsConf), nil
}
