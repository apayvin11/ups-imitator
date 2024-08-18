package apiserver

import (
	"io"
	"os"

	_ "github.com/alex11prog/ups-imitator/docs"
	"github.com/alex11prog/ups-imitator/internal/app/imitator"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type server struct {
	router   *gin.Engine
	imitator *imitator.Imitator
}

func newServer(imitator *imitator.Imitator) *server {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	gin.DefaultWriter = io.MultiWriter(os.Stdout)
	r.Use(gin.Recovery()) // to recover gin automatically

	s := &server{
		router:   r,
		imitator: imitator,
	}
	s.configureRouter()

	return s
}

func StartServer(bindAddr string, imitator *imitator.Imitator) error {
	s := newServer(imitator)
	return s.router.Run(bindAddr)
}

func (s *server) configureRouter() {
	s.router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	subRouter_imitator := s.router.Group("/imitator")
	subRouter_imitator.GET("/mode", s.handlerGetMode)
	subRouter_imitator.PUT("/mode", s.handlerUpdateMode)
	subRouter_imitator.GET("/ups", s.handlerGetAllUpsParams)
 	subRouter_imitator.PATCH("/ups/params", s.handlerUpdateUpsParams)
	subRouter_imitator.PATCH("/ups/:bat_id", s.handlerUpdateBattery) 
	subRouter_imitator.PATCH("/ups/alarms", s.handlerUpdateAlarms) 
}

func (s *server) errorResponse(c *gin.Context, code int, err error) {
	c.AbortWithStatusJSON(code, errorResponse{err.Error()})
}

type statusBody struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}
