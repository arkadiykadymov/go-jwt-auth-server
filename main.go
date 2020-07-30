package main

import (
	"dev/jwt-auth-server/service"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func handleRequest() {

	// Init Router
	router := gin.New()

	router.POST("/register", service.RegisterUser)
	router.POST("/login/:uuid", service.LoginUser)
	router.POST("/refresh-token", service.RefreshAccessToken)
	router.POST("/delete-refresh-token", service.DeleteRefreshToken)
	router.POST("/delete-refresh-tokens/:uuid", service.DeleteAllRefreshTokens)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	s.ListenAndServe()

}

func main() {
	handleRequest()
}
