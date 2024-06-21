package notification

import (
	"context"
	"errors"
	"fmt"
	"time"

	smtpclient "github.com/arfan21/vocagame/client/smtp"
	"github.com/arfan21/vocagame/internal/model"
	"github.com/arfan21/vocagame/pkg/constant"
	"github.com/arfan21/vocagame/pkg/logger"
	"github.com/dranikpg/gtrs"
	gomail "github.com/wneessen/go-mail"
)

func (n Notification) Consume() {
	ctx := context.Background()

	groupC1 := gtrs.NewGroupConsumer[model.Event](ctx, n.client, fmt.Sprintf("%s-group", model.StreamNotification), model.StreamNotification, model.StreamNotification, "0-0")

	defer func() {
		groupC1.Close()
	}()

	logger.Log(ctx).Info().Msg("start consuming redis email stream...")

	for {
		select {
		case <-ctx.Done():
			logger.Log(ctx).Info().Msg("stop consuming stream")
			return
		default:
			msg := <-groupC1.Chan()

			if msg.Err != nil {
				logger.Log(ctx).Error().Err(msg.Err).Msg("failed to consume message")
				continue
			}

			if msg.Data.Name == "" {
				groupC1.Ack(msg)
				continue
			}

			if msg.Data.RetryCount > 3 {
				logger.Log(ctx).Error().Msgf("failed to consume message: %s, with email: %s, retry count: %d, max retry reached", msg.Data.Name, msg.Data.Email, msg.Data.RetryCount)
				groupC1.Ack(msg)
				continue
			}

			data := msg.Data
			var content *gomail.Msg
			var err error
			switch data.Name {
			case model.EventTransactionNotification:
				content, err = n.handleTransactionNotification(ctx, data)
			default:
				logger.Log(ctx).Error().Msgf("unknown event name: %s", data.Name)
				groupC1.Ack(msg)
				continue
			}

			if err != nil {
				logger.Log(ctx).Error().Err(err).Msgf("failed to get email template: %s, with email: %s", data.Name, data.Email)
				groupC1.Ack(msg)

				if errors.Is(err, constant.ErrTransactionNotFound) {
					data.RetryCount++
					time.Sleep(time.Second * 5)
					// retry
					_, err = n.stream.Add(ctx, data)
					if err != nil {
						logger.Log(ctx).Error().Err(err).Msgf("failed to retry event: %s", data.Name)
					}
				}

				continue
			}

			err = n.smtpClient.SendEmail(ctx, data.Email, content)
			if err != nil {
				logger.Log(ctx).Error().Msgf("failed to send email: %s, with email: %s", data.Name, data.Email)
				continue
			}
			groupC1.Ack(msg)
		}
	}
}

func (n Notification) handleTransactionNotification(ctx context.Context, data model.Event) (res *gomail.Msg, err error) {
	logger.Log(ctx).Info().Msgf("handle transaction notification: %s", data.TransactionID)

	transactionData, err := n.transactionRepo.GetByID(ctx, data.TransactionID)
	if err != nil {
		err = fmt.Errorf("notification.handleTransactionNotification: failed to get transaction data: %w", err)
		return
	}

	content, err := n.smtpClient.GetNotifTemplate(smtpclient.NotifEmailBody{
		ToEmail:     data.Email,
		ProductName: transactionData.Product.Name,
		Status:      string(transactionData.Status),
		TotalPrice:  transactionData.TotalPrice.InexactFloat64(),
		Url:         data.Url,
	})
	if err != nil {
		err = fmt.Errorf("notification.handleTransactionNotification: failed to get email template: %w", err)
		return
	}

	res = content

	return
}
