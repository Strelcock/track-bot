package trackbot

import (
	"context"
	"fmt"

	"github.com/Strelcock/pb/bot/pb"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

var adminKeyboard = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Место для вашей админки", "Братан, ты че админ?"),
	),
)

func (b *TrackBot) startCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
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

func (b *TrackBot) addCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
	chatID := update.Message.Chat.ID
	if ctx.Value(numbers) == nil {
		text := "Введите номера посылок через запятую (,):"
		msg := tgbotapi.NewMessage(chatID, text)
		b.waitForInput[chatID] = true
		return msg, nil
	}

	resp, err := b.TrackClient.AddTrack(ctx, &pb.TrackRequest{
		Number: ctx.Value(numbers).([]string),
		User:   update.Message.From.ID,
	})

	if err != nil {
		return nilMsg, err
	}

	text := fmt.Sprintf("Добавлены заказы %v", resp.Status)
	msg := tgbotapi.NewMessage(chatID, text)
	b.waitForInput[chatID] = false
	return msg, nil
}

func (b *TrackBot) adminCommand(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
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

func (b *TrackBot) route(ctx context.Context, update tgbotapi.Update) (tgbotapi.MessageConfig, error) {
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
