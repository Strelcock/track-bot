package trackbot

import (
	"context"
	"log"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackBot struct {
	*tgbotapi.BotAPI
	Client pb.UserServiceClient
}

const (
	helpMsg = "/start - starts the bot;\n" +
		"/add_track - adds track number(s);\n" +
		"/stop - This stops notifications;\n" +
		"/help - help list;\n"
	uknownMsg = "Unknown command, use /help to list all possible commands"
)

var adminKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Место для вашей админки", "Братан, ты че админ?"),
	),
)

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
func (b *TrackBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	b.Hadnle(updates)

}

func (b *TrackBot) Hadnle(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		func() {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			//some checks
			if update.CallbackQuery != nil {
				callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
				if _, err := b.Request(callback); err != nil {
					log.Fatal(err)
				}
				msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
				if _, err := b.Send(msg); err != nil {
					log.Fatal(err)
				}
				return
			}

			if !update.Message.IsCommand() {
				return
			}

			if update.Message == nil {
				return
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

			case "admin":
				resp, err := b.Client.IsAdmin(ctx, &pb.AdminRequest{
					Id: update.Message.From.ID,
				})
				if err != nil {
					log.Fatal(err)
				}
				if resp.IsAdmin {
					msg.ReplyMarkup = adminKeyboard
					msg.Text = "Admin panel"
				} else {
					msg.Text = uknownMsg
				}

			default:
				msg.Text = uknownMsg
			}

			_, err := b.Send(msg)
			if err != nil {
				log.Fatal(err)
			}

			// for i := range 3 {
			// 	time.Sleep(time.Second * time.Duration(math.Pow(2, float64(i))))
			// 	if _, err := b.Send(msg); err != nil {
			// 		log.Printf("Cannot send message: %s", err.Error())
			// 	} else {
			// 		break
			// 	}
		}()
	}
}
