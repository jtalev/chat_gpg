package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetItemTypes(db *sql.DB) ([]models.ItemType, error) {
	q := `
	select * from item_types;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	itemTypes := []models.ItemType{}
	i := models.ItemType{}
	for rows.Next() {
		if err := rows.Scan(
			&i.UUID,
			&i.Type,
			&i.Description,
			&i.CreatedAt,
			&i.ModifiedAt,
		); err != nil {
			return nil, err
		}
		itemTypes = append(itemTypes, i)
	}

	return itemTypes, nil
}

func PostItemType(itemType models.ItemType, db *sql.DB) error {
	q := `
	insert into item_types(uuid, type, description)
	values ($1, $2, $3);	
	`

	_, err := db.Exec(q, itemType.UUID, itemType.Type, itemType.Description)
	if err != nil {
		return err
	}

	return nil
}
