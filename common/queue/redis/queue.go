package redis

import (
	"Panda/common/log"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

type Queue struct {
	topic   string
	body    string
	redis   *redis.Client
	message chan *message
}

type message struct {
	Value []byte
}

const (
	prefixQueueKey = "travel:queue"
	channelCount   = 1
)

func NewQueue(r *redis.Client) *Queue {
	return &Queue{redis: r}
}

func (q *Queue) getTopic(t string) string {
	return fmt.Sprintf("%s:%s", prefixQueueKey, t)
}

func (q *Queue) Producer(topic string, message interface{}) (err error) {
	if topic == "" {
		return errors.New("topic is invalid")
	}
	data, errs := json.Marshal(message)
	if errs != nil {
		log.Error("producer message marshal err:", err)
		return errs
	}
	err = q.redis.RPush(q.getTopic(topic), data).Err()
	if err != nil {
		log.Error("redis queue err:", err)
	}
	return
}

func (q *Queue) GetMessage() <-chan *message {
	return q.message
}

func (q *Queue) Consumer(topic string, ctx context.Context) *Queue {
	q.message = make(chan *message, channelCount)
	go func() {
		topic = q.getTopic(topic)
		for {
			select {
			case <-ctx.Done():
				return
			default:
				result, err := q.redis.BLPop(30*time.Second, topic).Result()
				if err == nil {
					q.message <- &message{[]byte(result[1])}
				}
			}
		}
	}()
	return q
}
