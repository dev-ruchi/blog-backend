package handlers

import (
	"fmt"

	"blog-backend/app"
	"blog-backend/models"

	"github.com/gin-gonic/gin"
)

func HandleAddUsers(context *gin.Context) {
	var user models.User

	// Bind the incoming JSON data to the user struct
	err := context.BindJSON(&user)
	if err != nil {
		fmt.Println(err)
		context.JSON(400, gin.H{
			"message": "Bad request",
		})
		return
	}

	// Query to insert the new user into the users table
	query := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`

	// Execute the query and scan the resulting ID into the user struct
	err = app.Db.QueryRow(query, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt).Scan(
		&user.Id,
	)

	// Check if any error occurred during the insertion
	if err != nil {
		fmt.Println(err)
		context.JSON(500, gin.H{
			"message": "Something went wrong",
		})
		return
	}

	// Return a successful response with the created user
	context.JSON(201, user)
}
