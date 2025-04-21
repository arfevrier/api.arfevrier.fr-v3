package youtube

import (
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
)

type Feed struct {
	XMLName xml.Name `xml:"feed"`
	Entry   []Entry  `xml:"entry"`
}

type Entry struct {
	XMLName   xml.Name `xml:"entry"`
	VideoID   string   `xml:"http://www.youtube.com/xml/schemas/2015 videoId"`
	Title     string   `xml:"title"`
	Published string   `xml:"published"`
	Author    Author   `xml:"author"`
}

type Author struct {
	XMLName xml.Name `xml:"author"`
	Name    string   `xml:"name"`
}

func generateChannelVideos(wg *sync.WaitGroup, chVid chan VideoContent, channelID string) {
	defer wg.Done()

	// Prepare HTTP request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s?channel_id=%s", videosFeed, channelID), nil)
	if err != nil {
		log.Println("Error on Unmarshal.\n[ERROR] -", err)
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Println("[ERROR] Error on generateChannelVideos request:", resp.StatusCode, err)
		return
	}
	defer resp.Body.Close()

	// Load video content
	data := Feed{}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("[ERROR] Error on Unmarshal: ", err)
		return
	}

	err = xml.Unmarshal(body, &data)
	if err != nil {
		log.Println("[ERROR] Error on Unmarshal: ", err)
		return
	}

	// Send the first 2 video
	var maxNewVideo int = 2
	for i, element := range data.Entry {
		if i >= maxNewVideo {
			break
		}
		chVid <- VideoContent{
			AuthorName:    element.Author.Name,
			PublishedDate: element.Published,
			VideoID:       element.VideoID,
			VideoTitle:    element.Title,
		}
	}
}
