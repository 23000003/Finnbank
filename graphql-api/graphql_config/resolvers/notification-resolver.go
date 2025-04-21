package resolvers

import (
	"finnbank/graphql-api/graphql_config/entities"
	sv "finnbank/graphql-api/services"
	"time"

	"github.com/graphql-go/graphql"
)

// func (s *StructGraphQLResolvers) GetNotificationQueryType(#) *graphql.Object

func (s *StructGraphQLResolvers) GetNotificationQueryType(NotifService *sv.NotificationService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "NotificationQuery",
		Fields: graphql.Fields{
			"getAllNotificationByUserId": &graphql.Field{
				Type: graphql.NewList(entities.GetNotificationEntityType()),
				Args: graphql.FieldConfigArgument{
					"notif_to_id": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					notifToID := p.Args["notif_to_id"].(string)
					notifications, err := NotifService.GetAllNotificationByUserId(notifToID)
					if err != nil {
						return nil, err
					}

					// optionally convert to []map[string]interface{} if your entity wants it that way
					result := make([]map[string]interface{}, len(notifications))
					for i, n := range notifications {
						result[i] = map[string]interface{}{
							"notif_id":        n.NotifID,
							"notif_type":      n.NotifType,
							"auth_id":         n.AuthID,
							"notif_to_id":     n.NotifToID,
							"notif_from_name": n.NotifFromName,
							"content":         n.Content,
							"is_read":         n.IsRead,
							"redirect_url":    n.RedirectURL,
							"date_notified":   n.DateNotified,
							"date_read":       n.DateRead,
						}
					}

					return result, nil
				},
			},
		},
	})
}

// func (s *StructGraphQLResolvers) GetNotificationMutationType(#) *graphql.Object

func (s *StructGraphQLResolvers) GetNotificationMutationType(NotifService *sv.NotificationService) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "NotificationMutation",
		Fields: graphql.Fields{
			"generateNotification": &graphql.Field{
				Type: entities.GetNotificationEntityType(),
				Args: graphql.FieldConfigArgument{
					"notif_type":      &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"auth_id":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"notif_to_id":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"notif_from_name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"content":         &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"redirect_url":    &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					notif := sv.Notification{
						NotifType:     p.Args["notif_type"].(string),
						AuthID:        p.Args["auth_id"].(string),
						NotifToID:     p.Args["notif_to_id"].(string),
						NotifFromName: p.Args["notif_from_name"].(string),
						Content:       p.Args["content"].(string),
						IsRead:        false,
						RedirectURL:   "",
						DateNotified:  time.Now(),
					}

					// Optional redirect_url (nullable)
					if val, ok := p.Args["redirect_url"].(string); ok {
						notif.RedirectURL = val
					}

					createdNotif, err := NotifService.GenerateNotification(notif)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"notif_id":        createdNotif.NotifID,
						"notif_type":      createdNotif.NotifType,
						"auth_id":         createdNotif.AuthID,
						"notif_to_id":     createdNotif.NotifToID,
						"notif_from_name": createdNotif.NotifFromName,
						"content":         createdNotif.Content,
						"is_read":         createdNotif.IsRead,
						"redirect_url":    createdNotif.RedirectURL,
						"date_notified":   createdNotif.DateNotified,
						"date_read":       createdNotif.DateRead,
					}, nil
				},
			},
		},
	})
}
