package outputs

type RestCharsObject struct {
	ID         uint64               `json:"id"`
	Components []RestCharsComponent `json:"components"`
}

func NewRestCharsObject(id uint64, components []RestCharsComponent) *RestCharsObject {
	return &RestCharsObject{id, components}
}
