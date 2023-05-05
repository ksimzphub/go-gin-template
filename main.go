package main

import (
	"context"
	"fmt"
	"go-gin-template/dao/mysql"
	"go-gin-template/logger"
	"go-gin-template/routes"
	"go-gin-template/settings"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	// Get the path to the config file on the command line
	if len(os.Args) < 2 {
		fmt.Println("need config file.eg: bluebell ./config.yaml")
		return
	}

	// load config file
	if err := settings.Init(os.Args[1]); err != nil {
		fmt.Printf("settings init failed, err: %v\n", err)
		return
	}

	// load logger
	if err := logger.Init(settings.Conf.LogConfig); err != nil {
		fmt.Printf("logger init failed, err: %v\n", err)
		return
	}
	defer zap.L().Sync()

	// connect mysql
	if err := mysql.Init(settings.Conf.MysqlConfig); err != nil {
		fmt.Printf("mysql connect failed, err: %v\n", err)
		return
	}
	defer mysql.Close()

	// connect redis
	// if err := redis.Init(settings.Conf.RedisConfig); err != nil {
	// 	fmt.Printf("redis connect failed, err: %v\n", err)
	// 	return
	// }
	// defer redis.Close()

	// registerd routes
	r := routes.Setup()

	// start server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", viper.GetInt("app.port")),
		Handler: r,
	}
	// Open a goroutine to start the service
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	// and set a 5-second timeout for the shutdown operation
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	zap.L().Info("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Gracefully close the service within 5 seconds (finish processing the unprocessed requests and then close the service)
	// and time out after 5 seconds
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("Server Shutdown: ", zap.Error(err))
	}
	zap.L().Info("Server exiting")
}
