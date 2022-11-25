package fasthttp

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/reuseport"
)

const (
	defaultShutdownTimeout = 15
)

type service struct {
	cfg config

	a app
	l logger
	m middleware

	listenHttp  net.Listener
	listenHttps net.Listener
}

// NewService return fasthttp service object.
func NewService(cfg config, a app, l logger) (*service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("incoming config object is nil")
	}
	if a == nil {
		return nil, fmt.Errorf("incomming app object is nil")
	}
	if l == nil {
		return nil, fmt.Errorf("incoming logger object is nil")
	}

	if cfg.GetShutdownTimeout() == 0 {
		cfg.SetShutdownTimeout(defaultShutdownTimeout)
	}

	listenHttp, err := reuseport.Listen("tcp4", ":"+strconv.Itoa(cfg.GetPortHTTP()))
	if err != nil {
		return nil, err
	}
	listenHttps, err := reuseport.Listen("tcp4", ":"+strconv.Itoa(cfg.GetPortHTTPS()))
	if err != nil {
		return nil, err
	}

	return &service{
		cfg: cfg,

		a: a,
		l: l,
		m: newMiddleware(l),

		listenHttp:  newGracefulListener(listenHttp, time.Second*time.Duration(cfg.GetShutdownTimeout())),
		listenHttps: newGracefulListener(listenHttps, time.Second*time.Duration(cfg.GetShutdownTimeout())),
	}, nil
}

// Start service.
func (s *service) Start(ctx context.Context) error {
	server := fasthttp.Server{
		Handler: s.m.base(func(ctx *fasthttp.RequestCtx) {
			// execute app route.
			if s.a.Route(ctx) {
				return
			}

			// default service route.
			switch string(ctx.Path()) {
			default:
				ctx.Error("unsupported path", fasthttp.StatusNotFound)
			}
		}),
		Logger: s.l,
	}
	errCh := make(chan error, 1)

	if s.cfg.GetPortHTTP() != 0 {
		go func(server *fasthttp.Server) {
			s.l.Infof("starting HTTP server", map[string]string{"addr": "http://" + s.listenHttp.Addr().String()})
			errCh <- server.Serve(s.listenHttp)
		}(&server)
	}
	if s.cfg.GetPortHTTPS() != 0 {
		cert, err := os.ReadFile(s.cfg.GetPemPath())
		if err != nil {
			errCh <- err
		}

		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)

		pair, err := tls.LoadX509KeyPair(s.cfg.GetCertPath(), s.cfg.GetKeyPath())
		if err != nil {
			errCh <- err
		}
		tlsConfig := &tls.Config{
			ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{pair},
			ClientCAs:    certPool,
		}
		go func(server *fasthttp.Server) {
			s.l.Infof("starting HTTPS server", map[string]string{"addr": "https://" + s.listenHttps.Addr().String()})
			lnTls := tls.NewListener(s.listenHttps, tlsConfig)
			errCh <- server.Serve(lnTls)
		}(&server)
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errCh:
			return err
		}
	}
}

// Shutdown service.
func (s *service) Shutdown() error {
	_ = s.listenHttps.Close()
	_ = s.listenHttp.Close()
	return nil
}

// GetName ...
func (s *service) GetName() string {
	return fmt.Sprintf("fasthttp_%s", s.a.Name())
}
