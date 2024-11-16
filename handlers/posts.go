package handlers

import (
	"fmt"

	"blog-backend/app"
	"blog-backend/models"

	"github.com/gin-gonic/gin"
)

	func HandleAddPosts(context *gin.Context) {
		var post models.Post

		err := context.BindJSON(&post)

		if err != nil {
			fmt.Println(err)
			context.JSON(400, gin.H{
				"message": "Bad request",
			})
			return
		}

		query := `
			INSERT INTO posts (title, body, user_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, title, body, user_id, created_at, updated_at`

		err = app.Db.QueryRow(query, post.Title, post.Body, post.UserId, post.CreatedAt, post.UpdatedAt).Scan(
			&post.Id,
			&post.Title,
			&post.Body,
			&post.UserId,
			&post.CreatedAt,
			&post.UpdatedAt,
		)

		if err != nil {
			fmt.Println(err)
			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
			return
		}

		context.JSON(201, post)
	}
