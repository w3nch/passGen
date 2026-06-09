package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/jedib0t/go-pretty/v6/table"
)

type passInfo struct {
	passwordLength      int64
	passwordSpecialChar bool
	passwordCapital     bool
	passwordNumbers     bool
}

const (
	lower   = "abcdefghijklmnopqrstuvwxyz"
	upper   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	digits  = "0123456789"
	special = "!@#$%^&*()_+-=[]{}|;:,.<>?"
)

func showMenu() {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	t.AppendHeader(table.Row{
		"Option",
		"Strength",
		"Length",
		"Special",
		"Capital",
		"Numbers",
	})

	t.AppendRows([]table.Row{
		{1, "Weak", 8, "No", "No", "Yes"},
		{2, "Medium", 12, "No", "Yes", "Yes"},
		{3, "Strong", 40, "Yes", "Yes", "Yes"},
	})

	t.Render()
}

func getPasswordProfile() passInfo {
	var choice int

	showMenu()

	fmt.Print("\nChoose a profile: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		return passInfo{
			passwordLength:      8,
			passwordSpecialChar: false,
			passwordCapital:     false,
			passwordNumbers:     true,
		}

	case 2:
		return passInfo{
			passwordLength:      12,
			passwordSpecialChar: false,
			passwordCapital:     true,
			passwordNumbers:     true,
		}

	case 3:
		return passInfo{
			passwordLength:      40,
			passwordSpecialChar: true,
			passwordCapital:     true,
			passwordNumbers:     true,
		}

	default:
		fmt.Println("Invalid option, using Strong.")
		return passInfo{
			passwordLength:      40,
			passwordSpecialChar: true,
			passwordCapital:     true,
			passwordNumbers:     true,
		}
	}
}

func collectEntropy() [32]byte {
	var data strings.Builder

	lastX, lastY := robotgo.Location()

	fmt.Println("\nMove your mouse around...")

	for range 100 {
		x, y := robotgo.Location()

		if x != lastX || y != lastY {
			data.WriteString(fmt.Sprintf(
				"%d:%d:%d|",
				x,
				y,
				time.Now().UnixNano(),
			))
		}

		lastX = x
		lastY = y

		time.Sleep(10 * time.Millisecond)
	}

	hash := sha256.Sum256([]byte(data.String()))

	fmt.Printf("\nEntropy Hash: %x\n", hash)

	return hash
}

func generatePassword(info passInfo, entropyHash [32]byte) (string, error) {
	charset := lower

	if info.passwordCapital {
		charset += upper
	}

	if info.passwordNumbers {
		charset += digits
	}

	if info.passwordSpecialChar {
		charset += special
	}

	if len(charset) == 0 {
		return "", fmt.Errorf("empty charset")
	}

	var password strings.Builder

	for i := int64(0); i < info.passwordLength; i++ {
		secureRandom, err := rand.Int(
			rand.Reader,
			big.NewInt(int64(len(charset))),
		)
		if err != nil {
			return "", err
		}

		entropyOffset := int64(entropyHash[i%32]) % int64(len(charset))

		index := (secureRandom.Int64() + entropyOffset) % int64(len(charset))

		password.WriteByte(charset[index])
	}

	return password.String(), nil
}

func main() {
	info := getPasswordProfile()

	entropyHash := collectEntropy()

	password, err := generatePassword(info, entropyHash)
	if err != nil {
		panic(err)
	}

	fmt.Println("\nGenerated Password:")
	fmt.Println(password)
}
