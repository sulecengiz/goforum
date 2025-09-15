package helpers

import (
	"github.com/gorilla/sessions"
	"net/http"
)

var SessionStore = sessions.NewCookieStore([]byte("site-secret-key-123"))

func SetUser(w http.ResponseWriter, r *http.Request, user interface{}) {
	session, _ := SessionStore.Get(r, "session")
	session.Values["user"] = user
	session.Save(r, w)
}

func GetUser(r *http.Request) interface{} {
	session, _ := SessionStore.Get(r, "session")
	return session.Values["user"]
}

func RemoveUser(w http.ResponseWriter, r *http.Request) {
	session, _ := SessionStore.Get(r, "session")
	delete(session.Values, "user")
	session.Save(r, w)
}
