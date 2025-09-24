package inputs

type GetAnswersInput struct {
	ProductID int64 `json:"product_id"`
	Limit     int   `json:"limit" binding:"required"`
	Offset    int   `json:"offset"`
}
