package main

import (
	"net/http"
	"hddtemp"
	"github.com/gin-gonic/gin"
)

func main() {

	hddtemp.Temp("/dev/sda");

	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H {
			"title": "Dashboard",
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}
