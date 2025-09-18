package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	"goforum/site/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Contact struct'ı, yönetici panelindeki iletişim sayfası için kontrolcüyü temsil eder.
type Contact struct{}

func (c Contact) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("contact/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Contacts"] = models.Contact{}.GetAll()
	data["Alert"] = map[string]interface{}{
		"is_alert": helpers.GetAlert(w, r) != "",
		"message":  helpers.GetAlert(w, r),
	}

	view.ExecuteTemplate(w, "index", data)
}

func (c Contact) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	id := params.ByName("id")

	// ID'yi uint'e çevir
	contactID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		helpers.SetAlert(w, r, "Geçersiz ID formatı")
		http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
		return
	}

	// Contact modelini kullan
	contact := models.Contact{}.Get(uint(contactID))
	if contact.ID == 0 {
		helpers.SetAlert(w, r, "Kayıt bulunamadı")
		http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
		return
	}

	// Kayıt başarıyla bulunduysa silme işlemini yapın.
	err = models.Contact{}.Delete(uint(contactID))
	if err != nil {
		helpers.SetAlert(w, r, "Kayıt silinirken bir hata oluştu: "+err.Error())
		http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
		return
	}

	helpers.SetAlert(w, r, "Kayıt başarıyla silindi")
	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}
