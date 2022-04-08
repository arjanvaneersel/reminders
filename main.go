package main

import (
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func MustGetenv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Printf("%q not found in ENV\n", key)
		os.Exit(1)
	}

	return v
}

func GetenvOr(key string, or string) string {
	v := os.Getenv(key)
	if v == "" {
		return or
	}

	return v
}

func sendMail(from, to, subject, user, password, host, port, msg string) error {
	body := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", from, to, subject, msg)
	auth := smtp.PlainAuth("", user, password, host)
	if err := smtp.SendMail(host+":"+port, auth, from, []string{to}, []byte(body)); err != nil {
		return err
	}

	return nil
}

func exitWithErr(err error) {
	fmt.Fprintf(os.Stderr, "%s", err)
	os.Exit(1)
}

func main() {
	if err := godotenv.Load(); err != nil {
		exitWithErr(err)
	}

	from := MustGetenv("FROM")
	user := GetenvOr("ACCOUNT", from)
	password := MustGetenv("PASSWD")
	host := MustGetenv("HOST")
	port := GetenvOr("PORT", "587")
	to := MustGetenv("TO")
	subject := MustGetenv("SUBJECT")
	file := GetenvOr("FILE", "message.txt")

	msg, err := os.ReadFile(file)
	if err != nil {
		exitWithErr(err)
	}

	if err := sendMail(from, to, subject, user, password, host, port, string(msg)); err != nil {
		exitWithErr(err)
	}
	fmt.Fprintf(os.Stdout, "Successfully sent initial mail to %s\n", to)
	c := 2

	t := time.NewTicker(1 * time.Hour)
	for {
		select {
		case <-t.C:
			if err := sendMail(from, to, subject, user, password, host, port, string(msg)); err != nil {
				exitWithErr(err)
			}
			fmt.Fprintf(os.Stdout, "Successfully sent mail %d to %s\n", c, to)
			c++
		}
	}
}
