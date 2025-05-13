package database // этот файл принадлежит пакету database - логическому модулю

// в GO всё состоит из пакетов и каждый модуль проекта это по сути пакет

import (
	"context" // либа для работы с запросами к БД
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool" // библиотека - драйвер для PostgreSQL которая представляет пул соединений
)

var DB *pgxpool.Pool // переменная с подключением к базе

// Функция для подключения к базе (возвращет ошибку, если соединение упадёт)
func Connect() error {
	dbUrl := "postgres://postgres:secret@localhost:5432/cs2bot" // Строка подключения

	var err error

	DB, err = pgxpool.New(context.Background(), dbUrl) // Создаём пул соединений
	if err != nil {
		return fmt.Errorf("Не удалось подключиться к БД: %w", err)
	}

	err = DB.Ping(context.Background()) // Пингуем базу данных, для проверки соединения
	if err != nil {
		return fmt.Errorf("БД не отвечает: %w", err)
	}

	fmt.Println("✅ Успешное подключение к PostgreSQL")
	return nil
}
