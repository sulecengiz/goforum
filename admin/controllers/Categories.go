package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	"goforum/admin/models"
	"html/template"
	"net/http"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Categories struct{}

func (category Categories) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	view, err := template.ParseFiles(helpers.Include("categories/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Categories"] = models.Category{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}

func (category Categories) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	categoryTitle := r.FormValue("category-title")
	categorySlug := slug.Make(categoryTitle)

	models.Category{Title: categoryTitle, Slug: categorySlug}.Add()
	helpers.SetAlert(w, r, "Kayıt başarıyla eklendi")
	http.Redirect(w, r, "/admin/kategoriler", http.StatusFound)
}

func (category Categories) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	categories := models.Category{}.Get(params.ByName("id"))
	categories.Delete()
	helpers.SetAlert(w, r, "Başarıyla silindi")
	http.Redirect(w, r, "/admin/kategoriler", http.StatusSeeOther)
}
