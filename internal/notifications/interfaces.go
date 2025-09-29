package notifications

import "context"

type INotifications interface {
	SendVerificationCode(ctx context.Context, to, code string) error
	SendLoginNotification(ctx context.Context, to string) error
}
