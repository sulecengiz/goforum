package config

import (
	admin "goblog/admin/controllers"
	site "goblog/site/controllers"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Routes() *httprouter.Router {
	r := httprouter.New()
	r.GET("/admin", admin.Dashboard{}.Index)
	r.GET("/admin/yeni-ekle", admin.Dashboard{}.NewItem)
	r.POST("/admin/add", admin.Dashboard{}.Add)
	r.GET("/admin/delete/:id", admin.Dashboard{}.Delete)
	r.GET("/admin/edit/:id", admin.Dashboard{}.Edit)
	r.POST("/admin/update/:id", admin.Dashboard{}.Update)

	r.GET("/admin/kategoriler", admin.Categories{}.Index)
	r.POST("/admin/kategoriler/add", admin.Categories{}.Add)
	r.GET("/admin/kategoriler/delete/:id", admin.Categories{}.Delete)

	r.GET("/admin/login", admin.Userops{}.Index)
	r.POST("/admin/do_login", admin.Userops{}.Login)
	r.GET("/admin/logout", admin.Userops{}.Logout)

	r.GET("/", site.Homepage{}.Index)
	r.GET("/yazilar/:slug", site.Homepage{}.Detail)
	r.POST("/yazilar/:slug/yorum-ekle", site.Homepage{}.AddComment) // Yorum ekleme için yeni POST rotası

	r.GET("/index", site.Homepage{}.About)
	r.GET("/contact", site.Homepage{}.Contact)

	r.POST("/contact/submit", site.ContactFormHandler)
	r.GET("/admin/contact", admin.Contact{}.Index)
	r.GET("/admin/contact/delete/:id", admin.Contact{}.Delete)

	r.GET("/admin/about", admin.AboutIndex)
	r.POST("/admin/about", admin.AboutUpdate)

	r.GET("/admin/comment/comments/:id", admin.Comment{}.PostComment)
	r.GET("/admin/comment/posts", admin.Comment{}.Posts)
	r.GET("/admin/comment/delete/:id", admin.Comment{}.Delete)

	// Statik dosyalar için yönlendirme

	r.ServeFiles("/admin/assets/*filepath", http.Dir("admin/assets"))
	r.ServeFiles("/uploads/*filepath", http.Dir("uploads"))
	r.ServeFiles("/assets/*filepath", http.Dir("site/assets"))
	return r
}
