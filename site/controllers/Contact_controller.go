package controllers

import (
	"goblog/site/models"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/julienschmidt/httprouter"
)

func ContactFormHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Geçersiz istek metodu.", http.StatusMethodNotAllowed)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || message == "" || phone == "" {
		http.Error(w, "Tüm zorunlu alanları doldurun.", http.StatusBadRequest)
		return
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		http.Error(w, "Geçerli bir e-posta adresi girin.", http.StatusBadRequest)
		return
	}

	phoneRegex := regexp.MustCompile("^[0-9]{3}[0-9]{3}[0-9]{2}[0-9]{2}$")
	if !phoneRegex.MatchString(phone) {
		http.Error(w, "Telefon numarası formatı geçerli değil (örn: 555-555-55-55).", http.StatusBadRequest)
		return
	}

	ContactForm := &models.ContactForm{
		Name:    name,
		Email:   email,
		Phone:   phone,
		Message: message,
	}

	err := ContactForm.Add()
	if err != nil {
		log.Printf("Veritabanına kaydetme hatası: %v", err)
		http.Error(w, "Veri kaydedilirken bir hata oluştu.", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Başarılı"))
}
