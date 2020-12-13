package server

import "github.com/gin-gonic/gin"

type Server struct {
	router *gin.Engine
}

func NewServer() *Server {
	router := gin.Default()

	return &Server{
		router: router,
	}
}

func (self *Server) Init() {
	self.router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "wang",
		})
	})
}

func (self *Server) TestRun() {
	self.router.Run()
}

func (self *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	self.router.Run()
}
