package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	"extract_load/models"
	"extract_load/utils"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

var (
	err         = godotenv.Load()
	dbUser      = os.Getenv("DB_USER")
	dbPassword  = os.Getenv("DB_PASSWORD")
	dbName      = os.Getenv("DB_NAME")
	dbHost      = os.Getenv("DB_HOST")
	dbPort      = os.Getenv("DB_PORT")
	insertTable = os.Getenv("INSERT_TABLE")

	cred_arr    = []string{dbUser, dbPassword, dbName, dbHost, dbPort}
	planetsList models.PlanetsList
	// insertTable = os.Getenv("INSERT_TABLE")
)

func main() {

	if err != nil {
		log.Fatal("Ошибка при загрузке файла .env: ", err)
	}

	var counts int = 10

	for _, value := range cred_arr {
		if value == "" {
			log.Fatal("Необходимо указать переменную окружения")
		}
	}

	fmt.Println("Запуск забора данных")

	for count := 1; count <= counts; count++ {
		url := fmt.Sprintf("https://rickandmortyapi.com/api/location?page=%s", strconv.Itoa(count))
		// Выполняем HTTP-запрос
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer resp.Body.Close()

		// Читаем ответ
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println(err)
			return
		}

		// Десериализуем JSON в структуру ApiResponse
		var apiResponse models.ApiResponse
		err = json.Unmarshal(body, &apiResponse)
		if err != nil {
			fmt.Println(err)
			return
		}

		planetsList = append(planetsList, apiResponse.Results...)

		if apiResponse.Info.Next == "" {
			fmt.Println("Curent page is", count)
			fmt.Println("No next page")
			break
		}

		// Если кол во страниц уже 10 но не сработал break
		// то увеличиваем кол во страниц на 10
		if count == counts {
			counts += 10
		}
	}

	// Обновляем количество жителей для каждого персонажа
	for i := range planetsList {
		planetsList[i].ResidentsCount = len(planetsList[i].Residents)
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		dbUser, dbPassword, dbName, dbHost, dbPort)

	db, err := utils.ConnectToDB(connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Выполняем транзакцию
	err = func() error {
		tx, err := db.Begin() // Начинаем транзакцию
		if err != nil {
			return err
		}

		// Обязательно откатываем транзакцию в случае ошибки
		defer func() {
			if err != nil {
				tx.Rollback()
			}
		}()

		truncateQuery := fmt.Sprintf("TRUNCATE %s", insertTable)
		fmt.Println(truncateQuery)
		_, err = tx.Exec(truncateQuery)

		if err != nil {
			return err
		}

		for _, value := range planetsList {
			// Выполняем вставку
			query := fmt.Sprintf("INSERT INTO %s (id, name, type, dimension, resident_cnt) VALUES ($1, $2, $3, $4, $5)", insertTable)
			_, err = tx.Exec(query, value.ID, value.Name, value.Type, value.Dimension, value.ResidentsCount)
			if err != nil {
				return err // Если ошибка, откатим транзакцию
			}
		}

		// Фиксируем транзакцию
		return tx.Commit()
	}()

	if err != nil {
		log.Fatalf("Ошибка при выполнении транзакции: %v", err)
	}

	fmt.Println("Данные успешно записаны в БД.")

}
