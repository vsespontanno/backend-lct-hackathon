package entity

// Структура прогресса медевежонка юзера
type Progress struct {
	UserID   int64 `json:"id"`
	Points   int   `json:"points"`
	Progress int   `json:"progress"`
}
