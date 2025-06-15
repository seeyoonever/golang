package steam

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

const CS22GameID = "730"

type Player struct {
	SteamID string `json:"steamid"`
	GameID  string `json:"gameid,omitempty"` // omitempty говорит GO о том, что если поля нет, то не страшно и нужно поставить просто пустую строку
	Persona string `json:"personaname"`
}

// Функция получениея информации об игроке, из которой мы вытащим играет ли он в кс2
type PlayerInfoResponse struct {
	Response struct {
		Players []Player `json:"players"`
	} `json:"response"`
}

// Функция для проверки играет ли пользователь в КС2 сейчас
func IsPlayingCS2(steamID string) (bool, string, error) {
	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		return false, "", fmt.Errorf("STEAM_API_KEY не установлен")
	}

	url := fmt.Sprintf("https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s", apiKey, steamID)

	resp, err := http.Get(url)
	if err != nil {
		return false, "", fmt.Errorf("ошибка при запросе к Steam API: %w", err)
	}
	defer resp.Body.Close()

	var result PlayerInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return false, "", fmt.Errorf("ошибка при декодировании JSON: %w", err)
	}

	if len(result.Response.Players) == 0 {
		return false, "", fmt.Errorf("пользователь не найден")
	}

	player := result.Response.Players[0]

	if player.GameID != "" && player.GameID == CS22GameID {
		return true, player.Persona, nil
	}

	return false, player.Persona, nil
}

// Новая функция, которая будет проверять статус игроков батчем
func GetPlayersStatuses(steamIDs []string) ([]Player, error) {
	apiKey := os.Getenv("STEAM_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("STEAM_API_KEY не установлен")
	}

	// Формируем строку с запятыми
	ids := strings.Join(steamIDs, ",")
	fmt.Println(ids)
	url := fmt.Sprintf(
		"https://api.steampowered.com/ISteamUser/GetPlayerSummaries/v2/?key=%s&steamids=%s",
		apiKey, ids)

	for attempt := 1; attempt <= 10; attempt++ {
		resp, err := http.Get(url)

		if err != nil {
			return nil, fmt.Errorf("ошибка при запросе к Steam API: %w", err)
		}

		log.Println(resp)

		if resp.StatusCode == 429 {
			log.Println("🚫 Получен 429 Too Many Requests. Ожидание 10 секунд перед повтором...")
			time.Sleep(10 * time.Second)
			continue
		}

		if resp.StatusCode != 200 {
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			return nil, fmt.Errorf("неожиданный код ответа %d: %s", resp.StatusCode, body)
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("ошибка при чтении тела ответа: %w", err)
		}

		var result PlayerInfoResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
		}

		return result.Response.Players, nil
	}

	// var result PlayerInfoResponse
	// if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
	// 	return nil, fmt.Errorf("ошибка при декодировании JSON: %w", err)
	// }

	return nil, fmt.Errorf("не удалось получить данные с Steam API после 10 попыток")
}
