package domain

type PurchaseOrder struct {
	UUID       string `json:"uuid"`
	EmployeeId string `json:"employee_id"`
	StoreId    string `json:"store_id"`
	JobId      int    `json:"job_id"`
	Date       string `json:"date"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type PurchaseOrderItem struct {
	UUID            string `json:"uuid"`
	PurchaseOrderId string `json:"purchase_order_id"`
	ItemName        string `json:"item_name"`
	ItemTypeId      string `json:"item_type_id"`
	Quantity        int    `json:"quantity"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

type ItemType struct {
	UUID       string `json:"uuid"`
	Type       string `json:"type"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type Store struct {
	UUID         string `json:"uuid"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Suburb       string `json:"suburb"`
	City         string `json:"city"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}
