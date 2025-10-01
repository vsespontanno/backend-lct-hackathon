package dto

import "black-pearl/backend-hackathon/internal/domain/prize/entity"
import sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"

type GetPetResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Exp   int    `json:"exp"`
}

type GetPrizesResp struct {
	Prizes []entity.Prize `json:"prizes"`
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
