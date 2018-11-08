package http

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/Waitfantasy/tmq/message/manager"
	"io/ioutil"
	"net/http"
)

type Config struct {
	Addr       string
	EnableTLS  bool
	CaFile     string
	CertFile   string
	KeyFile    string
	ClientAuth bool
}

type MqServer struct {
	cfg *Config
	api *api
}

func New(c *Config, manager *manager.Manager) *MqServer {
	return &MqServer{
		cfg: c,
		api: &api{
			manager: manager,
		},
	}
}

func (s *MqServer) Run() error {
	s.api.register()

	if s.cfg.EnableTLS {
		if s.cfg.ClientAuth {
			data, err := ioutil.ReadFile(s.cfg.CaFile)
			if err != nil {
				return fmt.Errorf("failed to read ca certificate: %v\n", err)
			}

			pool := x509.NewCertPool()
			if ok := pool.AppendCertsFromPEM(data); !ok {
				return fmt.Errorf("add ca certificate failed.\n")
			}

			server := &http.Server{
				Addr:    s.cfg.Addr,
				Handler: s.api.gin,
				TLSConfig: &tls.Config{
					ClientCAs:  pool,
					ClientAuth: tls.RequireAndVerifyClientCert,
				},
			}
			return server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		} else {
			server := &http.Server{
				Addr:    s.cfg.Addr,
				Handler: s.api.gin,
			}
			return server.ListenAndServeTLS(s.cfg.CertFile, s.cfg.KeyFile)
		}
	} else {
		server := &http.Server{
			Addr:    s.cfg.Addr,
			Handler: s.api.gin,
		}
		return server.ListenAndServe()
	}
}
