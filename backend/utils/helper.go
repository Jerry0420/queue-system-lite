package utils

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	gomail "gopkg.in/gomail.v2"
)

func CheckKeysInMap(data map[string]interface{}, keys []string) bool {
	for _, key := range keys {
		if _, ok := data[key]; !ok {
			return false
		}
	}
	return true
}

func InitDirPath(csvDirPath string) error {
	if _, err := os.Stat(csvDirPath); os.IsNotExist(err) {
		err := os.Mkdir(csvDirPath, 0777)
		if err != nil {
			return err
		}
	}
	return nil
}

func InitEmailDialer(emailServer string, emailPort int, emailUserName string, emailPassword string) *gomail.Dialer {
	dialer := gomail.NewDialer(
		emailServer,
		emailPort,
		emailUserName,
		emailPassword,
	)
	return dialer
}

func GenerateCSV(ctx context.Context, timeOut time.Duration, name string, content []byte) (filePath string, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	csvFilePath := filepath.Join("csvs", name+".csv")
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		return csvFilePath, err
	}
	defer csvFile.Close()

	err = os.Chmod(csvFilePath, 0777)
	if err != nil {
		return csvFilePath, err
	}

	csvWriter := csv.NewWriter(csvFile)

	var cotentMap [][]string
	json.Unmarshal(content, &cotentMap)
	err = csvWriter.WriteAll(cotentMap)
	if err != nil {
		return csvFilePath, err
	}
	csvWriter.Flush()

	err = csvWriter.Error()
	if err != nil {
		return csvFilePath, err
	}

	return csvFilePath, nil

}

func SendEmail(ctx context.Context, timeOut time.Duration, dialer *gomail.Dialer, fromEmail string, subject string, content string, email string, filePath string) (result bool, err error) {
	ctx, cancel := context.WithTimeout(ctx, timeOut)
	defer cancel()

	message := gomail.NewMessage()
	message.SetHeader("From", fromEmail)
	message.SetHeader("To", email)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html", content)

	if filePath != "" {
		message.Attach(filePath)
	}
	err = dialer.DialAndSend(message)
	if err != nil {
		fmt.Println(err)
		// TODO: keep it!
		return false, err
	}

	if filePath != "" {
		os.Remove(filePath)
	}

	return true, nil
}
