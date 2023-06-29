package repository

import (
	"context"

	"github.com/vynnydev/chatgpt-microservices-k8s/packages/microservices/chatservice/domain/entity"
)

type ChatRepository interface {
	CreateChat(ctx context.Context, chat *entity.Chat) error
	FindChatByID(ctx context.Context, chatID string) (*entity.Chat, error)
	SaveChat(ctx context.Context, chat *entity.Chat) error
}
