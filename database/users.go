// Заполнение таблицы users

package database

import (
	"context"
	"fmt"
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
}

// Добавление пользователя
func AddUser(telegramID int64, steamID string, info string) error {
	query := "INSERT INTO users (telegram_id, steam_id, info) VALUES ($1,$2,$3)"
	_, err := DB.Exec(context.Background(), query, telegramID, steamID, info)
	if err != nil {
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
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
	query := "SELECT id, telegram_id, steam_id, created_at, info FROM users"
	rows, err := DB.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.TelegramID, &user.SteamID, &user.CreatedAt, &user.Info)
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
