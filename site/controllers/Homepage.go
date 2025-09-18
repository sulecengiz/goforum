package controllers

import (
	"crypto/sha256"
	"fmt"
	"goforum/site/helpers"
	"goforum/site/models"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gosimple/slug"
	"github.com/julienschmidt/httprouter"
)

// Ortak fonksiyon haritası (profile vb. her yerde lazım olanlar)
func siteFuncMap() template.FuncMap {
	return template.FuncMap{
		"getCategory":  func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
		"getDate":      func(t time.Time) string { return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year()) },
		"upper":        strings.ToUpper,
		"getPostSlug":  func(id uint) string { return models.Post{}.Get(id).Slug },
		"getPostTitle": func(id uint) string { return models.Post{}.Get(id).Title },
		"html":         func(content string) template.HTML { return template.HTML(content) },
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
		"inSlice": func(slice []uint, item uint) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		"getUser": func(userID uint) models.User { return models.User{}.Get(userID) },
		"isAdmin": func(userID uint) bool { return userID == 1 }, // Admin kontrolü
	}
}

// Ortak template yükleyici
func loadSiteTemplates(name string, files ...string) (*template.Template, error) {
	partials := helpers.Include("templates") // layout, navbar, footer, vs
	all := append(partials, files...)
	return template.New(name).Funcs(siteFuncMap()).ParseFiles(all...)
}

type Homepage struct{}

// Anasayfa
func (homepage Homepage) Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Template setini oluştur (templates + homepage/list)
	templateFiles := helpers.Include("templates")
	listFiles := helpers.Include("homepage/list")
	templateFiles = append(templateFiles, listFiles...)

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
		"getDate":     func(t time.Time) string { return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year()) },
		"getUser":     func(id uint) models.User { return models.User{}.Get(id) },
		"upper":       strings.ToUpper,
		"isAdmin": func(authorID uint) bool {
			if authorID == 0 {
				return true
			}
			u := models.User{}.Get(authorID)
			return u.Role == 1
		},
		"firstInitial": func(s string) string {
			if s == "" {
				return "?"
			}
			r, _ := utf8.DecodeRuneInString(s)
			return strings.ToUpper(string(r))
		},
		"inSlice": func(slice []uint, item uint) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	}).ParseFiles(templateFiles...)
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	user, err := helpers.GetCurrentUser(r)
	userStruct := user
	if err != nil {
		log.Printf("GetCurrentUser hatası: %v", err)
		user = nil
	}

	// Kullanıcının kaydettiği postların ID'lerini al
	var savedPostIDs []uint
	if user != nil {
		savedPost := models.SavedPost{}
		savedPosts := savedPost.GetByUser(user.ID)
		for _, sp := range savedPosts {
			savedPostIDs = append(savedPostIDs, sp.PostID)
		}
	}

	categories := models.Category{}.GetAll()
	data := map[string]interface{}{
		"Posts":              models.Post{}.GetAll("approved = ?", true),
		"User":               userStruct,
		"Categories":         categories,
		"totalCategoryPosts": countCategoryPosts(),
		"SavedPostIDs":       savedPostIDs,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err = view.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template render hatası", http.StatusInternalServerError)
	}
}

func countComments(comments []*models.Comment) int {
	total := len(comments)
	for _, c := range comments {
		total += countComments(c.Replies)
	}
	return total
}

func countCategoryPosts() map[string]int {
	categories := models.Category{}.GetAll()
	counts := make(map[string]int, len(categories))

	for _, category := range categories {
		posts := models.Post{}.GetAll("category_id = ? AND approved = ?", category.ID, true)
		counts[category.Title] = len(posts)
	}
	return counts
}

// Kategori sayfası
func (homepage Homepage) Category(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slugParam := params.ByName("slug")
	category := models.Category{}.GetBySlug(slugParam)
	if category.ID == 0 {
		http.NotFound(w, r)
		return
	}
	posts := models.Post{}.GetAll("category_id = ? AND approved = ?", category.ID, true)

	// templates + homepage/list dahil et
	templateFiles := helpers.Include("templates")
	listFiles := helpers.Include("homepage/list")
	templateFiles = append(templateFiles, listFiles...)

	view, err := template.New("index").Funcs(template.FuncMap{
		"getCategory": func(categoryID int) string { return models.Category{}.Get(categoryID).Title },
		"getDate":     func(t time.Time) string { return fmt.Sprintf("%02d.%02d.%d", t.Day(), int(t.Month()), t.Year()) },
		"getUser":     func(id uint) models.User { return models.User{}.Get(id) },
		"upper":       strings.ToUpper,
		"isAdmin": func(authorID uint) bool {
			if authorID == 0 {
				return true
			}
			u := models.User{}.Get(authorID)
			return u.Role == 1
		},
		"firstInitial": func(s string) string {
			if s == "" {
				return "?"
			}
			r, _ := utf8.DecodeRuneInString(s)
			return strings.ToUpper(string(r))
		},
		"inSlice": func(slice []uint, item uint) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
	}).ParseFiles(templateFiles...)
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	user, _ := helpers.GetCurrentUser(r)

	// Kullanıcının kaydettiği postların ID'lerini al
	var savedPostIDs []uint
	if user != nil {
		savedPost := models.SavedPost{}
		savedPosts := savedPost.GetByUser(user.ID)
		for _, sp := range savedPosts {
			savedPostIDs = append(savedPostIDs, sp.PostID)
		}
	}

	data := map[string]interface{}{
		"Posts":              posts,
		"User":               user,
		"Categories":         models.Category{}.GetAll(),
		"totalCategoryPosts": countCategoryPosts(),
		"ActiveCategory":     category,
		"SavedPostIDs":       savedPostIDs,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := view.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template render hatası", http.StatusInternalServerError)
	}
}

// Yazı detay sayfası
func (homepage Homepage) Detail(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slugParam := params.ByName("slug")
	post := models.Post{}.Get("slug = ?", slugParam)

	user, _ := helpers.GetCurrentUser(r)
	if post.ID == 0 || post.Title == "" || (!post.Approved && (user == nil || user.Role != 1) && !(user != nil && user.ID == post.AuthorID)) {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("404 - Yazı bulunamadı!")); err != nil {
			log.Printf("Write error: %v", err)
		}
		return
	}

	comment := &models.Comment{}
	rootComments, err := comment.GetCommentTree(int(post.ID))
	if err != nil {
		log.Println("Yorumlar alınırken hata oluştu:", err)
	}

	// Beğeni bilgilerini yorumlara ekle
	annotateComments(rootComments, user)

	templateFiles := helpers.Include("templates")
	detailFiles := helpers.Include("homepage/detail")
	templateFiles = append(templateFiles, detailFiles...)

	view, err := template.New("index").Funcs(template.FuncMap{
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
		"upper": func(s string) string { return strings.ToUpper(s) },
		"inSlice": func(slice []uint, item uint) bool {
			for _, s := range slice {
				if s == item {
					return true
				}
			}
			return false
		},
		// Post içeriğini HTML olarak işlemek için helper
		"html": func(content string) template.HTML { return template.HTML(content) },
	}).ParseFiles(templateFiles...)

	if err != nil {
		log.Printf("Template parse hatası: %v\n", err)
		http.Error(w, "Sayfa yüklenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}

	// Kullanıcının kaydettiği postların ID'lerini al
	var savedPostIDs []uint
	if user != nil {
		savedPost := models.SavedPost{}
		savedPosts := savedPost.GetByUser(user.ID)
		for _, sp := range savedPosts {
			savedPostIDs = append(savedPostIDs, sp.PostID)
		}
	}

	data := map[string]interface{}{
		"Post":          post,
		"Comments":      rootComments,
		"TotalComments": countComments(rootComments),
		"User":          user,
		"SavedPostIDs":  savedPostIDs,
	}

	err = view.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Printf("Template render hatası: %v\n", err)
		http.Error(w, "Sayfa görüntülenirken bir hata oluştu", http.StatusInternalServerError)
		return
	}
}

// Yorum ağacını dolaşıp beğeni durumlarını toplu ekleyen optimize fonksiyon
func annotateComments(list []*models.Comment, user *models.User) {
	if len(list) == 0 {
		return
	}
	// Tüm yorum ID'lerini topla
	var all []*models.Comment
	var ids []uint
	var walk func(items []*models.Comment)
	walk = func(items []*models.Comment) {
		for _, c := range items {
			all = append(all, c)
			ids = append(ids, c.ID)
			if len(c.Replies) > 0 {
				walk(c.Replies)
			}
		}
	}
	walk(list)

	// Toplu like sayıları
	counts := models.Like{}.GetLikeCounts(ids)
	var liked map[uint]bool
	if user != nil {
		liked = models.Like{}.GetUserLikedCommentIDs(user.ID, ids)
	} else {
		liked = make(map[uint]bool)
	}

	// Değerleri yerleştir
	for _, c := range all {
		c.LikeCount = counts[c.ID]
		c.IsLiked = liked[c.ID]
	}
}

func (homepage Homepage) About(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	partials := helpers.Include("templates")
	aboutFiles := helpers.Include("about")
	all := append(partials, aboutFiles...)
	view, err := template.New("index").ParseFiles(all...)
	if err != nil {
		log.Printf("Template parse hatası (About): %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	user, _ := helpers.GetCurrentUser(r)
	about, err := models.About{}.Get()
	if err != nil {
		log.Println("Hakkında sayfası bilgisi alınırken hata oluştu:", err)
		http.Error(w, "Sunucu Hatası", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"About": about,
		"User":  user,
	}
	if err = view.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template render hatası", http.StatusInternalServerError)
	}
}

func (homepage Homepage) Contact(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	partials := helpers.Include("templates")
	contactFiles := helpers.Include("contact")
	all := append(partials, contactFiles...)
	view, err := template.New("index").ParseFiles(all...)
	if err != nil {
		log.Printf("Template parse hatası (Contact): %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}
	user, _ := helpers.GetCurrentUser(r)
	data := map[string]interface{}{
		"User": user,
	}
	if err = view.ExecuteTemplate(w, "index", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template render hatası", http.StatusInternalServerError)
	}
}

func (homepage Homepage) AddComment(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	slugParam := params.ByName("slug")
	post := models.Post{}.Get("slug = ?", slugParam)

	if post.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, err := w.Write([]byte("Yazı bulunamadı!")); err != nil {
			log.Printf("Write error: %v", err)
		}
		return
	}

	// Get current user
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		// If user is not logged in, redirect to login
		http.Redirect(w, r, "/login", http.StatusSeeOther)
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
		UserID:   &user.ID, // Set the UserID from the logged-in user
		Name:     name,
		Comment:  commentText,
		ParentID: parentID,
	}

	err = comment.AddComment()
	if err != nil {
		log.Println("Yorum eklenirken hata:", err)
	}

	http.Redirect(w, r, "/yazilar/"+slugParam, http.StatusSeeOther)
}

func (homepage Homepage) Login(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	view, err := loadSiteTemplates("layout", "site/views/userops/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":     "Giriş Yap",
		"BodyClass": "auth-body",
	}

	err = view.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (homepage Homepage) LoginPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))

	log.Printf("Giriş denemesi - Username: %s", username)

	user := &models.User{
		Username: username,
		Password: password,
	}

	if dbUser := user.Login(); dbUser != nil {
		log.Printf("Giriş başarılı - User ID: %d, Username: %s, Role: %d", dbUser.ID, dbUser.Username, dbUser.Role)

		session, _ := helpers.SessionStore.Get(r, "session")
		session.Values["userID"] = dbUser.ID
		session.Values["username"] = dbUser.Username
		session.Values["role"] = dbUser.Role

		// Session kaydetmeden önce değerleri kontrol et
		log.Printf("Session'a kaydedilecek değerler - userID: %v (tip: %T), username: %s, role: %d",
			dbUser.ID, dbUser.ID, dbUser.Username, dbUser.Role)

		err := session.Save(r, w)
		if err != nil {
			log.Printf("Session kaydetme hatası: %v", err)
		} else {
			log.Println("Session başarıyla kaydedildi")
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	log.Println("Giriş başarısız - Kullanıcı adı veya şifre hatalı")

	// Başarısız giriş
	view, _ := loadSiteTemplates("layout", "site/views/userops/login.html")
	if err := view.ExecuteTemplate(w, "layout", map[string]interface{}{
		"Title":     "Giriş Yap",
		"Error":     "Kullanıcı adı veya şifre hatalı",
		"BodyClass": "auth-body",
	}); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template render hatası", http.StatusInternalServerError)
	}
}

func (homepage Homepage) Register(w http.ResponseWriter, _ *http.Request, _ httprouter.Params) {
	view, err := loadSiteTemplates("layout", "site/views/userops/register.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":     "Kayıt Ol",
		"BodyClass": "auth-body",
	}

	err = view.ExecuteTemplate(w, "layout", data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (homepage Homepage) RegisterPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	username := r.FormValue("username")
	password := fmt.Sprintf("%x", sha256.Sum256([]byte(r.FormValue("password"))))
	email := r.FormValue("email")

	user := &models.User{
		Username: username,
		Password: password,
		Email:    email,
		Role:     0, // okuyucu rolü
	}

	if err := user.Register(); err != nil {
		view, _ := loadSiteTemplates("layout", "site/views/userops/register.html")
		if execErr := view.ExecuteTemplate(w, "layout", map[string]interface{}{
			"Title":     "Kayıt Ol",
			"Error":     "Kayıt işlemi başarısız: " + err.Error(),
			"BodyClass": "auth-body",
		}); execErr != nil {
			log.Printf("Template execute hatası: %v", execErr)
			http.Error(w, "Template render hatası", http.StatusInternalServerError)
		}
		return
	}

	// Kayıt başarılı, hemen giriş yaptıralım
	loginUser := &models.User{
		Username: username,
		Password: password,
	}

	if dbUser := loginUser.Login(); dbUser != nil {
		session, _ := helpers.SessionStore.Get(r, "session")
		session.Values["userID"] = dbUser.ID
		session.Values["username"] = dbUser.Username
		session.Values["role"] = dbUser.Role
		if err := session.Save(r, w); err != nil {
			log.Printf("Session kaydetme hatası: %v", err)
		}
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (homepage Homepage) Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, _ := helpers.SessionStore.Get(r, "session")

	// Tüm session değerlerini temizle
	delete(session.Values, "userID")
	delete(session.Values, "username")
	delete(session.Values, "role")

	// Ya da session'ı tamamen temizle
	session.Values = make(map[interface{}]interface{})

	if err := session.Save(r, w); err != nil {
		log.Printf("Session kaydetme hatası: %v", err)
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Profile sayfası
func (homepage Homepage) Profile(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcının giriş yapıp yapmadığını kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Kullanıcının bloglarını getir
	userPosts := models.Post{}.GetByAuthor(user.ID)
	approvedPosts := []models.Post{}
	pendingPosts := []models.Post{}
	for _, p := range userPosts {
		if p.Approved {
			approvedPosts = append(approvedPosts, p)
		} else {
			pendingPosts = append(pendingPosts, p)
		}
	}
	userComments := models.Comment{}.GetByUser(user.ID)

	// Kullanıcının kaydettiği blogları getir
	var savedPosts []models.SavedPost
	var savedBlogs []models.Post
	models.GetDB().Where("user_id = ?", user.ID).Find(&savedPosts)
	for _, sp := range savedPosts {
		var post models.Post
		if err := models.GetDB().First(&post, sp.PostID).Error; err == nil && post.ID != 0 {
			savedBlogs = append(savedBlogs, post)
		}
	}

	view, err := loadSiteTemplates("layout", "site/views/profile/profile.html")
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":         "Profil",
		"User":          user,
		"Posts":         userPosts,
		"ApprovedPosts": approvedPosts,
		"PendingPosts":  pendingPosts,
		"Categories":    models.Category{}.GetAll(),
		"SavedBlogs":    savedBlogs,
		"Comments":      userComments,
		// Eklendi:
		"UserComments": userComments,
		"PostCount":    len(userPosts),
		"CommentCount": len(userComments),
	}

	if err = view.ExecuteTemplate(w, "layout", data); err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
	}
}

// Yeni blog yazısı sayfası
func (homepage Homepage) NewPost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcının giriş yapıp yapmadığını kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	view, err := loadSiteTemplates("layout", "site/views/profile/new-post.html")
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Yeni Blog Yazısı",
		"User":       user,
		"Categories": models.Category{}.GetAll(),
	}

	err = view.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
	}
}

// Blog yazısı kaydetme
func (homepage Homepage) CreatePost(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Kullanıcının giriş yapıp yapmadığını kontrol et
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Form verilerini al
	title := r.FormValue("title")
	content := r.FormValue("content")
	description := r.FormValue("description")
	categoryIDStr := r.FormValue("category_id")

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		categoryID = 1 // Varsayılan kategori
	}

	// Slug oluştur (admin paneli gibi)
	slugStr := slug.Make(title)

	var picture_url string

	// Fotoğraf yükleme işlemi
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("ParseMultipartForm hatası: %v", err)
	}
	file, header, err := r.FormFile("blog-picture")
	if err != nil {
		// Fotoğraf yüklenmemişse varsayılan resim kullan
		picture_url = ""
	} else {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("File close hatası: %v", closeErr)
		}

		// Dosyayı kaydet
		f, err := os.OpenFile("uploads/"+header.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Printf("Dosya açılırken hata: %v", err)
			picture_url = ""
		} else {
			if closeErr := f.Close(); closeErr != nil {
				log.Printf("File close hatası: %v", closeErr)
			}
			_, err = io.Copy(f, file)
			if err != nil {
				log.Printf("Dosya kopyalanırken hata: %v", err)
				picture_url = ""
			} else {
				picture_url = "uploads/" + header.Filename
			}
		}
	}

	// Post oluştur
	post := models.Post{
		Title:       title,
		Slug:        slugStr,
		Description: description,
		Content:     content,
		CategoryID:  categoryID,
		AuthorID:    user.ID,
		Picture_url: picture_url,
	}

	post.Add()

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Blog düzenleme sayfası
func (homepage Homepage) EditPost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}

	post := models.Post{}.GetByIDAndAuthor(uint(postID), user.ID)
	if post.ID == 0 {
		http.Error(w, "Blog yazısı bulunamadı veya yetkiniz yok", http.StatusNotFound)
		return
	}

	view, err := loadSiteTemplates("layout", "site/views/profile/edit-post.html")
	if err != nil {
		log.Printf("Template parse hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":      "Blog Düzenle",
		"User":       user,
		"Post":       post,
		"Categories": models.Category{}.GetAll(),
	}

	err = view.ExecuteTemplate(w, "layout", data)
	if err != nil {
		log.Printf("Template execute hatası: %v", err)
		http.Error(w, "Template hatası", http.StatusInternalServerError)
	}
}

// Blog güncelleme
func (homepage Homepage) UpdatePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID64, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}
	postID := uint(postID64)

	existingPost := models.Post{}.GetByIDAndAuthor(postID, user.ID)
	if existingPost.ID == 0 {
		http.Error(w, "Blog yazısı bulunamadı veya yetkiniz yok", http.StatusNotFound)
		return
	}

	// Kullanıcı (admin olmayan) onaylı bir postu güncellediğinde approved=false ve approved_at=nil olacak şekilde tekrar onaya düşürme
	needResetApproval := existingPost.Approved && (user.Role != 1)

	// Form verileri + multipart
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		log.Printf("Multipart parse hata: %v", err)
	}
	title := r.FormValue("title")
	content := r.FormValue("content")
	description := r.FormValue("description")
	categoryIDStr := r.FormValue("category_id")
	categoryID, convErr := strconv.Atoi(categoryIDStr)
	if convErr != nil {
		categoryID = existingPost.CategoryID
	}

	// Slug üret
	slugStr := strings.ToLower(strings.ReplaceAll(title, " ", "-"))
	replacements := map[string]string{"ç": "c", "ğ": "g", "ı": "i", "ö": "o", "ş": "s", "ü": "u"}
	for k, v := range replacements {
		slugStr = strings.ReplaceAll(slugStr, k, v)
	}

	// Resim (opsiyonel)
	pictureURL := existingPost.Picture_url
	if file, header, ferr := r.FormFile("blog-picture"); ferr == nil {
		if closeErr := file.Close(); closeErr != nil {
			log.Printf("File close hatası: %v", closeErr)
		}
		out, oerr := os.OpenFile("uploads/"+header.Filename, os.O_CREATE|os.O_WRONLY, 0666)
		if oerr == nil {
			if closeErr := out.Close(); closeErr != nil {
				log.Printf("File close hatası: %v", closeErr)
			}
			if _, cerr := io.Copy(out, file); cerr == nil {
				pictureURL = "uploads/" + header.Filename
			}
		}
	}

	updated := models.Post{
		Title:       title,
		Slug:        slugStr,
		Description: description,
		Content:     content,
		CategoryID:  categoryID,
		Picture_url: pictureURL,
	}
	editErr := models.Post{}.EditPost(postID, updated)
	if editErr != nil {
		log.Printf("Blog güncelleme hatası: %v", editErr)
		http.Error(w, "Blog güncellenirken hata oluştu", http.StatusInternalServerError)
		return
	}
	if needResetApproval {
		// GORM struct update false değerini yazmayacağı için map kullanıyoruz
		models.GetDB().Model(&models.Post{}).Where("id = ?", postID).Updates(map[string]interface{}{
			"approved":    false,
			"approved_at": nil,
		})
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Blog silme
func (homepage Homepage) DeletePost(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	postID, err := strconv.ParseUint(params.ByName("id"), 10, 32)
	if err != nil {
		http.Error(w, "Geçersiz post ID", http.StatusBadRequest)
		return
	}

	err = models.Post{}.DeleteByAuthor(uint(postID), user.ID)
	if err != nil {
		log.Printf("Blog silme hatası: %v", err)
		http.Error(w, "Blog silinirken hata oluştu", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// Yorum beğenme / beğeni kaldırma
func (homepage Homepage) ToggleLike(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	user, err := helpers.GetCurrentUser(r)
	if err != nil || user == nil {
		w.WriteHeader(http.StatusUnauthorized)
		if _, writeErr := w.Write([]byte(`{"success":false,"error":"giris gerekli"}`)); writeErr != nil {
			log.Printf("Write error: %v", writeErr)
		}
		return
	}

	commentID64, err := strconv.ParseUint(params.ByName("commentId"), 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		if _, writeErr := w.Write([]byte(`{"success":false,"error":"gecersiz yorum id"}`)); writeErr != nil {
			log.Printf("Write error: %v", writeErr)
		}
		return
	}
	commentID := uint(commentID64)

	// Yorum var mı kontrol et
	var c models.Comment
	if err := models.GetDB().First(&c, commentID).Error; err != nil || c.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		if _, writeErr := w.Write([]byte(`{"success":false,"error":"yorum bulunamadi"}`)); writeErr != nil {
			log.Printf("Write error: %v", writeErr)
		}
		return
	}

	like := models.Like{UserID: user.ID, CommentID: commentID}
	isLiked, err := like.Toggle()
	if err != nil {
		log.Printf("ToggleLike hata user=%d comment=%d err=%v", user.ID, commentID, err)
		w.WriteHeader(http.StatusInternalServerError)
		// Kullanıcıya genel mesaj, log'da detay var
		if _, writeErr := w.Write([]byte(`{"success":false,"error":"begeni islem hatasi"}`)); writeErr != nil {
			log.Printf("Write error: %v", writeErr)
		}
		return
	}

	likeCount := models.Like{}.GetLikeCount(commentID)
	if _, writeErr := w.Write([]byte(fmt.Sprintf(`{"success":true,"likeCount":%d,"isLiked":%t}`, likeCount, isLiked))); writeErr != nil {
		log.Printf("Write error: %v", writeErr)
	}
}
