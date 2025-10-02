package inputs

type GetAnswersInput struct {
	ProductID int64 `json:"product_id" form:"product_id"`
	Limit     int   `json:"limit" binding:"required" form:"limit"`
	Offset    int   `json:"offset" form:"offset"`
}
