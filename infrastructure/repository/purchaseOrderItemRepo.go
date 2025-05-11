package infrastructure

import (
	"database/sql"

	models "github.com/jtalev/chat_gpg/domain/models"
)

func GetItemsByOrderUuid(uuid string, db *sql.DB) ([]models.PurchaseOrderItem, error) {
	q := `
	select * from purchase_order_item where purchase_order_id = ?
	`

	rows, err := db.Query(q, uuid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []models.PurchaseOrderItem{}
	i := models.PurchaseOrderItem{}
	for rows.Next() {
		if err := rows.Scan(
			&i.UUID,
			&i.PurchaseOrderId,
			&i.ItemName,
			&i.ItemTypeId,
			&i.Quantity,
			&i.CreatedAt,
			&i.ModifiedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	return items, nil
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
