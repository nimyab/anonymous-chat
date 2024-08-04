package websocket

import (
	"github.com/nimyab/anonymous-chat/internal/handlers/chat"
	"github.com/nimyab/anonymous-chat/internal/handlers/message"
	"sync"
)

type WebsocketService struct {
	chatService    *chat.ChatService
	messageService *message.MessageService

	userInSearch *UserQueue
}

type UserQueue struct {
	mu    sync.Mutex
	items []uint
}

func NewUserQueue() *UserQueue {
	return &UserQueue{items: make([]uint, 0, 100)}
}

func (q *UserQueue) PullTwoUserIds() []uint {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.items) < 2 {
		return nil
	}
	userIds := make([]uint, 2)
	copy(userIds, q.items[:2])
	copy(q.items, q.items[2:])
	return userIds
}

func (q *UserQueue) Push(userId uint) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, userId)
}
