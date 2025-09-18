package config

import (
	admin "goforum/admin/controllers"
	"goforum/admin/middleware"
	site "goforum/site/controllers"
	"goforum/site/helpers"
	"html/template"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes() *httprouter.Router {
	r := httprouter.New()

	// Statik dosyalar için yönlendirme
	r.ServeFiles("/admin/assets/*filepath", http.Dir("admin/assets"))
	r.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	r.ServeFiles("/assets/*filepath", http.Dir("site/assets"))

	// Site Routes (public)
	r.GET("/", site.Homepage{}.Index)
	r.GET("/login", site.Homepage{}.Login)
	r.POST("/login", site.Homepage{}.LoginPost)
	r.GET("/register", site.Homepage{}.Register)
	r.POST("/register", site.Homepage{}.RegisterPost)
	r.GET("/logout", site.Homepage{}.Logout)
	r.GET("/yazilar/:slug", site.Homepage{}.Detail)
	r.GET("/kategoriler/:slug", site.Homepage{}.Index)
	r.POST("/yazilar/:slug/yorum-ekle", site.Homepage{}.AddComment)
	r.GET("/about", site.Homepage{}.About)
	r.GET("/contact", site.Homepage{}.Contact)
	r.POST("/contact/submit", site.ContactFormHandler)
	r.GET("/kategori/:slug", site.Homepage{}.Category)

	// Profil sayfaları
	r.GET("/profile", site.Homepage{}.Profile)
	r.GET("/new-post", site.Homepage{}.NewPost)
	r.POST("/create-post", site.Homepage{}.CreatePost)

	// Blog yönetimi
	r.GET("/edit-post/:id", site.Homepage{}.EditPost)
	r.POST("/update-post/:id", site.Homepage{}.UpdatePost)
	r.GET("/delete-post/:id", site.Homepage{}.DeletePost)

	// Yorum beğeni sistemi (conflict fix: önceki /yazilar/yorum/... wildcard slug ile çakışıyordu)
	r.POST("/like-comment/:commentId", site.Homepage{}.ToggleLike)

	// SavedPost (blog kaydetme) işlemleri
	r.POST("/save-post/:postID", site.SavePost)
	r.POST("/unsave-post/:postID", site.UnsavePost)
	r.GET("/saved-posts", site.Homepage{}.GetSavedPosts)

	// Admin Routes
	r.GET("/admin/login", admin.Userops{}.Index)
	r.POST("/admin/do_login", admin.Userops{}.Login)
	r.GET("/admin/logout", admin.Userops{}.Logout)

	// Admin Panel Routes (Protected)
	r.GET("/admin", middleware.CheckAdmin(admin.Dashboard{}.Index))
	r.GET("/admin/yeni-ekle", middleware.CheckAdmin(admin.Dashboard{}.NewItem))
	r.POST("/admin/add", middleware.CheckAdmin(admin.Dashboard{}.Add))
	r.GET("/admin/delete/:id", middleware.CheckAdmin(admin.Dashboard{}.Delete))
	r.GET("/admin/edit/:id", middleware.CheckAdmin(admin.Dashboard{}.Edit))
	r.POST("/admin/update/:id", middleware.CheckAdmin(admin.Dashboard{}.Update))

	r.GET("/admin/kategoriler", middleware.CheckAdmin(admin.Categories{}.Index))
	r.POST("/admin/kategoriler/add", middleware.CheckAdmin(admin.Categories{}.Add))
	r.GET("/admin/kategoriler/delete/:id", middleware.CheckAdmin(admin.Categories{}.Delete))

	r.GET("/admin/contact", middleware.CheckAdmin(admin.Contact{}.Index))
	r.GET("/admin/contact/delete/:id", middleware.CheckAdmin(admin.Contact{}.Delete))

	r.GET("/admin/about", middleware.CheckAdmin(admin.AboutIndex))
	r.POST("/admin/about", middleware.CheckAdmin(admin.AboutUpdate))

	r.GET("/admin/comment/comments/:id", middleware.CheckAdmin(admin.Comment{}.PostComment))
	r.GET("/admin/comment/posts", middleware.CheckAdmin(admin.Comment{}.Posts))
	r.GET("/admin/comment/delete/:id", middleware.CheckAdmin(admin.Comment{}.Delete))

	r.GET("/admin/user", middleware.CheckAdmin(admin.Userops{}.Users))
	r.GET("/admin/user/:id/posts", middleware.CheckAdmin(admin.Userops{}.UserPosts))
	r.GET("/admin/post/toggle-approve/:id", middleware.CheckAdmin(admin.Userops{}.ToggleApprove))
	r.GET("/admin/post/detail/:id", middleware.CheckAdmin(admin.Userops{}.PostDetail))

	// 404 Handler
	r.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		partials := helpers.Include("templates") // layout, navbar, scripts, footer, head, vs.
		partials = append(partials, "site/views/404.html")
		view, err := template.New("layout").ParseFiles(partials...)
		if err != nil {
			http.Error(w, "Template parsing error", http.StatusInternalServerError)
			return
		}

		if err := view.ExecuteTemplate(w, "layout", map[string]interface{}{
			"Title": "404 - Sayfa Bulunamadı",
		}); err != nil {
			http.Error(w, "Template execution error", http.StatusInternalServerError)
		}
	})

	return r
}
