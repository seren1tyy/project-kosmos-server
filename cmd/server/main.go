package main

import (
	"fmt"
	"log"
	"net"
	"project-kosmos-server/internal/auth"
	"project-kosmos-server/internal/character"
	"project-kosmos-server/internal/config"
	"project-kosmos-server/internal/db"
	"project-kosmos-server/internal/items"
	"project-kosmos-server/internal/network"
	"project-kosmos-server/internal/session"
	"project-kosmos-server/internal/universe"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Загрузка конфига
	cfg := config.Load("config.json")
	log.Printf("Config loaded (port: %d)", cfg.Server.Port)

	// Подключение к БД
	database := db.NewMariaDB()
	if err := database.Connect(cfg.Database.DSN); err != nil {
		log.Fatalf("DB connect failed: %v", err)
	}
	defer database.Close()
	log.Printf("MariaDB pool initialized")

	// Получаем *sql.DB из интерфейса
	dbConn := database.(*db.MariaDB).GetDB()

	// Создаём сервисы
	itemSvc := items.NewService(dbConn)
	universeSvc := universe.NewService(dbConn)
	charSvc := character.NewService(dbConn, []byte(cfg.Server.CryptoKey), itemSvc, universeSvc)

	authSvc := &auth.Service{
		DB:         database,
		Key:        []byte(cfg.Server.CryptoKey),
		SessionMgr: session.NewManager(),
	}

	// Запуск TCP-сервера
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", addr, err)
	}
	defer ln.Close()
	log.Printf("Server listening on %s", addr)

	// 🔹 Основной цикл (без сложного graceful shutdown для дев-среды)
	// На Windows ln.Close() + Ctrl+C корректно убивает процесс через ОС
	for {
		conn, err := ln.Accept()
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && opErr.Err.Error() == "use of closed network connection" {
				return
			}
			continue
		}
		// Передаём оба сервиса
		go network.HandleClient(conn, authSvc, charSvc)
	}
}
