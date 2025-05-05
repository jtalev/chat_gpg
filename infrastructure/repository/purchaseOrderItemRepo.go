package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

type PurchaseOrderItem struct {
	UUID            string `json:"uuid"`
	PurchaseOrderId string `json:"purchase_order_id"`
	ItemName        string `json:"item_name"`
	ItemTypeId      string `json:"item_type_id"`
	Quantity        int    `json:"quantity"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

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
