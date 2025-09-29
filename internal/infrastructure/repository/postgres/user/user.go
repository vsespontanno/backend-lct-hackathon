package user

import (
	"black-pearl/backend-hackathon/internal/domain/user/entity"
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type UserRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *UserRepo) GetUserByID(ctx context.Context, id int) (entity.User, error) {
	builder := r.builder.
		Select("id").
		From("users").
		Where(sq.Eq{"id": id}).
		Suffix("ON CONFLICT (id) DO NOTHING")

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return entity.User{}, err
	}

	var user entity.User
	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	if err := row.Scan(&user.ID); err != nil {
		if err == sql.ErrNoRows {
			return entity.User{}, sql.ErrNoRows
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepo) CreateUser(ctx context.Context, id int) error {
	builder := r.builder.
		Insert("users").
		Columns("id").
		Values(id)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
