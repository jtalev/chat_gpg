package domain

type Role struct {
	UUID       string `json:"uuid"`
	EmployeeId string `json:"employee_id"`
	Role       string `json:"role"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}
