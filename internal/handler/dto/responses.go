package dto

import sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"

type GetPetReq struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Exp   int    `json:"exp"`
	Level int    `json:"lvl"`
}

type SectionWithItemsResp struct {
	ID    int64                            `json:"id"`
	Title string                           `json:"title"`
	Items []sectionItemsEntity.SectionItem `json:"items"`
}

type TaskResp struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}

type TheoryResp struct {
	ID      int64  `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}
