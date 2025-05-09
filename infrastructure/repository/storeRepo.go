package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetStores(db *sql.DB) ([]models.Store, error) {
	q := `
	select * from stores;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []models.Store{}
	s := models.Store{}
	for rows.Next() {
		if err := rows.Scan(
			&s.UUID,
			&s.BusinessName,
			&s.Email,
			&s.Phone,
			&s.Address,
			&s.Suburb,
			&s.City,
			&s.CreatedAt,
			&s.ModifiedAt,
		); err != nil {
			return nil, err
		}
		stores = append(stores, s)
	}
	return stores, nil
}

// TODO: complete this function so that the store names can be properly set in purchase order history
func GetStoreByUuid(uuid string, db *sql.DB) (models.Store, error) {
	q := `
	select * from stores where uuid = ?;
	`
	rows, err := db.Query(q, uuid)
	if err != nil {
		return models.Store{}, err
	}
	defer rows.Close()

	s := models.Store{}
	if rows.Next() {
		err := rows.Scan(
			&s.UUID,
			&s.BusinessName,
			&s.Email,
			&s.Phone,
			&s.Address,
			&s.Suburb,
			&s.City,
			&s.CreatedAt,
			&s.ModifiedAt,
		)
		if err != nil {
			return models.Store{}, err
		}
	}
	return s, nil
}

func PostStore(store models.Store, db *sql.DB) error {
	q := `
	insert into stores(uuid, business_name, email,
		phone, address, suburb, city)
	values ($1, $2, $3, $4, $5, $6, $7);
	`

	_, err := db.Exec(q, store.UUID, store.BusinessName, store.Email,
		store.Phone, store.Address, store.Suburb, store.City)
	if err != nil {
		return err
	}

	return nil
}
