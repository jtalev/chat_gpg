package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
)

var itemTypeList = []application.ItemType{
	{
		UUID: "123",
		Type: "Paint",
	},
	{
		UUID: "234",
		Type: "Accessory",
	},
}

var purchaseOrderItemList = []application.PurchaseOrderItem{
	{
		UUID:       "123",
		ItemName:   "15L expressions low sheen monument",
		ItemTypeId: "123",
		Quantity:   2,

		ItemTypes: itemTypeList,
	},
	{
		UUID:       "234",
		ItemName:   "10L enamel antique white",
		ItemTypeId: "123",
		Quantity:   1,

		ItemTypes: itemTypeList,
	},
	{
		UUID:       "345",
		ItemName:   "Haymes elite sash cutter",
		ItemTypeId: "234",
		Quantity:   6,

		ItemTypes: itemTypeList,
	},
}
var order = application.PurchaseOrder{}

func (h *Handler) ServePurchaseOrderView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order.Reset()
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

			order.PopulateStores(h.DB)

			data := application.PurchaseOrder{
				EmployeeId: employeeId,

				Stores: order.Stores,
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
			item := application.PurchaseOrderItem{
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

func (h *Handler) PostPurchaseOrder() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order.Reset()
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("error reading request body: %v", err)
				http.Error(w, "error reading request body, bad request", http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(body, &order)
			if err != nil {
				log.Printf("error unmarshalling json: %v", err)
				http.Error(w, "error unmarshalling json, bad request", http.StatusBadRequest)
				return
			}

			purchaseOrderErrors, err := order.PostPurchaseOrder(h.DB)
			if err != nil {
				log.Printf("error posting purchase order: %v", err)
				http.Error(w, "error posting purchase order, internal server error", http.StatusInternalServerError)
				return
			}
			order.Errors = purchaseOrderErrors

			tmpl, err := template.ParseFiles(
				purchaseOrderFormPath,
				purchaseOrderItemRowPath,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Error(w, "error parsing template files, internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "purchaseOrderForm", order)
			if err != nil {
				log.Printf("error executing purchaseOrderForm.html: %v", err)
				http.Error(w, "error executing purchaseOrderForm.html, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) PostItem() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var itemType application.ItemType
			if ok := h.DecodeJson(&itemType, w, r); !ok {
				return
			}
			log.Println(itemType)
		},
	)
}

func (h *Handler) PostStore() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var store application.Store
			if ok := h.DecodeJson(&store, w, r); !ok {
				return
			}
			err := application.PostStore(store, h.DB)
			if err != nil {
				log.Printf("error posting store: %v", err)
				http.Error(w, "error posting store, bad request", http.StatusBadRequest)
				return
			}
			err = a.ServeSingleTemplate(adminAddStoreModalPath, "addStoreModal", nil, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) ServeEmployeePOHistory() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order.Reset()
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Printf("error fetching employee ID: %v", err)
				http.Error(w, "error fetching employee ID, unauthorized", http.StatusUnauthorized)
				return
			}
			purchaseOrders, err := order.FetchEmployeeHistory(employeeId, h.DB)

			log.Println(purchaseOrders)
		},
	)
}
