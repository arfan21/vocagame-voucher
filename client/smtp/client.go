package smtpclient

import (
	"context"
	"embed"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/arfan21/vocagame/config"
	"github.com/arfan21/vocagame/pkg/logger"
	gomail "github.com/wneessen/go-mail"
)

//go:embed *.html
var embedEmailTemplate embed.FS

type NotifEmailBody struct {
	ToEmail     string
	ProductName string
	Status      string
	TotalPrice  float64
	Url         string
}

type SmtpClient struct{}

func New() *SmtpClient {
	logger.Log(context.Background()).Info().Msgf("smtp: username %v, from os %v", config.Get().Smtp.Username, os.Getenv("SMTP_USERNAME"))
	logger.Log(context.Background()).Info().Msgf("smtp: host %v, from os %v", config.Get().Smtp.Host, os.Getenv("SMTP_HOST"))
	logger.Log(context.Background()).Info().Msgf("smtp: port %v, from os %v", config.Get().Smtp.Port, os.Getenv("SMTP_PORT"))

	return &SmtpClient{}
}

func (s SmtpClient) initMailer() (*gomail.Msg, error) {
	mailer := gomail.NewMsg()

	smtpUsername := config.Get().Smtp.Username

	err := mailer.EnvelopeFrom(smtpUsername)
	if err != nil {
		err = fmt.Errorf("service: failed to set envelope from: %w", err)
		return nil, err
	}

	emailFrom := config.Get().Smtp.Username

	err = mailer.FromFormat("Vocagame noreply", emailFrom)
	if err != nil {
		err = fmt.Errorf("service: failed to set email from: %w", err)
		return nil, err
	}

	return mailer, nil
}

func (s SmtpClient) GetNotifTemplate(bodyEmail NotifEmailBody) (*gomail.Msg, error) {
	t, err := template.ParseFS(embedEmailTemplate, "notif.html")
	if err != nil {
		err = fmt.Errorf("service: failed to parse email template: %w", err)
		return nil, err
	}

	mailer, err := s.initMailer()
	if err != nil {
		err = fmt.Errorf("service: failed to init mailer: %w", err)
		return nil, err
	}

	err = mailer.To(bodyEmail.ToEmail)
	if err != nil {
		err = fmt.Errorf("service: failed to set email to: %w", err)
		return nil, err
	}

	mailer.Subject("Transaction Notification")

	mailer.AddAlternativeHTMLTemplate(t, bodyEmail)

	return mailer, nil
}

func (s SmtpClient) SendEmail(ctx context.Context, toEmail string, content *gomail.Msg) (err error) {
	logger.Log(ctx).Info().Msgf("smtp: sending otp to email %s", toEmail)

	portInt, err := strconv.Atoi(config.Get().Smtp.Port)
	if err != nil {
		err = fmt.Errorf("service: failed to convert port to int: %w", err)
		return err
	}

	mailOpts := []gomail.Option{
		gomail.WithPort(portInt),
		gomail.WithSMTPAuth(gomail.SMTPAuthPlain),
		gomail.WithUsername(config.Get().Smtp.Username),
		gomail.WithPassword(config.Get().Smtp.Password),
	}

	if config.Get().Smtp.Port == "465" {
		mailOpts = append(mailOpts, gomail.WithSSLPort(false))
	}

	client, err := gomail.NewClient(config.Get().Smtp.Host, mailOpts...)
	if err != nil {
		err = fmt.Errorf("service: failed to create new client: %w", err)
		return err
	}

	err = client.DialAndSend(content)
	if err != nil {
		err = fmt.Errorf("service: failed to dial and send email: %w", err)
		return err
	}

	return nil
}
