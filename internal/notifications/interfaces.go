package notifications

import "context"

type INotifications interface {
	SendRegistrationSuccess(ctx context.Context, to string) error
	SendLoginNotification(ctx context.Context, to string) error
}
