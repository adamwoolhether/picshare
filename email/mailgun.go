package email

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"net/url"
	"picapp/conf"
	"time"
)

const (
	resetURL = "https://test.adamwoolhether/reset"
	sender   = "adamwoolhether@gmail.com"
	subject  = "Welcome to PicApp"
	body     = `Hey!
	Welcome to my site. This web app is a working project to play with my Go skills.

	Feel free to create a gallery, upload photos, and share them with friends.

	Please reach out if you have suggestions for improvement.

	- Adam`

	htmlBody = `Hey!<br><br>
	Welcome to my site, <a href ="test.adamwoolhether.com">PicApp</a>.<br><br>
	This web app is a working project to play with my Go skills.<br><br>
	Feel free to create a gallery, upload photos, and share them with friends.<br><br>
	Please reach out if you have suggestions for improvement.<br><br>
	- Adam`

	resetSubject = "Uh oh, forgot your password?"

	resetText = `Hey!
	A password reset request was submitted for your account.
	If this was you, please follow the link below to update your password:

	%s
	
	If the link doesn't automatically populate the "Token" field, please paste in the following line:
	
	%s

	If this request wasn't submitted by you, please disregard this email.

	Thanks,
	Adam
	`

	resetHTML = `Hey!<br><br>
	A password reset request was submitted for your account.<br><br>
	If this was you, please follow the link below to update your password:<br><br>
	<a href="%s">%s</a><br><br>
	If the link doesn't automatically populate the "Token" field, please paste in the following line:<br><br>
	%s<br><br>
	If this request wasn't submitted by you, please disregard this email.<br><br>
	Thanks,
	Adam
	`
)

func SignUpEmail(emailAddy string) {
	cfg := conf.LoadConfig(true)
	mg := mailgun.NewMailgun(cfg.Mailgun.Domain, cfg.Mailgun.APIKey)

	msg := mg.NewMessage(sender, subject, body, emailAddy)
	msg.SetHtml(htmlBody)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
}

func ResetPw(emailAddy, token string) error {
	cfg := conf.LoadConfig(true)
	mg := mailgun.NewMailgun(cfg.Mailgun.Domain, cfg.Mailgun.APIKey)
	//TODO: Build the reset URL
	v := url.Values{}
	v.Set("token", token)
	resetURL := resetURL + "?" + v.Encode()

	resetBody := fmt.Sprintf(resetText, resetURL, token)
	msg := mg.NewMessage(sender, resetSubject, resetBody, emailAddy)
	resetBodyHTML := fmt.Sprintf(resetHTML, resetURL, resetURL, token)
	msg.SetHtml(resetBodyHTML)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	resp, id, err := mg.Send(ctx, msg)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return nil
}
