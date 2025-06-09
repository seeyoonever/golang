// Заполнение таблицы users

package database

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Структура - специальный тип данных который позволяет объединить несколько значений одного типа под одним имененем
// Это как класс в Python, но без методов
type User struct {
	ID         int
	TelegramID int64
	SteamID    string
	CreatedAt  time.Time
	Info       string
	ChatID     int64
}

// Добавление пользователя
func AddUser(telegramID int64, steamID string, info string, chatID int64) error {
	query := "INSERT INTO users (telegram_id, steam_id, info, chat_id) VALUES ($1,$2,$3,$4)"
	_, err := DB.Exec(context.Background(), query, telegramID, steamID, info, chatID)
	if err != nil {
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
	}

	// Вставка в таблицу player_status
	queryStatus := "INSERT INTO player_status (steam_id, in_game) VALUES ($1, FALSE)"
	_, err = DB.Exec(context.Background(), queryStatus, steamID)
	if err != nil {
		return fmt.Errorf("не удалось добавить статус игрока: %w", err)
	}

	fmt.Println("✅ Пользователь добавлен")
	return nil
}

// Проверка существования пользователя в базе
func UserExists(telegramID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE telegram_id=$1"

	var count int
	err := DB.QueryRow(context.Background(), query, telegramID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить наличие пользователя: %w", err)
	}

	return count > 0, nil

}

// Просмотр всех пользователей
func GetAllUsers() ([]User, error) {
	query := "SELECT id, telegram_id, steam_id, created_at, info, chat_id FROM users"
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.TelegramID, &user.SteamID, &user.CreatedAt, &user.Info, &user.ChatID)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func DeleteUser(telegramID int64) error {
	query := "DELETE FROM users WHERE telegram_id = $1"

	_, err := DB.Exec(context.Background(), query, telegramID)

	return err
}

// Функция которая забирает статус игрока из таблицы player_status
func GetUserStatus(steamID string) (bool, error) {
	var status bool
	query := "SELECT in_game FROM player_status WHERE steam_id = $1"

	err := DB.QueryRow(context.Background(), query, steamID).Scan(&status)
	if err != nil {
		log.Println("Ошибка запроса к БД для извлечения статуса игрока: ", err)
		return false, nil
	}

	return status, nil

}

// Функция для обновления статуса игрока в таблице player_status
func UpdatePlayerStatus(steamID string, status bool) error {
	query := `
			UPDATE player_status
			SET in_game = $1,
				last_checked = CURRENT_TIMESTAMP
			WHERE steam_id = $2
			`
	_, err := DB.Exec(context.Background(), query, status, steamID)
	if err != nil {
		return fmt.Errorf("не удалось обновить статус игрока: %w", err)
	}
	return nil
}
