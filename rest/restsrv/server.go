package restsrv

import (
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/mediacoin-pro/core/chain/bcstore"
	"github.com/mediacoin-pro/core/common/consts"
	"github.com/mediacoin-pro/core/common/xlog"
)

type Server struct {
	cfg *Config
	bc  *bcstore.ChainStorage
}

func StartServer(cfg *Config, bc *bcstore.ChainStorage) {
	s := NewService(cfg, bc)
	s.Start()
}

func NewService(cfg *Config, bc *bcstore.ChainStorage) *Server {
	return &Server{
		cfg: cfg,
		bc:  bc,
	}
}

func (s *Server) Start() {
	server := &http.Server{
		Addr:           s.cfg.HTTPConn,
		Handler:        s,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: int(64 * consts.MiB),
	}
	if err := server.ListenAndServe(); err != nil {
		xlog.Panic(err)
	}
}

func (s *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("http-PANIC: %v", r)
			xlog.Error.Printf("http> ServeHTTP-PANIC: %v\n%s", err, string(debug.Stack()))
			rw.WriteHeader(http.StatusInternalServerError)
			//http.Error(rw, err.Error(), http.StatusInternalServerError)
		}
	}()

	ctx := newContext(s, req, rw)
	ctx.Exec()
}
