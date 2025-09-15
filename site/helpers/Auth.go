package helpers

import (
	"errors"
	"goforum/site/models" // Make sure this path is correct
	"net/http"
)

// Assume SessionStore is defined elsewhere in helpers
// var SessionStore = sessions.NewCookieStore([]byte("your-secret-key"))

func GetCurrentUser(r *http.Request) (*models.User, error) {
	session, err := SessionStore.Get(r, "session")
	if err != nil {
		return nil, errors.New("session error")
	}

	userID, ok := session.Values["userID"]
	if !ok {
		return nil, errors.New("kullanıcı girişi yapılmamış")
	}

	var id uint
	switch v := userID.(type) {
	case uint:
		id = v
	case int:
		id = uint(v)
	case int64:
		id = uint(v)
	case float64: // JSON decode durumunda
		id = uint(v)
	default:
		return nil, errors.New("geçersiz kullanıcı ID'si")
	}

	var userModel models.User
	user := userModel.Get(id) // This fetches the full User struct from DB

	if user.ID == 0 {
		return nil, errors.New("kullanıcı bulunamadı")
	}
	return &user, nil
}
