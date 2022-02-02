package main

import (
	"context"
	"encoding/json"
	"github.com/scjtqs2/bot_adapter/event"
	"github.com/scjtqs2/bot_adapter/pb/entity"
	"github.com/tidwall/gjson"
)

func parseMsg(data string) {
	msg := gjson.Parse(data)
	switch msg.Get("post_type").String() {
	case "message": // 消息事件
		switch msg.Get("message_type").String() {
		case event.MESSAGE_TYPE_PRIVATE:
			var req event.MessagePrivate
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.MESSAGE_TYPE_GROUP:
			var req event.MessageGroup
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		}
	case "notice": // 通知事件
		switch msg.Get("notice_type").String() {
		case event.NOTICE_TYPE_FRIEND_ADD:
			var req event.NoticeFriendAdd
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_FRIEND_RECALL:
			var req event.NoticeFriendRecall
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_BAN:
			var req event.NoticeGroupBan
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_DECREASE:
			var req event.NoticeGroupDecrease
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_INCREASE:
			var req event.NoticeGroupIncrease
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_ADMIN:
			var req event.NoticeGroupAdmin
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_RECALL:
			var req event.NoticeGroupRecall
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_GROUP_UPLOAD:
			var req event.NoticeGroupUpload
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_POKE:
			var req event.NoticePoke
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_HONOR:
			var req event.NoticeHonor
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.NOTICE_TYPE_LUCKY_KING:
			var req event.NoticeLuckyKing
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.CUSTOM_NOTICE_TYPE_GROUP_CARD:
		case event.CUSTOM_NOTICE_TYPE_OFFLINE_FILE:
		}
	case "request": // 请求事件
		switch msg.Get("request_type").String() {
		case event.REQUEST_TYPE_FRIEND:
			var req event.RequestFriend
			_ = json.Unmarshal([]byte(msg.Raw), &req)
			apploveFriendRequest(req.Flag)
		case event.REQUEST_TYPE_GROUP:
			var req event.RequestGroup
			_ = json.Unmarshal([]byte(msg.Raw), &req)
			apploveGroupRequest(req.Flag, req.SubType)
		}
	case "meta_event": // 元事件
		switch msg.Get("meta_event_type").String() {
		case event.META_EVENT_LIFECYCLE:
			var req event.MetaEventLifecycle
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		case event.META_EVENT_HEARTBEAT:
			var req event.MetaEventHeartbeat
			_ = json.Unmarshal([]byte(msg.Raw), &req)
		}
	}
}

// apploveFriendRequest 通过好友请求
func apploveFriendRequest(flag string) {
	_, _ = bot_adapter_client.SetFriendAddRequest(context.TODO(), &entity.SetFriendAddRequestReq{Approve: true, Flag: flag})
}

// apploveGroupRequest 通过群组添加请求
func apploveGroupRequest(flag string, sub_type string) {
	_, _ = bot_adapter_client.SetGroupAddRequest(context.TODO(), &entity.SetGroupAddRequestReq{Approve: true, Flag: flag, SubType: sub_type})
}
