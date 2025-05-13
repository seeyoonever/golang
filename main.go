package main

import (
	"cs2_telegram_bot/database"
	"fmt"
	"log"
)

func main() {
	fmt.Println("Запуск бота...")
	err := database.Connect()

	if err != nil {
		log.Fatal("Ошибка подключения:", err)
	}
}
