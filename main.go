package main

import (
	"apr-backend/client"
	"apr-backend/internal/auth"
	"apr-backend/internal/controllers"
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	logger := log.Default()
	router := gin.New()
	router.Use(gin.Recovery())

	mysqlDb, err := sql.Open("mysql", "apr:bBfg3wV7mngCYvCIvX2Iqdv2X7oaivjzIOfk7KxY@/apr")
	if err != nil {
		panic(err.Error())
	}
	defer mysqlDb.Close()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sex", model.ValidateSex)
	}

	userRepo := db.NewUserRepo(mysqlDb)
	authServ := services.NewAuthService(userRepo)
	userServ := services.NewUserService(userRepo)
	authCtr := controllers.NewAuthController(authServ)
	userCtr := controllers.NewUserController(userServ)
	jwtGenerator, err := auth.NewJwtGenerator("/run/secrets/apr_rsa_private")
	if err != nil {
		return
	}

	router.Use(client.CheckAuth(jwtGenerator, client.Apr))
	authGroup := router.Group("/api/auth")
	{
		authGroup.POST("/login", authCtr.Login)
	}
	userGroup := router.Group("/api/user")
	{
		userGroup.POST("/", userCtr.RegisterUser)
	}

	srv := &http.Server{Addr: "0.0.0.0:7887", Handler: router}
	go func() {
		log.Println("server starting")
		if err := srv.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				log.Fatal(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// gracefully shutdown server
	logger.Println("service shutting down ...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}
	logger.Println("server stopped")
}
