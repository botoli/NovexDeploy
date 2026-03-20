package main

import (
	"localVercel/app"
	"localVercel/db"
	"localVercel/middleware"
	"log"
	"net/http"
	"os"

	_ "localVercel/docs"

	"github.com/joho/godotenv"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Подключаемся к базе данных
	if err := db.InitDB(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}

	application := app.New()

	mux := http.NewServeMux()
	application.RegisterRoutes(mux)

	// Swagger UI
	mux.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	// Статические файлы для превью билдов
	mux.Handle("/preview/", http.StripPrefix("/preview/", http.FileServer(http.Dir("./builds"))))

	// Оборачиваем в middleware: сначала CORS, потом логирование
	handler := middleware.CorsMiddleware(middleware.WithLogging(mux))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8888"
	}
	addr := ":" + port

	log.Printf("🚀 Server starting on http://localhost%s", addr)
	log.Printf("📚 Swagger UI: http://localhost%s/swagger/", addr)
	log.Printf("🔗 GitHub OAuth: http://localhost%s/auth/github/login", addr)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		log.Fatalf("❌ Server failed: %v", err)
	}
}