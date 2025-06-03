package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetItemSizes(db *sql.DB) ([]models.ItemSize, error) {
	q := `
	select * from item_size;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemSizes := []models.ItemSize{}
	i := models.ItemSize{}
	for rows.Next() {
		if err := rows.Scan(
			&i.UUID,
			&i.Size,
			&i.Description,
			&i.CreatedAt,
			&i.ModifiedAt,
		); err != nil {
			return nil, err
		}
		itemSizes = append(itemSizes, i)
	}
	return itemSizes, nil
}

func GetItemSizeByUuid(uuid string, db *sql.DB) (models.ItemSize, error) {
	q := `
	select * from item_size where uuid = ?;
	`

	i := models.ItemSize{}
	rows, err := db.Query(q, uuid)
	if err != nil {
		return i, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&i.UUID,
			&i.Size,
			&i.Description,
			&i.CreatedAt,
			&i.ModifiedAt,
		); err != nil {
			return i, err
		}
	}
	return i, nil
}

func PostItemSize(itemSize models.ItemSize, db *sql.DB) error {
	q := `
	insert into item_size(uuid, size, description)
	values ($1, $2, $3);
	`

	_, err := db.Exec(q, itemSize.UUID, itemSize.Size, itemSize.Description)
	if err != nil {
		return err
	}
	return nil
}

func PutItemSize(itemSize models.ItemSize, db *sql.DB) error {
	q := `
	update item_size
	set size = $1, description = $2, modified_at = CURRENT_TIMESTAMP
	where uuid = $3
	`

	_, err := db.Exec(q, itemSize.Size, itemSize.Description, itemSize.UUID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItemSize(uuid string, db *sql.DB) error {
	q := `
	delete from item_size where uuid = ?;
	`
	_, err := db.Exec(q, uuid)
	if err != nil {
		return err
	}

	return nil
}
