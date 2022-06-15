package main

import (
	"net"
	"net/http"
	"fmt"
	"strings"
	"bufio"
	"github.com/gin-gonic/gin"
	Mount "k8s.io/mount-utils"
	"strconv"
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


func main() {
	mount1 := "/mnt/ssd_wind/"
	mount2 := "/mnt/2tb_linux/"
	mount3 := "/mnt/2tb_wind/"

	router := gin.Default()
	router.Static("/static", "./views")
	router.LoadHTMLGlob("views/html/*")
	router.GET("/", func(c *gin.Context) {
		sda := DiskUsage(mount1)
		sdb := DiskUsage(mount2)
		sdc := DiskUsage(mount3)

		c.HTML(http.StatusOK, "index.html", gin.H {
			"title": "Dashboard",
			"sdaTemp": Temp(DeviceFromMount(mount1)),
			"sdbTemp": Temp(DeviceFromMount(mount2)),
			"sdcTemp": Temp(DeviceFromMount(mount3)),
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
	status := TrimFirst(reply)
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

func TrimFirst(s string) string {
	return s[1:]
}

func TrimLast(s string) string {
	return s[:len(s)-1]
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

func DeviceFromMount(path string) (device string) {
	interf := Mount.New("/bin/mount")
	device, _, err := Mount.GetDeviceNameFromMount(interf, path)
	device = TrimLast(device)
	if err != nil {
		device := "Error finding device from mountpoint: " + err.Error()
		fmt.Println(device)
	}
	return
}
