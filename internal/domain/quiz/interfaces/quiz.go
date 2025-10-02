package interfaces

import (
	"black-pearl/backend-hackathon/internal/domain/quiz/entity"
	"context"
)

type QuizRepoInterface interface {
	GetQuizByID(ctx context.Context, quizID int) (*entity.Quiz, error)
	InsertQuiz(ctx context.Context, quiz *entity.Quiz) error
}
