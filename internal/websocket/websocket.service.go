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

	searchUsers chan<- []uint
}

func NewUserQueue(searchUsers chan<- []uint) *UserQueue {
	return &UserQueue{items: make([]uint, 0, 100), searchUsers: searchUsers}
}

func (q *UserQueue) DeleteUserId(userId uint) int {
	q.mu.Lock()
	defer q.mu.Unlock()

	for i := 0; i < len(q.items); i++ {
		if q.items[i] == userId {
			q.items = append(q.items[:i], q.items[i+1:]...)
			return i
		}
	}
	return -1
}

func (q *UserQueue) Push(userId uint) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.items = append(q.items, userId)
	if len(q.items) >= 2 {
		go func() {
			userIds := make([]uint, 2)
			copy(userIds, q.items[:2])
			copy(q.items, q.items[2:])
			q.searchUsers <- userIds
		}()
	}
}
