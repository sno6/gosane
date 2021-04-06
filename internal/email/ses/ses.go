package ses

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/sno6/gosane/internal/email"
	"gopkg.in/gomail.v2"
)

const (
	sender = "sender@yoursite.com"
	region = "ap-southeast-2"
)

type SES struct {
	ses *ses.SES
}

func New() (*SES, error) {
	sesh, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	return &SES{
		ses: ses.New(sesh),
	}, nil
}

func (s *SES) SendTemplateEmail(toEmail string, template string, data *email.EmailData) error {
	templateData, err := json.Marshal(data.Content)
	if err != nil {
		return err
	}

	_, err = s.ses.SendTemplatedEmail(&ses.SendTemplatedEmailInput{
		Destination: &ses.Destination{
			CcAddresses: nil,
			ToAddresses: []*string{aws.String(toEmail)},
		},
		Source:       aws.String(sender),
		Template:     aws.String(template),
		TemplateData: aws.String(string(templateData)),
	})

	return err
}

func (s *SES) SendRawEmail(toEmail string, data *email.RawEmailData) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", sender)
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", data.Subject)
	msg.SetBody("text/html", data.Content)

	var raw bytes.Buffer
	msg.WriteTo(&raw)

	_, err := s.ses.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: raw.Bytes(),
		},
	})

	return err
}
