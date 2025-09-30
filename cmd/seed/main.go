package main

import (
	"black-pearl/backend-hackathon/internal/config"
	"black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"encoding/json"
	"github.com/lib/pq"
	"io/ioutil"
	"log"
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
	data, err := ioutil.ReadFile("seed/tasks.json")
	if err != nil {
		log.Fatalf("failed to read seed/tasks.json: %v", err)
		return
	}

	var tasks []entity.Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		log.Fatalf("failed to parse seed/tasks.json: %v", err)
		return
	}

	for _, task := range tasks {
		var exists bool
		err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id=$1)", task.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("failed to query task: %v", err)
		}
		if exists {
			continue
		}
		_, err := db.Exec(
			"INSERT INTO tasks (id, title, content, options, correct_answer, progress) VALUES ($1, $2, $3, $4, $5, $6)",
			task.ID, task.Title, task.Content, pq.Array(task.Options), task.CorrectAnswer, task.Progress)
		if err != nil {
			log.Fatalf("failed to insert task: %v", err)
		} else {
			log.Printf("inserted task: %v", task.ID)
		}
	}
}
