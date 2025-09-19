package models

type MediaType string

const (
	MediaImage MediaType = "Image"
	MediaVideo MediaType = "Video"
)

type Media struct {
	ID        uint64    `json:"id"`
	Url       string    `json:"url"`
	Type      MediaType `json:"type"`
	ProductID uint64    `json:"product_id"`
}

func NewMedia(id uint64, url string, t MediaType, prod_id uint64) *Media {
	return &Media{
		id, url, t, prod_id,
	}
}

func NewMediaFromInput(id uint64, prod_id uint64, input *InputMedia) *Media {
	return &Media{
		id, input.Url, input.Type, prod_id,
	}
}

type Medias []Media

func (ms Medias) GetIDs() []uint64 {
	ids := make([]uint64, len(ms))
	for i, m := range ms {
		ids[i] = m.ID
	}
	return ids
}

type InputMedia struct {
	Url  string    `json:"url"`
	Type MediaType `json:"type"`
}
