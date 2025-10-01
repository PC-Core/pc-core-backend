package models

import "time"

type CommentReactions struct {
	ReactionsAmount map[ReactionType]uint64 `json:"reactions_amount"`
	YourReaction    *ReactionType           `json:"your_reaction"`
}

type Comment struct {
	ID            int64 `json:"id"`
	*User         `json:"user"`
	Text          string           `json:"text"`
	Children      []Comment        `json:"children"`
	ChildrenCount uint64           `json:"children_count"`
	Rating        *int16           `json:"rating"`
	CreatedAt     *time.Time       `json:"created_at"`
	UpdatedAt     *time.Time       `json:"updated_at"`
	Medias        Medias           `json:"medias"`
	Reactions     CommentReactions `json:"reactions"`
	Deleted       bool             `json:"deleted"`
}

func NewComment(id int64, user *User, text string, children []Comment, chilren_count uint64, rating *int16, created_at *time.Time, updated_at *time.Time, medias []Media, reactions CommentReactions, deleted bool) *Comment {
	return &Comment{
		id,
		user,
		text,
		children,
		chilren_count,
		rating,
		created_at,
		updated_at,
		medias,
		reactions,
		deleted,
	}
}
