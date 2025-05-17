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

func AddUser(telegramID int64, steamID string, info string) error {
	query := "INSERT INTO users (telegram_id, steam_id, info) VALUES ($1,$2,$3)"
	_, err := DB.Exec(context.Background(), query, telegramID, steamID, info)
	if err != nil {
		return fmt.Errorf("не удалось добавить пользователя: %w", err)
	}
	fmt.Println("✅ Пользователь добавлен")
	return nil
}

func UserExists(telegramID int64) (bool, error) {
	query := "SELECT COUNT(*) FROM users WHERE telegram_id=$1"

	var count int
	err := DB.QueryRow(context.Background(), query, telegramID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("не удалось проверить наличие пользователя: %w", err)
	}

	return count > 0, nil

}
