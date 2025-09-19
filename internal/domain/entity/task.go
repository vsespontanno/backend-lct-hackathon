package entity

// Структура даваемых задач, после прочтения текста
type Task struct {
	ID            int64
	Title         string
	Content       string
	Options       []string
	CorrectAnswer string
	Progress      int
}
