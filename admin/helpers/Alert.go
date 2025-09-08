// helpers/Alert.go

package helpers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte("123123"))

func SetAlert(w http.ResponseWriter, r *http.Request, message string) error {
	session, err := store.Get(r, "go-alert")
	if err != nil {
		fmt.Println(err)
		return err
	}
	session.AddFlash(message)
	return session.Save(r, w)
}

func GetAlert(w http.ResponseWriter, r *http.Request) interface{} {
	session, err := store.Get(r, "go-alert")
	if err != nil {
		fmt.Println(err)
		return nil
	}
	flashes := session.Flashes()

	// Mesajı aldıktan sonra hemen çerezi kaydederek temizle
	session.Save(r, w)

	if len(flashes) > 0 {
		return flashes[0]
	}
	return nil
}
