package quiz

import (
	"context"
	"database/sql"

	"github.com/lib/pq"

	"black-pearl/backend-hackathon/internal/domain/quiz/entity"

	sq "github.com/Masterminds/squirrel"

	"fmt"
)

type quizRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewQuizRepo(db *sql.DB) *quizRepo {
	return &quizRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *quizRepo) GetQuizByID(ctx context.Context, quizID int64) (*entity.Quiz, error) {
	sqlStr := `SELECT id, title, content, options, correct_answer FROM quizs WHERE id = $1`

	var quiz entity.Quiz
	var options pq.StringArray

	row := r.db.QueryRowContext(ctx, sqlStr, quizID)
	err := row.Scan(&quiz.ID, &quiz.Title, &quiz.Content, &options, &quiz.CorrectAnswer)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("quiz not found")
		}
		return nil, err
	}

	quiz.Options = []string(options)
	return &quiz, nil
}

func (r *quizRepo) InsertQuiz(ctx context.Context, quiz *entity.Quiz) error {
	sqlStr := `INSERT INTO quizs (title, content, options, correct_answer) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, sqlStr, quiz.Title, quiz.Content, pq.Array(quiz.Options), quiz.CorrectAnswer)
	return err
}
