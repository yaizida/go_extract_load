package utils

import (
	"fmt"
	"sort"

	"extract_load/models"

	"github.com/jmoiron/sqlx"
)

func ConnectToDB(connStr string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("ошибка подключения к БД: %w", err)
	}
	return db, nil
}
func ReturnTopThree(planetsList models.PlanetsList) models.PlanetsList {
	// Копируем срез, чтобы не изменять оригинальный
	copiedPlanets := make(models.PlanetsList, len(planetsList))
	copy(copiedPlanets, planetsList)

	// Сортируем планеты по количеству жителей в порядке убывания
	sort.SliceStable(copiedPlanets, func(i, j int) bool {
		return copiedPlanets[i].ResidentsCount > copiedPlanets[j].ResidentsCount
	})

	// Возвращаем первые три планеты
	return copiedPlanets[:3]
}
