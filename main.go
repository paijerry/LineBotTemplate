// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/zabawaba99/firego"
)

var bot *linebot.Client
var fire *firego.Firebase

func init() {
	firego.TimeoutDuration = time.Minute
	fire = firego.New("https://haru-line.firebaseio.com/", nil)
	fire.Auth("SnDhy01FPnNOzOyjNQzgDksX8WI2")
}

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/callback", callbackHandler)
	port := os.Getenv("PORT")
	addr := fmt.Sprintf(":%s", port)
	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)

	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}
	log.Print("ParseRequest success")

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "test" {

					a1 := linebot.NewMessageTemplateAction("say hello", "hello")
					a2 := linebot.NewURITemplateAction("google", "https://www.google.com.tw/")

					tmp := linebot.NewButtonsTemplate("https://static.pexels.com/photos/126407/pexels-photo-126407.jpeg", "test", "你好嗎", a1, a2)
					msg := linebot.NewTemplateMessage("Test", tmp)
					if _, err = bot.ReplyMessage(event.ReplyToken, msg).Do(); err != nil {
						log.Print(err)
					}
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("To "+event.Source.UserID+":"+event.Source.GroupID+":"+event.Source.RoomID+":"+message.Text+" 是嗎!?")).Do(); err != nil {
					log.Print(err)
				}
				if err = fire.Set(map[string]string{time.Now().Format(time.RFC3339): message.Text}); err != nil {
					log.Print(err)
				}
			}
		}
	}
}
