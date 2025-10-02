package models

type KeyboardChars struct {
	ID            uint64   `json:"id"`
	Name          string   `json:"name"`
	TypeKeyBoards string   `json:"type_keyboards"`
	Switches      []string `json:"switches"`
	ReleaseYear   uint64   `json:"release_year"`
}

func NewKeyBoardChars(id uint64, name, type_keyboards string, switches []string, release_year uint64) *KeyboardChars{
	return &KeyboardChars{
		id, name, type_keyboards, switches, release_year,
	}
}
