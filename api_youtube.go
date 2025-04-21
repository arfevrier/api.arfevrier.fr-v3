package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"

	"api.arfevrier.fr/v3/crypto"
	"api.arfevrier.fr/v3/youtube"
	"github.com/gin-gonic/gin"
)

// > Input
type urlToken struct {
	Token string `uri:"token" binding:"required"`
}

// > Output
type SubscriptionsResult struct {
	VideosList []youtube.VideoContent `json:"subscriptions"`
}

// youtubeSubscriptions godoc
// @Summary Get YouTube subscriptions
// @Schemes
// @Description Fetch a list of YouTube subscription videos for a given token
// @Tags youtube
// @Accept json
// @Produce json
// @Param token path string true "Token"
// @Success 200 {object} SubscriptionsResult
// @Failure 400 {string} string "Invalid input"
// @Failure 500 {string} string "No content found for this token"
// @Router /youtube/subscriptions/{token} [get]
func (Api API) youtubeSubscriptions(c *gin.Context) {
	var urlToken urlToken

	// Get the token from REST url
	if err := c.ShouldBindUri(&urlToken); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Validate token
	r, err := regexp.Compile(`^[a-zA-Z0-9_=-]+$`)
	if err != nil || !r.MatchString(urlToken.Token) {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	// Decrypt token using Fernet
	urlToken.Token = crypto.Decrypt(urlToken.Token)
	if len(urlToken.Token) == 0 {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	// Start generate video list
	videosList := youtube.GetSubscriptionsVideosList(urlToken.Token)

	// Return json result if list contains video
	if videosList == nil {
		c.JSON(500, gin.H{"msg": "No content found for this token"})
	} else {
		c.JSON(http.StatusOK, SubscriptionsResult{
			VideosList: videosList,
		})
	}
}

// > Input
type urlTypeId struct {
	Type string `uri:"type" binding:"required"`
	ID   string `uri:"id" binding:"required"`
}

// youtubeDownload godoc
// @Summary Download YouTube video or audio
// @Schemes
// @Description Download a YouTube video or audio file by type and ID
// @Tags youtube
// @Accept json
// @Produce application/octet-stream
// @Param type path string true "Type (video or audio)"
// @Param id path string true "Video ID"
// @Success 200 {string} string "File stream"
// @Failure 400 {string} string "Invalid input"
// @Router /youtube/download/{type}/{id} [get]
func (Api API) youtubeDownload(c *gin.Context) {
	var urlTypeId urlTypeId

	// Get the token from REST url
	if err := c.ShouldBindUri(&urlTypeId); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Validate input type
	r, err := regexp.Compile(`^video|audio$`)
	if err != nil || !r.MatchString(urlTypeId.Type) {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	// Validate input video ID
	r, err = regexp.Compile(`^[a-zA-Z0-9_\-]{1,25}$`)
	if err != nil || !r.MatchString(urlTypeId.ID) {
		c.JSON(400, gin.H{"msg": "Invalid input"})
		return
	}

	// Set header depending of the
	if urlTypeId.Type == "video" {
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.mp4"`, urlTypeId.ID))
	} else {
		c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.webm"`, urlTypeId.ID))
	}

	c.Stream(func(w io.Writer) bool {
		fmt.Println("Download stream")
		return youtube.GenerateDownloadStream(w, urlTypeId.Type, urlTypeId.ID)
	})
}
