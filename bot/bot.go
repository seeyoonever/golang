package bot

import (
	"cs2_telegram_bot/database"
	"log"
	"os"
	"time"

	teleg "gopkg.in/telebot.v3"
)

var Bot *teleg.Bot

// Перемененная для отслеживания состояния пользователя
var awaitingSteamId = make(map[int64]bool)
var awaitingInfo = make(map[int64]string)

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

	// Обработчики команд
	Bot.Handle("/start", handleStart)
	Bot.Handle("/register", handleRegister)

	// Обработчик текстовых сообщений
	Bot.Handle(teleg.OnText, handleTextMessage)

	Bot.Start()

}

// Обработчик команды start
// teleg.Context - это объект который содержит всю информацию о сообщении пользователя
func handleStart(context teleg.Context) error {
	return context.Send("Ола Амиго! 👋 Я бот для отслеживания игроков CS2. Используй /register чтобы начать.")
}

// Функция проверки валидности SteamId
func isValidSteamId(input string) bool {
	if len(input) != 17 {
		return false
	}
	for _, ch := range input {
		if ch < 0 || ch > 9 {
			return false
		}
	}
	return true
}

func handleRegister(context teleg.Context) error {
	userID := context.Sender().ID
	userName := context.Sender().FirstName

	// Проверка на существование пользователя в базе
	found, err := database.UserExists(userID)
	if err != nil {
		log.Fatal(err)
		return context.Send("Произошла ошибка при проверке пользователя. Negativ!")
	}

	if found {
		return context.Send("Вы уже зарегистрированы. ✅ OK Let's Go")
	}

	// Записываем id пользователя в мапу, которая говорит, что от него мы ждём steamID
	awaitingSteamId[userID] = true
	return context.Send(userName, "отправь мне свой SteamID для регистрации")

}

func handleTextMessage(context teleg.Context) error {
	userID := context.Sender().ID
	userName := context.Sender().FirstName
	text := context.Text()

	// Проверка пользователя на ожидание SteamID
	if awaitingSteamId[userID] {
		// Валидация
		if !isValidSteamId(text) {
			return context.Send(userName, "Йоу, бро! Это не похоже на настоящий SteamID. Попробуй снова — только цифры, 17 символов.")
		}

		// Если steamID валиден мы просим пользователя ввести Info
		awaitingInfo[userID] = text
		delete(awaitingSteamId, userID)
		return context.Send("Отлично! Теперь введи дополнительную информацию (имя/ник)")
	}

	// Проверка пользователя на ожидание
	if steamID, ok := awaitingInfo[userID]; ok {
		info := text

		if len(info) > 20 {
			return context.Send(userName, "Слишком много текста. Дай краткую инфу о себе (имя/ник)")
		}

		err := database.AddUser(userID, steamID, info)
		if err != nil {
			return context.Send("Ошибка при добавлении пользователя. Negativ!")
		}

		// Удаляем пользователя из мапы ожидания SteamID
		delete(awaitingInfo, userID)

		return context.Send(userName, "Ты успешно зарегистрирован! Easy peasy lemon squeezy! 🎉")
	}

	// Если пользователь ничего не должен вводить
	return nil

}
