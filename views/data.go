package views

import "picapp/models"

const (
	AlertLvlError   = "danger"
	AlerLvlWarning  = "warning"
	AlertLvlInfo    = "info"
	AlertLvlSuccess = "success"

	// AlertMsgGeneric displays when any random/unexpected error is encountered by our backend.
	AlertMsgGeneric = "Something went wrong. Please try again and let me know if it persists."
)

// Alert is used to render boostrap alert messages in templates
type Alert struct {
	Level   string
	Message string
}

// Data is the top level structure that views expect data to come in
type Data struct {
	Alert *Alert
	User *models.User
	Yield interface{}
}

func (d *Data) SetAlert(err error) {
	if pErr, ok := err.(PublicError); ok { // example of a type assertion
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: pErr.Public(),
		}
	} else {
		d.Alert = &Alert{
			Level:   AlertLvlError,
			Message: AlertMsgGeneric,
		}
	}
}

func(d *Data) AlertError(msg string) {
	d.Alert = &Alert{
		Level:   AlertLvlError,
		Message: msg,
	}
}

type PublicError interface {
	error
	Public() string
}
