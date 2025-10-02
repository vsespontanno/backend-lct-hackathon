package main

import (
	cfg "black-pearl/backend-hackathon/internal/config"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	quizEntity "black-pearl/backend-hackathon/internal/domain/quiz/entity"
	sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	sectionsEntity "black-pearl/backend-hackathon/internal/domain/sections/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	theoryEntity "black-pearl/backend-hackathon/internal/domain/theory/entity"
	"black-pearl/backend-hackathon/internal/infrastructure/db"
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/lib/pq"
)

func main() {
	config, err := cfg.ReadConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	database, err := db.ConnectToPostgres(config.DB.Username, config.DB.Password, config.DB.DBName, config.DB.Host, config.DB.Port, config.DB.SSLMode)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Запускаем сиды (закомментируй те, для которых нет JSON файлов)
	seedQuizzes(database)
	seedTasks(database)
	seedSections(database)
	seedSectionItems(database)
	seedTheory(database)
	seedPrizes(database)

	log.Println("All seeding completed!")
}

// ------------------------- Seed Quizzes -------------------------
func seedQuizzes(db *sql.DB) {
	var quizzes []quizEntity.Quiz
	loadJSON("seed/quizs.json", &quizzes)

	for _, q := range quizzes {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM quiz WHERE id=$1)", q.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("query quiz exists failed: %v", err)
		}
		if exists {
			continue
		}
		_, err = db.Exec("INSERT INTO quiz (id, title, content, options, correct_answer) VALUES ($1,$2,$3,$4,$5)",
			q.ID, q.Title, q.Content, pq.Array(q.Options), q.CorrectAnswer)
		if err != nil {
			log.Fatalf("insert quiz failed: %v", err)
		}
	}
	log.Println("quizzes seeded successfully")
}

// ------------------------- Seed Tasks -------------------------
func seedTasks(db *sql.DB) {
	var tasks []taskEntity.Task
	loadJSON("seed/tasks.json", &tasks)

	for _, t := range tasks {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM tasks WHERE id=$1)", t.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("query task exists failed: %v", err)
		}
		if exists {
			continue
		}
		_, err = db.Exec("INSERT INTO tasks (id, title) VALUES ($1,$2)", t.ID, t.Title)
		if err != nil {
			log.Fatalf("insert task failed: %v", err)
		}
	}
	log.Println("tasks seeded successfully")
}

// ------------------------- Seed Sections -------------------------
func seedSections(db *sql.DB) {
	var sections []sectionsEntity.Sections
	loadJSON("seed/sections.json", &sections)

	for _, s := range sections {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sections WHERE id=$1)", s.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("query section exists failed: %v", err)
		}
		if exists {
			continue
		}
		_, err = db.Exec("INSERT INTO sections (id, title) VALUES ($1,$2)", s.ID, s.Title)
		if err != nil {
			log.Fatalf("insert section failed: %v", err)
		}
	}
	log.Println("sections seeded successfully")
}

// ------------------------- Seed Section Items -------------------------
func seedSectionItems(db *sql.DB) {
	var items []sectionItemsEntity.SectionItem
	loadJSON("seed/section_items.json", &items)

	for _, i := range items {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM sectionitems WHERE sectionID=$1)", i.SectionID).Scan(&exists)
		if err != nil {
			log.Fatalf("query section_item exists failed: %v", err)
		}
		if exists {
			continue
		}

		_, err = db.Exec(
			"INSERT INTO sectionitems (sectionID, isTest, title, itemID) VALUES ($1,$2,$3,$4)",
			i.SectionID, i.IsTest, i.Title, i.ItemID,
		)
		if err != nil {
			log.Fatalf("insert section_item failed: %v", err)
		}
	}

	log.Println("section_items seeded successfully")
}

// ------------------------- Seed Theory -------------------------
func seedTheory(db *sql.DB) {
	var theories []theoryEntity.Theory
	loadJSON("seed/theory.json", &theories)

	for _, t := range theories {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM theory WHERE id=$1)", t.ID).Scan(&exists)
		if err != nil {
			log.Fatalf("query theory exists failed: %v", err)
		}
		if exists {
			continue
		}
		_, err = db.Exec("INSERT INTO theory (id, title, content) VALUES ($1,$2,$3)", t.ID, t.Title, t.Content)
		if err != nil {
			log.Fatalf("insert theory failed: %v", err)
		}
	}
	log.Println("theory seeded successfully")
}

// ------------------------- Seed Prizes -------------------------
func seedPrizes(db *sql.DB) {
	var prizes []prizeEntity.Prize
	loadJSON("seed/prizes.json", &prizes)

	for _, p := range prizes {
		var exists bool
		err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM prizes WHERE title=$1)", p.Title).Scan(&exists)
		if err != nil {
			log.Fatalf("query prize exists failed: %v", err)
		}
		if exists {
			continue
		}
		_, err = db.Exec(
			"INSERT INTO prizes (title, descr, type) VALUES ($1,$2,$3)",
			p.Title, p.Description, p.Type,
		)
		if err != nil {
			log.Fatalf("insert prize failed: %v", err)
		}
	}
	log.Println("prizes seeded successfully")
}

// ------------------------- Load JSON helper -------------------------
func loadJSON(path string, target interface{}) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("failed to read %s: %v", path, err)
	}
	if err := json.Unmarshal(data, target); err != nil {
		log.Fatalf("failed to parse %s: %v", path, err)
	}
}
