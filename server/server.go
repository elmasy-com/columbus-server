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
	"github.com/elmasy-com/columbus-server/server/lookup"
	"github.com/elmasy-com/columbus-server/server/stat"

	"github.com/gin-gonic/gin"
)

func GinLog(param gin.LogFormatterParams) string {

	if param.StatusCode >= 200 && param.StatusCode < 300 && param.Latency < time.Second && config.LogErrorOnly {
		return ""
	}

	return fmt.Sprintf("%s - [%s] \"%s %s\" %d %d \"%s\" %s\n%s",
		param.ClientIP,
		param.TimeStamp.Format(time.RFC1123),
		param.Method,
		param.Path,
		param.StatusCode,
		param.BodySize,
		param.Request.UserAgent(),
		param.Latency,
		param.ErrorMessage,
	)
}

// ServerRun start the http server and block.
// The server can stopped with a SIGINT.
func Run() error {

	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()

	var (
		err    error
		router = gin.New()
		quit   = make(chan os.Signal, 1)
	)

	router.Use(gin.LoggerWithFormatter(GinLog))
	router.Use(gin.Recovery())

	router.SetTrustedProxies(config.TrustedProxies)

	router.GET("/api/lookup/:domain", lookup.GetApiLookup)
	router.GET("/api/starts/:domain", lookup.GetApiStarts)
	router.GET("/api/tld/:domain", lookup.GetApiTLD)
	router.GET("/api/history/:domain", lookup.GetApiHistory)

	// router.PUT("/insert/:domain", InsertPut)

	router.GET("/api/stat", stat.GetApiStat)
	router.GET("/stat", stat.GetStat)

	router.GET("/api/tools/tld/:fqdn", ToolsTLDGet)
	router.GET("/api/tools/domain/:fqdn", ToolsDomainGet)
	router.GET("/api/tools/subdomain/:fqdn", ToolsSubdomainGet)
	router.GET("/api/tools/isvalid/:fqdn", ToolsIsValidGet)

	// Permanent Redirect
	router.GET("/lookup/:domain", Redirect)
	router.GET("/tld/:domain", Redirect)
	router.GET("/tools/tld/:fqdn", Redirect)
	router.GET("/tools/domain/:fqdn", Redirect)
	router.GET("/tools/subdomain/:fqdn", Redirect)
	router.GET("/tools/isvalid/:fqdn", Redirect)

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
