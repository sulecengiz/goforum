// admin/controllers/About.go
package controllers

import (
	"fmt"
	"goblog/admin/helpers"
	"goblog/site/models"
	"html/template"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func AboutIndex(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}
	about, err := models.About{}.Get()
	if err != nil {
		log.Println("Hakkında sayfası bilgisi alınırken hata oluştu:", err)
		http.Error(w, "Sunucu Hatası", http.StatusInternalServerError)
		return
	}

	view, err := template.ParseFiles(helpers.Include("about/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["About"] = about
	data["Alert"] = helpers.GetAlert(w, r)

	view.ExecuteTemplate(w, "index", data)
}
func AboutUpdate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	content := r.FormValue("content")

	about, err := models.About{}.Get()
	if err != nil {
		log.Println("Hakkımda kaydı bulunamadı:", err)
		http.Redirect(w, r, "/admin/about", http.StatusFound)
		helpers.SetAlert(w, r, "Hakkımda kaydı bulunamadı")
		return
	}

	about.Content = content
	about.Update()
	helpers.SetAlert(w, r, "Hakkımda kaydedildi")
	http.Redirect(w, r, "/admin/about", http.StatusFound)
}
