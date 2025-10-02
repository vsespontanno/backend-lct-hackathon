package sectionitems

import (
	"black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	"context"
	"database/sql"
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

type SectionItemsRepo struct {
	db      *sql.DB
	builder sq.StatementBuilderType
}

func NewSectionItemsRepo(db *sql.DB) *SectionItemsRepo {
	return &SectionItemsRepo{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}

// Получить все элементы секции по sectionId
func (r *SectionItemsRepo) GetSectionItemsBySectionId(ctx context.Context, sectionId int) (*[]entity.SectionItem, error) {
	builder := r.builder.
		Select("sectionid", "istest", "title", "itemid").
		From("sectionitems").
		Where(sq.Eq{"sectionid": sectionId})

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := r.db.QueryContext(ctx, sqlStr, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []entity.SectionItem
	for rows.Next() {
		var item entity.SectionItem
		if err := rows.Scan(&item.SectionID, &item.IsTest, &item.Title, &item.ItemID); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no section items found for sectionId %d", sectionId)
	}

	return &items, nil
}

// Создать новый элемент секции
func (r *SectionItemsRepo) CreateSectionItem(ctx context.Context, item *entity.SectionItem) error {
	builder := r.builder.
		Insert("sectionitems").
		Columns("sectionid", "istest", "title", "itemid").
		Values(item.SectionID, item.IsTest, item.Title, item.ItemID)

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return err
	}

	_, err = r.db.ExecContext(ctx, sqlStr, args...)
	return err
}
