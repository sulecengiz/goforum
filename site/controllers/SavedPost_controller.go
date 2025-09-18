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

// SavePost - Blog yazısını kaydetme/kaldırma (toggle)
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

	// Kaydetme durumunu değiştir (toggle)
	isSaved, err := models.SavedPost{}.ToggleSave(user.ID, uint(postID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"message": "Bir hata oluştu",
		}); err != nil {
			log.Printf("JSON encode error: %v", err)
		}
		return
	}

	// Başarılı yanıt
	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"isSaved": isSaved,
		"message": func() string {
			if isSaved {
				return "Post kaydedildi"
			}
			return "Post kayıtlardan kaldırıldı"
		}(),
	}); err != nil {
		log.Printf("JSON encode error: %v", err)
	}
}

// UnsavePost - Blog yazısını kayıtlardan kaldırma (deprecated - SavePost toggle kullanıyor)
func UnsavePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	// Bu fonksiyon artık gerekli değil çünkü SavePost toggle mantığı kullanıyor
	// Ama mevcut route'lar için uyumluluk sağlamak adına SavePost'u çağırıyoruz
	SavePost(w, r, params)
}
