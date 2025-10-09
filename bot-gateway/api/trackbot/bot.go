package trackbot

import (
	"context"
	"log"
	"math"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackBot struct {
	*tgbotapi.BotAPI
	Client pb.UserServiceClient
}

const helpMsg = "/start - starts the bot;\n" +
	"/add_track - adds track number(s);\n" +
	"/stop - This stops notifications;\n" +
	"/help - help list;\n"

// New bot
func New(token string, client pb.UserServiceClient) (*TrackBot, error) {
	bot, err := tgbotapi.NewBotAPI("8286937197:AAFrfcaG_g_s1Sw5YZKUVgbtxyWbC9M8LWc")
	if err != nil {
		return nil, errCantCreate(err)
	}

	bot.Debug = true

	log.Printf("Authorized account %s", bot.Self.UserName)

	return &TrackBot{bot, client}, nil

}

// starts bot and handles commands
func (b *TrackBot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)

	for update := range updates {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		//some checks
		if update.Message == nil {
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		//router
		switch update.Message.Command() {

		case "start":
			resp, err := b.Client.CreateUser(ctx, &pb.UserRequest{
				Id:       update.Message.From.ID,
				Name:     update.Message.From.UserName,
				IsActive: true,
			})

			if err != nil {
				log.Fatal(err)
			}

			msg.Text = resp.Resp

		case "add_track":

		case "stop":

		case "help":
			msg.Text = helpMsg

		default:
			msg.Text = "Unknown command, use /help to list all possible commands"
		}

		for i := range 3 {
			time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
			if _, err := b.Send(msg); err != nil {
				log.Printf("Cannot send message %s", err.Error())
			} else {
				break
			}
		}

	}
	return nil
}
