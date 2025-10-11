package services

import (
	"github.com/nemdull/myapi/models"
	"github.com/nemdull/myapi/repositories"
)

// PostCommentHandlerで使用することを想定したサービス
// 引数の情報をもとに新しいコメントを作り、結果を返却
func (s *MyAppService) PostCommentService(comment models.Comment) (models.Comment, error) {
	newComment, err := repositories.InsertComment(s.DB, comment)
	if err != nil {
		return models.Comment{}, err
	}

	return newComment, nil
}
