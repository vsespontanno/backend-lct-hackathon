package dto

type SetPetNameReq struct {
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

type NewSectionReq struct {
	Title string `json:"title"`
}

type NewSectionItemReq struct {
	SectionID int    `json:"sectionId"`
	IsTest    bool   `json:"isTest"`
	Title     string `json:"title"`
	ItemID    int    `json:"itemId"`
}

type GetTheoryReq struct {
	ID int `json:"id"`
}

type NewTheoryReq struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

type SendXPReq struct {
	UserID int `json:"userID"`
	Exp    int `json:"exp"`
}
