package server

import (
	"Pixivel/internal/config"
	"Pixivel/internal/pixivel"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	router  *gin.Engine
	pixivel *pixivel.Pixivel
}

func Decorator(h gin.HandlerFunc, decors ...func(gin.HandlerFunc) gin.HandlerFunc) gin.HandlerFunc {
	for i := range decors {
		d := decors[len(decors)-1-i]
		h = d(h)
	}
	return h
}

func NewServer() *Server {
	router := gin.Default()
	handler := pixivel.GetHandler()
	return &Server{
		router:  router,
		pixivel: handler,
	}
}

func (self *Server) Init() {
	v1 := self.router.Group("/v1")
	{
		v1.GET("/illust/:id", self.pixivel.Cache(func(c *gin.Context) {
			id := config.Atoi(c.Param("id"))
			response := self.pixivel.SingleIllust(id)
			c.JSON(200, response)
		}, 60))
	}
}

func (self *Server) TestRun() {
	self.router.Run(":8000")
}

func (self *Server) Run() {
	//gin.SetMode(gin.ReleaseMode)
	srv := &http.Server{
		Addr:    ":8000",
		Handler: self.router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
		self.pixivel.Close()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
}
