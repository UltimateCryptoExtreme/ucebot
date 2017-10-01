package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/telegram-bot-api.v4"
)

type Exchange interface {
	Name() string
	Init()
	CurrencyMarket(currency string) string
}

type PriceSource interface {
	Price(market string) (Price, error)
}

type Price struct {
	Source     string
	Quote      string
	Base       string
	Price      float64
	Change24h  *float64
	BaseVolume float64
	Low        float64
	High       float64
	LowestAsk  float64
	HighestBid float64
}

var (
	cmdRegex     = regexp.MustCompile(`\/([^ @]*)[^ ]* (.*)`)
	invalidChars = regexp.MustCompile(`[^a-zA-Z0-9\-]`)
)

var bot *tgbotapi.BotAPI

var exchanges = []Exchange{
	&Poloniex{},
	&Bittrex{},
	&Cryptopia{},
	&Yobit{},
	&Liqui{},
}

func main() {
	for _, e := range exchanges {
		e.Init()
	}
	var err error
	bot, err = tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	//bot.Debug = true

	log.Printf("Bot \"%s\" started...", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 300

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		panic(err)
	}

	log.Println("Listening for updates.")

	for update := range updates {
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}
		cmd, text := parseCmd(update.Message.Text)
		switch cmd {
		case "p", "price":
			ps := prices(text)
			if len(ps) == 0 {
				text = "A veri bad error occured!"
			} else {
				text = formatPrices(ps)
			}
			sendMsg(update.Message.Chat.ID, text)
		}

	}
}

func parseCmd(msg string) (string, string) {
	finds := cmdRegex.FindStringSubmatch(msg)
	if len(finds) != 3 {
		return "", ""
	}
	return finds[1], finds[2]
}

func prices(market string) []Price {
	market = invalidChars.ReplaceAllLiteralString(market, "")
	market = strings.ToUpper(market)

	log.Printf("Price: %s", market)

	prices := make([]Price, 0)

	var wg sync.WaitGroup
	pc := make(chan Price)
	for _, e := range exchanges {
		ps, ok := e.(PriceSource)
		if !ok {
			log.Printf("Exchange %s can not be used for fetching prices.", e.Name())
			continue
		}
		wg.Add(1)
		go func() {
			if p, err := ps.Price(e.CurrencyMarket(market)); err == nil {
				pc <- p
			} else {
				log.Printf("%v", err)
			}
			wg.Done()
		}()
	}

	go func() {
		wg.Wait()
		close(pc)
	}()

	var maxVol float64
	for p := range pc {
		prices = append(prices, p)
		if p.BaseVolume > maxVol {
			maxVol = p.BaseVolume
		}
	}
	volLimit := maxVol * 0.10

	pricesFiltered := make([]Price, 0)
	for _, p := range prices {
		if p.BaseVolume > volLimit {
			pricesFiltered = append(pricesFiltered, p)
		}
	}

	return pricesFiltered
}

func formatPrices(prices []Price) string {
	s := ""
	for _, p := range prices {
		if s != "" {
			s += "\n\n"
		}
		s += formatPrice(p)
	}
	return s
}

func formatPrice(price Price) string {
	s := fmt.Sprintf("*%s:* `%.8f` %s",
		price.Source,
		price.Price,
		price.Base)
	if price.Change24h != nil {
		changeSign := ""
		if *price.Change24h > 0 {
			changeSign = "+"
		}
		s += fmt.Sprintf(" | `%s%.2f%%` (24h)",
			changeSign,
			*price.Change24h)
	}
	s += fmt.Sprintf("\n*Low:* `%.8f` | *High:* `%.8f`\n"+
		"*Ask:* `%.8f` | *Bid:* `%.8f`\n"+
		"*Vol:* `%f` %s",
		price.Low, price.High,
		price.LowestAsk, price.HighestBid,
		price.BaseVolume, price.Base)
	return s
}

func sendMsg(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "markdown"
	bot.Send(msg)
}
