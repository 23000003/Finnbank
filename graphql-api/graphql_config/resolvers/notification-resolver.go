package resolvers

import (
	"finnbank/graphql-api/graphql_config/entities"
	sv "finnbank/graphql-api/services"

	"github.com/graphql-go/graphql"
)

// "github.com/graphql-go/graphql"

// func (s *StructGraphQLResolvers) GetNotificationQueryType(#) *graphql.Object {

// }

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

// func (s *StructGraphQLResolvers) GetNotificationMutationType(#) *graphql.Object {

// }
