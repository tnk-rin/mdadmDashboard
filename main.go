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
	"syscall"
	"math"
)

type DiskStatus struct {
	Total uint64 `json:"Total"`
	Used uint64 `json:"Used"`
	Free uint64 `json:"Free"`
}

const (
	B  = 1
	KB = 1024
	MB = KB * KB
	GB = MB * KB
)


/*
 * TODO:	add config tab on main page where user can pick 
 *			a mounted folder, using https://pkg.go.dev/k8s.io/kubernetes/pkg/util/mount
 *			to find device names for the thermals
 */


func main() {
	router := gin.Default()
	router.Static("/static", "./views")
	router.LoadHTMLGlob("views/html/*")
	router.GET("/", func(c *gin.Context) {
		sda := DiskUsage("/mnt/ssd_wind/")
		sdb := DiskUsage("/mnt/2tb_linux/")
		sdc := DiskUsage("/mnt/2tb_wind/")

		c.HTML(http.StatusOK, "index.html", gin.H {
			"title": "Dashboard",
			"sdaTemp": Temp("/dev/sda"),
			"sdbTemp": Temp("/dev/sdb"),
			"sdcTemp": Temp("/dev/sdc"),
			"sdaUsed": (math.Round((float64(sda.Used)/float64(GB)) * 100) / 100),
			"sdbUsed": (math.Round((float64(sdb.Used)/float64(GB)) * 100) / 100),
			"sdcUsed": (math.Round((float64(sdc.Used)/float64(GB)) * 100) / 100),
			"sdaTotal": (math.Round((float64(sda.Total)/float64(GB)) * 100) / 100),
			"sdbTotal": (math.Round((float64(sdb.Total)/float64(GB)) * 100) / 100),
			"sdcTotal": (math.Round((float64(sdc.Total)/float64(GB)) * 100) / 100),

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

func DiskUsage(path string) (disk DiskStatus) {
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(path, &fs)

	if err != nil {
		fmt.Println(err)
		return
	}
	disk.Total = fs.Blocks * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.Total - disk.Free
	return
}
