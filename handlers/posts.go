package handlers

import (
	"database/sql"
	"fmt"
	"log"

	"blog-backend/app"
	"blog-backend/models"
	"net/http"
	"strconv"

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

	// Hardcode UserId
	post.UserId = 1

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

func HandleFetchPosts(context *gin.Context) {
	rows, err := app.Db.Query("SELECT * FROM posts")

	if err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	defer rows.Close()

	var posts []models.Post

	for rows.Next() {

		var post models.Post

		if err := rows.Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &post.CreatedAt, &post.UpdatedAt); err != nil {

			log.Fatal(err)

			context.JSON(500, gin.H{
				"message": "Something went wrong",
			})
		}

		posts = append(posts, post)
	}

	if err = rows.Err(); err != nil {

		log.Fatal(err)

		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
	}

	if posts == nil {
		context.JSON(200, []models.Post{})
		return
	}

	context.JSON(200, posts)
}

func HandleUpdatePosts(context *gin.Context) {
	// Get the post ID from URL parameters
	id := context.Param("id")
	if id == "" {
		context.JSON(400, gin.H{
			"message": "Invalid ID",
		})
		return
	}

	var post models.Post

	// Bind JSON data to the post model
	if err := context.BindJSON(&post); err != nil {
		context.JSON(400, gin.H{
			"message": "Invalid JSON payload",
		})
		return
	}

	// Update query
	query := `
		UPDATE posts 
		SET title = $1, body = $2, user_id = $3, updated_at = $4
		WHERE id = $5
		RETURNING id, title, body, user_id, created_at, updated_at`

	// Execute the query and update the database
	err := app.Db.QueryRow(query, post.Title, post.Body, post.UserId, post.UpdatedAt, id).Scan(
		&post.Id,
		&post.Title,
		&post.Body,
		&post.UserId,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		log.Println("Database error:", err)
		context.JSON(500, gin.H{
			"message": "Could not update post",
		})
		return
	}

	// Respond with the updated post details
	context.JSON(200, gin.H{
		"message": "Post updated successfully",
		"post":    post,
	})
}

func HandleDeletePosts(context *gin.Context) {

	query := `
      DELETE FROM posts WHERE id=$1;`

	_, err := app.Db.Query(query, context.Param("id"))

	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	context.Status(204)
}

func HandleFetchPost(context *gin.Context) {
	// Get the post ID from the URL parameters
	postIDStr := context.Param("id")

	// Convert postID to an integer
	postID, err := strconv.Atoi(postIDStr)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid post ID format",
		})
		return
	}

	// Prepare the SQL query
	query := "SELECT id, title, body, user_id, created_at, updated_at FROM posts WHERE id = $1"

	// Query the database for the post
	var post models.Post

	err = app.Db.QueryRow(query, postID).Scan(&post.Id, &post.Title, &post.Body, &post.UserId, &post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			// If no post is found, return a 404 Not Found response
			context.JSON(http.StatusNotFound, gin.H{
				"message": "Post not found",
			})
		} else {
			// Log the error and return a 500 Internal Server Error response
			log.Println("Error fetching post:", err)
			context.JSON(http.StatusInternalServerError, gin.H{
				"message": "Failed to fetch post",
			})
		}
		return
	}

	// Return the product as JSON
	context.JSON(http.StatusOK, post)
}
