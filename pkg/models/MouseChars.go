package models

type MouseChars struct {
	ID          uint64 `json:"id"`
	Name        string `json:"name"`
	TypeMouses  string `json:"type_mouses"`
	Dpi         uint64 `json:"dpi"`
	ReleaseYear uint64 `json:"release_year"`
}

func NewMouseChars(id uint64, name, type_mouses string, dpi, release_year uint64) *MouseChars{
	return &MouseChars{
		id, name, type_mouses, dpi, release_year,
	}
}

