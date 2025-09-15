package controllers

import (
	"encoding/json"
	"goforum/site/models"
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

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	name := strings.TrimSpace(r.FormValue("name"))
	email := strings.TrimSpace(r.FormValue("email"))
	phone := strings.TrimSpace(r.FormValue("phone"))
	message := strings.TrimSpace(r.FormValue("message"))

	if name == "" || email == "" || message == "" || phone == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Tüm zorunlu alanları doldurun."})
		return
	}

	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if !emailRegex.MatchString(email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Geçerli bir e-posta adresi girin."})
		return
	}

	// Telefonu sadece rakamlara indir
	digitsOnly := regexp.MustCompile("\\D").ReplaceAllString(phone, "")
	if len(digitsOnly) < 10 || len(digitsOnly) > 15 { // uluslararası esneklik
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Telefon numarası formatı geçerli değil."})
		return
	}

	ContactForm := &models.ContactForm{
		Name:    name,
		Email:   email,
		Phone:   digitsOnly,
		Message: message,
	}

	err := ContactForm.Add()
	if err != nil {
		log.Printf("Veritabanına kaydetme hatası: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]interface{}{"success": false, "error": "Veri kaydedilirken hata oluştu."})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{"success": true, "message": "Mesaj kaydedildi."})
}
