package models

type Post struct {
	Id        int    `json:"id"`        
	Title     string `json:"title"`    
	Body      string `json:"body"`     
	UserId    int    `json:"user_id"`   
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

