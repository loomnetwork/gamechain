package battleground

import (
	"fmt"
	"github.com/loomnetwork/gamechain/types/zb/zb_data"
	contract "github.com/loomnetwork/go-loom/plugin/contractpb"
)

func createBaseNotification(ctx contract.Context, currentNotifications []*zb_data.Notification, notificationType zb_data.NotificationType_Enum) *zb_data.Notification {
	var id int32 = 0
	if len(currentNotifications) > 0 {
		id = currentNotifications[len(currentNotifications)-1].Id + 1
	}

	return &zb_data.Notification{
		Id:        id,
		Type:      notificationType,
		CreatedAt: ctx.Now().Unix(),
		Seen:      false,
	}
}

func removeNotification(notifications []*zb_data.Notification, id int32) ([]*zb_data.Notification, error) {
	for index, notification := range notifications {
		if notification.Id == id {
			notifications = append(notifications[:index], notifications[index+1:]...)
			return notifications, nil
		}
	}

	return nil, fmt.Errorf("notification with id %d not found", id)
}
