package main

import (
	"net/http"
	"time"

	"api.arfevrier.fr/v3/webconnect"
	"github.com/gin-gonic/gin"

	// Swagger documentation generation
	_ "api.arfevrier.fr/v3/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var hub = webconnect.NewHub()

func CORS(c *gin.Context) {
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
}

// @title           API arfevrier.fr
// @version         3.0
// @host      api.arfevrier.fr
// @BasePath  /v3
// @Schemes   https
func main() {
	// --- Startup code for api.arfevrier.fr/v3 ---
	// Enable production run and create gin framework
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Run hub for websocket and webRTC connection
	go hub.Run()

	// Adding routes
	router.GET("/", func(c *gin.Context) { c.Redirect(http.StatusPermanentRedirect, "/v3/index.html") })
	for _, path := range []string{
		"/index.html",
		"/swagger-ui.css",
		"/swagger-ui-bundle.js",
		"/swagger-ui-standalone-preset.js",
		"/favicon-32x32.png",
		"/favicon-16x16.png",
		"doc.json",
	} {
		router.GET(path, ginSwagger.WrapHandler(swaggerfiles.Handler))
	}
	router.GET("/youtube/subscriptions/:token", CORS, Limiter{}.mw(1, 30*time.Second), API{}.youtubeSubscriptions)
	router.GET("/youtube/download/:type/:id", CORS, Limiter{}.mw(4, 30*time.Second), API{}.youtubeDownload)
	router.GET("/webconnect/new/:channel", CORS, Limiter{}.mw(4, 30*time.Second), API{}.newWS)
	router.GET("/webconnect/connect/:id", CORS, Limiter{}.mw(4, 30*time.Second), API{}.connectWS)
	router.GET("/bitcoin/price", CORS, Limiter{}.mw(1, 30*time.Second), API{}.bitcoinPrice)

	// Enable run on localhost port 1239, called by reverse proxy apache2
	router.Run("127.0.0.1:1239")
}
