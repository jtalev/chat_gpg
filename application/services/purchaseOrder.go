package application

import (
	"log"

	"github.com/google/uuid"
	domain "github.com/jtalev/chat_gpg/domain/models"
	models "github.com/jtalev/chat_gpg/domain/models"
)

type Store struct {
	StoreId   string
	StoreName string
}

type ItemType struct {
	UUID string `json:"uuid"`
	Type string `json:"type"`
}

type PurchaseOrderItem struct {
	UUID         string `json:"uuid"`
	ItemName     string `json:"item_name"`
	ItemTypeId   string `json:"item_type_id"`
	ItemTypeName string `json:"item_type_name"`
	Quantity     int    `json:"quantity"`

	ItemTypes []ItemType
}

type PurchaseOrder struct {
	UUID               string              `json:"uuid"`
	EmployeeId         string              `json:"employee_id"`
	StoreId            string              `json:"store_id"`
	JobId              int                 `json:"job_id"`
	Date               string              `json:"date"`
	PurchaseOrderItems []PurchaseOrderItem `json:"purchase_order_items"`

	Stores []Store
	Jobs   []models.Job
}

func mapPurchaseOrder(purchaseOrder PurchaseOrder) domain.PurchaseOrder {
	purchaseOrder.UUID = uuid.New().String()
	return domain.PurchaseOrder{
		UUID:       purchaseOrder.UUID,
		EmployeeId: purchaseOrder.EmployeeId,
		StoreId:    purchaseOrder.StoreId,
		JobId:      purchaseOrder.JobId,
		Date:       purchaseOrder.Date,
	}
}

func mapPurchaseOrderItems(purchaseOrder PurchaseOrder) []domain.PurchaseOrderItem {
	items := make([]domain.PurchaseOrderItem, len(purchaseOrder.PurchaseOrderItems))
	for i, item := range purchaseOrder.PurchaseOrderItems {
		item.UUID = uuid.New().String()
		items[i] = domain.PurchaseOrderItem{
			UUID:            item.UUID,
			PurchaseOrderId: purchaseOrder.UUID,
			ItemName:        item.ItemName,
			ItemTypeId:      item.ItemTypeId,
			Quantity:        item.Quantity,
		}
	}
	return items
}

func (o *PurchaseOrder) PostPurchaseOrder() error {
	order := mapPurchaseOrder(*o)
	o.UUID = order.UUID
	log.Println(order)
	items := mapPurchaseOrderItems(*o)
	log.Println(items)
	return nil
}
