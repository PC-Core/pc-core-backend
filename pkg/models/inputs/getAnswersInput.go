package inputs

type GetAnswersInput struct {
	ProductID int64 `json:"product_id"`
	Limit     int   `json:"limit"`
	Offset    int   `json:"offset"`
}
