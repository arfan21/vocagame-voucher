package notification

import (
	"context"
	"time"

	smtpclient "github.com/arfan21/vocagame/client/smtp"
	"github.com/arfan21/vocagame/internal/entity"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/dranikpg/gtrs"
	"github.com/redis/go-redis/v9"
	gomail "github.com/wneessen/go-mail"
)

type SMTPClient interface {
	SendEmail(ctx context.Context, toEmail string, content *gomail.Msg) (err error)
	GetNotifTemplate(bodyEmail smtpclient.NotifEmailBody) (*gomail.Msg, error)
}

type TransactionRepo interface {
	GetByID(ctx context.Context, id string) (res entity.Transaction, err error)
}

type Notification struct {
	client          *redis.Client
	smtpClient      SMTPClient
	transactionRepo TransactionRepo
	stream          gtrs.Stream[model.Event]
}

func New(client *redis.Client, smtpClient SMTPClient, transactionRepo TransactionRepo) *Notification {

	stream := gtrs.NewStream[model.Event](client, model.StreamNotification, &gtrs.Options{
		TTL:    time.Hour * 24 * 7,
		MaxLen: 1000,
		Approx: true,
	})

	return &Notification{client: client, smtpClient: smtpClient, transactionRepo: transactionRepo, stream: stream}
}
