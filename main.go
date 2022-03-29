package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/claudiu/gocron"
	"github.com/go-redis/redis/v8"
	gomail "gopkg.in/mail.v2" //go get gopkg.in/mail.v2
)

// ctx (global) for redis
var ctx = context.Background()

// check error from any tools
func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
		// fmt.Println(err)
	}
}

// GoMail - Send Email
func SendEmail(from string, to string, subject string, text string) {
	gml := gomail.NewMessage()

	gml.SetHeader("From", from)
	gml.SetHeader("To", to)
	gml.SetHeader("Subject", subject)
	gml.SetBody("text/plain", text)

	sender := gomail.NewDialer("smtp.gmail.com", 587, "cobapbp@gmail.com", "CobaPBP5656")

	err := sender.DialAndSend(gml)
	CheckError(err)

	fmt.Println("Email Sent")
}

// GoRedis - Set & Get from Redis
func SetRedis(rdb *redis.Client, key string, value string, expiration int) {
	err := rdb.Set(ctx, key, value, 0).Err()
	CheckError(err)
}

func GetRedis(rdb *redis.Client, key string) string {
	val, err := rdb.Get(ctx, key).Result()

	// if err == redis.Nil {
	// 	fmt.Println(key, "does not exist")
	// }

	CheckError(err)
	return val
}

// GoRoutine - Do SendEmail 2 times with different body text at the same time
func task(eng string, idn string) { //parse eng and idn text
	go SendEmail("cobapbp@gmail.com", "cobapbp@gmail.com", "Hello From LetsGyu - Reminder", eng) // sender, receiver, subject, body
	SendEmail("cobapbp@gmail.com", "cobapbp@gmail.com", "Hello From LetsGyu - Reminder", idn)
}

// GoRoutine - Do Multiple
// func Do(text string) {
// 	for i := 0; i < 2; i++ { // Multiple Do
// 		time.Sleep(5 * time.Second) // Time range for each Do
// 		SendEmail("cobapbp@gmail.com", "cobapbp@gmail.com", "Test", "Hai") // Do something
// 		fmt.Println("email sent")
// 	}
// }

func main() {

	// LetsGyu Reminder in English & Indonesia

	// GoRedis
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	SetRedis(rdb, "eng", "Don't forget to order braderr", 0) // set key and its value
	SetRedis(rdb, "idn", "Info mazeehhh - jangan lupa order :)", 0)

	eng := GetRedis(rdb, "eng") // get value with specific key
	idn := GetRedis(rdb, "idn")

	// Erase All Keys
	// rdb.FlushDB(ctx)

	// GoCRON
	gocron.Start()
	gocron.Every(10).Seconds().Do(task, eng, idn) // every 10 seconds do task with parameter eng and idn

	time.Sleep(40 * time.Second) // program will stop at 40 seconds

	gocron.Clear() // remove all gocron task
	fmt.Println("All task removed")

}
