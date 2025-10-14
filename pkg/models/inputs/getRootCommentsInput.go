package inputs

type GetRootCommentsInput struct {
	Limit  int `json:"limit" form:"limit" binding:"required"`
	Offset int `json:"offset" form:"offset"`
}
