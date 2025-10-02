package pet

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type PetRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewPetRepo(db *sql.DB) *PetRepo {
	return &PetRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

func (r *PetRepo) GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error) {
	builder := r.builder.
		Select("id", "petName", "age", "exp").
		From("pets").
		Where(sq.Eq{"user_id": userID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var pet entity.Pet
	row := r.db.QueryRowContext(ctx, sqlStr, args...)
	if err := row.Scan(&pet.ID, &pet.Name, &pet.Age, &pet.Exp); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("pet not found")
		}
		return nil, err
	}
	return &pet, nil
}

func (r *PetRepo) SetPetName(ctx context.Context, name string, userID int) error {
	builder := r.builder.
		Update("pets").
		Set("petName", name).
		Where(sq.Eq{"user_id": userID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}

	return err
}

func (r *PetRepo) CreatePet(ctx context.Context, userID int) error {
	builder := r.builder.
		Insert("pets").
		Columns("user_id", "petName", "age", "exp").
		Values(userID, "", 0, 0).
		Suffix("ON CONFLICT (id) DO NOTHING")

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return err
}

func (r *PetRepo) UpdateXP(ctx context.Context, xp int, userID int) error {
	builder := r.builder.
		Update("pets").
		Set("exp", sq.Expr("exp + ?", xp)).
		Where(sq.Eq{"user_id": userID})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	if err != nil {
		return err
	}
	return err
}
