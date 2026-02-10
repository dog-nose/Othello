package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dog-nose/othello-backend/config"
	"github.com/dog-nose/othello-backend/handler"
	"github.com/dog-nose/othello-backend/middleware"
	"github.com/dog-nose/othello-backend/repository"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("connected to database")

	repo := repository.NewMySQLRepository(db)
	h := handler.New(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/start-game", h.StartGame)
	mux.HandleFunc("/place-stone", h.PlaceStone)
	mux.HandleFunc("/end-game", h.EndGame)

	server := middleware.CORS(mux)

	addr := ":8080"
	fmt.Printf("server listening on %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, server))
}
