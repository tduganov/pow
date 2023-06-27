package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

const (
	Difficulty = 4 // Уровень сложности PoW (количество ведущих нулей в хэше)
)

var (
	Quotes = []string{
		"Кто не рискует тот не рискует",
		"Кто не ест тот не ест",
		"Мертвый ранненого тащит",
	}
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	log.Println("Сервер запущен и ожидает подключений...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	// Отправка клиенту проблемы PoW
	challenge := generateChallenge()
	conn.Write([]byte(challenge + "\n"))

	// Получение решения PoW от клиента
	solution, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Println(err)
		return
	}
	solution = strings.TrimSpace(solution)

	// Проверка решения PoW
	if !validateSolution(challenge, solution) {
		conn.Write([]byte("Неверное решение PoW\n"))
		log.Printf("Клиент %s предоставил неверное решение PoW\n", conn.RemoteAddr().String())
		return
	}

	// Генерация и отправка цитаты
	quote := getRandomQuote()
	conn.Write([]byte(quote + "\n"))
	log.Printf("Клиент %s получил цитату: %s\n", conn.RemoteAddr().String(), quote)
}

func generateChallenge() string {
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%d", timestamp)
}

func validateSolution(challenge, solution string) bool {
	hash := sha256.Sum256([]byte(challenge + solution))
	hashString := hex.EncodeToString(hash[:])
	for i := 0; i < Difficulty; i++ {
		if hashString[i] != '0' {
			return false
		}
	}
	return true
}

func getRandomQuote() string {
	// Генерация случайного числа на основе времени
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(Quotes))
	return Quotes[randomIndex]
}
