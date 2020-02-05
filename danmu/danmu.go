package danmu

import (
	"errors"
	"fmt"
	"github.com/x554462/danmu_douyu/client"
	"github.com/x554462/danmu_douyu/conf"
	"github.com/x554462/danmu_douyu/message"
	"github.com/x554462/danmu_douyu/util"
	"log"
	"os"
	"os/signal"
	"time"
)

const loopDuration = time.Second << 5

func checkRoom() error {
	log.Println("检查房间")
	resBody, err := client.HttpReq(client.HttpMethodGet, fmt.Sprintf(conf.RoomInfoUrl, conf.DefaultRoomId), nil)
	if err != nil {
		log.Println("http:", err)
		return err
	}
	var resMap map[string]interface{}
	err = util.JsonDecodeWithByte(resBody, &resMap)
	if err != nil {
		log.Fatal("错误: ", string(resBody))
	}
	errMsg := ""
	if resErr, ok := (resMap["error"]).(float64); ok {
		if resErr != 0 {
			log.Fatal("错误: ", "房间号不存在")
		} else {
			if data, ok := resMap["data"].(map[string]interface{}); ok {
				if roomStatus, ok := data["room_status"].(string); ok {
					if roomStatus == "2" {
						errMsg = "主播未开播"
					}
				} else {
					errMsg = "room status error"
				}
			} else {
				errMsg = "data error"
			}
		}
	} else {
		errMsg = "json error"
	}
	if errMsg != "" {
		log.Println("http:", errMsg)
		return errors.New(errMsg)
	}
	return nil
}

func pull() error {
	log.Println("连接弹幕服务器")
	wsClient := client.NewWebsocketClient(conf.Scheme, conf.SiteName, conf.Port)
	if err := wsClient.Connect(); err != nil {
		log.Println("dial:", err)
		return err
	}
	log.Println("拉取弹幕")
	wsClient.SendMsg(message.NewSendMsg(message.SendTypeLoginRoom, conf.DefaultRoomId).PackMsg())
	wsClient.SendMsg(message.NewSendMsg(message.SendTypeJoinRoom, conf.DefaultRoomId).PackMsg())
	wsClient.SetTickerFunc(func() bool {
		wsClient.SendMsg(message.NewSendMsg(message.SendTypeKeepLive, "").PackMsg())
		return true
	}, 45*time.Second)
	wsClient.OnReceive(func(bytes []byte) bool {
		if msg := message.Handle(string(bytes)); msg != nil {
			fmt.Printf("%+v\n", msg)
		}
		return true
	})
	wsClient.OnClose(func() {
		loop(loopDuration, false)
	})
	return nil
}

var interrupt = make(chan os.Signal, 1)

func loop(duration time.Duration, first bool) {
	timer := time.NewTimer(0)
	defer timer.Stop()
	for {
		select {
		case <-timer.C:
			if err := checkRoom(); err == nil {
				if err = pull(); err == nil {
					if first {
						signal.Notify(interrupt, os.Interrupt)
						<-interrupt
					}
					return
				}
			}
		case <-interrupt:
			return
		}
		timer.Reset(duration)
	}
}

func Run() {
	loop(loopDuration, true)
}
