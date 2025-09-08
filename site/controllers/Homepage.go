package controllers

import (
	"fmt"
	"goblog/site/helpers"
	"goblog/site/models"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

type Homepage struct{}

func (homepage Homepage) Index(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string {
			return models.Category{}.Get(categoryID).Title
		},
		"getDate": func(t time.Time) string {
			return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year())
		},
	}).ParseFiles(helpers.Include("homepage/list")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	view.ExecuteTemplate(w, "index", data)
}

func (homepage Homepage) Detail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slug := params.ByName("slug")
	post := models.Post{}.Get("slug = ?", slug)

	if post.ID == 0 || post.Title == "" {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 - Yazı bulunamadı!"))
		return
	}

	// Ana yorumları al ve yanıt ağacını oluştur
	comment := &models.Comment{}
	rootComments, err := comment.GetCommentTree(int(post.ID))
	if err != nil {
		log.Println("Yorumlar alınırken hata oluştu:", err)
	}

	// Template + FuncMap
	templateFiles := helpers.Include("templates")
	detailFiles := helpers.Include("homepage/detail")
	templateFiles = append(templateFiles, detailFiles...)

	view, err := template.New("index").Funcs(template.FuncMap{
		"getReplies": func(commentID uint) []*models.Comment {
			return nil // Artık buna ihtiyacımız yok, çünkü yorumlar zaten ağaç yapısında geliyor
		},
		"substr": func(s string, start, length int) string {
			if start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
		"upper": func(s string) string {
			if len(s) == 0 {
				return ""
			}
			return strings.ToUpper(s)
		},
	}).ParseFiles(templateFiles...)

	if err != nil {
		log.Printf("Template parse hatası: %v\n", err)
		http.Error(w, "Sayfa yüklenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Post":     post,
		"Comments": rootComments,
	}

	err = view.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Printf("Template render hatası: %v\n", err)
		http.Error(w, "Sayfa görüntülenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}
}
func (homepage Homepage) AddComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slug := params.ByName("slug")
	post := models.Post{}.Get("slug = ?", slug)

	if post.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Yazı bulunamadı!"))
		return
	}

	name := r.FormValue("name")
	commentText := r.FormValue("comment")
	parentIDStr := r.FormValue("parent_id")

	var parentID *uint
	if parentIDStr != "" {
		id, err := strconv.ParseUint(parentIDStr, 10, 64)
		if err == nil {
			u := uint(id)
			parentID = &u
		}
	}

	comment := &models.Comment{
		PostID:   post.ID,
		Name:     name,
		Comment:  commentText,
		ParentID: parentID,
	}

	err := comment.AddComment()
	if err != nil {
		log.Println("Yorum eklenirken hata:", err)
	}

	http.Redirect(w, r, "/yazilar/"+slug, http.StatusSeeOther)
}

func (homepage Homepage) About(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("/about")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	about, err := models.About{}.Get()
	if err != nil {
		log.Println("Hakkında sayfası bilgisi alınırken hata oluştu:", err)
		http.Error(w, "Sunucu Hatası", http.StatusInternalServerError)
		return
	}
	data := make(map[string]interface{})
	data["About"] = about
	view.ExecuteTemplate(w, "index", data)
}

func (homepage Homepage) Contact(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	view, err := template.ParseFiles(helpers.Include("/contact")...)
	if err != nil {
		fmt.Println(err)
		return
	}
	view.ExecuteTemplate(w, "index", nil)
}
