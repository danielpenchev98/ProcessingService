package main

import (
	"context"
	"danielpenchev98/http-job-processing-service/api/rest"
	"danielpenchev98/http-job-processing-service/pkg"
	"danielpenchev98/http-job-processing-service/pkg/algorithm"
	"danielpenchev98/http-job-processing-service/pkg/script"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/pkg/errors"

	"github.com/gin-gonic/gin"
)

const (
	hostParamName = "HOST"
	portParamName = "PORT"
)

type ServerConfig struct {
	Host string
	Port int
}

func main() {
	serverCfg, err := getServerConfig()
	if err != nil {
		log.Fatalf("Proble with the server config. Reason %s", err)
	}

	httpServer := createHttpServer(serverCfg.Host, serverCfg.Port)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(errors.Wrapf(err, "server listen-and-serve failed"))
		}
	}()

	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT)
	<-done

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		panic(errors.Wrapf(err, "failed to shutdown server"))
	}

	<-ctx.Done()
}

func getServerConfig() (ServerConfig, error) {
	portStr := os.Getenv(portParamName)
	if portStr == "" {
		return ServerConfig{}, errors.Errorf("Please set %s env variable", portParamName)
	}

	portNum, err := strconv.Atoi(portStr)
	if err != nil {
		return ServerConfig{}, errors.Errorf("The env variable %s has illegal port number", portParamName)
	}

	return ServerConfig{
		Host: os.Getenv(hostParamName),
		Port: portNum,
	}, nil
}

func createHttpServer(host string, port int) *http.Server {
	var router = gin.Default()

	taskScheduler := pkg.NewTaskScheduler(algorithm.TopologicalSort{}, algorithm.DependencyGraphCreator{})
	contentCretor := pkg.NewContentCreator(script.BashContentBuilderCreator{})
	processorEndpoint := rest.NewTaskProcessorEndpoint(taskScheduler, contentCretor)

	v1 := router.Group("/v1")
	{
		v1.GET("/healthcheck", rest.CheckHealth)
		v1.POST("/processing-plan", processorEndpoint.CreateProcessingPlan)
		v1.POST("/bash", processorEndpoint.GenerateBashContent)
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", host, port),
		Handler: router,
	}

	return httpServer
}
