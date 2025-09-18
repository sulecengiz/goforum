package controllers

import (
	"goforum/site/models"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// ContactFormHandler - İletişim formu işleyicisi
func ContactFormHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Form verilerini al
	name := r.FormValue("name")
	email := r.FormValue("email")
	subject := r.FormValue("subject")
	message := r.FormValue("message")

	// Basit validasyon
	if name == "" || email == "" || subject == "" || message == "" {
		http.Error(w, "Tüm alanlar doldurulmalıdır", http.StatusBadRequest)
		return
	}

	// Contact modelini oluştur ve kaydet
	contact := &models.Contact{
		Name:    name,
		Email:   email,
		Subject: subject,
		Message: message,
	}

	if err := contact.Save(); err != nil {
		log.Printf("İletişim formu kaydetme hatası: %v", err)
		http.Error(w, "Mesajınız gönderilemedi. Lütfen tekrar deneyin.", http.StatusInternalServerError)
		return
	}

	// Başarılı gönderim - iletişim sayfasına yönlendir
	http.Redirect(w, r, "/contact?success=1", http.StatusSeeOther)
}
