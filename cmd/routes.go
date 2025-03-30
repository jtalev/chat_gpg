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

	// safety routes
	mux.Handle("/safety/generate-incident-report", a.AuthorizeUser(h.GenerateIncidentReport()))
	mux.Handle("/safety/get-incident-report", a.AuthorizeAdmin(h.GetIncidentReport()))
	mux.Handle("/safety/delete-incident-report", a.AuthorizeAdmin(h.DeleteIncidentReport()))
}
