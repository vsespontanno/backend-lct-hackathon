package entity

type SectionItem struct {
	SectionID int    `json:"sectionID"`
	IsTest    bool   `json:"isTest"`
	Title     string `json:"title"`
	ItemID    int    `json:"itemID"`
}
