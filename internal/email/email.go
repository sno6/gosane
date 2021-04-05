package email

type Emailer interface {
	SendTemplateEmail(toEmail string, template string, data *EmailData) error
	SendRawEmail(toEmail string, data *RawEmailData) error
}

type EmailData struct {
	Subject string
	Content interface{}
}

type RawEmailData struct {
	Subject string
	Content string
}