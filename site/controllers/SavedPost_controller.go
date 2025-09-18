package controllers

import (
	"encoding/json"
	"goforum/site/helpers"
	"goforum/site/models"
	"log"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// SavePost - Blog yazısını kaydetme
func SavePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Kullanıcı giriş yapmış mı kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Giriş yapmalısınız",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Post ID'yi al
	postIDStr := params.ByName("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Geçersiz post ID",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Post var mı kontrol et
	post := models.Post{}.Get("id = ?", uint(postID))
	if post.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Post bulunamadı",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Postu kaydet
	savedPost := &models.SavedPost{}
	if err := savedPost.Save(user.ID, uint(postID)); err != nil {
		log.Printf("Post kaydetme hatası: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Post kaydedilemedi",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post başarıyla kaydedildi",
	}); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// UnsavePost - Kaydedilmiş blog yazısını kaldırma
func UnsavePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	// Kullanıcı giriş yapmış mı kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Giriş yapmalısınız",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Post ID'yi al
	postIDStr := params.ByName("postID")
	postID, err := strconv.ParseUint(postIDStr, 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Geçersiz post ID",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Kaydedilmiş postu kaldır
	savedPost := models.SavedPost{}
	if err := savedPost.Unsave(user.ID, uint(postID)); err != nil {
		log.Printf("Post kayıt kaldırma hatası: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Post kaydı kaldırılamadı",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Post kaydı kaldırıldı",
	}); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// GetSavedPosts - Kaydedilmiş blog yazılarını listeleme
func GetSavedPosts(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcı giriş yapmış mı kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Template yükle
	view, err := loadSiteTemplates("layout", "site/views/profile/saved-posts.html")
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	// Kaydedilmiş postları getir
	savedPost := models.SavedPost{}
	savedPosts := savedPost.GetByUser(user.ID)

	// Sadece post bilgilerini çıkar
	var posts []models.Post
	for _, sp := range savedPosts {
		posts = append(posts, sp.Post)
	}

	data := map[string]interface{}{
		"Title":      "Kaydedilmiş Yazılar",
		"User":       user,
		"SavedPosts": posts,
	}

	if err = view.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
	}
}
