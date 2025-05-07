package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func BatchPostPurchaseOrderItem(purchaseOrderItems []models.PurchaseOrderItem, db *sql.DB) error {
	q := `
	insert into purchase_order_item(uuid, purchase_order_id, item_name,
		item_type_id, quantity)
	values ($1, $2, $3, $4, $5);
	`

	for _, item := range purchaseOrderItems {
		_, err := db.Exec(q, item.UUID, item.PurchaseOrderId, item.ItemName,
			item.ItemTypeId, item.Quantity)
		if err != nil {
			return err
		}
	}

	return nil
}
