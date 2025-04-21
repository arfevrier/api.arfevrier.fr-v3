package youtube

type VideoContent struct {
	AuthorName    string `json:"author_name"`
	VideoID       string `json:"video_id"`
	VideoTitle    string `json:"video_title"`
	PublishedDate string `json:"published_date"`
}
