package service

import (
	"context"
	"fmt"
	"log"
	sms "notifyservice/email"
	proto "notifyservice/protos/notify"
)

type NotificationService struct {
	proto.UnimplementedNotificationServiceServer
}

func NewNotificationService() *NotificationService {
	return &NotificationService{}
}

func (s *NotificationService) SendNotification(ctx context.Context, req *proto.SendNotificationRequest) (*proto.SendNotificationResponse, error) {
	message := req.Message
	if message == "" {
		message = fmt.Sprintf("%s Service Notification", req.ServiceType)
	}

	err := sms.SendSMS(message, req.Message, req.Email, req.UserID)
	if err != nil {
		log.Printf("Failed to send notification: %v", err)
		return nil, fmt.Errorf("failed to send notification: %v", err)
	}

	return &proto.SendNotificationResponse{Status: "Success"}, nil
}
