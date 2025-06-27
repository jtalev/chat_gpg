package main

import (
	"context"
	"net/http"

	"github.com/jtalev/chat_gpg/handlers"
	"github.com/jtalev/chat_gpg/infrastructure/auth"
)

func add_routes(mux *http.ServeMux, ctx context.Context, h *handlers.Handler, a *infrastructure.Auth) {
	fileServer := http.FileServer(http.Dir("../ui/static"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	// login/logout requests
	mux.HandleFunc("/login", h.ServeLoginView)
	mux.Handle("/logout", h.LogoutHandler())
	mux.Handle("/authenticate-user", h.LoginHandler())

	mux.HandleFunc("/error", h.ServeErrorView)

	mux.Handle("/dashboard", a.AuthorizeUser(h.ServeDashboardView()))
	mux.Handle("/jobs", a.AuthorizeUser(h.ServeJobsView()))
	mux.Handle("/timesheets", a.AuthorizeUser(h.ServeTimesheetsView()))
	mux.Handle("/leave", a.AuthorizeUser(h.ServeLeaveView()))
	mux.Handle("/reports", a.AuthorizeAdmin(h.ServeReportsView()))
	mux.Handle("/admin", a.AuthorizeAdmin(h.ServeAdminView()))
	mux.Handle("/account", a.AuthorizeUser(h.ServeAccountView()))
	mux.Handle("/safety", a.AuthorizeUser(h.ServeSafetyView()))
	mux.Handle("/purchase-order", a.AuthorizeUser(h.ServePurchaseOrderView()))

	// leave requests
	mux.Handle("/leave/get", a.AuthorizeUser(h.GetLeaveRequests()))
	mux.Handle("/leave/get-by-employee", a.AuthorizeUser(h.GetLeaveRequestsByEmployee()))
	mux.Handle("/leave/post", a.AuthorizeUser(h.PostLeaveRequest()))
	mux.Handle("/leave/put", a.AuthorizeUser(h.PutLeaveRequest()))
	mux.Handle("/leave/delete", a.AuthorizeUser(h.DeleteLeaveRequest()))
	mux.Handle("/leave/form", a.AuthorizeUser(h.RenderLeaveFormTab()))
	mux.Handle("/leave/history", a.AuthorizeUser(h.RenderLeaveHistoryTab()))
	mux.Handle("/leave/leave-request-modal", a.AuthorizeUser(h.EmployeeLeaveModal()))

	// job requests
	mux.Handle("/jobs/get", a.AuthorizeUser(h.GetJobs()))
	mux.Handle("/jobs/get-by-id", a.AuthorizeUser(h.GetJobById()))
	mux.Handle("/jobs/get-by-name", a.AuthorizeUser(h.GetJobByName()))
	mux.Handle("/jobs/post", a.AuthorizeAdmin(h.PostJob()))
	mux.Handle("/jobs/put", a.AuthorizeAdmin(h.PutJob()))
	mux.Handle("/jobs/delete", a.AuthorizeAdmin(h.DeleteJob()))

	// employee requests
	mux.Handle("/employee/delete", a.AuthorizeAdmin(h.DeleteEmployee()))
	mux.Handle("/employee/post", a.AuthorizeAdmin(h.PostEmployee()))
	mux.Handle("/employee/put", a.AuthorizeAdmin(h.PutEmployee()))

	// timesheet requests
	mux.Handle("/timesheets/get", a.AuthorizeUser(h.GetTimesheets()))
	mux.Handle("/timesheets/get-by-id", a.AuthorizeUser(h.GetTimesheetById()))
	mux.Handle("/timesheets/put", a.AuthorizeUser(h.PutTimesheet()))
	mux.Handle("/timesheets/get-timesheet-table", a.AuthorizeUser(h.GetTimesheetTable()))
	mux.Handle("/timesheets/job-select-modal", a.AuthorizeUser(h.RenderJobSelectModal()))

	// timesheet week requests
	mux.Handle("/timesheet-week/init-timesheet-week", a.AuthorizeUser(h.InitTimesheetWeek()))
	mux.Handle("/timesheet-week/get-by-employee", a.AuthorizeUser(h.GetTimesheetWeekByEmployee()))
	mux.Handle("/timesheet-week/delete", a.AuthorizeUser(h.DeleteTimesheetWeek()))

	// reports requests
	mux.Handle("/reports/timesheet-report", a.AuthorizeAdmin(h.GetEmployeeTimesheetReport()))
	mux.Handle("/reports/prev-timesheet-report", a.AuthorizeAdmin(h.PrevEmployeeTimesheetReport()))
	mux.Handle("/reports/next-timesheet-report", a.AuthorizeAdmin(h.NextEmployeeTimesheetReport()))
	mux.Handle("/reports/init-timesheet-report", a.AuthorizeAdmin(h.InitTimesheetReport()))
	mux.Handle("/reports/init-job-report", a.AuthorizeAdmin(h.InitJobReportData()))
	mux.Handle("/reports/get-job-report", a.AuthorizeAdmin(h.GetJobReport()))

	// admin requests
	mux.Handle("/admin/render-employee-tab", a.AuthorizeAdmin(h.RenderEmployeeTab()))
	mux.Handle("/admin/render-job-tab", a.AuthorizeAdmin(h.RenderJobTab()))
	mux.Handle("/admin/render-leave-tab", a.AuthorizeAdmin(h.RenderLeaveTab()))
	mux.Handle("/admin/add-job-modal", a.AuthorizeAdmin(h.AddJobModal()))
	mux.Handle("/admin/put-modal", a.AuthorizeAdmin(h.PutJobModal()))
	mux.Handle("/admin/leave-request-modal", a.AuthorizeAdmin(h.LeaveRequestModal()))
	mux.Handle("/admin/leave-finalise", a.AuthorizeAdmin(h.LeaveRequestFinalise()))
	mux.Handle("/admin/add-employee-modal", a.AuthorizeAdmin(h.AddEmployeeModal()))
	mux.Handle("/admin/employee-put-modal", a.AuthorizeAdmin(h.PutEmployeeModal()))
	mux.Handle("/admin/render-safety-tab", a.AuthorizeAdmin(h.RenderSafetyTab()))

	mux.Handle("/admin/render-purchase-order-tab", a.AuthorizeAdmin(h.RenderPurchaseOrderTab()))
	mux.Handle("/admin/serve-purchase-order-history", a.AuthorizeAdmin(h.ServePurchaseOrderHistory()))

	mux.Handle("/admin/serve-stores-template", a.AuthorizeAdmin(h.ServeStores()))
	mux.Handle("/admin/serve-add-store-modal", a.AuthorizeAdmin(h.ServeAddStoreModal()))
	mux.Handle("/admin/serve-put-store-modal", a.AuthorizeAdmin(h.ServePutStoreModal()))

	mux.Handle("/admin/serve-item-sizes-template", a.AuthorizeAdmin(h.ServeItemSizes()))
	mux.Handle("/admin/serve-add-size-modal", a.AuthorizeAdmin(h.ServeAddSizeModal()))
	mux.Handle("/admin/serve-put-size-modal", a.AuthorizeAdmin(h.ServePutSizeModal()))

	mux.Handle("/admin/serve-item-types-template", a.AuthorizeAdmin(h.ServeItemTypes()))
	mux.Handle("/admin/serve-add-item-modal", a.AuthorizeAdmin(h.ServeAddItemModal()))
	mux.Handle("/admin/serve-put-item-modal", a.AuthorizeAdmin(h.ServePutItemModal()))

	mux.Handle("/admin/store/post", a.AuthorizeAdmin(h.PostStore()))
	mux.Handle("/admin/store/put", a.AuthorizeAdmin(h.PutStore()))
	mux.Handle("/admin/store/delete", a.AuthorizeAdmin(h.DeleteStore()))

	mux.Handle("/admin/item-type/post", a.AuthorizeAdmin(h.PostItemType()))
	mux.Handle("/admin/item-type/put", a.AuthorizeAdmin(h.PutItemType()))
	mux.Handle("/admin/item-type/delete", a.AuthorizeAdmin(h.DeleteItemType()))

	mux.Handle("/admin/item-size/post", a.AuthorizeAdmin(h.PostItemSize()))
	mux.Handle("/admin/item-size/put", a.AuthorizeAdmin(h.PutItemSize()))
	mux.Handle("/admin/item-size/delete", a.AuthorizeAdmin(h.DeleteItemSize()))

	mux.Handle("/admin/safety/serve-incident-report-content", a.AuthorizeAdmin(h.AdminServeIncidentReportContent()))

	// safety routes
	mux.Handle("/safety/serve-incident-report-form", a.AuthorizeUser(h.ServeIncidentReportForm()))
	mux.Handle("/safety/generate-incident-report", a.AuthorizeUser(h.GenerateIncidentReport()))
	mux.Handle("/safety/get-incident-report", a.AuthorizeAdmin(h.GetIncidentReport()))
	mux.Handle("/safety/delete-incident-report", a.AuthorizeAdmin(h.DeleteIncidentReport()))
	mux.Handle("/safety/put-incident-report-html", a.AuthorizeAdmin(h.PutIncidentReportHtml()))
	mux.Handle("/safety/put-incident-report", a.AuthorizeAdmin(h.PutIncidentReport()))

	mux.Handle("/safety/serve-swm-user-content", a.AuthorizeUser(h.ServeSwmUserContent()))
	mux.Handle("/safety/swms/get-list-html", a.AuthorizeUser(h.GetSwmsListHtml()))
	mux.Handle("/safety/swms/generate-swms-pdf", a.AuthorizeForemanManagement(h.GenerateSwmsPdf()))
	mux.Handle("/safety/swms/serve-swms-pdf", a.AuthorizeUser(h.ServeSwmsPdf()))
	mux.Handle("/safety/swms/serve-swms-form", a.AuthorizeForemanManagement(h.ServeSwmsForm()))
	mux.Handle("/safety/swms/post", a.AuthorizeForemanManagement(h.PostSwms()))
	mux.Handle("/safety/swms/delete", a.AuthorizeForemanManagement(h.DeleteSwms()))
	mux.Handle("/safety/swms/serve-swms-form-put", a.AuthorizeForemanManagement(h.ServeSwmsFormPut()))
	mux.Handle("/safety/swms/put", a.AuthorizeForemanManagement(h.PutSwms()))

	// purchase order
	mux.Handle("/purchase-order/serve-item-row", a.AuthorizeUser(h.ServeItemRow()))
	mux.Handle("/purchase-order/serve-purchase-order-employee-history", a.AuthorizeUser(h.ServeEmployeePOHistory()))
	mux.Handle("/purchase-order/serve-form", a.AuthorizeUser(h.ServePurchaseOrderForm()))
	mux.Handle("/purchase-order/serve-purchase-order", a.AuthorizeUser(h.ServePurchaseOrder()))

	mux.Handle("/purchase-order/post", a.AuthorizeUser(h.PostPurchaseOrder()))
	mux.Handle("/purchase-order/delete", a.AuthorizeUser(h.DeletePurchaseOrder()))

	// job notes
	mux.Handle("/job-notes/serve-jobnote-tiles", a.AuthorizeUser(h.ServeJobnoteTiles()))
	mux.Handle("/job-notes/serve-note-form", a.AuthorizeUser(h.ServeNoteForm()))

	mux.Handle("/job-notes/archive", a.AuthorizeUser(h.ArchiveJobNote()))
	mux.Handle("/job-notes/get", a.AuthorizeUser(h.GetJobNotes()))
	mux.Handle("/job-notes/post", a.AuthorizeUser(h.PostNote()))
	mux.Handle("/job-notes/put", a.AuthorizeUser(h.PutNote()))
	mux.Handle("/job-notes/delete", a.AuthorizeUser(h.DeleteNote()))
}
