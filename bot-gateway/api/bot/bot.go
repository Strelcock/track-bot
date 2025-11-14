package bot

import (
	"context"
	"log"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type tracksKey string

const numbers tracksKey = "numbers"

const (
	add  = "Добавить посылку"
	stop = "Остановить рассылку"
	info = "Где моя посылка?"
	help = "Помощь"
)

const (
	startCmd = "start"
	admin    = "admin"
	addCmd   = "add_track"
	stopCmd  = "stop"
	helpCmd  = "help"
)

var commandList = []string{add, stop, info, help}

type botMap struct {
	mu           sync.RWMutex
	waitForInput map[int64]bool
}

type Bot struct {
	*tgbotapi.BotAPI
	UserClient  pb.UserServiceClient
	TrackClient pb.TrackServiceClient
	botMap      *botMap
	commands    tgbotapi.ReplyKeyboardMarkup
}

// New bot
func New(token string, userClient pb.UserServiceClient, trackClient pb.TrackServiceClient) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errCantCreate(err)
	}

	bot.Debug = true

	log.Printf("Authorized account %s", bot.Self.UserName)
	waitForInput := make(map[int64]bool)
	commandKeyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(add),
			tgbotapi.NewKeyboardButton(stop),
			tgbotapi.NewKeyboardButton(info),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(help),
		),
	)
	input := &botMap{waitForInput: waitForInput}

	return &Bot{bot, userClient, trackClient, input, commandKeyboard}, nil

}

// starts bot and handles commands
func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.GetUpdatesChan(u)
	b.Hadnle(updates)

}

func (b *Bot) Hadnle(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		go func() {
			var tracks = []string{}

			if update.Message == nil {
				return
			}

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
				err := b.HandleSlashCommands(timerCtx, update)
				if err != nil {
					log.Print(err)
				}
				return
			}

			if slices.Contains(commandList, strings.TrimSpace(update.Message.Text)) {

				err := b.HandleNonSlashCommands(timerCtx, update)
				if err != nil {
					log.Print(err)
				}
				return
			} else {

				tracks = strings.Split(update.Message.Text, ",")
				ctx := context.WithValue(timerCtx, numbers, tracks)
				if b.botMap.waitForInput[update.Message.Chat.ID] {
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
			}

		}()
	}
}

func (b *Bot) HandleSlashCommands(ctx context.Context, update tgbotapi.Update) error {
	//router
	msg, err := b.routeSlashCommands(ctx, update)
	if err != nil {
		return err
	}

	_, err = b.Send(msg)

	if err != nil {
		return err
	}
	return nil
}

func (b *Bot) HandleNonSlashCommands(ctx context.Context, update tgbotapi.Update) error {
	//router
	msg, err := b.routeNonSlashCommands(ctx, update)
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
