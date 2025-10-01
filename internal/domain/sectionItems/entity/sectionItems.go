package entity

type SectionItem struct {
	SectionID int64  `json:"sectionID"`
	IsTest    bool   `json:"isTest"`
	Title     string `json:"title"`
	ItemID    int64  `json:"itemID"`
}
