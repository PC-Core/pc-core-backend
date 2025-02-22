package inputs

type LoginUserInput struct {
	Email    string `form:"email" binding:"required,email" json:"email"`
	Password string `form:"password" binding:"required" json:"password"`
}
