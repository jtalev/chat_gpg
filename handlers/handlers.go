package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/gorilla/sessions"
	application "github.com/jtalev/chat_gpg/application/services"
	"github.com/jtalev/chat_gpg/application/services/jobnotes"
	"github.com/jtalev/chat_gpg/application/services/safety"
	domain "github.com/jtalev/chat_gpg/domain/models"
	infrastructure "github.com/jtalev/chat_gpg/infrastructure/repository"
	"github.com/jtalev/chat_gpg/internal/task_queue"

	"go.uber.org/zap"
)

type Handler struct {
	DB     *sql.DB
	Store  *sessions.CookieStore
	Logger *zap.SugaredLogger

	TaskProducer         *task_queue.TaskProducer
	LeaveService         *application.LeaveService
	PurchaseOrderService *application.PurchaseOrder
	jobnotes             *jobnotes.Jobnotes

	noteRepo jobnotes.NoteRepo
}

func (h *Handler) StartWorkers() error {
	oneTimeQueue := make(chan task_queue.Task, 30)
	for i := 0; i < 4; i++ {
		w := task_queue.InitWorker(context.Background(), h.DB)
		go w.StartWorkerLoop(oneTimeQueue)
	}

	scheduledQueue := make(chan task_queue.Task, 10)
	for i := 0; i < 2; i++ {
		w := task_queue.InitWorker(context.Background(), h.DB)
		go w.StartWorkerLoop(scheduledQueue)
	}

	tp := task_queue.TaskProducer{
		Db:             h.DB,
		OneTimeQueue:   oneTimeQueue,
		ScheduledQueue: scheduledQueue,
	}

	err := tp.InitDbBackupTask()
	if err != nil {
		return err
	}

	err = tp.InitQueues()
	if err != nil {
		log.Printf("failed to initialize persisted tasks: %v", err)
		return err
	}

	h.TaskProducer = &tp
	return nil
}

func (h *Handler) RegisterRepos() {
	jobnotesRepo := infrastructure.JobnotesRepo{
		Db: h.DB,
	}
	h.noteRepo = &jobnotesRepo
}

func (h *Handler) RegisterServices() {
	leaveService := application.LeaveService{
		TaskProducer: h.TaskProducer,
	}

	purchaseOrderService := application.PurchaseOrder{
		TaskProducer: h.TaskProducer,
	}

	jobnotes := jobnotes.Jobnotes{
		Repo: h.noteRepo,
	}

	h.LeaveService = &leaveService
	h.PurchaseOrderService = &purchaseOrderService
	h.jobnotes = &jobnotes
}

const (
	layoutPath                     = "../ui/layouts/layout.html"
	navPath                        = "../ui/templates/nav.html"
	dashboardPath                  = "../ui/views/dashboard.html"
	dashboardTilePath              = "../ui/templates/dashboardTile.html"
	jobsPath                       = "../ui/views/jobs.html"
	timesheetsPath                 = "../ui/views/timesheets.html"
	timesheetTablePath             = "../ui/templates/timesheetTable.html"
	timesheetNavContainerPath      = "../ui/templates/timesheetNavContainer.html"
	timesheetHeadPath              = "../ui/templates/timesheetHead.html"
	existingTimesheetRowPath       = "../ui/templates/existingTimesheetRow.html"
	jobSelectModalPath             = "../ui/templates/jobSelectModal.html"
	leavePath                      = "../ui/views/leave.html"
	leaveHistoryPath               = "../ui/templates/leaveHistory.html"
	leaveFormPath                  = "../ui/templates/leaveForm.html"
	employeeLeaveRequestPath       = "../ui/templates/employeeLeaveRequest.html"
	employeeLeaveModalPath         = "../ui/templates/employeeLeaveModal.html"
	reportsPath                    = "../ui/views/reports.html"
	reportEmployeeLeaveRequestPath = "../ui/templates/reportEmployeeLeaveRequest.html"
	timesheetReportPath            = "../ui/templates/timesheetReport.html"
	employeeTimesheetReportPath    = "../ui/templates/employeeTimesheetReport.html"
	initJobReportPath              = "../ui/templates/initJobReport.html"
	jobReportPath                  = "../ui/templates/jobReport.html"
	jobReportEmployeeTable         = "../ui/templates/jobReportEmployeeTable.html"
	adminPath                      = "../ui/views/admin.html"
	adminEmployeeTabPath           = "../ui/templates/adminEmployeeTab.html"
	adminEmployeeListPath          = "../ui/templates/adminEmployeeList.html"
	adminEmployeeListRowPath       = "../ui/templates/adminEmployeeListRow.html"
	adminAddEmployeeModalPath      = "../ui/templates/adminAddEmployeeModal.html"
	adminPutEmployeeModalPath      = "../ui/templates/adminPutEmployeeModal.html"
	adminJobListPath               = "../ui/templates/adminJobList.html"
	adminJobTabPath                = "../ui/templates/adminJobTab.html"
	addJobModalPath                = "../ui/templates/addJobModal.html"
	putJobModalPath                = "../ui/templates/putJobModal.html"
	adminLeaveTabPath              = "../ui/templates/adminLeaveTab.html"
	adminLeaveRequestPath          = "../ui/templates/adminLeaveRequest.html"
	adminLeaveRequestModalPath     = "../ui/templates/adminLeaveModal.html"
	adminSafetyTabPath             = "../ui/templates/adminSafetyTab.html"
	adminIncidentReportListPath    = "../ui/templates/adminIncidentReportList.html"
	accountPath                    = "../ui/views/account.html"
	safetyPath                     = "../ui/views/safety.html"
	incidentReportFormPath         = "../ui/templates/incidentReportForm.html"
	iframePdfPath                  = "../ui/templates/iframePdf.html"
	putIncidentReportFormPath      = "../ui/templates/putIncidentReportForm.html"
	swmUserContentPath             = "../ui/templates/swmUserContent.html"
	swmListPath                    = "../ui/templates/swmList.html"
	swmsFormPath                   = "../ui/templates/swmsForm.html"
	updateSwmsFormPath             = "../ui/templates/updateSwmsForm.html"

	purchaseOrderPath                = "../ui/views/purchaseOrder.html"
	purchaseOrderFormPath            = "../ui/templates/purchaseOrderForm.html"
	purchaseOrderItemRowPath         = "../ui/templates/purchaseOrderItemRow.html"
	viewPurchaseOrderModalPath       = "../ui/templates/viewPurchaseOrderModal.html"
	employeePurchaseOrderHistoryPath = "../ui/templates/employeePurchaseOrderHistory.html"
	adminPurchaseOrderViewPath       = "../ui/templates/adminPurchaseOrderView.html"
	adminPurchaseOrderHistoryPath    = "../ui/templates/adminPurchaseOrderHistory.html"
	adminItemTypesPath               = "../ui/templates/adminItemTypes.html"
	adminItemSizesPath               = "../ui/templates/adminItemSizes.html"
	adminStoresPath                  = "../ui/templates/adminStores.html"
	adminAddItemModalPath            = "../ui/templates/addItemModal.html"
	adminAddSizeModalPath            = "../ui/templates/addSizeModal.html"
	adminAddStoreModalPath           = "../ui/templates/addStoreModal.html"

	jobnoteTilesPath      = "../ui/templates/jobnoteTiles.html"
	jobNotesPath          = "../ui/templates/jobNotes.html"
	archivedJobNotesPath  = "../ui/templates/archivedJobNotes.html"
	archivedPaintNotePath = "../ui/templates/archivedPaintNote.html"
	archivedTaskNotePath  = "../ui/templates/archivedTaskNote.html"
	archivedImageNotePath = "../ui/templates/archivedImageNote.html"
	paintNotePath         = "../ui/templates/paintNote.html"
	taskNotePath          = "../ui/templates/taskNote.html"
	imageNotePath         = "../ui/templates/imageNote.html"
	paintNoteFormPath     = "../ui/templates/paintNoteForm.html"
	taskNoteFormPath      = "../ui/templates/taskNoteForm.html"
	imageNoteFormPath     = "../ui/templates/imageNoteForm.html"
)

func renderTemplate(
	w http.ResponseWriter,
	r *http.Request,
	component, title string,
	componentData interface{},
) {
	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		http.Error(w, "unable to retrieve is_admin", http.StatusUnauthorized)
		return
	}

	data := struct {
		Title     string
		Component string
		IsAdmin   bool
		Data      interface{}
	}{
		Title:     title,
		Component: component,
		IsAdmin:   isAdmin,
		Data:      componentData,
	}

	tmpl, err := template.ParseFiles(
		layoutPath,
		navPath,
		dashboardPath,
		dashboardTilePath,
		jobsPath,
		timesheetsPath,
		timesheetTablePath,
		timesheetNavContainerPath,
		timesheetHeadPath,
		existingTimesheetRowPath,
		jobSelectModalPath,
		leavePath,
		leaveHistoryPath,
		leaveFormPath,
		reportsPath,
		timesheetReportPath,
		employeeTimesheetReportPath,
		adminPath,
		adminEmployeeTabPath,
		adminEmployeeListPath,
		accountPath,
		safetyPath,
		incidentReportFormPath,
		purchaseOrderPath,
		purchaseOrderFormPath,
		purchaseOrderItemRowPath,
		jobnoteTilesPath,
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = tmpl.ExecuteTemplate(w, "layout", data)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func parseRequestValues(keys []string, r *http.Request) ([]string, error) {
	out := make([]string, len(keys))

	err := r.ParseForm()
	if err != nil {
		return out, err
	}

	for i := range keys {
		val := r.FormValue(keys[i])
		out[i] = val
	}

	return out, nil
}

func decodeJSON(payload interface{}, r *http.Request) error {
	err := json.NewDecoder(r.Body).Decode(payload)
	return err
}

func getEmployeeId(w http.ResponseWriter, r *http.Request) (string, error) {
	employeeId, ok := r.Context().Value("employee_id").(string)
	if !ok {
		return "", errors.New("Error retrieving employee_id")
	}
	return employeeId, nil
}

func getIsAdmin(w http.ResponseWriter, r *http.Request) (bool, error) {
	isAdmin, ok := r.Context().Value("is_admin").(bool)
	if !ok {
		return false, errors.New("Error retrieving is_admin")
	}
	return isAdmin, nil
}

func responseJSON(w http.ResponseWriter, data any, sugar *zap.SugaredLogger) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(data); err != nil {
		sugar.Errorf("Error encoding leave requests: %v", err)
		http.Error(w, "failed to fetch leave requests", http.StatusInternalServerError)
		return
	}
}

func executePartialTemplate(filepath string, name string, data interface{}, w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		return err
	}
	err = tmpl.ExecuteTemplate(w, name, data)
	if err != nil {
		return err
	}
	return nil
}

type HtmlServer interface {
	ServeSingleTemplate(templatePath, templateName string, data interface{}, w http.ResponseWriter) error
}

func (h *Handler) ServeSingleTemplate(templatePath, templateName string, data interface{}, w http.ResponseWriter) error {
	err := executePartialTemplate(templatePath, templateName, data, w)
	if err != nil {
		log.Printf("error executing %s: %v", templatePath, err)
		http.Error(w, "error exectuting %s, internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handler) ServeMultiTemplate(templatePaths []string, parentTemplateName string, data interface{}, w http.ResponseWriter) error {
	tmpl, err := template.ParseFiles(templatePaths...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		http.Error(w, "error parsing template files, internal server error", http.StatusInternalServerError)
		return err
	}

	err = tmpl.ExecuteTemplate(w, parentTemplateName, data)
	if err != nil {
		log.Printf("error executing purchaseOrderForm.html: %v", err)
		http.Error(w, "error executing purchaseOrderForm.html, internal server error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (h *Handler) DecodeJson(out any, w http.ResponseWriter, r *http.Request) bool {
	if err := json.NewDecoder(r.Body).Decode(out); err != nil {
		log.Printf("error decoding json: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return false
	}
	return true
}

func (h *Handler) ServeLoginView(w http.ResponseWriter, r *http.Request) {
	login_path := filepath.Join("..", "ui", "views", "login.html")
	tmpl := template.Must(template.ParseFiles(login_path))
	tmpl.Execute(w, nil)
}

func (h *Handler) ServeErrorView(w http.ResponseWriter, r *http.Request) {
	errorPath := filepath.Join("..", "ui", "views", "error.html")
	tmpl := template.Must(template.ParseFiles(errorPath))
	tmpl.Execute(w, nil)
}

func (h *Handler) ServeAccountView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			component := "account"
			title := "Account - GPG"
			renderTemplate(w, r, component, title, nil)
		},
	)
}

func (h *Handler) ServeDashboardView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			data, err := getDashboardData(w, r)
			if err != nil {
				log.Printf("Error getting dashboard data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			component := "dashboard"
			title := "Dashboard - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

func (h *Handler) ServeSafetyView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, err := getEmployeeId(w, r)
			if err != nil {
				log.Printf("error getting employee id: %v", err)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			employee, err := infrastructure.GetEmployeeByEmployeeId(employeeId, h.DB)
			if err != nil {
				log.Printf("error getting employee: %v", err)
				http.Error(w, "Not found", http.StatusNotFound)
				return
			}

			data := safety.IncidentReportValues{
				ReporterId: employeeId,
				Reporter:   fmt.Sprintf("%s %s", employee.FirstName, employee.LastName),
			}

			component := "safety"
			title := "Safety - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}

type LeaveViewData struct {
	LeaveRequests []domain.LeaveRequest
	LeaveFormDto  application.LeaveFormDto
}

func (h *Handler) ServeLeaveView() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			employeeId, ok := r.Context().Value("employee_id").(string)
			if !ok {
				h.Logger.Warn("Error getting employee_id from context")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			employee, err := infrastructure.GetEmployeeByEmployeeId(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting employee data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			leaveFormDto := application.LeaveFormDto{
				EmployeeId: employee.EmployeeId,
				FirstName:  employee.FirstName,
				LastName:   employee.LastName,
				DateErr:    "",
			}

			leaveRequests, err := infrastructure.GetLeaveRequestsByEmployee(employeeId, h.DB)
			if err != nil {
				h.Logger.Errorf("Error getting leave page data: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			data := LeaveViewData{
				LeaveFormDto:  leaveFormDto,
				LeaveRequests: leaveRequests,
			}

			component := "leave"
			title := "Leave - GPG"
			renderTemplate(w, r, component, title, data)
		},
	)
}
