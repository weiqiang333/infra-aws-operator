package telegram

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"infra-aws-operator/pkg/utils/str"
)

func SendMessage(botToken, chatId, content string) error {
	content = str.EscapeStrings(content, []string{"_"})
	data := fmt.Sprintf(`{
        "parse_mode": "Markdown",
		"chat_id": "%s",
		"text": "%s",
    }`, chatId, content)
	bodys := strings.NewReader(data)
	fmt.Println(bodys)
	resp, err := http.Post(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken), "application/json", bodys)
	if err != nil {
		log.Println("Failed http.Post", http.StatusInternalServerError, err.Error())
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		body, _ := ioutil.ReadAll(resp.Body)
		msg := fmt.Sprintf("Failed http.Post telegram SendMessage StatusCode is %v, body %s", resp.StatusCode, string(body))
		log.Println(msg)
		return fmt.Errorf(msg)
	}
	log.Println("INFO http.Post telegram SendMessage", resp.StatusCode)
	return nil
}
