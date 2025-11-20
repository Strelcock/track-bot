package bot

import (
	"context"
	"log"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/Strelcock/pb/bot/pb"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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

var commandKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(add),
		tgbotapi.NewKeyboardButton(stop),
		tgbotapi.NewKeyboardButton(info),
	),
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton(help),
	),
)

var commandList = []string{add, stop, info, help}

type Bot struct {
	*tgbotapi.BotAPI
	UserClient  pb.UserServiceClient
	TrackClient pb.TrackServiceClient
	waitMap     *waitMap
	infoMap     *infoMap
}

// New bot
func New(token string, userClient pb.UserServiceClient, trackClient pb.TrackServiceClient) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, errCantCreate(err)
	}

	bot.Debug = true

	log.Printf("Authorized account %s", bot.Self.UserName)

	waitMap := newWaitMap()
	infoMap := NewInfoMap()

	return &Bot{bot, userClient, trackClient, waitMap, infoMap}, nil

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
			timerCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			//handle callbacks
			if update.CallbackQuery != nil {
				err := b.HandleCallback(update)
				if err != nil {
					log.Print(err)
				}
				return
			}

			if update.Message == nil {
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
			}

			if b.waitMap.waitForInput[update.Message.Chat.ID] {
				msg, err := b.addCommand(timerCtx, update)
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
	log.Printf("\n\ncallback %s\n\n", update.CallbackQuery.Data)

	callbackData := update.CallbackQuery.Data
	if strings.HasPrefix(callbackData, "page:") {
		pageStr := strings.TrimPrefix(callbackData, "page:")
		page, err := strconv.Atoi(pageStr)
		if err != nil {
			log.Println("handleCallback: ", err)
			return err
		}
		keyboard := generateReplyKeyboard(b.infoMap.Get(update.CallbackQuery.From.ID), page)
		msg := tgbotapi.NewEditMessageTextAndMarkup(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID, update.CallbackQuery.Message.Text, keyboard)

		if _, err := b.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
