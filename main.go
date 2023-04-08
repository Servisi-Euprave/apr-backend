package main

import (
	"apr-backend/internal/auth"
	"apr-backend/internal/controllers"
	"apr-backend/internal/db"
	"apr-backend/internal/model"
	"apr-backend/internal/services"
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
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

	dbUsr := "apr"
	dbPass, err := ioutil.ReadFile("../db_pass.key")
	sqlConStr := fmt.Sprintf("%s:%s@/apr", dbUsr, strings.TrimSpace(string(dbPass)))
	mysqlDb, err := sql.Open("mysql", sqlConStr)
	if err != nil {
		logger.Println(err.Error())
		return
	}
	defer mysqlDb.Close()

	err = mysqlDb.Ping()
	if err != nil {
		logger.Println(err.Error())
		return
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("sex", model.ValidateSex)
	}

	userRepo := db.NewUserRepo(mysqlDb)
	authServ := services.NewAuthService(userRepo)
	userServ := services.NewUserService(userRepo)

	jwtGenerator, err := auth.NewJwtGenerator("../private.pem")
	authCtr := controllers.NewAuthController(authServ, jwtGenerator)
	userCtr := controllers.NewUserController(userServ)
	if err != nil {
		logger.Println(err.Error())
		return
	}

	// router.Use(client.CheckAuth(jwtGenerator, client.Apr))
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
