package main

import (
	admin_models "goblog/admin/models"
	"goblog/config"
	site_models "goblog/site/models"
	"net/http"
)

func main() {
	admin_models.Post{}.Migrate()
	admin_models.User{}.Migrate()
	admin_models.Category{}.Migrate()
	site_models.ContactForm{}.Migrate()
	site_models.About{}.Migrate()
	site_models.Comment{}.Migrate()
	http.ListenAndServe(":8080", config.Routes())
}
