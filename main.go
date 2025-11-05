package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"strconv"
)

type Bank struct {
	Name    string
	BinFrom int
	BinTo   int
}

func loadBankData(path string) ([]Bank, error) {
	file, err := os.Open("banks.txt") // открываем файл
	if err != nil { // err — ошибку, если файл не удалось открыть.
		return nil, err
	}

	defer file.Close() // если мы открываем файл, то мы обязаны его в конце закрыть

	scanner := bufio.NewScanner(file) // Создаём сканер, который будет читать файл построчно

	var banks []Bank
	// Создаём пустой срез banks.
	// Сюда мы будем добавлять банки, считанные из файла

	for scanner.Scan() { // цикл, который читает все строки
		line := scanner.Text() // текущая строка из файла пусть станет текстом
		parts := strings.Split(line, ",") 
		// парсим строку по запятым
		// "Sberbank,400000,499999" → ["Sberbank", "400000", "499999"]
		if len(parts) != 3 { // Проверяем, что строка действительно состоит из трёх частей
			continue
		}

		BinFrom, err := strconv.Atoi(parts[1]) // переводим второй индекс парса из строки в число
		if err != nil { // проверка ошибки
			return nil, err
		}

		BinTo, err := strconv.Atoi(parts[2]) // переводим третий индекс парса из строки в число
		if err != nil { // проверка ошибки
			return nil, err
		}

		bank := Bank{ // создаем структуру с нашими данными
			Name:    parts[0],
			BinFrom: BinFrom,
			BinTo:   BinTo,
		}

		banks = append(banks, bank) // Добавляем этот банк в наш список banks
	}

	return banks, nil // возвращаем значение
}


func identifyBank(bin int, banks []Bank) string {
	for _, bank := range banks { // Начинаем цикл по всем банкам в списке banks
		if bin >= bank.BinFrom && bin <= bank.BinTo { // проверка диапазона BinFrom -> BinTo
			return bank.Name // возвращаем имя банка, обращаемся к полю Name
		}
	}

	return "Неизвестный банк"
}

func extractBIN(cardNumber string) int {
	if len(cardNumber) < 6 { // if меньше 6 цифр - ошибка
		return 0
	}

	sixNumber := cardNumber[:6] // срезаем первые 6 цифр
	sixNumberAtoi, err := strconv.Atoi(sixNumber) // переводим их в числа
	if err != nil { // проверка ошибки
		return 0
	}
	return sixNumberAtoi
}


// StringToDigit проверяет все ли руны в строке это числа
// добавляет каждую цифру в массив и возвращает его
func StringToDigit(s string) []int {
	var digits []int // создаем пустой массив из чисел
	for _, r := range s { // проходимся по всей строке
		if r < '0' || r > '9' { // если руна в строке не число, возвращаем пустой массив
			return nil
		}
		digits = append(digits, int(r-'0')) // если все ок, добавляет в пустой массив значение
	}
	return digits
}


// validateLuhn проходится по массиву цифр справа налево, удваивает каждое второе число
// при необходимости вычитает 9 и проверяет итоговую сумму. sum % 10 == 0
func validateLuhn(cardNumber string) bool {
	arrayDigits := StringToDigit(cardNumber) // создание массива с помощью ф-и StringToDigit
	if arrayDigits == nil { // проверка ошибки
		return false
	}

	var sum int // сумма
	double := false // флаг для умножения на 2

	for i := len(arrayDigits) - 1; i >= 0; i-- { // итерация по массиву с конца, с шагом 2
		d := arrayDigits[i]
		if double { // если флаг true, то умножаем на 2
			d *= 2
			if d > 9 { // если получившееся число > 9, вычитаем 9
				d -= 9
			}
		}
		sum += d // складываем все результаты
		double = !double // обновляем флаг, чтобы он стал false
	}

	return sum % 10 == 0
}



// GetUserInput считывает и возвращает строку с терминала
func getUserInput() string {
    reader := bufio.NewReader(os.Stdin)
    input, _ := reader.ReadString('\n')
    return strings.TrimSpace(input)
}


// validateInput проверяет номер карты на корректность >=13 && <= 19
func validateInput(cardNumber string) bool { 	
	digits := StringToDigit(cardNumber)
    if digits == nil {
        return false
    }
    return len(digits) >= 13 && len(digits) <= 19
}



func main() {
	fmt.Println("Добро пожаловать в программу валидации карт!")
	banks, err := loadBankData("banks.txt") // загрузка всех банков из txt
	if err != nil { // проверка ошибки
		fmt.Println("Ошибка при загрузке данных банков: ", err)
		os.Exit(1)
	}
	fmt.Println("Данные банков успешно загружены!")

	for { // бессконечный цикл с проверками
		cardNumber := getUserInput() // ввод в терминале

		if cardNumber == "exit" { // if ввели exit, программа завершена
			fmt.Println("Программа завершена")
			break
		}
		if cardNumber == "" { // if ввели "", пустая строка
			fmt.Println("Пустая строка")
			continue
		}

		if !validateInput(cardNumber) { 
			fmt.Println("Ошибка: номер карты должен содержать только цифры и иметь длину от 13 до 19 символов.")
			continue
		}

		if !validateLuhn(cardNumber) {
			fmt.Println("Невалидный номер карты!")
			continue
		}

		bin := extractBIN(cardNumber)
		bank := identifyBank(bin, banks)

		fmt.Println("Номер карты валиден")
		fmt.Println("Банк:", bank)
	}
}
