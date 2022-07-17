package main

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/gocolly/colly"
)

func getBankuzInfo() string {
	c := colly.NewCollector()

	var curr []string
	c.OnHTML(".left-side .other-bank-course-block-mob span", func(e *colly.HTMLElement) {
		if e.Attr("class") == "semibold-text" {
			curr = append(curr, e.Text)
		}
	})

	c.Visit("https://bank.uz/currency/kzt")
	return curr[1] + "\n" + "https://bank.uz/currency/kzt"
}

func getKaseInfo() []string {
	c := colly.NewCollector()

	var curr []string
	c.OnHTML(".center-column span", func(e *colly.HTMLElement) {
		if e.Attr("class") == "currency-round__round" {
			curr = append(curr, strings.TrimSpace(e.Text))
		}
	})

	c.Visit("https://kase.kz/ru/currency/")
	return curr
}

func getCurrency(msg string) string {
	var result string
	switch msg {
	case "usd":
		result += getKaseInfo()[0] + "\n" + "https://kase.kz/ru/currency/"
	case "rub":
		result += getKaseInfo()[3] + "\n" + "https://kase.kz/ru/currency/"
	case "eur":
		result += getKaseInfo()[7] + "\n" + "https://kase.kz/ru/currency/"
	case "cny":
		result += getKaseInfo()[1] + "\n" + "https://kase.kz/ru/currency/"
	case "uzs":
		result += getBankuzInfo()
	}
	return result
}

func main() {
	var currency []string = []string{"usd", "rub", "eur", "uzs", "cny"}
	bot, err := tgbotapi.NewBotAPI("***************")
	if err != nil {
		fmt.Println(err)
	}

	bot.Debug = true

	log.Printf("Auth on accaount %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

			for i := 0; i < len(currency); i++ {
				if strings.ToLower(msg.Text) == currency[i] {
					msg.ReplyToMessageID = update.Message.MessageID
					msg.Text = getCurrency(currency[i])

					bot.Send(msg)
					break
				}
			}
		}
	}
}
