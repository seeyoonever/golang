package steam

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const cs2GameID = "730"

// Функция получениея информации об игроке, из которой мы вытащим играет ли он в кс2
type PlayerInfoResponse struct {
	Response struct {
		Players []struct {
			SteamID string `json:"steamid"`
			GameID  string `json:"gameid.omitempty"`
			Persona string `json:"personaname"`
		} `json:"players"`
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

	if player.GameID != "" && player.GameID == cs2GameID {
		return true, player.Persona, nil
	}

	return false, player.Persona, nil
}
