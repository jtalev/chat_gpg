package application

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/google/uuid"
	domain "github.com/jtalev/chat_gpg/domain/models"
	models "github.com/jtalev/chat_gpg/domain/models"
	repo "github.com/jtalev/chat_gpg/infrastructure/repository"
	"github.com/jtalev/chat_gpg/internal/task_queue"
)

type Store struct {
	UUID         string `json:"uuid"`
	BusinessName string `json:"business_name"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	Address      string `json:"address"`
	Suburb       string `json:"suburb"`
	City         string `json:"city"`

	Errors     models.StoreErrors
	ModalTitle string
}

type ItemType struct {
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
	Description string `json:"description"`

	Errors     models.ItemTypeErrors
	ModalTitle string
}

type PurchaseOrderItem struct {
	UUID       string `json:"uuid"`
	ItemName   string `json:"item_name"`
	ItemTypeId string `json:"item_type_id"`
	ItemType   string
	Quantity   int `json:"quantity"`

	ItemTypes []ItemType
}

type PurchaseOrder struct {
	UUID               string              `json:"uuid"`
	EmployeeId         string              `json:"employee_id"`
	StoreId            string              `json:"store_id"`
	JobId              int                 `json:"job_id"`
	Date               string              `json:"date"`
	PurchaseOrderItems []PurchaseOrderItem `json:"purchase_order_items"`

	Employee string
	Store    models.Store
	Job      string

	Stores []Store
	Jobs   []models.Job

	Errors domain.PurchaseOrderErrors

	TaskProducer *task_queue.TaskProducer
}

func (p *PurchaseOrder) mapPurchaseOrder() domain.PurchaseOrder {
	p.UUID = uuid.New().String()
	return domain.PurchaseOrder{
		UUID:       p.UUID,
		EmployeeId: p.EmployeeId,
		StoreId:    p.StoreId,
		JobId:      p.JobId,
		Date:       p.Date,
	}
}

func (p *PurchaseOrder) mapPurchaseOrderItems() []domain.PurchaseOrderItem {
	items := make([]domain.PurchaseOrderItem, len(p.PurchaseOrderItems))
	for i, item := range p.PurchaseOrderItems {
		item.UUID = uuid.New().String()
		items[i] = domain.PurchaseOrderItem{
			UUID:            item.UUID,
			PurchaseOrderId: p.UUID,
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

func GetItemTypes(db *sql.DB) ([]models.ItemType, error) {
	itemTypes, err := repo.GetItemTypes(db)
	if err != nil {
		return nil, err
	}
	return itemTypes, nil
}

func (p *PurchaseOrder) PopulateItemTypes(db *sql.DB) error {
	itemTypeList, err := GetItemTypes(db)
	if err != nil {
		return err
	}
	outItemTypes := make([]ItemType, len(itemTypeList))
	for i, itemType := range itemTypeList {
		outItemTypes[i].UUID = itemType.UUID
		outItemTypes[i].Type = itemType.Type
		outItemTypes[i].Description = itemType.Description
	}
	if len(p.PurchaseOrderItems) == 0 {
		p.PurchaseOrderItems = append(p.PurchaseOrderItems, PurchaseOrderItem{
			ItemTypes: outItemTypes,
		})
		return nil
	}
	for i := range p.PurchaseOrderItems {
		p.PurchaseOrderItems[i].ItemTypes = outItemTypes
	}
	return nil
}

func GetStores(db *sql.DB) ([]models.Store, error) {
	storeList, err := repo.GetStores(db)
	if err != nil {
		log.Printf("error getting stores: %v", err)
		return nil, err
	}
	return storeList, nil
}

func (p *PurchaseOrder) PopulateStores(db *sql.DB) error {
	storeList, err := GetStores(db)
	if err != nil {
		return err
	}
	outStores := make([]Store, len(storeList))
	for i, store := range storeList {
		outStores[i].UUID = store.UUID
		outStores[i].BusinessName = store.BusinessName
		outStores[i].Address = store.Address
	}
	p.Stores = outStores
	return nil
}

func (p *PurchaseOrder) initRequiredViewData(db *sql.DB) error {
	jobs, err := getJobs(db)
	if err != nil {
		return err
	}
	p.Jobs = jobs
	p.PopulateItemTypes(db)
	p.PopulateStores(db)
	return nil
}

func validatePurchaseOrderAndItems(purchaseOrder domain.PurchaseOrder, purchaseOrderItems []domain.PurchaseOrderItem) (domain.PurchaseOrderErrors, domain.PurchaseOrderItemErrors) {
	purchaseOrderErrors := purchaseOrder.Validate()
	itemErrors := domain.PurchaseOrderItemErrors{}
	if len(purchaseOrderItems) == 0 {
		purchaseOrderErrors.IsSuccessful = false
		purchaseOrderErrors.PurchaseOrderItemsErr = "*must include at least 1 item to submit"
		return purchaseOrderErrors, itemErrors
	}
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

func generateEmailBody(purchaseOrder models.PurchaseOrder, purchaseOrderItems []models.PurchaseOrderItem, db *sql.DB) (string, error) {
	employee, err := repo.GetEmployeeByEmployeeId(purchaseOrder.EmployeeId, db)
	if err != nil {
		return "", err
	}
	store, err := repo.GetStoreByUuid(purchaseOrder.StoreId, db)
	if err != nil {
		return "", err
	}
	job, err := repo.GetJobById(purchaseOrder.JobId, db)
	if err != nil {
		return "", err
	}
	emailBody := fmt.Sprintf(
		"Purchase Order - Geelong Paint Group\n"+
			"Date: %s\n"+
			"Store: %s\n"+
			"Ordered By: %s %s\n"+
			"Reference: %s\n"+
			"\n"+
			"Items:\n",
		purchaseOrder.Date,
		store.BusinessName,
		employee.FirstName, employee.LastName,
		job.Address,
	)
	for _, item := range purchaseOrderItems {
		emailBody += fmt.Sprintf("%v x %s\n", item.Quantity, item.ItemName)
	}
	return emailBody, nil
}

func generateEmailPayload(purchaseOrder models.PurchaseOrder, purchaseOrderItems []models.PurchaseOrderItem, db *sql.DB) (task_queue.EmailHandler, error) {
	e := task_queue.EmailHandler{}
	store, err := repo.GetStoreByUuid(purchaseOrder.StoreId, db)
	if err != nil {
		return e, err
	}
	emailBody, err := generateEmailBody(purchaseOrder, purchaseOrderItems, db)
	if err != nil {
		return e, err
	}
	env := os.Getenv("ENV")
	recipientEmail := ""
	if env == "development" {
		recipientEmail = "j.talev@outlook.com"
	} else if env == "production" {
		recipientEmail = store.Email
	}
	e = task_queue.CreateEmailPayload(
		"Admin@GeelongPaintGroup",
		"admin@geelongpaintgroup.com.au",
		store.BusinessName,
		recipientEmail,
		"Purchase Order - Geelong Paint Group",
		emailBody,
		"",
	)
	return e, nil
}

func (p *PurchaseOrder) PostPurchaseOrder(db *sql.DB) (domain.PurchaseOrderErrors, error) {
	purchaseOrder := p.mapPurchaseOrder()
	purchaseOrderItems := p.mapPurchaseOrderItems()

	purchaseOrderErrors, _ := validatePurchaseOrderAndItems(purchaseOrder, purchaseOrderItems)

	// the data initialised here is required to allow view features
	err := p.initRequiredViewData(db)
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

	e, err := generateEmailPayload(purchaseOrder, purchaseOrderItems, db)
	if err != nil {
		return purchaseOrderErrors, err
	}

	p.TaskProducer.Enqueue("one_time", "send_email", e)

	purchaseOrderErrors.SuccessMsg = "Purchase order submitted successfully."
	return purchaseOrderErrors, nil
}

func (p *PurchaseOrder) DeletePurchaseOrder(db *sql.DB) error {
	err := repo.DeletePurchaseOrder(p.UUID, db)
	if err != nil {
		return err
	}
	return nil
}

func filterEmployeeOrders(orders []models.PurchaseOrder, employeeId string) []models.PurchaseOrder {
	filteredOrders := []models.PurchaseOrder{}
	for _, order := range orders {
		if order.EmployeeId == employeeId {
			filteredOrders = append(filteredOrders, order)
		}
	}
	return filteredOrders
}

func mapStoreNameToOrder(inOrders []models.PurchaseOrder, outOrders []PurchaseOrder, db *sql.DB) {
	if len(inOrders) != len(outOrders) {
		log.Println("inOrders and outOrders should be same length")
	}
	var store models.Store
	var err error
	for i, order := range inOrders {
		store, err = repo.GetStoreByUuid(order.StoreId, db)
		if err != nil {
			log.Println("error getting job")
			return
		}
		outOrders[i].Store = store
	}
}

func mapJobNameToOrder(inOrders []models.PurchaseOrder, outOrders []PurchaseOrder, db *sql.DB) {
	if len(inOrders) != len(outOrders) {
		log.Println("inOrders and outOrders should be same length")
	}
	var job models.Job
	var err error
	for i, order := range inOrders {
		job, err = repo.GetJobById(order.JobId, db)
		if err != nil {
			log.Println("error getting job")
			return
		}
		outOrders[i].Job = job.Name
	}
}

func mapUUIDToOrder(inOrders []models.PurchaseOrder, outOrders []PurchaseOrder, db *sql.DB) {
	if len(inOrders) != len(outOrders) {
		log.Println("inOrders and outOrders should be same length")
	}
	for i, order := range inOrders {
		outOrders[i].UUID = order.UUID
	}
}

func (p *PurchaseOrder) FetchEmployeeHistory(employeeId string, db *sql.DB) ([]PurchaseOrder, error) {
	purchaseOrders, err := repo.GetPurchaseOrders(db)
	if err != nil {
		return nil, err
	}
	filteredOrders := filterEmployeeOrders(purchaseOrders, employeeId)
	outOrders := make([]PurchaseOrder, len(filteredOrders))
	mapStoreNameToOrder(filteredOrders, outOrders, db)
	mapJobNameToOrder(filteredOrders, outOrders, db)
	mapUUIDToOrder(filteredOrders, outOrders, db)

	return outOrders, nil
}

func (p *PurchaseOrder) GetPurchaseOrders(db *sql.DB) ([]PurchaseOrder, error) {
	purchaseOrders, err := repo.GetPurchaseOrders(db)
	if err != nil {
		return nil, err
	}
	outOrders := make([]PurchaseOrder, len(purchaseOrders))
	mapStoreNameToOrder(purchaseOrders, outOrders, db)
	mapJobNameToOrder(purchaseOrders, outOrders, db)
	mapUUIDToOrder(purchaseOrders, outOrders, db)
	return outOrders, nil
}

func mapStore(store Store, uuid string) models.Store {
	return models.Store{
		UUID:         uuid,
		BusinessName: store.BusinessName,
		Email:        store.Email,
		Phone:        store.Phone,
		Address:      store.Address,
		Suburb:       store.Suburb,
		City:         store.City,
	}
}

func PostStore(store Store, db *sql.DB) (Store, error) {
	store.ModalTitle = "Add Store"
	uuid := uuid.New().String()
	modelStore := mapStore(store, uuid)
	errors := modelStore.Validate()
	if !errors.IsSuccessful {
		store.Errors = errors
		return store, nil
	}
	err := repo.PostStore(modelStore, db)
	if err != nil {
		return Store{}, err
	}
	store.Errors.SuccessMsg = "Store submitted successfully."
	return store, nil
}

func PutStore(store Store, db *sql.DB) (Store, error) {
	store.ModalTitle = "Update Store"
	modelStore := mapStore(store, store.UUID)
	errors := modelStore.Validate()
	if !errors.IsSuccessful {
		store.Errors = errors
		return store, nil
	}
	err := repo.PutStore(modelStore, db)
	if err != nil {
		return store, err
	}
	store.Errors.SuccessMsg = "Store updated successfully."
	return store, nil
}

func DeleteStore(uuid string, db *sql.DB) error {
	err := repo.DeleteStore(uuid, db)
	if err != nil {
		return err
	}
	return nil
}

func mapItemType(itemType ItemType, uuid string) models.ItemType {
	return models.ItemType{
		UUID:        uuid,
		Type:        itemType.Type,
		Description: itemType.Description,
	}
}

func PostItemType(itemType ItemType, db *sql.DB) (ItemType, error) {
	itemType.ModalTitle = "Add Item Type"
	uuid := uuid.New().String()
	modelItemType := mapItemType(itemType, uuid)
	existingTypes, err := repo.GetItemTypes(db)
	if err != nil {
		return itemType, err
	}
	itemType.Errors = modelItemType.Validate(existingTypes)
	if !itemType.Errors.IsSuccessful {
		return itemType, nil
	}
	err = repo.PostItemType(modelItemType, db)
	if err != nil {
		return ItemType{}, err
	}
	itemType.Errors.SuccessMsg = "Item type submitted successfully."
	return itemType, nil
}

func PutItemType(itemType ItemType, db *sql.DB) (ItemType, error) {
	itemType.ModalTitle = "Update Item Type"
	modelItemType := mapItemType(itemType, itemType.UUID)
	existingTypes, err := repo.GetItemTypes(db)
	if err != nil {
		return itemType, err
	}
	errors := modelItemType.Validate(existingTypes)
	if !errors.IsSuccessful {
		itemType.Errors = errors
		return itemType, nil
	}
	err = repo.PutItemType(modelItemType, db)
	if err != nil {
		return itemType, err
	}
	itemType.Errors.SuccessMsg = "Item type updated successfully."
	return itemType, nil
}

func DeleteItemType(uuid string, db *sql.DB) error {
	err := repo.DeleteItemType(uuid, db)
	if err != nil {
		return err
	}
	return nil
}

func (p *PurchaseOrder) mapModelToPurchaseOrder(model models.PurchaseOrder) {
	p.EmployeeId = model.EmployeeId
	p.UUID = model.UUID
	p.Date = model.Date
	p.StoreId = model.StoreId
}

func (p *PurchaseOrder) mapStoreToOrder(db *sql.DB) error {
	store, err := repo.GetStoreByUuid(p.StoreId, db)
	if err != nil {
		return err
	}
	p.Store = store
	return nil
}

func (p *PurchaseOrder) mapEmployeeName(db *sql.DB) error {
	employee, err := repo.GetEmployeeByEmployeeId(p.EmployeeId, db)
	if err != nil {
		return err
	}
	p.Employee = fmt.Sprintf("%s %s", employee.FirstName, employee.LastName)
	return nil
}

func (p *PurchaseOrder) mapItems(db *sql.DB) error {
	items, err := repo.GetItemsByOrderUuid(p.UUID, db)
	if err != nil {
		return err
	}
	p.PurchaseOrderItems = make([]PurchaseOrderItem, len(items))
	for i, item := range items {
		p.PurchaseOrderItems[i].ItemName = item.ItemName
		p.PurchaseOrderItems[i].ItemTypeId = item.ItemTypeId
		p.PurchaseOrderItems[i].Quantity = item.Quantity
	}
	return nil
}

func (p *PurchaseOrder) mapItemTypeToItems(db *sql.DB) error {
	for i, item := range p.PurchaseOrderItems {
		itemType, err := repo.GetItemTypeByUuid(item.ItemTypeId, db)
		if err != nil {
			return err
		}
		p.PurchaseOrderItems[i].ItemType = itemType.Type
	}
	return nil
}

func GetPurchaseOrder(uuid string, db *sql.DB) (PurchaseOrder, error) {
	var p PurchaseOrder
	pm, err := repo.GetPurchaseOrderByUuid(uuid, db)
	if err != nil {
		return p, err
	}

	p.mapModelToPurchaseOrder(pm)
	err = p.mapStoreToOrder(db)
	p.mapEmployeeName(db)
	p.mapItems(db)
	p.mapItemTypeToItems(db)
	if err != nil {
		return p, err
	}
	return p, nil
}
