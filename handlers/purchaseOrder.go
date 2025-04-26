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
	UUID       string `json:"uuid"`
	ItemName   string `json:"item_name"`
	ItemTypeId string `json:"item_type_id"`
	Quantity   int    `json:"quantity"`
}

var purchaseOrderItemList = []purchaseOrderItem{
	{
		UUID:       "123",
		ItemName:   "15L expressions low sheen monument",
		ItemTypeId: "123",
		Quantity:   2,
	},
	{
		UUID:       "234",
		ItemName:   "10L enamel antique white",
		ItemTypeId: "123",
		Quantity:   1,
	},
	{
		UUID:       "345",
		ItemName:   "Haymes elite sash cutter",
		ItemTypeId: "234",
		Quantity:   6,
	},
}

type purchaseOrder struct {
	UUID               string `json:"uuid"`
	EmployeeId         string `json:"employee_id"`
	StoreId            string `json:"store_id"`
	JobId              int    `json:"job_id"`
	Date               string `json:"date"`
	PurchaseOrderItems []purchaseOrderItem

	Stores    []store
	Jobs      []models.Job
	ItemTypes []itemType
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

				Stores:    storeList,
				Jobs:      jobs,
				ItemTypes: itemTypeList,
			}

			component := "purchaseOrder"
			title := "Purchase Order - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}
