package controllers

import (
	"fmt"
	"goblog/admin/helpers"
	"goblog/site/models"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Contact struct'ı, yönetici panelindeki iletişim sayfası için kontrolcüyü temsil eder.
type Contact struct{}

func (c Contact) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}
	view, err := template.ParseFiles(helpers.Include("contact/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Contacts"] = models.ContactForm{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}

func (c Contact) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}

	id := params.ByName("id")

	// Önce kaydı bulmaya çalışın.
	contact, err := models.ContactForm{}.Get(id)
	if err != nil {
		// Kayıt bulunamadıysa veya başka bir hata oluştuysa
		helpers.SetAlert(w, r, "Kayıt bulunamadı veya bir hata oluştu: "+err.Error())
		http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
		return
	}

	// Kayıt başarıyla bulunduysa silme işlemini yapın.
	err = contact.Delete()
	if err != nil {
		helpers.SetAlert(w, r, "Kayıt silinirken bir hata oluştu: "+err.Error())
		http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
		return
	}

	helpers.SetAlert(w, r, "Kayıt başarıyla silindi.")
	http.Redirect(w, r, "/admin/contact", http.StatusSeeOther)
}
