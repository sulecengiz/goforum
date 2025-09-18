package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	"goforum/admin/models"
	"html/template"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

type Dashboard struct{}

func (dashboard Dashboard) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Dashboard.Index başladı")

	if !helpers.IsAdminLoggedIn(r) {
		fmt.Println("Admin login yok, redirect")
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	fmt.Println("Admin login başarılı")

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string {
			return models.Category{}.Get(categoryID).Title
		},
	}).ParseFiles(helpers.Include("dashboard/list")...)
	if err != nil {
		fmt.Println("Template parse error:", err)
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	alert := helpers.GetAlert(w, r) // string
	data["Alert"] = map[string]interface{}{
		"is_alert": alert != "",
		"message":  alert,
	}

	err = view.ExecuteTemplate(w, "index", data)
	if err != nil {
		fmt.Println("Template execute error:", err)
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Template başarıyla render edildi")
}

func (dashboard Dashboard) NewItem(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Fonksiyon başladı:", r.URL.Path)

	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(helpers.Include("dashboard/add")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Categories"] = models.Category{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (dashboard Dashboard) Add(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Fonksiyon başladı:", r.URL.Path)

	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	title := r.FormValue("forum-title")
	slug := slug.Make(title)
	description := r.FormValue("forum-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("forum-category"))
	content := r.FormValue("forum-content")

	//Upload
	r.ParseMultipartForm(10 << 20)
	file, header, err := r.FormFile("forum-picture")
	if err != nil {
		fmt.Println(err)
		return
	}
	f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = io.Copy(f, file)
	// Upload End
	if err != nil {
		fmt.Println(err)
		return
	}
	models.Post{
		Title:       title,
		Slug:        slug,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: "uploads/" + header.Filename,
		Approved:    true, // admin eklediği için otomatik onaylı
	}.Add()
	helpers.SetAlert(w, r, "Kayıt Başarıyla Eklendi")
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (dashboard Dashboard) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Fonksiyon başladı:", r.URL.Path)

	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	post := models.Post{}.Get(params.ByName("id"))
	post.Delete()
	http.Redirect(w, r, "/admin", http.StatusSeeOther)
}

func (dashboard Dashboard) Edit(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Fonksiyon başladı:", r.URL.Path)

	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(helpers.Include("dashboard/edit")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Post"] = models.Post{}.Get(params.ByName("id"))
	data["Categories"] = models.Category{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (dashboard Dashboard) Update(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	fmt.Println("Fonksiyon başladı:", r.URL.Path)

	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	post := models.Post{}.Get(params.ByName("id"))
	title := r.FormValue("forum-title")
	slug := slug.Make(title)
	description := r.FormValue("forum-desc")
	categoryID, _ := strconv.Atoi(r.FormValue("forum-category"))
	content := r.FormValue("forum-content")
	is_selected := r.FormValue("is_selected")
	var picture_url string

	if is_selected == "1" {
		//Upload
		r.ParseMultipartForm(10 << 20)
		file, header, err := r.FormFile("forum-picture")
		if err != nil {
			fmt.Println(err)
			return
		}
		f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		_, err = io.Copy(f, file)
		// Upload End
		picture_url = "uploads/" + header.Filename
		os.Remove(post.Picture_url)
	} else {
		picture_url = post.Picture_url
	}

	post.Updates(models.Post{
		Title:       title,
		Slug:        slug,
		Description: description,
		CategoryID:  categoryID,
		Content:     content,
		Picture_url: picture_url,
	})
	http.Redirect(w, r, "/admin/edit/"+params.ByName("id"), http.StatusSeeOther)
}
