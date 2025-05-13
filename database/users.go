// Заполнение таблицы users

package database

import (
	"context"
	"fmt"
	"time"
)

type User struct {
	ID         int
	TelegramID int64
	SteamID    string
	CreatedAt  time.Time
	Info       string
}

func AddUser(telegramID int64, steamID string, info string) error {
	query := "INSERT INTO users (telegram_id, steam_id, info) VALUES ($1,$2,$3)"
	_, err := DB.Exec(context.Background(), query, telegramID, steamID)
	if err != nil {
		return fmt.Errorf("Не удалось добавить пользователя: %w", err)
	}
	fmt.Println("✅ Пользователь добавлен")
	return nil
}
