package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func portInUse(portNumber int) int {
	res := -1
	var outBytes bytes.Buffer
	cmdStr := fmt.Sprintf("netstat -ano -p tcp | findstr 0.0.0.0:%d", portNumber)
	cmd := exec.Command("cmd", "/c", cmdStr)
	cmd.Stdout = &outBytes
	cmd.Run()
	resStr := outBytes.String()
	r := regexp.MustCompile(`\s\d+\s`).FindAllString(resStr, -1)
	if len(r) > 0 {
		pid, err := strconv.Atoi(strings.TrimSpace(r[0]))
		if err != nil {
			res = -1
		} else {
			res = pid
		}
	}
	return res
}

var quit = make(chan os.Signal)

func main() {
	//fmt.Println(pid)
	//if pid != -1 {
	//	c := exec.Command("taskkill.exe", "/F", "/PID", strconv.Itoa(pid))
	//	if err := c.Start(); err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	//time.Sleep(2*time.Second)
	pid := portInUse(8080)
	if pid != -1 {
		http.Get("http://localhost:8080/shutdown")
	}
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		//time.Sleep(5 * time.Second)
		c.String(http.StatusOK, "Welcome Gin Server6")
	})
	router.GET("/shutdown", func(c *gin.Context) {
		quit <- syscall.SIGINT
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		var pid int
		for pid != -1 {
			pid = portInUse(8080)
		}
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
