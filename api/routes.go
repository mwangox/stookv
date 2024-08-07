package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net"
	"stoo-kv/config"
	"stoo-kv/internal/store"
)

func InitializeRoutes(storage store.Store, cfg *config.Config) error {
	gin.SetMode(cfg.Application.ServerLogLevel)
	r := gin.Default()
	r.Use(cors.Default())

	if err := r.SetTrustedProxies(nil); err != nil {
		return errors.Wrapf(err, "failed to set trusted proxies")
	}
	handler := NewHandler(storage, cfg)
	r.GET("/stoo-kv/:namespace/:profile/:key", handler.GetHandler)
	r.GET("/stoo-kv/:namespace/:profile", handler.GetByNamespaceAndProfileHandler)
	//r.GET("/stoo-kv", handler.GetAllHandler)
	r.POST("/stoo-kv/:namespace/:profile", handler.SetHandler)
	r.POST("/stoo-kv/secrets/:namespace/:profile", handler.SetSecretHandler)
	r.DELETE("/stoo-kv/:namespace/:profile", handler.DeleteHandler)
	r.POST("/stoo-kv/encrypt", handler.EncryptHandler)
	if cfg.Application.EnableDecryptEndpoint {
		r.POST("/stoo-kv/decrypt", handler.DecryptHandler)
	}
	return r.Run(net.JoinHostPort(cfg.Application.ServerBindingHost, cfg.Application.ServerPort))
}
