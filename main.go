package main

import (
	"cs2_telegram_bot/bot"
	"cs2_telegram_bot/database"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	// Получаем переменные из .env
	env_err := godotenv.Load()
	if env_err != nil {
		log.Fatal("❌ Не удалось загрузить .env файл")
	}
	fmt.Println("✅ Токен Telegram получен")

	// Получаем токен ТГ
	token := os.Getenv("TELEGRAM_TOKEN")
	if token == "" {
		log.Fatal("❌ Переменная TELEGRAM_TOKEN не найдена")
	}

	// Получение строки подключения к базе данных
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("❌ Переменная DATABASE_URL не найдена")
	}
	fmt.Println("✅ URL базы данных получен")

	fmt.Println("⚠️  Запуск бота...")

	// Подключаемся к базе данных
	connect_err := database.Connect()

	if connect_err != nil {
		log.Fatalf("❌ Ошибка подключения к базе: %v", connect_err)
	}
	fmt.Println("✅ Успешное подключение к базе данных")

	// Запуск бота
	bot.StartBot()
	fmt.Println("✅ Бот успешно запущен")

}

func StartStatusChecker() {

	go func() {
		for {
			// Тайм аут
			time.Sleep(1 * time.Minute)

			// Получаем всех пользователей
			users, err := database.GetAllUsers()
			if err != nil {
				log.Println("Ошибка при получении пользователей:", err)
				continue
			}

			for _, user := range users {
				go checkUserStatus(user.TelegramID, user.SteamID)
			}
		}
	}()
}

func checkUserStatus(telegramID int64, steamID string) {

}
