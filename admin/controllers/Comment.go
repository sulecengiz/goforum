package controllers

import (
	"fmt"
	"goblog/admin/helpers"
	"goblog/site/models"
	"html/template"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Comment struct'ı, yönetici panelindeki iletişim sayfası için kontrolcüyü temsil eder.
type Comment struct{}

func (c Comment) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}
	view, err := template.ParseFiles(helpers.Include("comment/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Comments"] = models.Comment{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}

func (c Comment) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}

	id := params.ByName("id")

	// Önce kaydı bulmaya çalışın.
	Comment, err := models.Comment{}.Get(id)
	if err != nil {
		// Kayıt bulunamadıysa veya başka bir hata oluştuysa
		helpers.SetAlert(w, r, "Kayıt bulunamadı veya bir hata oluştu: "+err.Error())
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	err = Comment.Delete()
	if err != nil {
		helpers.SetAlert(w, r, "Kayıt silinirken bir hata oluştu: "+err.Error())
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	helpers.SetAlert(w, r, "Kayıt başarıyla silindi.")
	http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
}

func (c Comment) PostComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}

	// URL'den post ID'sini al
	postIDStr := params.ByName("id")
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		helpers.SetAlert(w, r, "Geçersiz post ID")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(helpers.Include("comment/comments")...)
	if err != nil {
		fmt.Println(err)
		return
	}

	post := models.Post{}.Get(postIDStr)

	// Admin paneli için flat liste alıyoruz (ParentID ve ParentName ile birlikte)
	comments, err := (&models.Comment{}).GetCommentsByPostID(postID)
	if err != nil {
		fmt.Println(err)
		return
	}

	data := make(map[string]interface{})
	data["Post"] = post
	data["Comments"] = comments // ParentID ve ParentName ile reply’ler görünecek
	data["Alert"] = helpers.GetAlert(w, r)

	view.ExecuteTemplate(w, "index", data)
}

// Post listesini gösteren fonksiyon
func (c Comment) Posts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.CheckUser(w, r) {
		return
	}
	view, err := template.ParseFiles(helpers.Include("comment/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	data["Alert"] = helpers.GetAlert(w, r)
	view.ExecuteTemplate(w, "index", data)
}
