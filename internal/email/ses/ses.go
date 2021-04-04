package ses

import (
	"bytes"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"gopkg.in/gomail.v2"
)

const (
	sender = "sender@mysite.com"
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

func (s *SES) SendEmail(toEmail string, template string, content interface{}) error {
	templateData, err := json.Marshal(content)
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

func (s *SES) SendRawEmail(toEmail string, subject, content string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", "sender@mysite.com")
	msg.SetHeader("To", toEmail)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", content)

	var raw bytes.Buffer
	msg.WriteTo(&raw)

	_, err := s.ses.SendRawEmail(&ses.SendRawEmailInput{
		RawMessage: &ses.RawMessage{
			Data: raw.Bytes(),
		},
	})

	return err
}
