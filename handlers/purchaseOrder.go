package handlers

import (
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
	models "github.com/jtalev/chat_gpg/domain/models"
)

type store struct {
	StoreId   string
	StoreName string
}

var storeList = []store{
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

type itemType struct {
	UUID string `json:"uuid"`
	Type string `json:"type"`
}

var itemTypeList = []itemType{
	{
		UUID: "123",
		Type: "Paint",
	},
	{
		UUID: "234",
		Type: "Accessory",
	},
}

type purchaseOrderItem struct {
	UUID         string `json:"uuid"`
	ItemName     string `json:"item_name"`
	ItemTypeId   string `json:"item_type_id"`
	ItemTypeName string `json:"item_type_name"`
	Quantity     int    `json:"quantity"`

	ItemTypes []itemType
}

var purchaseOrderItemList = []purchaseOrderItem{
	{
		UUID:         "123",
		ItemName:     "15L expressions low sheen monument",
		ItemTypeId:   "123",
		ItemTypeName: "Paint",
		Quantity:     2,

		ItemTypes: itemTypeList,
	},
	{
		UUID:         "234",
		ItemName:     "10L enamel antique white",
		ItemTypeId:   "123",
		ItemTypeName: "Paint",
		Quantity:     1,

		ItemTypes: itemTypeList,
	},
	{
		UUID:         "345",
		ItemName:     "Haymes elite sash cutter",
		ItemTypeId:   "234",
		ItemTypeName: "Accessory",
		Quantity:     6,

		ItemTypes: itemTypeList,
	},
}

type purchaseOrder struct {
	UUID               string `json:"uuid"`
	EmployeeId         string `json:"employee_id"`
	StoreId            string `json:"store_id"`
	JobId              int    `json:"job_id"`
	Date               string `json:"date"`
	PurchaseOrderItems []purchaseOrderItem

	Stores []store
	Jobs   []models.Job
}

func (h *Handler) ServePurchaseOrderView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Printf("error fetching employee id: %v", err)
				http.Error(w, "error fetching employee id, unauthorized", http.StatusUnauthorized)
				return
			}

			jobs, err := application.GetJobs(h.DB)
			if err != nil {
				log.Printf("error fetching jobs: %v", err)
				http.Error(w, "error fetching jobs, internal server error", http.StatusInternalServerError)
				return
			}

			data := purchaseOrder{
				EmployeeId:         employeeId,
				StoreId:            "345",
				JobId:              3,
				PurchaseOrderItems: purchaseOrderItemList,

				Stores: storeList,
				Jobs:   jobs,
			}

			component := "purchaseOrder"
			title := "Purchase Order - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeItemRow() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			item := purchaseOrderItem{
				ItemTypes: itemTypeList,
			}
			err := executePartialTemplate(purchaseOrderItemRowPath, "purchaseOrderItemRow", item, w)
			if err != nil {
				log.Printf("error executing purchaseOrderItemRow.html: %v", err)
				http.Error(w, "error executing purchaseOrderItemRow.html, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}
