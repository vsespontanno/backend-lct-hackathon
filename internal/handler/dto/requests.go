package dto

type SetPetNameReq struct {
	Name   string `json:"name"`
	UserID int    `json:"userID"`
}

type NewSectionReq struct {
	Title string `json:"title"`
}

type NewSectionItemReq struct {
	SectionID int64  `json:"sectionId"`
	IsTest    bool   `json:"isTest"`
	Title     string `json:"title"`
	ItemID    int64  `json:"itemId"`
}

type GetTheoryReq struct {
	ID int64 `json:"id"`
}

type NewTheoryReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
