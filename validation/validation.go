package validation

import (
	"fmt"
	"net/http"
)

type validation_result struct {
	Is_valid bool
	Msg      string
}

func Validate_login(username, password string) validation_result {
	result := validation_result{
		Is_valid: false,
		Msg:      "come back to this when more is working",
	}
	return result
}

func Handle_login_validation() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				fmt.Printf("handle_login_validation error parsing form: %v", err)
				http.Error(w, "error parsing form", http.StatusBadRequest)
			}
			username := r.FormValue("username")
			password := r.FormValue("password")
			result := Validate_login(username, password)
			var html string
			if !result.Is_valid {
				html = fmt.Sprintf(`<p class="err" id="err">*%s</p>`, result.Msg)
				if _, err := w.Write([]byte(html)); err != nil {
					fmt.Printf("handle_login_validation error writing response: %v", err)
					http.Error(w, "error writing response", http.StatusInternalServerError)
					return
				}
			}
		},
	)
}
