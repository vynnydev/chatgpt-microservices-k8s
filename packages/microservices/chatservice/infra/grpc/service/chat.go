package service

import (
	"github.com/vynnydev/chatgpt-microservices-k8s/packages/microservices/chatservice/infra/grpc/pb"
	"github.com/vynnydev/chatgpt-microservices-k8s/packages/microservices/chatservice/usecases/chatcompletionstream"
)

type ChatService struct {
	pb.UnimplementedChatServiceServer
	ChatCompletionStreamusecases chatcompletionstream.ChatCompletionusecases
	ChatConfigStream             chatcompletionstream.ChatCompletionConfigInputDTO
	StreamChannel                chan chatcompletionstream.ChatCompletionOutputDTO
}

func NewChatService(chatCompletionStreamusecases chatcompletionstream.ChatCompletionusecases, chatConfigStream chatcompletionstream.ChatCompletionConfigInputDTO, streamChannel chan chatcompletionstream.ChatCompletionOutputDTO) *ChatService {
	return &ChatService{
		ChatCompletionStreamusecases: chatCompletionStreamusecases,
		ChatConfigStream:             chatConfigStream,
		StreamChannel:                streamChannel,
	}
}

func (c *ChatService) ChatStream(req *pb.ChatRequest, stream pb.ChatService_ChatStreamServer) error {
	chatConfig := chatcompletionstream.ChatCompletionConfigInputDTO{
		Model:                c.ChatConfigStream.Model,
		ModelMaxTokens:       c.ChatConfigStream.ModelMaxTokens,
		Temperature:          c.ChatConfigStream.Temperature,
		TopP:                 c.ChatConfigStream.TopP,
		N:                    c.ChatConfigStream.N,
		Stop:                 c.ChatConfigStream.Stop,
		MaxTokens:            c.ChatConfigStream.MaxTokens,
		InitialSystemMessage: c.ChatConfigStream.InitialSystemMessage,
	}
	input := chatcompletionstream.ChatCompletionInputDTO{
		UserMessage: req.GetUserMessage(),
		UserID:      req.GetUserId(),
		ChatID:      req.GetChatId(),
		Config:      chatConfig,
	}

	ctx := stream.Context()
	go func() {
		for msg := range c.StreamChannel {
			stream.Send(&pb.ChatResponse{
				ChatId:  msg.ChatID,
				UserId:  msg.UserID,
				Content: msg.Content,
			})
		}
	}()

	_, err := c.ChatCompletionStreamusecases.Execute(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
