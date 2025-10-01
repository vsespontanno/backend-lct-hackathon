package entity

// Структура даваемых задач, после прочтения текста
type Quiz struct {
	ID            int64    `json:"id"`
	Title         string   `json:"title"`
	Content       string   `json:"content"`
	Options       []string `json:"options"`
	CorrectAnswer string   `json:"correctAnswer"`
}
