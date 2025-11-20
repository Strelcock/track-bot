package bot

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func generateReplyKeyboard(nums []string, curPage int) tgbotapi.InlineKeyboardMarkup {
	pages := float64(len(nums)) / 6.0
	var endOffset int

	if curPage*6+6 > len(nums) {
		endOffset = len(nums) - 6*curPage
	} else {
		endOffset = 6
	}

	pageBtns := [][]tgbotapi.InlineKeyboardButton{}

	slicePage := nums[curPage*6 : curPage*6+endOffset]

	cutPages := CutSlice(slicePage, 3)

	for i := range cutPages {
		row := tgbotapi.NewInlineKeyboardRow()
		for _, num := range cutPages[i] {
			row = append(row, tgbotapi.NewInlineKeyboardButtonData(num, num))
		}
		pageBtns = append(pageBtns, row)
	}

	if pages > 1 && curPage == 0 {
		pageBtns = append(pageBtns, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Далее -->", fmt.Sprintf("page:%d", curPage+1)),
		))
	}

	if curPage > 0 && float64(curPage+1) < pages {
		pageBtns = append(pageBtns, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Далее -->", fmt.Sprintf("page:%d", curPage+1)),
		))
		pageBtns = append(pageBtns, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад <--", fmt.Sprintf("page:%d", curPage-1)),
		))
	}

	if curPage > 0 && float64(curPage+1) >= pages {
		pageBtns = append(pageBtns, tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Назад <--", fmt.Sprintf("page:%d", curPage-1)),
		))
	}
	replyKeyboard := tgbotapi.NewInlineKeyboardMarkup(pageBtns...)
	return replyKeyboard
}

func CutSlice(slice []string, cut int) [][]string {
	res := [][]string{}
	var offset int
	if len(slice) < cut {
		return [][]string{slice}
	}
	if len(slice)-cut < cut {

	}
	for i := 0; i < len(slice); i += cut {
		row := []string{}
		if len(slice) < i+cut {
			offset = len(slice)
		} else {
			offset = i + cut
		}
		for j := i; j < offset; j++ {
			row = append(row, slice[j])

		}
		res = append(res, row)
	}
	return res
}
