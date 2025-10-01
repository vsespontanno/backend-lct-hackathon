package main

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/domain/quiz/entity"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/lib/pq"
)

func main() {
	cfg, err := config.ReadConfig()
	if err != nil {
		log.Printf("failed to connect to config: %v", err)
		return
	}
	//	log.Printf("DB_HOST=%s DB_PORT=%s DB_USER=%s DB_NAME=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName)
	db, err := db.ConnectToPostgres(cfg.DB.Username, cfg.DB.Password, cfg.DB.DBName, cfg.DB.Host, cfg.DB.Port, cfg.DB.SSLMode)
	if err != nil {
		log.Printf("failed to connect to database: %v", err)
		return
	}
	data, err := ioutil.ReadFile("seed/quizs.json")
	if err != nil {
		log.Fatalf("failed to read seed/quizs.json: %v", err)
		return
	}

	var quizs []entity.Quiz
	err = json.Unmarshal(data, &quizs)
	if err != nil {
		log.Fatalf("failed to parse seed/quizs.json: %v", err)
		return
	}

	for _, quiz := range quizs {
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM quizs WHERE id=$1)", quiz.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to query quiz: %v", err)
		}
		if exists {
			continue
		}
		_, err := db.Exec(
			"INSERT INTO quizs (id, title, content, options, correct_answer) VALUES ($1, $2, $3, $4, $5)",
			quiz.ID, quiz.Title, quiz.Content, pq.Array(quiz.Options), quiz.CorrectAnswer)
		if err != nil {
			log.Fatalf("failed to insert quiz: %v", err)
		} else {
			log.Printf("inserted quiz: %v", quiz.ID)
		}
	}
}
