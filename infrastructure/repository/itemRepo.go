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

func GetItemTypeByUuid(uuid string, db *sql.DB) (models.ItemType, error) {
	q := `
	select * from item_types where uuid = ?;
	`

	i := models.ItemType{}
	rows, err := db.Query(q, uuid)
	if err != nil {
		return i, err
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(
			&i.UUID,
			&i.Type,
			&i.Description,
			&i.CreatedAt,
			&i.ModifiedAt,
		); err != nil {
			return i, err
		}
	}

	return i, nil
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

func PutItemType(itemType models.ItemType, db *sql.DB) error {
	q := `
	update item_types
	set type = $1, description = $2, modified_at = CURRENT_TIMESTAMP
	where uuid = $3
	`

	_, err := db.Exec(q, itemType.Type, itemType.Description, itemType.UUID)
	if err != nil {
		return err
	}
	return nil
}

func DeleteItemType(uuid string, db *sql.DB) error {
	q := `
	delete from item_types where uuid = ?;
	`
	_, err := db.Exec(q, uuid)
	if err != nil {
		return err
	}

	return nil
}
