package controllers

import (
	"fmt"
	"goforum/admin/helpers"
	"goforum/site/models"
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
)

// Comment struct'ı, yönetici panelindeki iletişim sayfası için kontrolcüyü temsil eder.
type Comment struct{}

func (c Comment) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	helpers.SetAlert(w, r, "Yorum yönetimi: lütfen sol menüden 'Post Yorumları'nı seçin.")
	http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
}

// Post listesini gösterir (yorum görüntüleme için seçim ekranı)
func (c Comment) Posts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	view, err := template.ParseFiles(helpers.Include("comment/list")...)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["Posts"] = models.Post{}.GetAll()
	alertMsg := helpers.GetAlert(w, r)
	data["Alert"] = map[string]interface{}{
		"is_alert": alertMsg != "",
		"message":  alertMsg,
	}
	view.ExecuteTemplate(w, "index", data)
}

// Belirli bir post'un yorumlarını gösterir
func (c Comment) PostComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}

	postIDStr := params.ByName("id")
	postID64, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		helpers.SetAlert(w, r, "Geçersiz Post ID")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}
	postID := uint(postID64)

	// Post bilgisini al
	post := models.Post{}.Get("id = ?", postID)
	if post.ID == 0 {
		helpers.SetAlert(w, r, "Post bulunamadı")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	// Yorum ağacını getir
	cm := &models.Comment{}
	rootComments, err := cm.GetCommentTree(int(post.ID))
	if err != nil {
		helpers.SetAlert(w, r, "Yorumlar alınırken hata oluştu")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	// Düz listeye aç (indent için Depth / IndentPx alanları)
	type FlatComment struct {
		ID        uint
		PostID    uint
		ParentID  *uint
		Name      string
		Comment   string
		Short     string
		CreatedAt time.Time
		Depth     int
		IndentPx  int
	}

	var flat []FlatComment
	var walk func(list []*models.Comment, depth int)
	walk = func(list []*models.Comment, depth int) {
		for _, cmt := range list {
			words := strings.Fields(cmt.Comment)
			shortText := cmt.Comment
			if len(words) > 4 {
				shortText = strings.Join(words[:4], " ") + "..."
			}
			flat = append(flat, FlatComment{
				ID:        cmt.ID,
				PostID:    cmt.PostID,
				ParentID:  cmt.ParentID,
				Name:      cmt.Name,
				Comment:   cmt.Comment,
				Short:     shortText,
				CreatedAt: cmt.CreatedAt,
				Depth:     depth,
				IndentPx:  depth * 20,
			})
			if len(cmt.Replies) > 0 {
				walk(cmt.Replies, depth+1)
			}
		}
	}
	walk(rootComments, 0)

	view, err := template.ParseFiles(helpers.Include("comment/comments")...)
	if err != nil {
		fmt.Println("Template parse hatası:", err)
		helpers.SetAlert(w, r, "Şablon yüklenemedi")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	alertMsg := helpers.GetAlert(w, r)
	data := map[string]interface{}{
		"Post":     post,
		"Comments": flat,
		"Alert": map[string]interface{}{
			"is_alert": alertMsg != "",
			"message":  alertMsg,
		},
	}
	view.ExecuteTemplate(w, "index", data)
}

// Yorum silme (tek yorum). İstenirse alt yorumlar ek olarak silinebilir.
func (c Comment) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	if !helpers.IsAdminLoggedIn(r) {
		http.Redirect(w, r, "/admin/login", http.StatusSeeOther)
		return
	}
	idStr := params.ByName("id")
	id64, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		helpers.SetAlert(w, r, "Geçersiz yorum id")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}
	cid := uint(id64)

	// Yorumu bul (redirect için PostID lazım)
	var found models.Comment
	if err := models.GetDB().First(&found, cid).Error; err != nil || found.ID == 0 {
		helpers.SetAlert(w, r, "Yorum bulunamadı")
		http.Redirect(w, r, "/admin/comment/posts", http.StatusSeeOther)
		return
	}

	comment := &models.Comment{}
	deleteErr := comment.DeleteComment(int(cid))
	if deleteErr != nil {
		helpers.SetAlert(w, r, "Yorum silinirken hata")
	} else {
		helpers.SetAlert(w, r, "Yorum silindi")
	}

	http.Redirect(w, r, "/admin/comment/comments/"+strconv.FormatUint(uint64(found.PostID), 10), http.StatusSeeOther)
}
