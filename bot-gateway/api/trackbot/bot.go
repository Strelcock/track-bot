package trackbot

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tracksKey string

const numbers tracksKey = "numbers"

type TrackBot struct {
	*tgbotapi.BotAPI
	UserClient  pb.UserServiceClient
	TrackClient pb.TrackServiceClient
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
func New(token string, userClient pb.UserServiceClient, trackClient pb.TrackServiceClient) (*TrackBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errCantCreate(err)
	}

	bot.Debug = true

	log.Printf("Authorized account %s", bot.Self.UserName)

	return &TrackBot{bot, userClient, trackClient}, nil

}

// starts bot and handles commands
func (b *TrackBot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	b.Hadnle(updates)

}

func (b *TrackBot) Hadnle(updates tgbotapi.UpdatesChannel) {
	var tracks = []string{}
	for update := range updates {
		func() {
			valueCtx := context.WithValue(context.Background(), numbers, tracks)
			ctx, cancel := context.WithTimeout(valueCtx, 2*time.Second)
			defer cancel()

			//handle callbacks
			if update.CallbackQuery != nil {
				err := b.HandleCallback(update)
				if err != nil {
					log.Print(err)
				}
				return
			}

			//handle commands
			if update.Message.IsCommand() {
				tracks = []string{}
				err := b.HandleCommands(ctx, update)
				if err != nil {
					log.Print(err)
				}
				return
			}

			if !update.Message.IsCommand() {
				tracks = strings.Split(update.Message.Text, ",")
				// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
				// b.Send(msg)
			}

			if update.Message == nil {
				return
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

func (b *TrackBot) HandleCommands(ctx context.Context, update tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
	//router
	switch update.Message.Command() {

	case "start":
		resp, err := b.UserClient.CreateUser(ctx, &pb.UserRequest{
			Id:       update.Message.From.ID,
			Name:     update.Message.From.UserName,
			IsActive: true,
		})

		if err != nil {
			return err
		}

		msg.Text = resp.Resp

	case "add_track":
		if len(ctx.Value(numbers).([]string)) == 0 {
			msg.Text = "Введите номера посылок через запятую (,):"
			break
		}

		resp, err := b.TrackClient.AddTrack(ctx, &pb.TrackRequest{
			Number: ctx.Value(numbers).([]string),
			User:   update.Message.From.ID,
		})
		if err != nil {
			return err
		}

		msg.Text = fmt.Sprintf("Добавлены заказы %v", resp.Number)

	case "stop":

	case "help":
		msg.Text = helpMsg

	case "admin":
		resp, err := b.UserClient.IsAdmin(ctx, &pb.AdminRequest{
			Id: update.Message.From.ID,
		})
		if err != nil {
			return err
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
		return err
	}
	return nil
}

func (b *TrackBot) HandleCallback(update tgbotapi.Update) error {
	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, update.CallbackQuery.Data)
	if _, err := b.Request(callback); err != nil {
		return err
	}

	msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Data)
	if _, err := b.Send(msg); err != nil {
		return err
	}

	return nil
}
