package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"regexp"

	"api.arfevrier.fr/v3/signal"
	"api.arfevrier.fr/v3/webconnect"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// > Input
type urlChannel struct {
	Channel string `uri:"channel" binding:"required"`
}

// newWS godoc
// @Summary Create a new WebSocket connection
// @Schemes
// @Description Create a new WebSocket connection for a given channel
// @Tags webconnect
// @Accept json
// @Produce json
// @Param channel path string true "Channel"
// @Param localdesc query string false "Local description for WebRTC"
// @Success 200 {string} string "WebSocket connection established"
// @Failure 400 {string} string "Invalid input"
// @Router /webconnect/new/{channel} [get]
func (Api API) newWS(c *gin.Context) {
	var urlChannel urlChannel

	// Get the token from REST url
	if err := c.ShouldBindUri(&urlChannel); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Validate input type
	r, err := regexp.Compile(`^.*$`)
	if err != nil || !r.MatchString(urlChannel.Channel) {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	if localDesc, ok := c.GetQuery("localdesc"); ok {
		newRtc := webconnect.NewRtc(hub, localDesc)
		hub.RegisterRtc <- newRtc
		c.JSON(200, gin.H{"websocket": fmt.Sprintf("%d", rand.Intn(1000)), "webrtc": signal.Encode(newRtc.Conn.LocalDescription())})
		return
	}

	c.JSON(200, gin.H{"websocket": "withoutRTC"})
}

// > Var
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// > Input
type urlId struct {
	Id string `uri:"id" binding:"required"`
}

// connectWS godoc
// @Summary Connect to an existing WebSocket
// @Schemes
// @Description Connect to an existing WebSocket by ID
// @Tags webconnect
// @Accept json
// @Produce json
// @Param id path string true "WebSocket ID"
// @Success 200 {string} string "Connection established"
// @Failure 400 {string} string "Invalid input"
// @Router /webconnect/connect/{id} [get]
func (Api API) connectWS(c *gin.Context) {
	var urlId urlId

	// Get the token from REST url
	if err := c.ShouldBindUri(&urlId); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Validate input type
	r, err := regexp.Compile(`^.*$`)
	if err != nil || !r.MatchString(urlId.Id) {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	// --- Start websocket connection upgrade
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	// Disconnect if client already connected
	for client := range hub.Clients {
		fmt.Printf("|> Websocket: Current client ID: %s\n", client.Id)
		if client.Id == urlId.Id {
			return
		}
	}

	// Add new client to the hub and run
	newClient := webconnect.NewClient(urlId.Id, hub, conn)
	hub.Register <- newClient
	newClient.Run()
	hub.Unregister <- newClient
}
