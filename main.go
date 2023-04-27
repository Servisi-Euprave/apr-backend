// Package classification APR API.
//
// This application is used to register companies.
//
//	Schemes: http
//	Host: apr
//	BasePath: /
//	Version: 1.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	SecurityDefinitions:
//	bearerAuth:
//	     type: apiKey
//	     in: header
//	     name: Authorization
//
// swagger:meta
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

	"github.com/gin-contrib/cors"
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

	mysqlAddr, ok := os.LookupEnv("MYSQL_ADDR")
	if !ok {
		logger.Fatal("MYSQL_ADDR env variable was not set")
	}
	dbPassFile, ok := os.LookupEnv("DB_PASS_FILE")
	if !ok {
		logger.Fatal("DB_PASS env variable was not set")
	}
	rsaKeyFile, ok := os.LookupEnv("RSA_KEY_FILE")
	if !ok {
		logger.Fatal("RSA_KEY_FILE env variable was not set")
	}

	dbPass, err := ioutil.ReadFile(dbPassFile)
	sqlConStr := fmt.Sprintf("%s:%s@tcp(%s)/apr", dbUsr, strings.TrimSpace(string(dbPass)), mysqlAddr)
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

	userRepo := db.NewPersonRepo(mysqlDb)
	comRepo := db.NewCompanyRepository(mysqlDb, userRepo)
	authServ := services.NewAuthService(comRepo)

	privateKey, err := auth.ReadRSAPrivateKeyFromFile(rsaKeyFile)
	if err != nil {
		logger.Println(err.Error())
		return
	}
	jwtGenerator := auth.NewJwtGenerator(privateKey)
	authCtr := controllers.NewAuthController(authServ, jwtGenerator)

	comServ := services.NewCompanyService(comRepo)
	comCtr := controllers.NewCompanyController(comServ, jwtGenerator)

	nstjRepo := db.NewNstjRepository(mysqlDb)
	nstjService := services.NewNstjService(nstjRepo)
	nstjCtr := controllers.NewNstjController(nstjService)

	router.POST("/api/auth/login/", authCtr.Login)
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200", "http://localhost:4201", "http://localhost:4202"},
		AllowMethods:     []string{"POST", "GET"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	comGroup := router.Group("/api/company/")
	{
		comGroup.POST("/", comCtr.CreateCompany)
		comGroup.GET("/", comCtr.FindCompanies)
		comGroup.GET("/:pib", comCtr.FindOne)
	}
	nstjGroup := router.Group("/api/nstj/")
	{
		nstjGroup.GET("/", nstjCtr.FindAll)
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
