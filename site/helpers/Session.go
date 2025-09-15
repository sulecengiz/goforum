package helpers

import (
	"github.com/gorilla/sessions"
)

var SessionStore = sessions.NewCookieStore([]byte("goforum-secret-key"))
