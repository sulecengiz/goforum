package helpers

import (
	"encoding/gob"
	"goforum/admin/models"
	"net/http"

	"github.com/gorilla/sessions"
)

var AdminSessionStore = sessions.NewCookieStore([]byte("admin-secret-key-123"))

func init() {
	// Role tipini gob'a kaydet
	gob.Register(models.Role(0))

	AdminSessionStore.Options = &sessions.Options{
		Path:     "/admin", // sadece admin URL’lerinde geçerli
		MaxAge:   3600 * 8, // 8 saat
		HttpOnly: true,
	}
}

// Admin session işlemleri
func GetAdminSession(r *http.Request) (*sessions.Session, error) {
	return AdminSessionStore.Get(r, "admin_session")
}

func IsAdminLoggedIn(r *http.Request) bool {
	session, err := GetAdminSession(r)
	if err != nil {
		return false
	}

	if session.Values["adminID"] != nil && session.Values["adminRole"] == models.RoleAdmin {
		return true
	}
	return false
}

func SaveAdminSession(w http.ResponseWriter, r *http.Request, userID int, username string) error {
	session, _ := GetAdminSession(r)
	session.Values["adminID"] = userID
	session.Values["adminUsername"] = username
	session.Values["adminRole"] = models.RoleAdmin
	return session.Save(r, w)
}

func ClearAdminSession(w http.ResponseWriter, r *http.Request) error {
	session, _ := GetAdminSession(r)
	session.Values = make(map[interface{}]interface{})
	return session.Save(r, w)
}
