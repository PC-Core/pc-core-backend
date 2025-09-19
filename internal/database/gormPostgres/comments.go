package gormpostgres

import (
	"time"

	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type TargetCommentGroup string

const (
	TCG_ROOT     TargetCommentGroup = "root"
	TCG_NON_ROOT TargetCommentGroup = "nonroot"
	TCG_ALL      TargetCommentGroup = "all"
)

func (*GormPostgresController) loadCommentIds(comments []DbComment) []int64 {
	result := make([]int64, 0, len(comments))

	for _, comm := range comments {
		result = append(result, comm.ID)
	}

	return result
}

func (c *GormPostgresController) loadReactionsForComments(ids []int64) (map[int64][]DbCommentReaction, errors.PCCError) {
	var reactions []DbCommentReaction
	err := c.db.Where("comment_id = ANY(?)", pq.Array(ids)).Find(&reactions).Error

	result := make(map[int64][]DbCommentReaction)

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	for _, reaction := range reactions {
		result[reaction.CommentID] = append(result[reaction.CommentID], reaction)
	}

	return result, nil
}

func (*GormPostgresController) getCommentReactions(dbreactions []DbCommentReaction, userID *int64) models.CommentReactions {
	reactions := make(map[models.ReactionType]uint64)
	var yourReaction *models.ReactionType = nil

	for _, reaction := range dbreactions {
		if prev, ok := reactions[reaction.Type]; ok {
			reactions[reaction.Type] = prev + 1
		} else {
			reactions[reaction.Type] = 1
		}

		if userID != nil && *userID == reaction.UserID {
			yourReaction = &reaction.Type
		}
	}

	return models.CommentReactions{ReactionsAmount: reactions, YourReaction: yourReaction}
}

func (c *GormPostgresController) LoadMediasForComment(mediaIDs []int64) (DbMedias, errors.PCCError) {
	var medias DbMedias

	err := c.db.Where("id = ANY(?)", pq.Array(mediaIDs)).Find(&medias).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	return medias, nil
}

func getFiltersForCommentGroup(group TargetCommentGroup) string {
	switch group {
	case TCG_ROOT:
		return " AND answer_on IS NULL"
	case TCG_NON_ROOT:
		return " AND answer_on IS NOT NULL"
	case TCG_ALL:
		fallthrough
	default:
		return ""
	}
}

func (c *GormPostgresController) loadComments(product_id int64, userID *int64, target TargetCommentGroup) ([]DbComment, []models.Comment, errors.PCCError) {
	filter := getFiltersForCommentGroup(target)

	var comments []DbComment
	result := make([]models.Comment, 0)
	commentReactions := make(map[int64]models.CommentReactions)

	err := c.db.Preload("User").Preload("Product").Where("product_id = ?"+filter, product_id).Find(&comments).Error

	if err != nil {
		return nil, nil, gormerrors.GormErrorCast(err)
	}

	commentIds := c.loadCommentIds(comments)
	dbreactions, perr := c.loadReactionsForComments(commentIds)

	if perr != nil {
		return nil, nil, perr
	}

	for k, v := range dbreactions {
		commentReactions[k] = c.getCommentReactions(v, userID)
	}

	for _, comment := range comments {
		medias, perr := c.LoadMediasForComment(comment.MediaIDs)

		if perr != nil {
			return nil, nil, perr
		}

		result = append(result, *models.NewComment(comment.ID, comment.User.IntoUser(), comment.CommentText, []models.Comment{}, comment.Rating, &comment.CreatedAt, comment.UpdatedAt, medias.IntoMedias(), commentReactions[comment.ID], comment.Deleted))
	}

	return comments, result, nil
}

func (c *GormPostgresController) GetRootCommentsForProduct(product_id int64, userID *int64) ([]models.Comment, errors.PCCError) {
	_, result, err := c.loadComments(product_id, userID, TCG_ALL)
	return result, err
}

func buildTree(comment *models.Comment, idToChildrenMap map[int64][]int64, idToCommentMap map[int64]*models.Comment) {
	childIds := idToChildrenMap[comment.ID]

	for _, childId := range childIds {
		childComment := idToCommentMap[childId]
		buildTree(childComment, idToChildrenMap, idToCommentMap)
		comment.Children = append(comment.Children, *childComment)
	}
}

func (c *GormPostgresController) GetAnswersOnComment(product_id int64, userID *int64, comment_id int64) ([]models.Comment, errors.PCCError) {
	dbComments, upgrouppedComments, err := c.loadComments(product_id, userID, TCG_ALL)

	if err != nil {
		return nil, err
	}

	idToCommentMap := make(map[int64]*models.Comment)
	idToChildrenIdMap := make(map[int64][]int64)

	for i, comment := range upgrouppedComments {
		dbComment := dbComments[i]
		idToCommentMap[comment.ID] = &comment

		if dbComment.AnswerOn != nil {
			parentID := *dbComment.AnswerOn
			idToChildrenIdMap[parentID] = append(idToChildrenIdMap[parentID], dbComment.ID)
		}
	}

	targetComment := idToCommentMap[comment_id]

	buildTree(targetComment, idToChildrenIdMap, idToCommentMap)

	return targetComment.Children, nil
}

func (c *GormPostgresController) CheckUserOwnCommentByID(commentID int64, userID int64) errors.PCCError {
	var comment DbComment

	err := c.db.Where("id = ? AND user_id = ?", commentID, userID).Find(&comment).Error

	if err != nil {
		return gormerrors.GormErrorCastUserOwn(err)
	}

	return nil
}

func (c *GormPostgresController) AddComment(input *inputs.AddCommentInput, userID int64, product_id int64) (int64, errors.PCCError) {
	comment := DbComment{
		ID:          0,
		UserID:      userID,
		ProductID:   product_id,
		CommentText: input.Text,
		AnswerOn:    input.Answer,
		Rating:      input.Rating,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		MediaIDs:    []int64{},
	}

	err := c.db.Create(&comment).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	return comment.ID, nil
}

func (c *GormPostgresController) EditComment(newText string, commentID int64, userID int64) (int64, errors.PCCError) {
	if err := c.CheckUserOwnCommentByID(commentID, userID); err != nil {
		return -1, err
	}

	var comment DbComment

	err := c.db.First(&comment, commentID).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	now := time.Now()

	comment.CommentText = newText
	comment.UpdatedAt = &now

	err = c.db.Save(&comment).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	return commentID, nil
}

func (c *GormPostgresController) DeleteComment(commentID int64, userID int64) (int64, errors.PCCError) {
	if err := c.CheckUserOwnCommentByID(commentID, userID); err != nil {
		return -1, err
	}

	var comment DbComment

	err := c.db.First(&comment, commentID).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	comment.CommentText = ""
	comment.Deleted = true

	err = c.db.Save(&comment).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	return commentID, nil
}

func (c *GormPostgresController) CreateReaction(commentID int64, userID int64, ty models.ReactionType) errors.PCCError {
	reaction := DbCommentReaction{
		UserID:    userID,
		CommentID: commentID,
		Type:      ty,
	}

	err := c.db.Create(&reaction).Error

	if err != nil {
		return gormerrors.GormErrorCast(err)
	}

	return nil
}

func (c *GormPostgresController) DeleteReaction(existing *DbCommentReaction) errors.PCCError {
	err := c.db.Delete(existing).Error

	if err != nil {
		return gormerrors.GormErrorCast(err)
	}

	return nil
}

func (c *GormPostgresController) SetReaction(commentID int64, userID int64, ty models.ReactionType) (int64, errors.PCCError) {

	var existing DbCommentReaction

	err := c.db.First(&existing, "comment_id = ? AND user_id = ?", &commentID, &userID).Error

	if err == gorm.ErrRecordNotFound {
		return commentID, c.CreateReaction(commentID, userID, ty)
	}

	if existing.Type == ty {
		return commentID, c.DeleteReaction(&existing)
	}

	existing.Type = ty

	err = c.db.Save(&existing).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	return commentID, nil
}
