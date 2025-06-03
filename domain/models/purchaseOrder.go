package domain

import "log"

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
	ItemSizeId      string `json:"item_size_id"`
	Quantity        int    `json:"quantity"`
	CreatedAt       string `json:"created_at"`
	ModifiedAt      string `json:"modified_at"`
}

type PurchaseOrderItemErrors struct {
	ItemNameErr   string
	QuantityErr   string
	ItemTypeIdErr string
	ItemSizeIdErr string

	IsSuccessful bool
}

func (p *PurchaseOrderItem) Validate() PurchaseOrderItemErrors {
	errors := &PurchaseOrderItemErrors{IsSuccessful: true}
	p.validateItemName(errors).validateQuantity(errors).validateItemTypeId(errors).validateItemSizeId(errors)
	log.Println(errors.IsSuccessful)
	return *errors
}

func (p *PurchaseOrderItem) validateItemName(errors *PurchaseOrderItemErrors) *PurchaseOrderItem {
	if p.ItemName == "" {
		errors.ItemNameErr = "*required"
		errors.IsSuccessful = false
	}
	return p
}

func (p *PurchaseOrderItem) validateQuantity(errors *PurchaseOrderItemErrors) *PurchaseOrderItem {
	if p.Quantity == 0 {
		errors.QuantityErr = "*required"
		errors.IsSuccessful = false
	}
	return p
}

func (p *PurchaseOrderItem) validateItemTypeId(errors *PurchaseOrderItemErrors) *PurchaseOrderItem {
	if p.ItemTypeId == "" {
		errors.ItemTypeIdErr = "*required"
		errors.IsSuccessful = false
	}
	return p
}

func (p *PurchaseOrderItem) validateItemSizeId(errors *PurchaseOrderItemErrors) *PurchaseOrderItem {
	if p.ItemSizeId == "" {
		errors.ItemSizeIdErr = "*required"
		errors.IsSuccessful = false
	}
	return p
}

type ItemSize struct {
	UUID        string `json:"uuid"`
	Size        string `json:"size"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
}

type ItemSizeErrors struct {
	SizeErr        string
	DescriptionErr string

	IsSuccessful bool
	SuccessMsg   string
}

func (i *ItemSize) Validate(existingSizes []ItemSize) ItemSizeErrors {
	errors := &ItemSizeErrors{IsSuccessful: true}
	i.validateSize(errors, existingSizes).validateDescription(errors)
	return *errors
}

func (i *ItemSize) validateSize(errors *ItemSizeErrors, existingSizes []ItemSize) *ItemSize {
	if i.Size == "" {
		errors.SizeErr = "*required"
		errors.IsSuccessful = false
		return i
	}
	return i
}

func (i *ItemSize) validateDescription(errors *ItemSizeErrors) *ItemSize {
	if i.Description == "" {
		errors.DescriptionErr = "*required"
		errors.IsSuccessful = false
	}
	return i
}

type ItemType struct {
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	ModifiedAt  string `json:"modified_at"`
}

type ItemTypeErrors struct {
	TypeErr        string
	DescriptionErr string

	IsSuccessful bool
	SuccessMsg   string
}

func (i *ItemType) Validate(existingTypes []ItemType) ItemTypeErrors {
	errors := &ItemTypeErrors{IsSuccessful: true}
	i.validateType(errors, existingTypes).validateDescription(errors)
	return *errors
}

func (i *ItemType) validateType(errors *ItemTypeErrors, existingTypes []ItemType) *ItemType {
	if i.Type == "" {
		errors.TypeErr = "*required"
		errors.IsSuccessful = false
		return i
	}
	return i
}

func (i *ItemType) validateDescription(errors *ItemTypeErrors) *ItemType {
	if i.Description == "" {
		errors.DescriptionErr = "*required"
		errors.IsSuccessful = false
	}
	return i
}

type Store struct {
	UUID         string `json:"uuid"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Suburb       string `json:"suburb"`
	City         string `json:"city"`
	AccountCode  string `json:"account_code"`
	CreatedAt    string `json:"created_at"`
	ModifiedAt   string `json:"modified_at"`
}

type StoreErrors struct {
	BusinessNameErr string
	EmailErr        string
	PhoneErr        string
	AddressErr      string
	SuburbErr       string
	CityErr         string
	AccountCodeErr  string

	IsSuccessful bool
	SuccessMsg   string
}

func (s *Store) Validate() StoreErrors {
	errors := &StoreErrors{IsSuccessful: true}
	s.validateBusinessName(errors).
		validateEmail(errors).
		validatePhone(errors).
		validateAddress(errors).
		validateSuburb(errors).
		validateCity(errors).
		validateAccountCode(errors)
	return *errors
}

func (s *Store) validateBusinessName(errors *StoreErrors) *Store {
	if s.BusinessName == "" {
		errors.BusinessNameErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validateEmail(errors *StoreErrors) *Store {
	if s.Email == "" {
		errors.EmailErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validatePhone(errors *StoreErrors) *Store {
	if s.Phone == "" {
		errors.PhoneErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validateSuburb(errors *StoreErrors) *Store {
	if s.Suburb == "" {
		errors.SuburbErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validateAddress(errors *StoreErrors) *Store {
	if s.Address == "" {
		errors.AddressErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validateCity(errors *StoreErrors) *Store {
	if s.City == "" {
		errors.CityErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}

func (s *Store) validateAccountCode(errors *StoreErrors) *Store {
	if s.AccountCode == "" {
		errors.AccountCodeErr = "*required"
		errors.IsSuccessful = false
	}
	return s
}
