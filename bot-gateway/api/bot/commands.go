package bot

import (
	"context"
	"log"
	"strings"

	"github.com/Strelcock/pb/bot/pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	helpMsg = "/start - starts the bot;\n" +
		"/add_track - adds track number(s);\n" +
		"/stop - This stops notifications;\n" +
		"/help - help list;\n"
	unknownMsg = "Unknown command, use /help to list all possible commands"
)

var nilMsg = tgbotapi.MessageConfig{}

var adminKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Место для вашей админки", "Братан, ты че админ?"),
	),
)

func (b *Bot) startCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	resp, err := b.UserClient.CreateUser(ctx, &pb.UserRequest{
		Id:       update.Message.From.ID,
		Name:     update.Message.From.UserName,
		IsActive: true,
	})

	if err != nil {
		return nilMsg, err
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, resp.Resp)
	log.Println("STARTUEM ", update.Message.Text)
	msg.ReplyMarkup = b.commands
	log.Println("ESLI NE OPEN TO HOOEVO")
	return msg, nil
}

func (b *Bot) addCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatID := update.Message.Chat.ID

	if update.Message.IsCommand() || update.Message.Text == add {
		text := "Введите номера посылок через запятую\n- чтобы ничего не добавлять"
		msg := tgbotapi.NewMessage(chatID, text)
		b.botMap.mu.Lock()
		defer b.botMap.mu.Unlock()
		b.botMap.waitForInput[chatID] = true

		return msg, nil
	}
	msgs := strings.Split(update.Message.Text, ",")

	msg := tgbotapi.NewMessage(chatID, "")

	if msgs[0] == "-" {
		msg.ReplyMarkup = b.commands
		msg.Text = "Работаем дальше"
		b.botMap.mu.Lock()
		defer b.botMap.mu.Unlock()
		b.botMap.waitForInput[chatID] = false

		return msg, nil
	}

	for _, num := range msgs {
		resp, err := b.TrackClient.AddTrack(ctx, &pb.TrackRequest{
			Number: num,
			User:   update.Message.From.ID,
		})

		if err != nil {
			log.Println(err, resp)
			errMsg := strings.Split(err.Error(), "=") //take an error info from rpc
			msg.Text += errMsg[2]                     // error
			continue
		}

		msg.Text += resp.Status

	}

	b.botMap.mu.Lock()
	b.botMap.waitForInput[chatID] = false
	b.botMap.mu.Unlock()
	msg.ReplyMarkup = b.commands
	return msg, nil
}

func (b *Bot) adminCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatID := update.Message.Chat.ID
	msg := tgbotapi.NewMessage(chatID, "")

	resp, err := b.UserClient.IsAdmin(ctx, &pb.AdminRequest{
		Id: update.Message.From.ID,
	})

	if err != nil {
		return nilMsg, err
	}

	if resp.IsAdmin {
		msg.ReplyMarkup = adminKeyboard
		msg.Text = "Admin panel"
	} else {
		msg.Text = unknownMsg
	}

	return msg, nil
}

func (b *Bot) routeSlashCommands(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatId := update.Message.Chat.ID
	switch update.Message.Command() {
	case startCmd:

		return b.startCommand(ctx, update)
	case admin:
		return b.adminCommand(ctx, update)
	case addCmd:
		return b.addCommand(ctx, update)

	case stopCmd:
		msg := tgbotapi.NewMessage(chatId, "Функция пока не реализована, гуляем")
		return msg, nil

	case helpCmd:
		msg := tgbotapi.NewMessage(chatId, helpMsg)
		return msg, nil
	default:
		msg := tgbotapi.NewMessage(chatId, unknownMsg)
		return msg, nil
	}
}

func (b *Bot) routeNonSlashCommands(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatId := update.Message.Chat.ID
	switch update.Message.Text {

	case add:
		return b.addCommand(ctx, update)

	case stop:
		msg := tgbotapi.NewMessage(chatId, "Функция пока не реализована, гуляем")
		return msg, nil

	case info:
		msg := tgbotapi.NewMessage(chatId, "Функция пока не реализована, гуляем")
		return msg, nil

	case help:
		msg := tgbotapi.NewMessage(chatId, helpMsg)
		return msg, nil

	}
	return nilMsg, nil
}
