package bot

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tracksKey string

const numbers tracksKey = "numbers"

type Bot struct {
	*tgbotapi.BotAPI
	UserClient   pb.UserServiceClient
	TrackClient  pb.TrackServiceClient
	waitForInput map[int64]bool
}

const (
	helpMsg = "/start - starts the bot;\n" +
		"/add_track - adds track number(s);\n" +
		"/stop - This stops notifications;\n" +
		"/help - help list;\n"
	unknownMsg = "Unknown command, use /help to list all possible commands"
)

var nilMsg = tgbotapi.MessageConfig{}

// New bot
func New(token string, userClient pb.UserServiceClient, trackClient pb.TrackServiceClient) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errCantCreate(err)
	}

	bot.Debug = true

	log.Printf("Authorized account %s", bot.Self.UserName)
	input := make(map[int64]bool)

	return &Bot{bot, userClient, trackClient, input}, nil

}

// starts bot and handles commands
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	b.Hadnle(updates)

}

func (b *Bot) Hadnle(updates tgbotapi.UpdatesChannel) {
	var tracks = []string{}
	for update := range updates {
		go func() {

			timerCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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
				err := b.HandleCommands(timerCtx, update)
				if err != nil {
					log.Print(err)
				}
				return
			}

			if !update.Message.IsCommand() {
				tracks = strings.Split(update.Message.Text, ",")
				ctx := context.WithValue(timerCtx, numbers, tracks)
				if b.waitForInput[update.Message.Chat.ID] {
					msg, err := b.addCommand(ctx, update)
					if err != nil {
						log.Print(err)
						return
					}
					_, err = b.Send(msg)
					if err != nil {
						log.Print(err)
					}
					return
				}
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

func (b *Bot) HandleCommands(ctx context.Context, update tgbotapi.Update) error {
	//router
	msg, err := b.route(ctx, update)
	if err != nil {
		return err
	}

	_, err = b.Send(msg)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) HandleCallback(update tgbotapi.Update) error {
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
