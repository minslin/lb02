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
	"strings"
	"bytes"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

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

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				var outmsg bytes.Buffer

				if strings.Compare(message.Text, "溫馨提醒") == 0 {
					outmsg.WriteString("<<<溫馨提醒>>>\r\n因為這個群很吵 -->\r\n右上角 可以 關閉提醒\r\n\r\n[同學會] 投票進行中 -->\r\n右上角 筆記本 可以進行投票\r\n\r\n[通訊錄] 需要大家的協助 -->\r\n右上角 筆記本 請更新自己的聯絡方式")
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outmsg.String())).Do(); err != nil {
					log.Print(err)
					}
				}
				
				if strings.HasSuffix(message.Text, "還是那麼帥") {
					outmsg.WriteString(message.Text)
					outmsg.WriteString("+1")
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outmsg.String())).Do(); err != nil {
					log.Print(err)
					}
				}
				
				if strings.Compare(message.Text, "PPAP") == 0 {
					outmsg.WriteString(GetPPAPText())
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(outmsg.String())).Do(); err != nil {
					log.Print(err)
					}
				}
			}
		}
	}
}

func GetPPAPText() string {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(100)
	switch i % 5 {
	case 0:
		return "I have a pencil,\r\nI have an Apple,\r\nApple pencil.\r\nI have a watch,\r\nI have an Apple,\r\nApple watch."
	case 1:
		return "順帶一提，請不要把Apple Pencil刺進水果裡，不管是蘋果還是鳳梨。"
	case 2:
		return "我懂了，這是以書寫工具與種類食物為題的饒舌歌。"
	case 3:
		return "我不太清楚PPAP是什麼，但你可以問我AAPL的相關資訊。"
	case 4:
		return "我是不會接著唱的！"
	}
}
