package dto

type SetPetNameReq struct {
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

type SendXPReq struct {
	UserID int `json:"userID"`
	Exp    int `json:"exp"`
}
