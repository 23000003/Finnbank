package services

import (
	"finnbank/services/common/utils"
	"github.com/gin-gonic/gin"
)

type NotificationService struct {
	log *utils.Logger
	url string
}

func NewNotificationService(log *utils.Logger) *NotificationService {
	return &NotificationService{
		log: log,
		url: "http://localhost:8083/graphql/notification",
	}
}

func (a *NotificationService) GetUserNotifications(*gin.Context)  {
}

func (a *NotificationService) GenerateNotification(*gin.Context)  {
}

func (a *NotificationService) MarkAsReadNotification(*gin.Context)  {
}

func (a *NotificationService) DeleteNotification(*gin.Context)  {
}