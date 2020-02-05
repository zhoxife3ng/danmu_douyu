package message

import (
	"regexp"
	"strconv"
)

const (
	RecvTypeGift         = "spbc"        // 全站礼物广播
	RecvTypeChatMsg      = "chatmsg"     // 弹幕消息
	RecvTypeUserEnter    = "uenter"      // 进入房间
	RecvTypeShareRoom    = "srres"       // 分享房间
	RecvTypeUserLevelUp  = "upgrade"     // 用户等级
	RecvTypeSuperChatMsg = "ssd"         // 超级弹幕
	RecvTypeBanned       = "newblackres" // 禁言
	RecvTypeError        = "error"       // 自定义error
)

var (
	regexpUid, _              = regexp.Compile("/uid@=([0-9]+)/")
	regexpType, _             = regexp.Compile("type@=(.*?)/")
	regexpTypeGift, _         = regexp.Compile("/sn@=(.*?)/dn@=(.*?)/gn@=(.*?)/gc@=(.*?)/")
	regexpTypeChatMsg, _      = regexp.Compile("/nn@=([^/]*?)/txt@=([^/]*?)/")
	regexpTypeUserEnter, _    = regexp.Compile("/nn@=(.*?)/")
	regexpTypeShareRoom, _    = regexp.Compile("/nickname@=(.*?)/")
	regexpTypeUserLevelUp, _  = regexp.Compile("/nn@=(.*?)/level@=(.*?)/")
	regexpTypeSuperChatMsg, _ = regexp.Compile("/content@=(.*?)/")
	regexpTypeBanned, _       = regexp.Compile("/snic@=(.*?)/dnic@=(.*?)/")
)

type RecvMsg struct {
	Type string
	Data []string
	Uid  int64
}

func Handle(msg string) *RecvMsg {
	match := regexpType.FindStringSubmatch(msg)
	if len(match) < 1 {
		return nil
	}
	typo := match[1]

	message := &RecvMsg{Type: typo}

	match = regexpUid.FindStringSubmatch(msg)
	if len(match) > 1 {
		if uid, err := strconv.ParseInt(match[1], 10, 64); err == nil {
			message.Uid = uid
		}
	}

	switch typo {
	case RecvTypeGift:
		if m := regexpTypeGift.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeChatMsg:
		if m := regexpTypeChatMsg.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeUserEnter:
		if m := regexpTypeUserEnter.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeShareRoom:
		if m := regexpTypeShareRoom.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeUserLevelUp:
		if m := regexpTypeUserLevelUp.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeSuperChatMsg:
		if m := regexpTypeSuperChatMsg.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	case RecvTypeBanned:
		if m := regexpTypeBanned.FindStringSubmatch(msg); len(m) > 1 {
			message.Data = m[1:]
		}
	default:
		return nil
		//message.Data = []string{msg}
	}
	return message
}
