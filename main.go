package main

import (
	"net"
	"net/http"
	"fmt"
	"strings"
	"bufio"
	"github.com/gin-gonic/gin"
	"strconv"
	"unicode/utf8"
)

func main() {
	router := gin.Default()

	router.Static("/static", "./views")

	router.LoadHTMLGlob("views/html/*")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H {
			"title": "Dashboard",
			"sdaTemp": Temp("/dev/sda"),
			"sdbTemp": Temp("/dev/sdb"),
			"sdcTemp": Temp("/dev/sdc"),
		})
	})

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	router.Run()
}

func Temp(drive string) int {
	c, err := net.Dial("tcp", "localhost:7634")
	if err != nil {
		fmt.Println(err)
		return -1
	}

	reply, err := bufio.NewReader(c).ReadString('\n')
	status := Trim(reply)
	s := strings.Split(status, "||")
	t := 0
	for _, j := range s {
		d := strings.Split(j, "|")
		if d[0] == drive {
			t, _ = strconv.Atoi(d[2])
		}
	}

	return t
}

func Trim(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}
