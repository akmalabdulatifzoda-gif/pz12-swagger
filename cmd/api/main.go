// Package main Notes API server.
//
//	@title						Notes API
//	@version					1.0
//	@description				Учебный REST API для заметок (CRUD). ПЗ №12 — Swagger/OpenAPI.
//	@description				Методы POST, PATCH и DELETE требуют Bearer-токен, если сервер запущен с API_TOKEN.
//	@contact.name				Абдулатифзода Акмал, ЭФМО-0126
//	@contact.url				https://github.com/akmalabdulatifzoda-gif/pz12-swagger
//	@license.name				MIT
//	@license.url				https://opensource.org/licenses/MIT
//	@BasePath					/api/v1
//	@schemes					http
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Введите токен в формате: Bearer ВАШ_ТОКЕН (для демо: Bearer secret)
package main

import (
	"log"
	"net/http"
	"os"

	_ "github.com/akmalabdulatifzoda-gif/pz12-swagger/docs" // сгенерированная спецификация (swag init)
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/core/service"
	httpx "github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/http/handlers"
	"github.com/akmalabdulatifzoda-gif/pz12-swagger/internal/repo"
)

func main() {
	// Сборка слоёв: repository → service → handler → router
	noteRepo := repo.NewNoteRepoMem()
	svc := service.NewNoteService(noteRepo)
	h := &handlers.Handler{Svc: svc}

	token := os.Getenv("API_TOKEN")
	r := httpx.NewRouter(h, token)

	addr := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		addr = ":" + p
	}
	if token != "" {
		log.Println("Bearer auth enabled for POST/PATCH/DELETE")
	}
	log.Println("Server started at", addr)
	log.Println("Swagger UI: http://localhost" + addr + "/docs/index.html")
	log.Println("ReDoc:      http://localhost" + addr + "/redoc")
	log.Fatal(http.ListenAndServe(addr, r))
}
