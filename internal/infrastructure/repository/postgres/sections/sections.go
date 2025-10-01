package sections

import (
	"black-pearl/backend-hackathon/internal/domain/sections/entity"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type SectionsRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewSectionsRepo(db *sql.DB) *SectionsRepo {
	return &SectionsRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Получить все секции
func (r *SectionsRepo) GetSections(ctx context.Context) (*[]entity.Sections, error) {
	builder := r.builder.
		Select("id", "title").
		From("sections")

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []entity.Sections
	for rows.Next() {
		var s entity.Sections
		if err := rows.Scan(&s.ID, &s.Title); err != nil {
			return nil, err
		}
		sections = append(sections, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &sections, nil
}

// Создать новую секцию
func (r *SectionsRepo) CreateSection(ctx context.Context, s *entity.Sections) error {
	builder := r.builder.
		Insert("sections").
		Columns("title").
		Values(s.Title).
		Suffix("RETURNING id")

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
