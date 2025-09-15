// middleware/Auth.go
package middleware

import (
	"goforum/admin/helpers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func CheckAdmin(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		if !helpers.IsAdminLoggedIn(r) {
			http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
			return
		}
		next(w, r, ps)
	}
}
