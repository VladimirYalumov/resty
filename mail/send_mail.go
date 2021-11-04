package mail

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"goRestApi_main/redis"
	"html/template"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

const SUBJECT_AUTH = "Confirm sign up on deewave.online"
const CODE_NUMBERS_COUNT = 4
const CODE_LIFETIME = "10"

type Code struct {
	Num1 string
	Num2 string
	Num3 string
	Num4 string
}

type Mail struct {
	senderId string
	toIds    string
	subject  string
	body     string
}

type SmtpServer struct {
	host string
	port string
}

func (s *SmtpServer) ServerName() string {
	return s.host + ":" + s.port
}

func (mail *Mail) BuildMessage() string {
	message := ""
	message += fmt.Sprintf("From: %s\r\n", mail.senderId)
	if len(mail.toIds) > 0 {
		message += fmt.Sprintf("To: %s\r\n", strings.Join([]string{mail.toIds}, ";"))
	}
	message += fmt.Sprintf("Subject: %s\r\n", mail.subject)
	message += fmt.Sprintf("Content-Type: text/html; charset=\"utf-8\"\r\n")
	message += "\r\n" + mail.body
	return message
}

type Config struct {
	SenderName     string
	SenderPassword string
	Smtp           string
}

var c Config

func InitMailConfig(smtp string, user string, password string) {
	c.Smtp = smtp
	c.SenderName = user
	c.SenderPassword = password
}

func SendAuthMessage(email string) (bool, error) {
	if c.Smtp == "" || c.SenderName == "" || c.SenderPassword == "" {
		// for testing
		return true, nil
	}
	mail := Mail{}
	mail.senderId = c.SenderName
	mail.toIds = email
	mail.subject = SUBJECT_AUTH

	redis.RedisClient.Set("name", "Elliot", 0)

	smtpServer := SmtpServer{host: c.Smtp, port: "465"}

	auth := smtp.PlainAuth("", mail.senderId, c.SenderPassword, smtpServer.host)

	tlsconfig := &tls.Config{InsecureSkipVerify: true, ServerName: smtpServer.host}
	conn, err := tls.Dial("tcp", smtpServer.ServerName(), tlsconfig)
	if err != nil {
		return false, err
	}
	client, err := smtp.NewClient(conn, smtpServer.host)
	if err != nil {
		return false, err
	}
	if err = client.Auth(auth); err != nil {
		return false, err
	}
	if err = client.Mail(mail.senderId); err != nil {
		return false, err
	}
	if err = client.Rcpt(mail.toIds); err != nil {
		return false, err
	}

	w, err := client.Data()

	if err != nil {
		return false, err
	}

	code, success := BuildCode(email)
	if !success {
		return false, nil
	}

	mail.body = createBody(code)
	messageBody := mail.BuildMessage()

	if _, err = w.Write([]byte(messageBody)); err != nil {
		return false, err
	}

	if err = w.Close(); err != nil {
		return false, err
	}
	client.Quit()

	return true, nil
}

func createBody(numbers []string) string {
	tpl, err := template.ParseFiles("mail/auth.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	code := Code{
		Num1: numbers[0],
		Num2: numbers[1],
		Num3: numbers[2],
		Num4: numbers[3],
	}
	buf := &bytes.Buffer{}
	err = tpl.Execute(buf, code)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func BuildCode(email string) ([]string, bool) {
	rand.Seed(time.Now().Unix())
	var code [CODE_NUMBERS_COUNT]int
	var result []string
	for i := 0; i < CODE_NUMBERS_COUNT; i++ {
		code[i] = rand.Intn(9)
		result = append(result, strconv.Itoa(code[i]))
	}

	authCodeCounter := redis.RedisClient.Get(redis.CreateKey(redis.REDIS_EMAIL_AUTH_CODE_COUNT, email)).Val()

	if authCodeCounter == "" {
		authCodeCounter = "1"
	} else {
		var authCodeCounterInt int
		authCodeCounterInt, _ = strconv.Atoi(authCodeCounter)
		if authCodeCounterInt != 1 && authCodeCounterInt != 2 && authCodeCounterInt != 3 {
			authCodeCounter = "3"
			return []string{}, false
		} else {
			authCodeCounterInt++
			authCodeCounter = strconv.Itoa(authCodeCounterInt)
		}
	}
	timelimit, _ := time.ParseDuration(fmt.Sprintf("%sm", CODE_LIFETIME))
	redis.RedisClient.Set(redis.CreateKey(redis.REDIS_EMAIL_AUTH_CODE_COUNT, email), authCodeCounter, timelimit)
	redis.RedisClient.Set(redis.CreateKey(redis.REDIS_EMAIL_AUTH_CODE, email), strings.Join(result, ""), timelimit)

	return result, true
}
