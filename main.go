package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Bank describes an issuer by the prefix its card numbers start with.
type Bank struct {
	Name   string
	Prefix string
}

// Matches reports whether the card number belongs to this bank.
func (b Bank) Matches(cardNumber string) bool {
	return strings.HasPrefix(cardNumber, b.Prefix)
}

const defaultCard = "4000123456789017"

func main() {
	banks, err := loadBankData("banks.txt")
	if err != nil {
		fmt.Println("Ошибка: не удалось загрузить список банков:", err)
		return
	}
	fmt.Printf("Загружено банков: %d\n", len(banks))

	cardNumber := defaultCard
	if len(os.Args) > 1 {
		cardNumber = os.Args[1]
	}

	fmt.Println("Карта", cardNumber, "валидна по Луну:", LuhnCheck(cardNumber))

	if bank := DetectBank(cardNumber, banks); bank != nil {
		fmt.Println("Банк:", bank.Name)
	} else {
		fmt.Println("Банк не определён")
	}
}

// loadBankData reads "Name,Prefix" lines from path.
func loadBankData(path string) ([]Bank, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open banks file: %w", err)
	}
	defer f.Close()

	banks := make([]Bank, 0, 8)
	sc := bufio.NewScanner(f)
	lineNo := 0
	for sc.Scan() {
		lineNo++
		raw := strings.TrimSpace(sc.Text())
		if raw == "" {
			continue
		}
		fields := strings.Split(raw, ",")
		if len(fields) != 2 {
			return nil, fmt.Errorf("line %d: expected Name,Prefix, got %q", lineNo, raw)
		}
		banks = append(banks, Bank{
			Name:   strings.TrimSpace(fields[0]),
			Prefix: strings.TrimSpace(fields[1]),
		})
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read banks file: %w", err)
	}
	return banks, nil
}

// DetectBank returns the first bank whose prefix matches, or nil.
func DetectBank(cardNumber string, banks []Bank) *Bank {
	for i := range banks {
		if banks[i].Matches(cardNumber) {
			return &banks[i]
		}
	}
	return nil
}

// LuhnCheck validates a card number with the Luhn algorithm.
func LuhnCheck(cardNumber string) bool {
	if len(cardNumber) < 2 {
		return false
	}
	sum := 0
	for i, pos := len(cardNumber)-1, 0; i >= 0; i, pos = i-1, pos+1 {
		d := int(cardNumber[i] - '0')
		if d < 0 || d > 9 {
			return false
		}
		if pos%2 == 1 {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
	}
	return sum%10 == 0
}
