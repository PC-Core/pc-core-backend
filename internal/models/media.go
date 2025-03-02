package models

type MediaType string

const (
	MediaImage MediaType = "Image"
	MediaVideo MediaType = "Video"
)

type Media struct {
	ID   uint64    `json:"id"`
	Url  string    `json:"url"`
	Type MediaType `json:"type"`
}

func NewMedia(id uint64, url string, t MediaType) *Media {
	return &Media{
		id, url, t,
	}
}

func NewMediaFromInput(id uint64, input *InputMedia) *Media {
	return &Media{
		id, input.Url, input.Type,
	}
}

type InputMedia struct {
	Url  string    `json:"url"`
	Type MediaType `json:"type"`
}
