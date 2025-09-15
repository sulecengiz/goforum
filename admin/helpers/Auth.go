// helpers/auth.go
package helpers

import (
	"goforum/admin/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckUser(r *http.Request) bool {
	sess, err := GetAdminSession(r)
	if err != nil {
		return false
	}

	user, exists := sess.Values["user"]
	if !exists || user == nil {
		return false
	}
	return true
}

func IsAuthenticated(c *gin.Context) bool {
	sess, err := GetAdminSession(c.Request)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return false
	}

	user, exists := sess.Values["user"]
	return exists && user != nil
}

func IsAdmin(c *gin.Context) bool {
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		return false
	}

	var user models.User
	user = user.Get(userID)
	return user.Role == models.RoleAdmin
}

func IsPostOwner(userID int, postID int) bool {
	var post models.Post
	post = post.Get(postID)
	return post.AuthorID == userID
}

func CanEditPost(c *gin.Context, postID int) bool {
	userID, exists := c.Get("userID")
	if !exists || userID == nil {
		return false
	}

	if IsAdmin(c) {
		return true
	}

	return false
}
