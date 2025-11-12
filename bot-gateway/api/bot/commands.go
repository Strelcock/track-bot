package bot

import (
	"context"
	"log"
	"strings"

	"github.com/Strelcock/pb/bot/pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	return msg, nil
}

func (b *Bot) addCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatID := update.Message.Chat.ID
	if ctx.Value(numbers) == nil {
		text := "Введите номера посылок через запятую (,):"
		msg := tgbotapi.NewMessage(chatID, text)
		b.botMap.mu.Lock()
		defer b.botMap.mu.Unlock()
		b.botMap.waitForInput[chatID] = true
		return msg, nil
	}

	msg := tgbotapi.NewMessage(chatID, "")
	for _, num := range ctx.Value(numbers).([]string) {
		resp, err := b.TrackClient.AddTrack(ctx, &pb.TrackRequest{
			Number: num,
			User:   update.Message.From.ID,
		})

		if err != nil {
			log.Println(err, resp)
			errMsg := strings.Split(err.Error(), "=")
			msg.Text += errMsg[2]
			continue
		}

		msg.Text += resp.Status

	}

	b.botMap.mu.Lock()
	b.botMap.waitForInput[chatID] = false
	b.botMap.mu.Unlock()
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

func (b *Bot) route(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatId := update.Message.Chat.ID
	switch update.Message.Command() {

	case "start":
		return b.startCommand(ctx, update)

	case "add_track":
		return b.addCommand(ctx, update)

	case "stop":

	case "help":
		msg := tgbotapi.NewMessage(chatId, helpMsg)
		return msg, nil

	case "admin":
		return b.adminCommand(ctx, update)
	default:
		msg := tgbotapi.NewMessage(chatId, unknownMsg)
		return msg, nil
	}
	return nilMsg, nil
}
