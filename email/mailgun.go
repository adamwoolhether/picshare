package email

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"time"
)

//const (
//	welcomeSubject = "Welcome to PicApp!"
//
//	welcomeBody = `Welcome to my site. This is a demonstration of my ability with Go!
//
//	Feel free to create a gallery, upload photos, and share them with friends.
//
//	Please reach out if you have suggestions for improvement.
//
//	- Adam`
//
//	htmlBody = `Hey!<br/>
//	Welcome to my site, <a href ="test.adamwoolhether.com">PicApp</a>.<br><br>
//	This is a demonstration of my ability with Go!<br><br>
//	Feel free to create a gallery, upload photos, and share them with friends.<br><br>
//	Please reach out if you have suggestions for improvement.<br><br>
//	- Adam`
//)
//
//type Client struct {
//	from string
//	mg   mailgun.Mailgun
//}
//
//type ClientConfig func(*Client)
//
//func NewClient(opts ...ClientConfig) *Client {
//	client := Client{
//		from: "adamwoolhether@gmail.com",
//	}
//	for _, opt := range opts {
//		opt(&client)
//	}
//	return &client
//}
//
//func WithSender(name, email string) ClientConfig {
//	return func(c *Client) {
//		c.from = buildEmail(name, email)
//
//	}
//}
//
//func WithMailGun(domain, apiKey string) ClientConfig {
//	return func(c *Client) {
//		mg := mailgun.NewMailgun(domain, apiKey)
//		c.mg = mg
//	}
//}
//
//func buildEmail(name, email string) string {
//	if name == "" {
//		return email
//	}
//	return fmt.Sprintf("%s <%s>", name, email)
//}
//
//func (c *Client) Welcome(name, email string) error {
//
//	mg := mailgun.NewMailgun("***REMOVED***", "***REMOVED***")
//	msg := mg.NewMessage(c.from, welcomeSubject, welcomeBody, buildEmail(name, email))
//	//msg := c.mg.NewMessage(c.from, welcomeSubject, welcomeBody, buildEmail(name, email))
//	msg.SetHtml(htmlBody)
//
//	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
//	defer cancel()
//	_, _, err := c.mg.Send(ctx, msg)
//	return err
//}


func SignUpEmail(domain, apiKey, emailAddy string) {
	mg := mailgun.NewMailgun(domain, apiKey)

	sender := "adamwoolhether@gmail.com"
	subject := "Welcome to PicApp"
	body := `Hey!
	Welcome to my site. This is a demonstration of my ability with Go!

	Feel free to create a gallery, upload photos, and share them with friends.

	Please reach out if you have suggestions for improvement.

	- Adam`

	htmlBody := `Hey!<br/>
	Welcome to my site, <a href ="test.adamwoolhether.com">PicApp</a>.<br><br>
	This is a demonstration of my ability with Go!<br><br>
	Feel free to create a gallery, upload photos, and share them with friends.<br><br>
	Please reach out if you have suggestions for improvement.<br><br>
	- Adam`

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
