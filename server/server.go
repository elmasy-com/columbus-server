package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/elmasy-com/columbus-server/config"
	"github.com/gin-gonic/gin"
)

// ServerRun start the http server and block.
// The server can stopped with a SIGINT.
func Run() error {

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	var (
		err    error
		router = gin.Default()
		quit   = make(chan os.Signal, 1)
	)

	router.SetTrustedProxies(config.TrustedProxies)

	router.GET("/lookup/:domain", LookupGet)
	router.PUT("/insert/:domain", InsertPut)
	router.GET("/openapi.yaml", StaticOpenApiYamlGet)

	router.GET("/user", UserGet)
	router.PUT("/user", UserPut)
	router.DELETE("/user", UserDelete)
	router.PATCH("/user", UserPatch)
	router.GET("/users", UsersGet)
	router.PATCH("user/other", UserOtherPatch)

	router.GET("/stat", StatGet)

	srv := &http.Server{
		Addr:    config.Address,
		Handler: router,
	}

	go func() {
		if config.SSLCert != "" && config.SSLKey != "" {
			err = srv.ListenAndServeTLS(config.SSLCert, config.SSLKey)
		} else {
			err = srv.ListenAndServe()
		}
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Fprintf(os.Stderr, "HTTP Server failed: %s\n", err)
			os.Exit(1)
		}
	}()

	signal.Notify(quit, os.Interrupt, syscall.SIGINT)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return srv.Shutdown(ctx)
}
