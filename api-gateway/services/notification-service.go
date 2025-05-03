package services

import (
	"bytes"
	"encoding/json"
	t "finnbank/api-gateway/types"
	"finnbank/common/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationService struct {
	log *utils.Logger
	url string
}

func newNotificationService(log *utils.Logger) *NotificationService {
	return &NotificationService{
		log: log,
		url: "http://localhost:8083/graphql/notification",
	}
}

func (a *NotificationService) GetAllNotificationByUserId(ctx *gin.Context) {
	id := ctx.Param("id");
	limit, err := strconv.Atoi(ctx.Param("limit"))

	if err != nil {
		a.log.Info("Error converting limit to int: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit format"})
		return
	}

	query := fmt.Sprintf(`{
		getAllNotificationByUserId(notif_to_id: "%s", limit: %d) {
			notif_id
			notif_type
			notif_from_name
			content
			is_read
			date_notified
		}
	}`, id, limit)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.GetAllNotificationsGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetAllNotificationByUserId})
}

func (a *NotificationService) GetAllUnreadNotificationByUserId(ctx *gin.Context) {
	id := ctx.Param("id");

	query := fmt.Sprintf(`{
		getAllUnreadNotificationByUserId(notif_to_id: "%s") {
			total_notification
			unread_notification
		}
	}`, id)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.GetAllUnreadNotificationGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetAllUnreadNotificationByUserId})
}

func (a *NotificationService) GetNotificationByUserId(ctx *gin.Context) {
	id := ctx.Param("id");

	query := fmt.Sprintf(`{
		getNotificationById(notif_id: "%s") {
			notif_id
			notif_type
			notif_from_name
			content
			is_read
			date_notified
		}
	}`, id)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.GetNotificationGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": data.Data.GetNotificationById})
}

func (a *NotificationService) GenerateNotification(ctx *gin.Context) {
	var req t.CreateNotificationRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := fmt.Sprintf(`mutation {
		generateNotification(notif_type: "%s", notif_to_id: "%s", notif_from_name: "%s", content: "%s") {
			notif_id
		}
	}`, req.NotifType, req.NotifToID, req.NotifFromName, req.Content)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}
	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	defer resp.Body.Close()

	a.log.Info("Response: %v", resp.Body)

	var data t.CreateOpenedAccountsGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	a.log.Info("%v ======= DATA", data)

	if data.Errors != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Notified successfully"})
}

func (a *NotificationService) ReadNotificationByUserId(ctx *gin.Context) {
	id := ctx.Param("id");

	query := fmt.Sprintf(`mutation {
		readNotificationByUserId(notif_id: "%s") {
			notif_id
		}
	}`, id)

	qlrequestBody := map[string]interface{}{
		"query": query,
	}

	qlrequestBodyJSON, _ := json.Marshal(qlrequestBody)

	resp, err := http.Post(a.url, "application/json", bytes.NewBuffer(qlrequestBodyJSON))
	if err != nil {
		a.log.Info("Request Error: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	var data t.GetNotificationGraphQLResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		a.log.Info("Error decoding response: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if data.Errors != nil {
		a.log.Info("GraphQL Errors: %v", data.Errors)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": data.Errors})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"data": "Notification read successfully"})
}
