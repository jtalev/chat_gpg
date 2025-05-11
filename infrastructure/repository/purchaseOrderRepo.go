package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetPurchaseOrders(db *sql.DB) ([]models.PurchaseOrder, error) {
	q := `
	select * from purchase_order;
	`

	rows, err := db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	purchaseOrders := []models.PurchaseOrder{}
	p := models.PurchaseOrder{}
	for rows.Next() {
		if err := rows.Scan(
			&p.UUID,
			&p.EmployeeId,
			&p.StoreId,
			&p.JobId,
			&p.Date,
			&p.CreatedAt,
			&p.ModifiedAt,
		); err != nil {
			return nil, err
		}
		purchaseOrders = append(purchaseOrders, p)
	}
	return purchaseOrders, nil
}

func GetPurchaseOrderByUuid(uuid string, db *sql.DB) (models.PurchaseOrder, error) {
	q := `
	select * from purchase_order where uuid = ?;
	`

	rows, err := db.Query(q, uuid)
	if err != nil {
		return models.PurchaseOrder{}, err
	}
	defer rows.Close()

	var p models.PurchaseOrder
	if rows.Next() {
		if err := rows.Scan(
			&p.UUID,
			&p.EmployeeId,
			&p.StoreId,
			&p.JobId,
			&p.Date,
			&p.CreatedAt,
			&p.ModifiedAt,
		); err != nil {
			return p, err
		}
	}
	return p, nil
}

func PostPurchaseOrder(purchaseOrder models.PurchaseOrder, db *sql.DB) error {
	q := `
	insert into purchase_order (uuid, employee_id, store_id,
		job_id, order_date)
	values ($1, $2, $3, $4, $5);
	`

	_, err := db.Exec(q, purchaseOrder.UUID, purchaseOrder.EmployeeId, purchaseOrder.StoreId,
		purchaseOrder.JobId, purchaseOrder.Date)
	if err != nil {
		return err
	}

	return nil
}
