package http_api

import (
	"errors"

	"github.com/nsqio/nsq/internal/protocol"
)

//url参数获取的接口
type getter interface {
	Get(key string) (string, error)
}

//通过url参数获取topic名字和channel名字
func GetTopicChannelArgs(rp getter) (string, string, error) {
	topicName, err := rp.Get("topic")
	if err != nil {
		return "", "", errors.New("MISSING_ARG_TOPIC")
	}

	if !protocol.IsValidTopicName(topicName) {
		return "", "", errors.New("INVALID_ARG_TOPIC")
	}

	channelName, err := rp.Get("channel")
	if err != nil {
		return "", "", errors.New("MISSING_ARG_CHANNEL")
	}

	if !protocol.IsValidChannelName(channelName) {
		return "", "", errors.New("INVALID_ARG_CHANNEL")
	}

	return topicName, channelName, nil
}
