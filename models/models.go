package models

type ValidationResult struct {
	Key     string
	IsValid bool
	Msg     string
}

type Employee struct {
	ID          int    `json:"id"`
	EmployeeId  string `json:"employee_id"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	IsAdmin     bool   `json:"is_admin"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type EmployeeAuth struct {
	AuthId       int    `json:"auth_id"`
	EmployeeId   string `json:"employee_id"`
	Username     string `json:"username"`
	PasswordHash string `json:"password_hash"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

type LeaveRequest struct {
	RequestId  int    `json:"request_id"`
	EmployeeId string `json:"employee_id"`
	FirstName  string `json:"first_name"`
	LastName   string `json:"last_name"`
	Type       string `json:"leave_type"`
	From       string `json:"from_date"`
	To         string `json:"to_date"`
	Note       string `json:"note"`
	IsApproved bool   `json:"is_approved"`
}

type Job struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Number    int    `json:"number"`
	Address   string `json:"address"`
	Suburb    string `json:"suburb"`
	PostCode  string `json:"post_code"`
	City      string `json:"city"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Timesheet struct {
	ID         int    `json:"id"`
	EmployeeId string `json:"employee_id"`
	JobId      int    `json:"job_id"`
	WeekStart  string `json:"week_start"`
	Date       string `json:"date"`
	Hours      int    `json:"hours"`
	Minutes    int    `json:"minutes"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
