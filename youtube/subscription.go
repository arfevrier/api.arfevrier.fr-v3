package youtube

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type listSubscriptions struct {
	Items         []Subscription `json:"items"`
	NextPageToken string         `json:"nextPageToken"`
}

type Subscription struct {
	Snippet Snippet `json:"snippet"`
}

type Snippet struct {
	ResourceId resourceId `json:"resourceId"`
}

type resourceId struct {
	ChannelId string `json:"channelId"`
}

func generateSubscriptions(chSub chan string, token string, pageToken string) {
	// Prepare HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s&pageToken=%s", subscriptionsURL, pageToken), nil)
	if err != nil {
		log.Println("Error on Unmarshal.\n[ERROR] -", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("[ERROR] Error on generateSubscriptions request:", resp.StatusCode, err)
		close(chSub)
		return
	}
	defer resp.Body.Close()

	// Load subscriptions content
	data := listSubscriptions{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Error on Unmarshal: ", err)
		close(chSub)
		return
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Println("[ERROR] Error on Unmarshal: ", err)
		close(chSub)
		return
	}

	for _, element := range data.Items {
		chSub <- element.Snippet.ResourceId.ChannelId
	}

	// Recursive load with the next page token

	if data.NextPageToken != "" {
		generateSubscriptions(chSub, token, data.NextPageToken)
	}

	// If we are the first function called
	if pageToken == "" {
		close(chSub)
	}
}
