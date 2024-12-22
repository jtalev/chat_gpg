package handlers

type DashboardData struct {
	IsAdmin bool
}

func getDashboardData() []DashboardData {
	data := []DashboardData{
		{},
	}
	return data
}
