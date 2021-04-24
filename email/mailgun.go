package email

import (
	"context"
	"fmt"
	"github.com/mailgun/mailgun-go/v4"
	"log"
	"picapp/conf"
	"time"
)

const (
	sender  = "adamwoolhether@gmail.com"
	subject = "Welcome to PicApp"
	body    = `Hey!
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
)

func SignUpEmail(emailAddy string) {
	mgCfg := conf.LoadConfig(true)
	mg := mailgun.NewMailgun(mgCfg.Mailgun.Domain, mgCfg.Mailgun.APIKey)

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
