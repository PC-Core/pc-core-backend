package inputs

type AddCommentInput struct {
	Text   string `json:"text"`
	Rating *int16 `json:"rating"`
	//Medias    models.Medias `json:"medias"`
	Answer *int64 `json:"answer"`
}
