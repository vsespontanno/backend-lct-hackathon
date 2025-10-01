package service

import (
	"black-pearl/backend-hackathon/internal/domain/quiz/entity"
	"black-pearl/backend-hackathon/internal/domain/quiz/interfaces"
	"context"

	"go.uber.org/zap"
)

type QuizService struct {
	repo   interfaces.QuizRepoInterface
	logger *zap.SugaredLogger
}

func NewQuizService(repo interfaces.QuizRepoInterface, logger *zap.SugaredLogger) *QuizService {
	return &QuizService{repo: repo, logger: logger}
}

// Временные затычки
func (s *QuizService) GetQuiz(ctx context.Context, quizID int64) (*entity.Quiz, error) {
	quiz, err := s.repo.GetQuizByID(ctx, quizID)
	if err != nil {
		s.logger.Errorw("failed to get quiz", "error", err, "stage", "GetQuiz.GetQuizByID")
		return nil, err
	}
	return quiz, err
}
