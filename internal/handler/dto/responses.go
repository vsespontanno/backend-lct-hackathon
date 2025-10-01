package dto

import "black-pearl/backend-hackathon/internal/domain/prize/entity"

type GetPetResp struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Exp   int    `json:"exp"`
	Level int    `json:"lvl"`
}

type GetPrizesResp struct {
	Prizes []entity.Prize `json:"prizes"`
}
