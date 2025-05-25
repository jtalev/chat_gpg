package handlers

import (
	"encoding/json"
	"html/template"
	"io"
	"log"
	"net/http"

	application "github.com/jtalev/chat_gpg/application/services"
)

func (h *Handler) ServePurchaseOrderView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order := application.PurchaseOrder{}
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

func (h *Handler) ServePurchaseOrderForm() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order := application.PurchaseOrder{}
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

			templatePaths := []string{purchaseOrderFormPath, purchaseOrderItemRowPath}
			err = h.ServeMultiTemplate(templatePaths, "purchaseOrderForm", data, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) ServeEmployeePOHistory() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order := application.PurchaseOrder{}
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Printf("error fetching employee ID: %v", err)
				http.Error(w, "error fetching employee ID, unauthorized", http.StatusUnauthorized)
				return
			}
			purchaseOrders, err := order.FetchEmployeeHistory(employeeId, h.DB)
			if err != nil {
				log.Printf("error fetching employee purchase order history: %v", err)
				http.Error(w, "error fetching employee purchase order history, status unauthorized", http.StatusUnauthorized)
				return
			}
			err = h.ServeSingleTemplate(employeePurchaseOrderHistoryPath, "employeePurchaseOrderHistory", purchaseOrders, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) ServePurchaseOrder() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			vals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			purchaseOrder, err := application.GetPurchaseOrder(vals[0], h.DB)
			if err != nil {
				log.Printf("error getting purchase order, :%v", err)
				http.Error(w, "error getting purchase order, internal server error", http.StatusInternalServerError)
				return

			}

			err = h.ServeSingleTemplate(viewPurchaseOrderModalPath, "viewPurchaseOrderModal", purchaseOrder, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) ServeItemRow() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			order := application.PurchaseOrder{}
			log.Println(order)
			order.PopulateItemTypes(h.DB)
			item := application.PurchaseOrderItem{}
			if len(order.PurchaseOrderItems) > 0 {
				item.ItemTypes = order.PurchaseOrderItems[0].ItemTypes
			}
			// log.Println(item)
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
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("error reading request body: %v", err)
				http.Error(w, "error reading request body, bad request", http.StatusBadRequest)
				return
			}

			err = json.Unmarshal(body, &h.PurchaseOrderService)
			if err != nil {
				log.Printf("error unmarshalling json: %v", err)
				http.Error(w, "error unmarshalling json, bad request", http.StatusBadRequest)
				return
			}

			purchaseOrderErrors, err := h.PurchaseOrderService.PostPurchaseOrder(h.DB)
			if err != nil {
				log.Printf("error posting purchase order: %v", err)
				http.Error(w, "error posting purchase order, internal server error", http.StatusInternalServerError)
				return
			}
			h.PurchaseOrderService.Errors = purchaseOrderErrors

			tmpl, err := template.ParseFiles(
				purchaseOrderFormPath,
				purchaseOrderItemRowPath,
			)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				http.Error(w, "error parsing template files, internal server error", http.StatusInternalServerError)
				return
			}

			err = tmpl.ExecuteTemplate(w, "purchaseOrderForm", h.PurchaseOrderService)
			if err != nil {
				log.Printf("error executing purchaseOrderForm.html: %v", err)
				http.Error(w, "error executing purchaseOrderForm.html, internal server error", http.StatusInternalServerError)
				return
			}
		},
	)
}

func (h *Handler) DeletePurchaseOrder() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			p := application.PurchaseOrder{}
			rVals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}
			p.UUID = rVals[0]

			err = p.DeletePurchaseOrder(h.DB)
			if err != nil {
				log.Printf("error deleting purchase order: %v", err)
				http.Error(w, "error deleting purchase order, not found", http.StatusNotFound)
				return
			}
		},
	)
}

func (h *Handler) PostItemType() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var itemType application.ItemType
			if ok := h.DecodeJson(&itemType, w, r); !ok {
				return
			}
			outItemType, err := application.PostItemType(itemType, h.DB)
			if err != nil {
				log.Printf("error posting item type: %v", err)
				http.Error(w, "error posting item type, bad request", http.StatusBadRequest)
				return
			}
			err = h.ServeSingleTemplate(adminAddItemModalPath, "addItemModal", outItemType, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) PutItemType() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var itemType application.ItemType
			if ok := h.DecodeJson(&itemType, w, r); !ok {
				return
			}

			outItemType, err := application.PutItemType(itemType, h.DB)
			if err != nil {
				log.Printf("error updating store: %v", err)
				http.Error(w, "error updating store, bad request", http.StatusBadRequest)
				return
			}

			err = h.ServeSingleTemplate(adminAddItemModalPath, "addItemModal", outItemType, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) DeleteItemType() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rVals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			err = application.DeleteItemType(rVals[0], h.DB)
			if err != nil {
				log.Printf("error deleting item type: %v", err)
				http.Error(w, "error deleting item type, not found", http.StatusNotFound)
				return
			}
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
			outStore, err := application.PostStore(store, h.DB)
			if err != nil {
				log.Printf("error posting store: %v", err)
				http.Error(w, "error posting store, bad request", http.StatusBadRequest)
				return
			}
			err = h.ServeSingleTemplate(adminAddStoreModalPath, "addStoreModal", outStore, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) PutStore() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var store application.Store
			if ok := h.DecodeJson(&store, w, r); !ok {
				return
			}

			outStore, err := application.PutStore(store, h.DB)
			if err != nil {
				log.Printf("error updating store: %v", err)
				http.Error(w, "error updating store, bad request", http.StatusBadRequest)
				return
			}

			err = h.ServeSingleTemplate(adminAddStoreModalPath, "addStoreModal", outStore, w)
			if err != nil {
				return
			}
		},
	)
}

func (h *Handler) DeleteStore() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			rVals, err := parseRequestValues([]string{"uuid"}, r)
			if err != nil {
				log.Printf("error parsing request values: %v", err)
				http.Error(w, "error parsing request values, bad request", http.StatusBadRequest)
				return
			}

			err = application.DeleteStore(rVals[0], h.DB)
			if err != nil {
				log.Printf("error deleting store: %v", err)
				http.Error(w, "error deleting store, not found", http.StatusNotFound)
				return
			}
		},
	)
}
