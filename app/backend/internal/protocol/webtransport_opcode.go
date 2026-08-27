package protocol

import "api/internal/server/wts"

type RequestMessageType = wts.MessageType

const (
	RequestTypeNone RequestMessageType = iota
	RequestTypePing
	RequestTypeHello
	RequestTypeChat
)

type ResponseMessageType = wts.MessageType

const (
	ResponseTypeNone ResponseMessageType = iota
	ResponseTypeLog
	ResponseTypePong
	ResponseTypeLeave
	ResponseTypeWelcome
	ResponseTypeMembers
	ResponseTypeJoin
	ResponseTypeChat
	ResponseTypeRoomInfo
)
