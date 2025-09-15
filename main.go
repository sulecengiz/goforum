package main

import (
	admin_models "goforum/admin/models"
	"goforum/config"
	site_models "goforum/site/models"
	"net/http"
)

func main() {
	site_models.ConnectDB()

	admin_models.Post{}.Migrate()
	admin_models.User{}.Migrate()
	admin_models.Category{}.Migrate()

	site_models.Category{}.Migrate()
	site_models.Post{}.Migrate()
	site_models.Comment{}.Migrate()
	site_models.Like{}.Migrate()
	site_models.ContactForm{}.Migrate()
	site_models.About{}.Migrate()

	// Eski admin (author_id=0 veya admin kullanıcı) postlarını onayla
	site_models.BackfillAdminPosts()

	http.ListenAndServe(":8080", config.Routes())
}
