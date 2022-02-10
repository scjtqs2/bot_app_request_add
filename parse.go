package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/scjtqs2/bot_adapter/coolq"
	"github.com/scjtqs2/bot_adapter/event"
	"github.com/scjtqs2/bot_adapter/pb/entity"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"strconv"
	"time"
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
		case event.NOTICE_TYPE_GROUP_DECREASE: // 群成员减少
			var req event.NoticeGroupDecrease
			_ = json.Unmarshal([]byte(msg.Raw), &req)
			groupDecrease(req)
		case event.NOTICE_TYPE_GROUP_INCREASE: // 群成员增加
			var req event.NoticeGroupIncrease
			_ = json.Unmarshal([]byte(msg.Raw), &req)
			groupIncrease(req)
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
	_, _ = botAdapterClient.SetFriendAddRequest(context.TODO(), &entity.SetFriendAddRequestReq{Approve: true, Flag: flag})
}

// apploveGroupRequest 通过群组添加请求
func apploveGroupRequest(flag string, sub_type string) {
	_, _ = botAdapterClient.SetGroupAddRequest(context.TODO(), &entity.SetGroupAddRequestReq{Approve: true, Flag: flag, SubType: sub_type})
}

// groupIncrease 回应群成员增加
func groupIncrease(req event.NoticeGroupIncrease) {
	at := coolq.EnAtCode(strconv.FormatInt(req.UserID, 10))
	img := coolq.EnImageCode(fmt.Sprintf("http://q1.qlogo.cn/g?b=qq&nk=%d&s=100", req.UserID), 1)
	groupName, err := getGroupName(req.GroupID)
	if err != nil {
		log.Errorf("获取群信息错误：err:%v", err)
		return
	}
	nowDate := time.Now().Format("2006-01-02 15:04:05")
	message := fmt.Sprintf("尊敬的 %s 您好~ \n %s \n 欢迎加入%s~ \n 本群的群号是 %d~ \n 已经赞你了哦！\n  请看公告遵守相关规定~ \n 当前时间 %s",
		at,
		img,
		groupName,
		req.GroupID,
		nowDate,
	)
	_, _ = botAdapterClient.SendGroupMsg(context.TODO(), &entity.SendGroupMsgReq{
		GroupId: req.GroupID,
		Message: []byte(message),
	})
}

// groupDecrease 回应群成员减少
func groupDecrease(req event.NoticeGroupDecrease) {
	// at := coolq.EnAtCode(strconv.FormatInt(req.UserID, 10))
	img := coolq.EnImageCode(fmt.Sprintf("http://q1.qlogo.cn/g?b=qq&nk=%d&s=100", req.UserID), 1)
	groupName, err := getGroupName(req.GroupID)
	if err != nil {
		log.Errorf("获取群信息错误：err:%v", err)
		return
	}
	nowDate := time.Now().Format("2006-01-02 15:04:05")
	var message string
	switch req.SubType {
	case "kick": // 被踢
		operatorName, err := getMemberNickName(req.GroupID, req.OperatorID)
		if err != nil {
			log.Errorf("获取群成员信息失败,group:%d,userid:%d,err:%v", req.GroupID, req.OperatorID, err)
			return
		}
		message = fmt.Sprintf("用户 %d \n %s \n 被管理员 %s 移除了群 %s \n 当前时间： %s", req.UserID, img, operatorName, groupName, nowDate)
	case "leave": // 主动退群
		message = fmt.Sprintf("用户 %d \n %s \n 主动离开了群 %s \n 当前时间： %s", req.UserID, img, groupName, nowDate)
	default:
		return
	}
	_, _ = botAdapterClient.SendGroupMsg(context.TODO(), &entity.SendGroupMsgReq{
		GroupId: req.GroupID,
		Message: []byte(message),
	})
}

// getGroupName 通过群号码查群名称
func getGroupName(groupID int64) (string, error) {
	groupInfo, err := botAdapterClient.GetGroupInfo(context.TODO(), &entity.GetGroupInfoReq{GroupId: groupID})
	if err != nil {
		return "", err
	}
	return groupInfo.GroupName, nil
}

// getMemberNickName 通过qq号和群号，查 群成员昵称
func getMemberNickName(groupID, userID int64) (string, error) {
	memberinfo, err := botAdapterClient.GetGroupMemberInfo(context.TODO(), &entity.GetGroupMemberInfoReq{GroupId: groupID, UserId: userID, NoCache: false})
	if err != nil {
		return "", err
	}
	return memberinfo.Nickname, nil
}
