package bot

import (
	"cs2_telegram_bot/database"
	"log"
	"os"
	"time"

	teleg "gopkg.in/telebot.v3"
)

var Bot *teleg.Bot

// Инициализация и запуск бота
func StartBot() {
	settings := teleg.Settings{ // Структура конфигурации бота
		Token:  os.Getenv("TELEGRAM_TOKEN"),
		Poller: &teleg.LongPoller{Timeout: 10 * time.Second}, // Обязательное поле Поллера, я выбрал лонг поллинг с таймаутом в 10 сек
	}

	var err error
	Bot, err = teleg.NewBot(settings)
	if err != nil {
		log.Fatal(err)
	}

	Bot.Handle("/start", handleStart)
	Bot.Handle("/register", handleRegister)

	Bot.Start()

}

// Обработчик команды start
// teleg.Context - это объект который содержит всю информацию о сообщении пользователя
func handleStart(context teleg.Context) error {
	return context.Send("Ола Амиго! 👋 Я бот для отслеживания игроков CS2. Используй /register чтобы начать.")
}

func handleRegister(context teleg.Context) error {
	userID := context.Sender().ID

	// Проверка на существование пользователя в базе
	found, err := database.UserExists(userID)
	if err != nil {
		log.Fatal(err)
		return context.Send("Произошла ошибка при проверке пользователя. Negativ!")
	}

	if found {
		return context.Send("Вы уже зарегистрированы. ✅ OK Let's Go")
	}

	err = database.AddUser(userID, "1234567", "Test Info")
	if err != nil {
		return context.Send("Ошибка при добавлении пользователя. Negativ!")
	}

	return context.Send("Вы успешно зарегестрированы в системе! Easy peasy lemon squeezy!")

}
