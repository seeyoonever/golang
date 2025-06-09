package bot

import (
	"cs2_telegram_bot/database"
	"cs2_telegram_bot/steam"
	"cs2_telegram_bot/utils"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	teleg "gopkg.in/telebot.v3"
)

var Bot *teleg.Bot

// Перемененная для отслеживания состояния пользователя
var awaitingSteamId = make(map[int64]bool)
var awaitingInfo = make(map[int64]string)
var awaitingDelete = make(map[int64]bool)

// var adminID = os.Getenv("MY_TELEGRAM_ID")

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

	StartStatusChecker(Bot)

	// Обработчик текстовых сообщений
	Bot.Handle(teleg.OnText, handleTextMessage)

	// Обработчики команд
	Bot.Handle("/start", handleStart)
	Bot.Handle("/register", handleRegister)

	// Обработчик команды администрирования - получение списка всех пользователей
	Bot.Handle("/admin", handleAdmin)

	// Обработчик команды удаления пользвателей (только для админа)
	Bot.Handle("/delete", handleDelete)

	// Обработчик получения всех игроков в cs2
	Bot.Handle("/status", handleStatus2)

	Bot.Start()

	log.Println("Бот запущен")
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
		if ch < '0' || ch > '9' {
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
		log.Println("Ошибка при проверке пользователя:", err)
		return context.Send("Произошла ошибка при проверке пользователя. Negativ!")
	}

	if found {
		return context.Send("Вы уже зарегистрированы. ✅ OK Let's Go")
	}

	// Записываем id пользователя в мапу, которая говорит, что от него мы ждём steamID
	utils.SetWithTTL(awaitingSteamId, userID, true, 5*time.Minute)

	return context.Send(userName + " отправь мне свой SteamID для регистрации")

}

func handleTextMessage(context teleg.Context) error {
	userID := context.Sender().ID
	userName := context.Sender().FirstName
	chatID := context.Chat().ID
	text := context.Text()

	log.Println("Получено сообщение:", context.Text(), "от", context.Sender().ID)

	// Проверка пользователя на ожидание SteamID
	if awaitingSteamId[userID] {
		// Валидация
		if !isValidSteamId(text) {
			return nil
			// return context.Send(userName + " Йоу, бро! Это не похоже на настоящий SteamID. Попробуй снова — только цифры, 17 символов.")
		}

		// Если steamID валиден мы просим пользователя ввести Info
		utils.SetWithTTL(awaitingInfo, userID, text, 5*time.Minute)

		delete(awaitingSteamId, userID)
		return context.Send("Отлично! Теперь введи дополнительную информацию (имя/ник)\nВнимание следующее сообщение будет принято как информация о пользователе")
	}

	// Проверка пользователя на ожидание
	if steamID, ok := awaitingInfo[userID]; ok {
		info := text

		if len(info) > 40 {
			return context.Send(userName + " Слишком много текста. Дай краткую инфу о себе (имя/ник)")
		}

		err := database.AddUser(userID, steamID, info, chatID)

		// Удаляем пользователя из мапы ожидания SteamID
		delete(awaitingInfo, userID)

		if err != nil {
			log.Println("Ошибка при проверке пользователя:", err)
			return context.Send("Ошибка при добавлении пользователя. Negativ!\nСкорее всего такой SteamID уже зарегистрирован.\nНачни заново с /register")
		}

		return context.Send(userName + " Ты успешно зарегистрирован! Easy peasy lemon squeezy! 🎉")
	}

	if awaitingDelete[userID] {
		delete(awaitingDelete, userID)

		tgID, err := strconv.ParseInt(text, 10, 64)
		if err != nil {
			fmt.Println("Error (обработка строки telegramID в число): ", err)
			return context.Send("❌ Некорректный Telegram ID. Попробуйте снова.")
		}

		err = database.DeleteUser(tgID)
		if err != nil {
			log.Println("Ошибка при удалении:", err)
			return context.Send("❌ Ошибка при удалении пользователя")
		}

		return context.Send("✅ Пользователь удалён")

	}

	// Если пользователь ничего не должен вводить
	return nil

}

func handleAdmin(context teleg.Context) error {
	userID := context.Sender().ID

	context.Send(userID)

	log.Println(userID)

	// admin, err := strconv.ParseInt(adminID, 10, 64)
	// if err != nil {
	// 	fmt.Println("Error: ", err)
	// }

	// Проверка на админа
	if userID != 330686271 {
		return context.Send("DENIED ❌ У вас нет доступа к этой команде.")
	}

	users, err := database.GetAllUsers()
	if err != nil {
		log.Println("Ошибка при получении пользователей: ", err)
	}

	if len(users) == 0 {
		return context.Send("Список пользователей пуст")
	}

	message := "📋 Список пользователей:\n\n"
	for _, user := range users {
		message += fmt.Sprintf("🧍 TelegramID: %d\n🎮 SteamID: %s\n📝 Info: %s\n\n", user.TelegramID, user.SteamID, user.Info)
	}

	return context.Send(message)

}

// Функция для удаления пользователей (делаю исключительно для админов, соответсвенно в ней нет защиты от дурака и ожидания ответа)
func handleDelete(context teleg.Context) error {
	userID := context.Sender().ID

	context.Send(userID)

	log.Println(userID)

	// admin, err := strconv.ParseInt(adminID, 10, 64)
	// if err != nil {
	// 	fmt.Println("Error (обработка строки telegramID в число): ", err)
	// }

	// Проверка на админа
	if userID != 330686271 {
		return context.Send("DENIED ❌ У вас нет доступа к этой команде.")
	}

	utils.SetWithTTL(awaitingDelete, userID, true, 5*time.Minute) // Включаем режим ожидания

	return context.Send("Введите telegramID кого хотите удалить")

}

func handleStatus(context teleg.Context) error {
	userID := context.Sender().ID

	// Проверка на существование пользователя в базе
	found, err := database.UserExists(userID)
	if err != nil {
		log.Fatal(err)
		return context.Send("Произошла ошибка при проверке пользователя. Negativ!")
	}

	if !found {
		return context.Send("🙅 Вы не зарегистрированы. Используйте /register")
	}

	// Получаем всех пользователей
	users, err := database.GetAllUsers()
	if err != nil {
		log.Println("Ошибка при получении пользователей: ", err)
		return context.Send("❌ Ошибка при получении списка пользователей")
	}

	var onlineList []string

	for _, user := range users {
		playing, persona, err := steam.IsPlayingCS2(user.SteamID)
		if err != nil {
			log.Println("Ошибка при получении статуса пользователя в CS2", err)
			continue // пропускаем, если что-то пошло не так
		}
		if playing {
			info := fmt.Sprintf("- %s (%s)", user.Info, persona)
			onlineList = append(onlineList, info)
		}
	}

	if len(onlineList) == 0 {
		return context.Send("😴 Никто из зарегистрированных пользователей сейчас не играет в CS2")
	}

	response := "🎮 Сейчас играют в CS2:\n" + strings.Join(onlineList, "\n")

	return context.Send(response)
}

// Функция получения статусов игроков через БАТЧ v1.1
func handleStatus2(context teleg.Context) error {
	userID := context.Sender().ID

	// Проверка на существование пользователя в базе
	found, err := database.UserExists(userID)
	if err != nil {
		log.Fatal(err)
		return context.Send("Произошла ошибка при проверке пользователя. Negativ!")
	}

	if !found {
		return context.Send("🙅 Вы не зарегистрированы. Используйте /register")
	}

	// Получаем всех пользователей
	users, err := database.GetAllUsers()
	if err != nil {
		log.Println("Ошибка при получении пользователей: ", err)
		return context.Send("❌ Ошибка при получении списка пользователей")
	}

	// Готовим пакет
	const batchSize = 100
	var responses []string

	for i := 0; i < len(users); i += batchSize {
		end := i + batchSize
		if end > len(users) {
			end = len(users)
		}

		batch := users[i:end]
		var steamIDs []string
		steamIDtoInfo := make(map[string]string)

		for _, user := range batch {
			steamIDs = append(steamIDs, user.SteamID)
			steamIDtoInfo[user.SteamID] = user.Info
		}

		players, err := steam.GetPlayersStatuses(steamIDs)
		if err != nil {
			log.Println("Ошибка при получении статусов игроков:", err)
			return context.Send("Ошибка при запросе к Steam API")
		}

		for _, player := range players {
			if player.GameID != "" && player.GameID == steam.CS22GameID {
				info := steamIDtoInfo[player.SteamID]
				responses = append(responses, fmt.Sprintf("🟢 %s (%s) играет в CS2", info, player.Persona))
			}
		}
	}

	if len(responses) == 0 {
		return context.Send("Никто из зарегестрированных пользователей сейчас не играет в CS2")
	}

	return context.Send(strings.Join(responses, "\n"))

}

func StartStatusChecker(bot *teleg.Bot) {

	go func() {
		for {
			log.Println("ШЕДУЛЕР ЗАПУЩЕН")
			// Получаем всех пользователей
			users, err := database.GetAllUsers()
			if err != nil {
				log.Println("Ошибка при получении пользователей:", err)
				continue
			}
			log.Println("Пользователи получены - ", users)
			// Готовим пакет
			const batchSize = 100

			for i := 0; i < len(users); i += batchSize {
				end := i + batchSize
				if end > len(users) {
					end = len(users)
				}

				batch := users[i:end]
				var steamIDs []string
				steamIDtoInfo := make(map[string]string)
				steamIDtoChatID := make(map[string]int64)

				for _, user := range batch {
					steamIDs = append(steamIDs, user.SteamID)
					steamIDtoInfo[user.SteamID] = user.Info
					steamIDtoChatID[user.SteamID] = user.ChatID
				}

				// Получаем игроков
				players, err := steam.GetPlayersStatuses(steamIDs)
				if err != nil {
					log.Println("Ошибка при получении статусов игроков:", err)
				}
				log.Println("Игроки получены - ", players)
				// Обрабатываем каждого игрока
				for _, player := range players {

					steamID := player.SteamID
					status := player.GameID == steam.CS22GameID

					prevStatus, err := database.GetUserStatus(steamID)
					if err != nil {
						log.Println("Ошибка получения статуса игрока из таблицы player_status", err)
						continue
					}
					log.Println("Статус игрока - ", player.Persona, prevStatus)

					if !prevStatus && status {
						//был не в игре, стал в игре
						chatID := steamIDtoChatID[steamID]
						log.Println("ЧАТ", chatID)
						_, err := bot.Send(&teleg.Chat{ID: chatID}, fmt.Sprintf("🎮 %s Ебашит в CS2!", player.Persona))
						if err != nil {
							log.Println("Ошибка при отправке сообщения: ", err)
						}
					}

					// Обновляем статус в базе данных
					err = database.UpdatePlayerStatus(steamID, status)
					if err != nil {
						log.Println("Ошибка при обновлении статуса: ", err)
					}

				}

			}

			// Тайм аут
			time.Sleep(1 * time.Minute)

		}
	}()
}
