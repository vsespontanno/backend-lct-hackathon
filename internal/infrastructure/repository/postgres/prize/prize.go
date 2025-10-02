package prize

import (
	"black-pearl/backend-hackathon/internal/domain/prize/entity"
	"context"
	"database/sql"
)

type PrizeRepo struct {
	db *sql.DB
}

func NewPrizeRepo(db *sql.DB) *PrizeRepo {
	return &PrizeRepo{
		db: db,
	}
}

func (r *PrizeRepo) GetMyPrizes(ctx context.Context, userID int) (*[]entity.Prize, error) {
	query := `SELECT p.title, p.descr, p.type FROM prizes p INNER JOIN user_prizes up ON p.id = up.prize_id WHERE up.user_id = $1`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []entity.Prize
	for rows.Next() {
		var p entity.Prize
		if err := rows.Scan(&p.Title, &p.Description, &p.Type); err != nil {
			return nil, err
		}
		prizes = append(prizes, p)
	}
	return &prizes, nil
}

func (r *PrizeRepo) GetAvailablePrizes(ctx context.Context, id int) (*[]entity.Prize, error) {
	query := `SELECT p.title, p.descr, p.type FROM prizes p LEFT JOIN user_prizes up ON p.id = up.prize_id AND up.user_id = $1 WHERE up.user_id IS NULL`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prizes []entity.Prize
	for rows.Next() {
		var p entity.Prize
		if err := rows.Scan(&p.Title, &p.Description, &p.Type); err != nil {
			return nil, err
		}
		prizes = append(prizes, p)
	}
	return &prizes, nil
}
