package models

type PlanetsList []Planets

// Определяем структуру для хранения данных персонажа
type Planets struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Dimension      string   `json:"dimension"`
	Residents      []string `json:"residents"`       // Массив жителей
	ResidentsCount int      `json:"residents_count"` // Поле для хранения количества жителей
	URL            string   `json:"url"`
	Created        string   `json:"created"`
}

// Определяем структуру для общего ответа API
type ApiResponse struct {
	Info    Info      `json:"info"`
	Results []Planets `json:"results"`
}

// Структура для информации о страницах
type Info struct {
	Count int     `json:"count"`
	Pages int     `json:"pages"`
	Next  string  `json:"next"`
	Prev  *string `json:"prev"` // Используем указатель, чтобы учитывать null
}
