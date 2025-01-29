package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/quic-go/quic-go/http3"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/socket.io/v2/socket"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("./web/html/*")
	r.StaticFile("/favicon.ico", "./web/favicon.ico")
	r.Static("/js", "./web/js")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	r.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/js", "/css", "/favicon.ico"},
	}))

	opts := socket.DefaultServerOptions()
	opts.SetTransports(types.NewSet("polling", "websocket", "webtransport")) // added webtransport
	// opts.SetTransports(types.NewSet("polling", "websocket"))
	opts.SetAllowEIO3(true)
	opts.SetCors(&types.Cors{
		Origin:      "*",
		Credentials: false,
	})

	server := socket.NewServer(nil, opts)

	r.GET("/socket.io/*any", gin.WrapH(server.ServeHandler(opts)))
	r.POST("/*any", gin.WrapH(server.ServeHandler(opts)))

	SetupSocketHandler(server)

	if _, err := os.Stat("./cert/cert.pem"); err == nil {
		log.Println(fmt.Sprint("Starting HTTP/3 server on :443"))
		err = http3.ListenAndServeTLS(":443", "./cert/cert.pem", "./cert/cert.key", r.Handler())
		if err != nil {
			log.Fatalf("HTTP/3 server failed: %v", err)
		}
		// r.RunTLS(":443", "./cert/cert.pem", "./cert/cert.key")
	} else {
		log.Println(fmt.Sprint("Starting HTTP/2 server on :80"))
		r.Run()
	}
}

func SetupSocketHandler(server *socket.Server) {
	namespace := server.Of(regexp.MustCompile(`/\w+`), nil)

	namespace.On("connection", func(clients ...interface{}) {
		utils.Log().Success("On Connect")
		client := clients[0].(*socket.Socket)

		client.On("message", func(msg ...any) {
			log.Println("Received message:", msg)
			client.Emit("message back", msg...)
		})
	})
}
