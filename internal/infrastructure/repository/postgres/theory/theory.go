package theory

import (
	"black-pearl/backend-hackathon/internal/domain/theory/entity"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type TheoryRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewTheoryRepo(db *sql.DB) *TheoryRepo {
	return &TheoryRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Получить теорию по ID
func (r *TheoryRepo) GetTheoryByID(ctx context.Context, id int64) (*entity.Theory, error) {
	builder := r.builder.
		Select("id", "title", "content").
		From("theory").
		Where(sq.Eq{"id": id})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var t entity.Theory
	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	if err := row.Scan(&t.ID, &t.Title, &t.Content); err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}
	return &t, nil
}

// Создать новую запись теории
func (r *TheoryRepo) CreateTheory(ctx context.Context, t *entity.Theory) error {
	builder := r.builder.
		Insert("theory").
		Columns("id", "title", "content").
		Values(t.ID, t.Title, t.Content)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
