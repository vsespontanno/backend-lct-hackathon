package dto

type GetPetReq struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Exp   int    `json:"exp"`
	Level int    `json:"lvl"`
}
