package email

type Emailer interface {
	SendEmail(toEmail string, template string, content interface{}) error
}
