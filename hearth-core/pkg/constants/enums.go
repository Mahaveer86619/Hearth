package constants

type WebSocketMessageType string

const (
	MessageTypeHeartbeat WebSocketMessageType = "heartbeat"
	MessageTypeData      WebSocketMessageType = "data"
)

type RedisChannels string

const (
	RedisChannelLiveLogs RedisChannels = "live_logs"
)
