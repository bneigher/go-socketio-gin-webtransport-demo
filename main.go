package main

import (
	"log"
	"net/http"
	"os"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/zishang520/engine.io/v2/engine"
	"github.com/zishang520/engine.io/v2/types"
	"github.com/zishang520/engine.io/v2/utils"
	"github.com/zishang520/engine.io/v2/webtransport"
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

	// WebTransport start
	// WebTransport uses udp, so you need to enable the new service.
	customServer := types.NewWebServer(nil)

	// A certificate is required and cannot be a self-signed certificate.
	wts := customServer.ListenWebTransportTLS(":443", "./cert/cert.pem", "./cert/cert.key", nil, nil)

	// if _, err := os.Stat("./cert/cert.pem"); err == nil {
	//  wts := customServer.ListenWebTransportTLS(":443", "./cert/cert.pem", "./cert/cert.key", nil, nil)
	// } else {
	// 	wts := customServer.Listen(":8080", nil)
	// }

	// Here is the core logic of the WebTransport handshake.
	customServer.HandleFunc(server.Path()+"/", func(w http.ResponseWriter, r *http.Request) {
		if webtransport.IsWebTransportUpgrade(r) {
			// You need to call server.ServeHandler(nil) before this, otherwise you cannot get the Engine instance.
			server.Engine().(engine.Server).OnWebTransportSession(types.NewHttpContext(w, r), wts)
		} else {
			customServer.DefaultHandler.ServeHTTP(w, r)
		}
	})
	// WebTransport end

	if _, err := os.Stat("./cert/cert.pem"); err == nil {
		r.RunTLS(":443", "./cert/cert.pem", "./cert/cert.key")
	} else {
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
