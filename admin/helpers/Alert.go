// helpers/Alert.go

package helpers

import (
	"net/http"
)

func SetAlert(w http.ResponseWriter, r *http.Request, message string) error {
	session, err := GetAdminSession(r)
	if err != nil {
		return err
	}
	session.Values["alert"] = message
	return session.Save(r, w)
}

func GetAlert(w http.ResponseWriter, r *http.Request) string {
	session, err := GetAdminSession(r)
	if err != nil {
		return ""
	}

	msg, ok := session.Values["alert"].(string)
	if !ok {
		return ""
	}

	delete(session.Values, "alert")
	session.Save(r, w)
	return msg
}
