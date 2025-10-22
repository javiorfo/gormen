package response

import "hex-arch-fiber/application/model"

type UserResponse struct {
	User model.User `json:"user"`
}

type UsersResponse struct {
	PageInfo PageInfo     `json:"pageInfo"`
	Users    []model.User `json:"user"`
}

type PageInfo struct {
	Number string `json:"number"`
	Size   string `json:"size"`
	Total  int64  `json:"total"`
}
