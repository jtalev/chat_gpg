package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

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
