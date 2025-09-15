package controllers

import (
	"crypto/sha256"
	"fmt"
	"goforum/admin/helpers"
	"goforum/admin/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Userops struct{}

func (userops Userops) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Eğer zaten admin girişi yapılmışsa direkt admin paneline yönlendir
	if helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(
		"admin/views/templates/auth_header.html",
		"admin/views/userops/login.html",
		"admin/views/templates/auth_footer.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := view.ExecuteTemplate(w, "login", nil); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (userops Userops) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))

	// Admin hesabı kontrolü
	if username == "admin" && password == fmt.Sprintf("%x", sha256.Sum256([]byte("123123"))) {
		// Yeni session yönetimi fonksiyonunu kullan
		if err := helpers.SaveAdminSession(w, r, 1, username); err != nil {
			log.Printf("Session save error: %v", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		if err := helpers.SetAlert(w, r, "Admin olarak giriş yaptınız"); err != nil {
			log.Printf("Set alert error: %v", err)
		}
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(
		"admin/views/templates/auth_header.html",
		"admin/views/userops/login.html",
		"admin/views/templates/auth_footer.html",
	)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := view.ExecuteTemplate(w, "login", "Kullanıcı adı veya şifre hatalı"); err != nil {
		log.Printf("Template execution error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (userops Userops) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Yeni session temizleme fonksiyonunu kullan
	if err := helpers.ClearAdminSession(w, r); err != nil {
		log.Printf("Session clear error: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err := helpers.SetAlert(w, r, "Admin panelinden çıkış yapıldı"); err != nil {
		log.Printf("Set alert error: %v", err)
	}
	http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
}

func (userops Userops) Users(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(request) {
		http.Redirect(writer, request, "/admin/login", http.StatusSeeOther)
		return
	}

	view, err := template.New("index").Funcs(template.FuncMap{
		"pendingCount": func(uid uint) int64 {
			var count int64
			models.GetDB().Model(&models.Post{}).Where("author_id = ? AND approved = ?", uid, false).Count(&count)
			return count
		},
	}).ParseFiles(helpers.Include("users/list")...)
	if err != nil {
		log.Println("Template parse error:", err)
		http.Error(writer, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Users": models.User{}.GetAll(),
	}
	alert := helpers.GetAlert(writer, request)
	data["Alert"] = map[string]interface{}{"is_alert": alert != "", "message": alert}

	if err := view.ExecuteTemplate(writer, "index", data); err != nil {
		log.Println("Template execute error:", err)
		http.Error(writer, "Template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (userops Userops) UserPosts(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	uidStr := params.ByName("id")
	uid, _ := strconv.Atoi(uidStr)

	user := models.User{}.Get(uid)
	posts := models.Post{}.GetAll("author_id = ?", uid)

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
	}).ParseFiles(helpers.Include("users/posts")...)
	if err != nil {
		log.Println("Template parse error:", err)
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	alert := helpers.GetAlert(w, r)
	data := map[string]interface{}{
		"UserObj": user,
		"Posts":   posts,
		"Alert": map[string]interface{}{
			"is_alert": alert != "",
			"message":  alert,
		},
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println("Template exec error:", err)
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (userops Userops) ToggleApprove(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	id := params.ByName("id")
	post := models.Post{}.Get(id)
	if post.ID == 0 {
		http.NotFound(w, r)
		return
	}
	newState := !post.Approved
	if newState {
		now := time.Now()
		models.GetDB().Model(&post).Updates(map[string]interface{}{"approved": true, "approved_at": &now})
	} else {
		models.GetDB().Model(&post).Updates(map[string]interface{}{"approved": false, "approved_at": nil})
	}
	helpers.SetAlert(w, r, fmt.Sprintf("Post #%d onay durumu: %v", post.ID, newState))
	ref := r.Referer()
	if ref == "" {
		ref = "/admin"
	}
	http.Redirect(w, r, ref, http.StatusSeeOther)
}

func (userops Userops) PostDetail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	id := params.ByName("id")
	post := models.Post{}.Get(id)
	if post.ID == 0 {
		http.NotFound(w, r)
		return
	}
	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
	}).ParseFiles(helpers.Include("users/postdetail")...)
	if err != nil {
		log.Println("Template parse error:", err)
		http.Error(w, "Template parse error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	alert := helpers.GetAlert(w, r)
	data := map[string]interface{}{
		"Post":  post,
		"Alert": map[string]interface{}{"is_alert": alert != "", "message": alert},
	}
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Println("Template exec error:", err)
		http.Error(w, "Template exec error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
