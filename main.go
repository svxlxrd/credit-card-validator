package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

// loadBankData загружает данные банков из файла и возвращает слайс структур Bank
func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var banks []Bank

	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) != 3 {
			continue
		}

		binFrom, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}

		binTo, err := strconv.Atoi(parts[2])
		if err != nil {
			return nil, err
		}

		banks = append(banks, Bank{
			Name:    parts[0],
			BinFrom: binFrom,
			BinTo:   binTo,
		})
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return banks, nil
}

// identifyBank возвращает имя банка по bin или "Неизвестный банк" если не найден
func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks {
		if bin >= bank.BinFrom && bin <= bank.BinTo {
			return bank.Name
		}
	}
	return "Неизвестный банк"
}

// extractBIN извлекает первые 6 цифр номера карты и возвращает их как число
func extractBIN(cardNumber string) int {
	if len(cardNumber) < 6 {
		return -1
	}

	bin, err := strconv.Atoi(cardNumber[:6])
	if err != nil {
		return -1
	}
	return bin
}

// parseDigits проверяет, что строка состоит только из цифр и возвращает слайс int
func parseDigits(s string) []int {
	var digits []int
	for _, r := range s {
		if r < '0' || r > '9' {
			return nil
		}
		digits = append(digits, int(r-'0'))
	}
	return digits
}

// validateLuhn проверяет номер карты по алгоритму Луна
func validateLuhn(cardNumber string) bool {
	digits := parseDigits(cardNumber)
	if digits == nil {
		return false
	}

	sum := 0
	shouldDouble := false

	for i := len(digits) - 1; i >= 0; i-- {
		d := digits[i]
		if shouldDouble {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		shouldDouble = !shouldDouble
	}

	return sum%10 == 0
}

// getUserInput считывает строку с терминала и удаляет пробелы
func getUserInput() string {
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// validateInput проверяет, что номер карты содержит только цифры и имеет длину 13 - 19
func validateInput(cardNumber string) bool {
	digits := parseDigits(cardNumber)
	return len(digits) >= 13 && len(digits) <= 19
}

func main() {
	fmt.Println("Добро пожаловать в программу валидации карт!")
	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка при загрузке данных банков:", err)
		os.Exit(1)
	}
	fmt.Println("Данные банков успешно загружены!")
	fmt.Println("Введите номер карты или exit для выхода:")

	for {
		cardNumber := getUserInput()

		if cardNumber == "exit" {
			fmt.Println("Программа завершена")
			break
		}
		if cardNumber == "" {
			fmt.Println("Пустая строка, попробуйте снова")
			continue
		}
		if !validateInput(cardNumber) {
			fmt.Println("Ошибка: номер карты должен содержать только цифры и иметь длину от 13 до 19")
			continue
		}
		if !validateLuhn(cardNumber) {
			fmt.Println("Невалидный номер карты")
			continue
		}

		bin := extractBIN(cardNumber)
		bank := identifyBank(bin, banks)

		fmt.Println("Номер карты валиден")
		fmt.Println("Банк:", bank)
	}
}
