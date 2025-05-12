package main

import (
	"fmt"
)

// Глобальный срез
var favoritLangs []string

// Глобальная память
var memory = make(map[string]string)

func main() {
	fmt.Println("Привет! Я твой персональный помощник.")

	name := askName()
	greetUser(name)

	for {
		fmt.Println("\nДоступные команды:")
		fmt.Println("1. joke — рассказать шутку")
		fmt.Println("2. age  — ввести свой возраст")
		fmt.Println("3. hello — приветствие")
		fmt.Println("4. favlang - любимые языки")
		fmt.Println("5. listlangs - список твоих любимых языков")
		fmt.Println("6. remember - запомнить что-то")
		fmt.Println("7. recall - показать информацию о чём-либо")
		fmt.Println("8. showmem - показать все записи")
		fmt.Println("--- exit — выйти")

		var command string
		fmt.Print("Введи команду: ")
		fmt.Scanln(&command)

		if command == "exit" {
			fmt.Println("Пока, " + name + "!")
			break // Выход из цикла
		}

		handleCommand(command)
	}
}

// Спросить имя
func askName() string {
	fmt.Print("Как тебя зовут? ")
	var name string
	fmt.Scanln(&name)
	return name
}

// Приветствие
func greetUser(name string) {
	fmt.Println("Рад познакомиться,", name+"!")
}

// Обработка команд
func handleCommand(cmd string) {
	if cmd == "joke" {
		fmt.Println("Почему гусь перешёл дорогу? Чтобы избежать null pointer exception.")
	} else if cmd == "age" {
		askAge()
	} else if cmd == "hello" {
		fmt.Println("Привет-привет! Рад тебя видеть 😊")
	} else if cmd == "favlang" {
		favlang()
	} else if cmd == "listlangs" {
		listLangs()
	} else if cmd == "remember" {
		rememberSomething()
	} else if cmd == "recall" {
		recallSomething()
	} else if cmd == "showmem" {
		showMemory()
	} else {
		fmt.Println("Неизвестная команда. Попробуй снова.")
	}
}

// Обработка любимого языка
func favlang() {
	var lang string
	flag := true
	fmt.Println("Введите любимый язык программирования")
	fmt.Scanln(&lang)

	for _, language := range favoritLangs {
		if lang == language {
			fmt.Println("Такой язык ты уже указывал")
			flag = false
		}
	}

	if flag == true {
		favoritLangs = append(favoritLangs, lang)
		fmt.Println("Добавил ", lang, " в список твоих любимых языков")
	}

}

// Список любимых языков
func listLangs() {
	if len(favoritLangs) == 0 {
		fmt.Println("Ты ещё ничего не добавил, наверное не любишь программирование")
	} else {
		fmt.Println("Твои любимые языки: ")
		for i, lang := range favoritLangs {
			fmt.Printf("%d. %s\n", i+1, lang)
		}
	}
}

// Возраст
func askAge() {
	var age int
	fmt.Print("Сколько тебе лет? ")
	fmt.Scanln(&age)

	if age < 20 {
		fmt.Println("Ого, да ты совсем пиздюк")
	} else if age < 30 {
		fmt.Println("Ты молодой и пиздатый")
	} else {
		fmt.Println("Старый, но не бесполезный!")
	}

}

// Функция запоминания
func rememberSomething() {
	var key, value string

	fmt.Println("Что мне запомнить, ВВЕДИ КЛЮЧ")
	fmt.Scanln(&key)

	fmt.Println("Что мне запомнить по этому ключу?")
	fmt.Scanln(&value)

	memory[key] = value
	fmt.Println("Запомнил:", key, " - ", value)
}

// Функция воспоминания по ключу
func recallSomething() {
	var key string

	if len(memory) == 0 {
		fmt.Println("Тебе нечего вспоминать")
		return
	}

	fmt.Println("Что тебе показать? ВВЕДИ КЛЮЧ")
	for k, _ := range memory {
		fmt.Printf("- %s\n", k)
	}
	fmt.Scanln(&key)

	value, exist := memory[key]
	if exist {
		fmt.Println("Ты говорил, что: ", key, " - ", value)
	} else {
		fmt.Println("Я ничего об этом не знаю")
	}

}

// Функция показа всей запомненной информации
func showMemory() {
	if len(memory) == 0 {
		fmt.Println("Я пока ничего не запоминал")
	} else {
		fmt.Println("Вот что я знаю амиго")
		for key, value := range memory {
			fmt.Printf("- %s: %s\n", key, value)
		}
	}
}
