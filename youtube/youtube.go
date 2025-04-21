package youtube

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
)

const (
	subscriptionsURL string = "https://www.googleapis.com/youtube/v3/subscriptions?part=snippet&mine=true&maxResults=50"
	videosFeed       string = "https://www.youtube.com/feeds/videos.xml"
	watchUrl         string = "https://www.youtube.com/watch"
)

func GetSubscriptionsVideosList(token string) []VideoContent {
	var videosList []VideoContent

	// Routine init
	var wgSub sync.WaitGroup
	var wgVid sync.WaitGroup
	chSub := make(chan string)
	chVid := make(chan VideoContent)

	// Generate the subscribers list
	wgVid.Add(1)
	go func() {
		defer wgVid.Done()
		for elem := range chVid {
			videosList = append(videosList, elem)
		}
	}()
	go generateSubscriptions(chSub, token, "")
	for elem := range chSub {
		wgSub.Add(1)
		go generateChannelVideos(&wgSub, chVid, elem)
	}
	wgSub.Wait()
	close(chVid)
	wgVid.Wait()

	// Return list of videos
	return videosList
}

func GenerateDownloadStream(w io.Writer, linkType string, linkID string) bool {
	var cmd *exec.Cmd

	if linkType == "video" {
		cmd = exec.Command("/home/ubuntu/.local/bin/yt-dlp", "-f", "b", fmt.Sprintf("%s?v=%s", watchUrl, linkID), "-o", "-")
	} else {
		cmd = exec.Command("/home/ubuntu/.local/bin/yt-dlp", "-f", "ba", fmt.Sprintf("%s?v=%s", watchUrl, linkID), "-o", "-")
	}
	cmd.Stdout = w
	cmd.Stderr = os.Stdout
	if err := cmd.Run(); err != nil {
		fmt.Println("[api-v3] YouTube GenerateDownloadStream() error:", err)
	}
	return false
}
