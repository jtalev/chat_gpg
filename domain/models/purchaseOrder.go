package domain

type PurchaseOrder struct {
	UUID       string `json:"uuid"`
	EmployeeId string `json:"employee_id"`
	StoreId    string `json:"store_id"`
	JobId      int    `json:"job_id"`
	Date       string `json:"order_date"`
	CreatedAt  string `json:"created_at"`
	ModifiedAt string `json:"modified_at"`
}

type PurchaseOrderErrors struct {
	EmployeeIdErr         string
	StoreIdErr            string
	JobIdErr              string
	DateErr               string
	PurchaseOrderItemsErr string

	IsSuccessful bool
	SuccessMsg   string
}

func (p *PurchaseOrder) Validate() PurchaseOrderErrors {
	errors := PurchaseOrderErrors{IsSuccessful: true}
	errors = p.validateEmployeeId(errors)
	errors = p.validateStoreId(errors)
	errors = p.validateJobId(errors)
	errors = p.validateDate(errors)
	return errors
}

func (p *PurchaseOrder) validateEmployeeId(errors PurchaseOrderErrors) PurchaseOrderErrors {
	if p.EmployeeId == "" {
		errors.EmployeeIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (p *PurchaseOrder) validateStoreId(errors PurchaseOrderErrors) PurchaseOrderErrors {
	if p.StoreId == "" {
		errors.StoreIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (p *PurchaseOrder) validateJobId(errors PurchaseOrderErrors) PurchaseOrderErrors {
	if p.JobId == 0 {
		errors.JobIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (p *PurchaseOrder) validateDate(errors PurchaseOrderErrors) PurchaseOrderErrors {
	if p.Date == "" {
		errors.DateErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
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

type PurchaseOrderItemErrors struct {
	ItemNameErr   string
	QuantityErr   string
	ItemTypeIdErr string

	IsSuccessful bool
}

func (p *PurchaseOrderItem) Validate() PurchaseOrderItemErrors {
	errors := PurchaseOrderItemErrors{IsSuccessful: true}
	errors = p.validateItemName(errors)
	errors = p.validateQuantity(errors)
	errors = p.validateItemTypeId(errors)
	return errors
}

func (p *PurchaseOrderItem) validateItemName(errors PurchaseOrderItemErrors) PurchaseOrderItemErrors {
	if p.ItemName == "" {
		errors.ItemNameErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (p *PurchaseOrderItem) validateQuantity(errors PurchaseOrderItemErrors) PurchaseOrderItemErrors {
	if p.Quantity == 0 {
		errors.QuantityErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

func (p *PurchaseOrderItem) validateItemTypeId(errors PurchaseOrderItemErrors) PurchaseOrderItemErrors {
	if p.ItemTypeId == "" {
		errors.ItemTypeIdErr = "*required"
		errors.IsSuccessful = false
		return errors
	}
	return errors
}

type ItemType struct {
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
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
