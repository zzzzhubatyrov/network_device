package main

import (
	"log"
	"network/internal/handlers/v1"
	"network/internal/models"
	"network/internal/repository"
	"network/internal/service"
	storage "network/internal/storage/sqlite"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	webview "github.com/webview/webview_go"
)

func main() {
	// Инициализация базы данных
	db, err := storage.NewSQLiteStorage("test.db")
	if err != nil {
		log.Fatal(err)
	}

	migrator := db.Migrator()

	// if err := migrator.DropTable(
	// 	&models.Router{},
	// 	&models.Port{},
	// 	&models.RouterConnection{},
	// ); err != nil {
	// 	log.Fatal(err)
	// }

	if err := migrator.AutoMigrate(
		&models.Router{},
		&models.Port{},
		&models.RouterConnection{},
	); err != nil {
		log.Fatal(err)
	}

	// Настройка API сервера
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New(
		cors.Config{
			AllowOrigins: "*",
			AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		},
	))

	// Инициализация хендлеров
	handler := handlers.NewHandler(
		service.NewService(
			repository.NewRepository(db),
		),
	)
	handler.InitRoute(app)

	app.Static("/", "./web/dist")

	// Запуск API сервера в горутине
	go func() {
		if err := app.Listen(":5500"); err != nil {
			log.Fatal(err)
		}
	}()

	// Даем серверу время на запуск
	time.Sleep(100 * time.Millisecond)

	// Создание и настройка webview
	w := webview.New(true)
	defer w.Destroy()

	w.SetTitle("Network Management")
	w.SetSize(1200, 800, webview.HintNone)
	// w.SetSize(800, 600, webview.HintFixed)

	// Добавляем обработчик для связи с Go
	// w.Bind("getRouters", func() string {
	// return `{"status": "ok"}`
	// })

	// Загружаем HTML напрямую
	w.Navigate("http://localhost:5500")

	// Запуск webview
	w.Run()
}
