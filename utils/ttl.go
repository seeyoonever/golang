package utils

import (
	"log"
	"time"
)

// Уничерсальная функция добавления ключа с последующим удалением
func SetWithTTL[T any](m map[int64]T, key int64, value T, ttl time.Duration) {
	m[key] = value
	log.Println(m)

	go func() {
		time.Sleep(ttl)
		delete(m, key)
		log.Println("Флаги очищены")
	}()
}
