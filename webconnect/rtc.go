package webconnect

import (
	"fmt"
	"time"

	"api.arfevrier.fr/v3/signal"
	"github.com/pion/webrtc/v3"
)

// Client is a middleman between the websocket connection and the hub.
type Rtc struct {
	hub *Hub
	// The websocket connection.
	Conn *webrtc.PeerConnection
	// Buffered channel of outbound messages.
	send chan []byte
}

func NewRtc(newHub *Hub, webOffer string) *Rtc {
	creationRtc := &Rtc{
		hub:  newHub,
		send: make(chan []byte, 10),
	}
	var err error

	// Prepare the configuration
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			{
				URLs: []string{"stun:stun.l.google.com:19302"},
			},
		},
	}

	// Create a new RTCPeerConnection
	creationRtc.Conn, err = webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	// defer func() {
	// 	if cErr := creationRtc.conn.Close(); cErr != nil {
	// 		fmt.Printf("cannot close peerConnection: %v\n", cErr)
	// 	}
	// }()

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	creationRtc.Conn.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		fmt.Printf("|> RTC: Peer Connection State has changed: %s\n", s.String())
		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			fmt.Println("|> RTC: Peer Connection has gone to failed exiting")
			creationRtc.hub.unregisterRtc <- creationRtc
			return
		}
	})

	// Register data channel creation handling
	creationRtc.Conn.OnDataChannel(func(d *webrtc.DataChannel) {
		fmt.Printf("|> RTC: New DataChannel %s %d\n", d.Label(), d.ID())
		// Register channel opening handling
		d.OnOpen(func() {
			fmt.Printf("|> RTC: Data channel '%s'-'%d' open. Random messages now be sent\n", d.Label(), d.ID())
			var count int = 0
			for {
				select {
				case message := <-creationRtc.send:
					// Send the message as text
					sendErr := d.SendText(string(message))
					if sendErr != nil {
						creationRtc.hub.unregisterRtc <- creationRtc
						return
					}
				case <-time.After(time.Millisecond * 16):
					count += 1
					sendErr := d.SendText(fmt.Sprintf("%06d:%s", count, time.Now()))
					if sendErr != nil {
						fmt.Println("|> RTC: Send error, close RTC client")
						creationRtc.hub.unregisterRtc <- creationRtc
						return
					}
				}
			}
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			fmt.Printf("|> RTC: Message from DataChannel '%s': '%s'\n", d.Label(), string(msg.Data))
		})
	})

	// Wait for the offer to be pasted
	offer := webrtc.SessionDescription{}
	signal.Decode(webOffer, &offer)

	// Set the remote SessionDescription
	err = creationRtc.Conn.SetRemoteDescription(offer)
	if err != nil {
		panic(err)
	}

	// Create an answer
	answer, err := creationRtc.Conn.CreateAnswer(nil)
	if err != nil {
		panic(err)
	}

	// Create channel that is blocked until ICE Gathering is complete
	gatherComplete := webrtc.GatheringCompletePromise(creationRtc.Conn)

	// Sets the LocalDescription, and starts our UDP listeners
	err = creationRtc.Conn.SetLocalDescription(answer)
	if err != nil {
		panic(err)
	}

	// Block until ICE Gathering is complete, disabling trickle ICE
	// we do this because we only can exchange one signaling message
	// in a production application you should exchange ICE Candidates via OnICECandidate
	<-gatherComplete

	return creationRtc
}
