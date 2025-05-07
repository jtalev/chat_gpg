package application

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/google/uuid"
	domain "github.com/jtalev/chat_gpg/domain/models"
	models "github.com/jtalev/chat_gpg/domain/models"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
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
	UUID       string `json:"uuid"`
	ItemName   string `json:"item_name"`
	ItemTypeId string `json:"item_type_id"`
	Quantity   int    `json:"quantity"`

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

	Errors domain.PurchaseOrderErrors
}

func (o *PurchaseOrder) Reset() {
	o = &PurchaseOrder{}
}

func (o *PurchaseOrder) mapPurchaseOrder() domain.PurchaseOrder {
	o.UUID = uuid.New().String()
	return domain.PurchaseOrder{
		UUID:       o.UUID,
		EmployeeId: o.EmployeeId,
		StoreId:    o.StoreId,
		JobId:      o.JobId,
		Date:       o.Date,
	}
}

func (o *PurchaseOrder) mapPurchaseOrderItems() []domain.PurchaseOrderItem {
	items := make([]domain.PurchaseOrderItem, len(o.PurchaseOrderItems))
	for i, item := range o.PurchaseOrderItems {
		item.UUID = uuid.New().String()
		items[i] = domain.PurchaseOrderItem{
			UUID:            item.UUID,
			PurchaseOrderId: o.UUID,
			ItemName:        item.ItemName,
			ItemTypeId:      item.ItemTypeId,
			Quantity:        item.Quantity,
		}
	}
	return items
}

func getJobs(db *sql.DB) ([]domain.Job, error) {
	jobs, err := repo.GetJobs(db)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

var itemTypeList = []ItemType{
	{
		UUID: "123",
		Type: "Paint",
	},
	{
		UUID: "234",
		Type: "Accessory",
	},
}

// TODO: once item repo is functional, get item types
func (o *PurchaseOrder) populateItemTypes() {
	for i := range o.PurchaseOrderItems {
		o.PurchaseOrderItems[i].ItemTypes = itemTypeList
	}
}

var storeList = []Store{
	{
		StoreId:   "123",
		StoreName: "Haymes Geelong West",
	},
	{
		StoreId:   "234",
		StoreName: "Dulux Geelong West",
	},
	{
		StoreId:   "345",
		StoreName: "Haymes Ocean Grove",
	},
}

// TODO: once stores repo is functional, get stores
func (o *PurchaseOrder) populateStores() {
	o.Stores = storeList
}

func (o *PurchaseOrder) initRequiredViewData(db *sql.DB) error {
	jobs, err := getJobs(db)
	if err != nil {
		return err
	}
	o.Jobs = jobs
	o.populateItemTypes()
	o.populateStores()
	return nil
}

func validatePurchaseOrderAndItems(purchaseOrder domain.PurchaseOrder, purchaseOrderItems []domain.PurchaseOrderItem) (domain.PurchaseOrderErrors, domain.PurchaseOrderItemErrors) {
	purchaseOrderErrors := purchaseOrder.Validate()
	itemErrors := domain.PurchaseOrderItemErrors{}
	for _, item := range purchaseOrderItems {
		itemErrors = item.Validate()
		if itemErrors.IsSuccessful == false {
			purchaseOrderErrors.IsSuccessful = false
			purchaseOrderErrors.PurchaseOrderItemsErr = "*item fields must not be empty"
		}
	}
	return purchaseOrderErrors, itemErrors
}

func parallelPostPurchaseOrder(purchaseOrder domain.PurchaseOrder, purchaseOrderItems []domain.PurchaseOrderItem, db *sql.DB) chan error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.PostPurchaseOrder(purchaseOrder, db); err != nil {
			errChan <- fmt.Errorf("purchase order post failed: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := repo.BatchPostPurchaseOrderItem(purchaseOrderItems, db); err != nil {
			errChan <- fmt.Errorf("purchase order items post failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	return errChan
}

func (o *PurchaseOrder) PostPurchaseOrder(db *sql.DB) (domain.PurchaseOrderErrors, error) {
	purchaseOrder := o.mapPurchaseOrder()
	purchaseOrderItems := o.mapPurchaseOrderItems()

	purchaseOrderErrors, _ := validatePurchaseOrderAndItems(purchaseOrder, purchaseOrderItems)

	// the data initialised here is required to allow view features
	err := o.initRequiredViewData(db)
	if err != nil {
		return purchaseOrderErrors, err
	}

	if !purchaseOrderErrors.IsSuccessful {
		return purchaseOrderErrors, nil
	}

	errChan := parallelPostPurchaseOrder(purchaseOrder, purchaseOrderItems, db)
	for err := range errChan {
		if err != nil {
			purchaseOrderErrors.IsSuccessful = false
			purchaseOrderErrors.SuccessMsg = err.Error()
			return purchaseOrderErrors, err
		}
	}

	// send email to store
	e := EmailSender{
		SenderName:       "admin",
		SenderEmail:      "j.talev@outlook.com",
		RecipientName:    "Josh Talev",
		RecipientEmail:   "j.talev@outlook.com",
		Subject:          "test",
		PlainTextContent: "hello from chat_gpg",
	}
	err = e.SendEmail()
	if err != nil {
		purchaseOrderErrors.IsSuccessful = false
		purchaseOrderErrors.SuccessMsg = err.Error()
		return purchaseOrderErrors, err
	}

	purchaseOrderErrors.SuccessMsg = "Purchase order submitted successfully."
	return purchaseOrderErrors, nil
}

func (o *PurchaseOrder) GetPurchaseOrders(db *sql.DB) ([]models.PurchaseOrder, error) {
	purchaseOrders, err := repo.GetPurchaseOrders(db)
	if err != nil {
		return nil, err
	}
	return purchaseOrders, nil
}
