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

	mux.Handle("/dashboard", a.AuthMiddleware(h.ServeDashboardView()))
	mux.Handle("/jobs", a.AuthMiddleware(h.ServeJobsView()))
	mux.Handle("/timesheets", a.AuthMiddleware(h.ServeTimesheetsView()))
	mux.Handle("/leave", a.AuthMiddleware(h.ServeLeaveView()))
	mux.Handle("/reports", a.AuthMiddleware(h.ServeReportsView()))
	mux.Handle("/admin", a.AuthMiddleware(h.ServeAdminView()))
	mux.Handle("/account", a.AuthMiddleware(h.ServeAccountView()))

	// leave requests
	mux.Handle("/leave/get", a.AuthMiddleware(h.GetLeaveRequests()))
	mux.Handle("/leave/get-by-employee", a.AuthMiddleware(h.GetLeaveRequestsByEmployee()))
	mux.Handle("/leave/post", a.AuthMiddleware(h.PostLeaveRequest()))
	mux.Handle("/leave/put", a.AuthMiddleware(h.PutLeaveRequest()))
	mux.Handle("/leave/delete", a.AuthMiddleware(h.DeleteLeaveRequest()))
	mux.Handle("/leave/form", a.AuthMiddleware(h.RenderLeaveFormTab()))
	mux.Handle("/leave/history", a.AuthMiddleware(h.RenderLeaveHistoryTab()))

	// job requests
	mux.Handle("/jobs/get", a.AuthMiddleware(h.GetJobs()))
	mux.Handle("/jobs/get-by-id", a.AuthMiddleware(h.GetJobById()))
	mux.Handle("/jobs/get-by-name", a.AuthMiddleware(h.GetJobByName()))
	mux.Handle("/jobs/post", a.AuthMiddleware(h.PostJob()))
	mux.Handle("/jobs/put", a.AuthMiddleware(h.PutJob()))
	mux.Handle("/jobs/delete", a.AuthMiddleware(h.DeleteJob()))

	// timesheet requests
	mux.Handle("/timesheets/get", a.AuthMiddleware(h.GetTimesheets()))
	mux.Handle("/timesheets/get-by-id", a.AuthMiddleware(h.GetTimesheetById()))
	mux.Handle("/timesheets/put", a.AuthMiddleware(h.PutTimesheet()))
	mux.Handle("/timesheets/get-timesheet-table", a.AuthMiddleware(h.GetTimesheetTable()))
	mux.Handle("/timesheets/job-select-modal", a.AuthMiddleware(h.RenderJobSelectModal()))

	// timesheet week requests
	mux.Handle("/timesheet-week/init-timesheet-week", a.AuthMiddleware(h.InitTimesheetWeek()))
	mux.Handle("/timesheet-week/get-by-employee", a.AuthMiddleware(h.GetTimesheetWeekByEmployee()))
	mux.Handle("/timesheet-week/delete", a.AuthMiddleware(h.DeleteTimesheetWeek()))

	// reports requests
	mux.Handle("/reports/timesheet-report", a.AuthMiddleware(h.GetEmployeeTimesheetReport()))

	// admin requests
	mux.Handle("/admin/render-job-tab", a.AuthMiddleware(h.RenderJobTab()))
}
