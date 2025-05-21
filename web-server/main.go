package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func handleGetStudents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	conn, err := connectToDb()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	students := queryData(conn)
	err = json.NewEncoder(w).Encode(students)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func connectToDb() (*pgx.Conn, error) {
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	connectionString := fmt.Sprintf("postgres://%s:%s@postgres:5432/%s", dbUser, dbPassword, dbName)
	conn, err := pgx.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func queryData(conn *pgx.Conn) []Student {
	rows, err := conn.Query(context.Background(), "SELECT id, name FROM students")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var students []Student = make([]Student, 0)
	for rows.Next() {
		var id int
		var name string
		err := rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}
		students = append(students, Student{id, name})
	}
	return students
}

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatalf("Specify PORT env: %s\n", err)
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /students", handleGetStudents)

	mux.Handle("/metrics", promhttp.Handler())
	
	log.Printf("Server started at port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), mux))
}
