package handlers

import (
	"fmt"
	"net/http"
)

type LeaveRequest struct {
	Name 		string  `json:"name"`
	Type 		string  `json:"type"`
	From 		string  `json:"from"`
	To   		string  `json:"to"`
	Note 		string  `json:"note"`
	IsApproved 	bool 	`json:"is_approved"`
}

func getLeaveRequests() []LeaveRequest {
	data := []LeaveRequest {
		{
			Name: 	"Slid",
			Type: 	"annual",
			From: 	"07/12/2024",
			To: 	"10/12/2024",
			Note: 	"Need a break Need a break Need a break Need a break Need a break",
			IsApproved: true,
		},
		{
			Name: 	"Slid",
			Type: 	"sick",
			From: 	"04/12/2024",
			To: 	"05/12/2024",
			Note: 	"Sick",
			IsApproved: true,
		},
		{
			Name: 	"Slid",
			Type: 	"stress",
			From: 	"10/11/2024",
			To: 	"15/11/2024",
			Note: 	"I'm stressed out",
			IsApproved: false,
		},
		{
			Name: 	"Slid",
			Type: 	"stress",
			From: 	"10/11/2024",
			To: 	"15/11/2024",
			Note: 	"I'm stressed out",
			IsApproved: false,
		},
		{
			Name: 	"Slid",
			Type: 	"stress",
			From: 	"10/11/2024",
			To: 	"15/11/2024",
			Note: 	"I'm stressed out",
			IsApproved: false,
		},
		{
			Name: 	"Slid",
			Type: 	"stress",
			From: 	"10/11/2024",
			To: 	"15/11/2024",
			Note: 	"I'm stressed out",
			IsApproved: false,
		},
		{
			Name: 	"Slid",
			Type: 	"stress",
			From: 	"10/11/2024",
			To: 	"15/11/2024",
			Note: 	"I'm stressed out",
			IsApproved: false,
		},
	}
	return data
}

func ServeLeaveView(w http.ResponseWriter, r *http.Request) {
	data := getLeaveRequests()
	component := "leave"
	title := "Leave - GPG"
	renderTemplate(w, component, title, data)
}

func Submit_leave_request(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "SUBMITTED")
}
