package gormpostgres

import (
	"time"

	gormerrors "github.com/PC-Core/pc-core-backend/internal/database/gormPostgres/gormErrors"
	"github.com/PC-Core/pc-core-backend/internal/errors"
	"github.com/PC-Core/pc-core-backend/pkg/models"
	"github.com/PC-Core/pc-core-backend/pkg/models/inputs"
	"github.com/PC-Core/pc-core-backend/pkg/models/outputs"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type commentRow struct {
	ID          int64         `gorm:"column:id;primaryKey"`
	UserID      int64         `gorm:"column:user_id"`
	ProductID   int64         `gorm:"column:product_id"`
	CommentText string        `gorm:"column:comment_text"`
	AnswerOn    *int64        `gorm:"column:answer_on"`
	Rating      *int16        `gorm:"column:rating"`
	CreatedAt   time.Time     `gorm:"column:created_at"`
	UpdatedAt   *time.Time    `gorm:"column:updated_at"`
	MediaIDs    pq.Int64Array `gorm:"column:media_ids;type:bigint[]"`
	Deleted     bool          `gorm:"column:is_deleted"`

	UserID_A  int64           `gorm:"column:user.id"`
	UserName  string          `gorm:"column:user.name"`
	UserEmail string          `gorm:"column:user.email"`
	UserRole  models.UserRole `gorm:"column:user.role"`
}

type TargetCommentGroup string

type LoadedComments struct {
	DbComments []DbComment
	Comments   []models.Comment
	TotalCount int64
}

func (*GormPostgresController) loadCommentIds(comments []DbComment) []int64 {
	result := make([]int64, 0, len(comments))

	for _, comm := range comments {
		result = append(result, comm.ID)
	}

	return result
}

func (c *GormPostgresController) getAnswersCount(parent_id int64) (int64, errors.PCCError) {
	var count int64

	err := c.db.Model(&DbComment{}).Where("answer_on = ?", parent_id).Count(&count).Error

	if err != nil {
		return 0, gormerrors.GormErrorCast(err)
	}

	return count, nil

	// 	var count int64

	// 	cte := `
	// WITH RECURSIVE comment_tree AS (
	//     SELECT id, answer_on
	//     FROM comments
	//     WHERE answer_on = ?

	//     UNION ALL

	//     SELECT c.id, c.answer_on
	//     FROM comments c
	//     INNER JOIN comment_tree ct ON c.answer_on = ct.id
	// )
	// SELECT COUNT(*) FROM comment_tree;
	// `

	// 	err := c.db.Raw(cte, parent_id).Scan(&count).Error
	// 	if err != nil {
	// 		return -1, gormerrors.GormErrorCast(err)
	// 	}

	// return count, nil
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

func (c *GormPostgresController) getRootCommentsCount(product_id int64) (int64, errors.PCCError) {
	var count int64

	err := c.db.Model(&DbComment{}).Where("product_id = ? AND answer_on IS NULL", product_id).Count(&count).Error

	if err != nil {
		return -1, gormerrors.GormErrorCast(err)
	}

	return count, nil
}

func (c *GormPostgresController) loadRootComments(product_id int64, userID *int64, limit int, offset int) (*LoadedComments, errors.PCCError) {
	var comments []DbComment
	commentReactions := make(map[int64]models.CommentReactions)

	err := c.db.Preload("User").Preload("Product").Order("created_at DESC").Limit(limit).Offset(offset).Where("product_id = ? AND answer_on is NULL", product_id).Find(&comments).Error

	if err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	commentIds := c.loadCommentIds(comments)
	dbreactions, perr := c.loadReactionsForComments(commentIds)

	if perr != nil {
		return nil, perr
	}

	for k, v := range dbreactions {
		commentReactions[k] = c.getCommentReactions(v, userID)
	}

	rootCount, perr := c.getRootCommentsCount(product_id)

	if perr != nil {
		return nil, perr
	}

	result, perr := c.dbCommentsIntoComments(comments, rootCount, userID)

	for i := range result.Comments {
		cmt := &result.Comments[i]

		cnt, perr := c.getAnswersCount(cmt.ID)
		if perr != nil {
			continue
		}

		cmt.ChildrenCount = uint64(cnt)
	}

	if perr != nil {
		return nil, perr
	}

	return result, nil
}

func (c *GormPostgresController) GetRootCommentsForProduct(product_id int64, userID *int64, limit int, offset int) (*outputs.CommentsOutput, errors.PCCError) {
	result, err := c.loadRootComments(product_id, userID, limit, offset)

	if err != nil {
		return nil, err
	}

	return &outputs.CommentsOutput{
		Comments: result.Comments,
		Amount:   result.TotalCount,
	}, err
}

func buildTree(comment *models.Comment, idToChildrenMap map[int64][]int64, idToCommentMap map[int64]*models.Comment) {
	childIds := idToChildrenMap[comment.ID]

	for _, childId := range childIds {
		childComment := idToCommentMap[childId]

		if childComment == nil {
			continue
		}

		buildTree(childComment, idToChildrenMap, idToCommentMap)
		comment.Children = append(comment.Children, *childComment)
	}
}

func (c *GormPostgresController) dbCommentsIntoComments(comments []DbComment, all_count int64, userID *int64) (*LoadedComments, errors.PCCError) {
	result := make([]models.Comment, 0)
	commentReactions := make(map[int64]models.CommentReactions)

	counts := make(map[int64]uint64)
	for _, c := range comments {
		if c.AnswerOn != nil {
			counts[*c.AnswerOn]++
		}
	}

	commentIds := c.loadCommentIds(comments)
	dbreactions, perr := c.loadReactionsForComments(commentIds)

	if perr != nil {
		return nil, perr
	}

	for k, v := range dbreactions {
		commentReactions[k] = c.getCommentReactions(v, userID)
	}

	for _, comment := range comments {
		medias, perr := c.LoadMediasForComment(comment.MediaIDs)

		if perr != nil {
			return nil, perr
		}

		result = append(result, *models.NewComment(comment.ID, comment.User.IntoUser(), comment.CommentText, []models.Comment{}, counts[comment.ID], comment.Rating, &comment.CreatedAt, comment.UpdatedAt, medias.IntoMedias(), commentReactions[comment.ID], comment.Deleted))
	}

	return &LoadedComments{comments, result, all_count}, nil
}

func (c *GormPostgresController) loadAnswersV2(parent_id int64, limit int, offset int) ([]DbComment, errors.PCCError) {
	var raw []commentRow

	query := `
WITH RECURSIVE comment_tree AS (
    SELECT 
        c.id,
        c.user_id,
        c.product_id,
        c.answer_on,
        c.comment_text,
        c.created_at,
        0 AS depth
    FROM comments c
    WHERE c.id = ?

    UNION ALL

    SELECT 
        c.id,
        c.user_id,
        c.product_id,
        c.answer_on,
        c.comment_text,
        c.created_at,
        ct.depth + 1
    FROM comments c
    INNER JOIN comment_tree ct ON c.answer_on = ct.id
)
SELECT 
    ct.id,
    ct.product_id,
    ct.answer_on,
    ct.comment_text,
    ct.created_at,
    ct.depth,
    u.id AS "user.id",
    u.name AS "user.name",
    u.email AS "user.email",
    u.role AS "user.role"
FROM comment_tree ct
JOIN users u ON u.id = ct.user_id
WHERE ct.depth = 0 OR ct.id IN (
    SELECT id FROM comment_tree WHERE depth > 0
    ORDER BY depth, created_at
    LIMIT ? OFFSET ?
)
ORDER BY depth, ct.created_at;
`

	if err := c.db.Raw(query, parent_id, limit, offset).Scan(&raw).Error; err != nil {
		return nil, gormerrors.GormErrorCast(err)
	}

	var comments []DbComment
	for _, r := range raw {
		c := DbComment{
			ID:          r.ID,
			UserID:      r.UserID,
			CommentText: r.CommentText,
			CreatedAt:   r.CreatedAt,
			ProductID:   r.ProductID,
			AnswerOn:    r.AnswerOn,
			Rating:      r.Rating,
			UpdatedAt:   r.UpdatedAt,
			MediaIDs:    r.MediaIDs,
			Deleted:     r.Deleted,

			User: DbUser{
				ID: int(r.UserID), Name: r.UserName, Email: r.UserEmail, Role: r.UserRole,
			},
		}
		comments = append(comments, c)
	}

	return comments, nil
}

func (c *GormPostgresController) GetAnswersOnComment(product_id int64, userID *int64, comment_id int64, limit int, offset int) (*outputs.CommentsOutput, errors.PCCError) {
	ans, err := c.loadAnswersV2(comment_id, limit, offset)

	if err != nil {
		return nil, err
	}

	res, err := c.dbCommentsIntoComments(ans, int64(len(ans)), userID)

	if err != nil {
		return nil, err
	}

	idToCommentMap := make(map[int64]*models.Comment)
	idToChildrenIdMap := make(map[int64][]int64)

	for i, comment := range res.Comments {
		dbComment := res.DbComments[i]
		idToCommentMap[comment.ID] = &comment

		if dbComment.AnswerOn != nil {
			parentID := *dbComment.AnswerOn
			idToChildrenIdMap[parentID] = append(idToChildrenIdMap[parentID], dbComment.ID)
		}
	}

	targetComment, ok := idToCommentMap[comment_id]

	if !ok {
		return &outputs.CommentsOutput{
			Comments: []models.Comment{},
			Amount:   res.TotalCount,
		}, nil
	}

	buildTree(targetComment, idToChildrenIdMap, idToCommentMap)

	return &outputs.CommentsOutput{
		Comments: targetComment.Children,
		Amount:   res.TotalCount,
	}, nil
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
