package dto

type SetPetNameReq struct {
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}
