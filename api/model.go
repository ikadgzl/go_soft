package main

type Todo struct {
	ID          int    `json:"id"`
	Title       string `json:"title" validate:"required"`
	IsCompleted bool   `json:"is_completed" gorm:"default:false"`
}
