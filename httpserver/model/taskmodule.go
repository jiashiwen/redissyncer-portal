package model

type TaskCreateBody struct {
	Name     string `json:"name"`
	Password int64  `json:"password"`
}
