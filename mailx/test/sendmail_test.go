package test

import (
	_ "embed"
	"fmt"
	"time"

	"net/mail"
	"os"
	"testing"

	"github.com/bougou/gopkg/mailx"
)

//go:embed testmail.html
var contentHTML string

var contentPlain = "this is hello"

func TestSendHTMLNative(t *testing.T) {

	to := []mail.Address{
		{
			Name:    "",
			Address: os.Getenv("TO_MAIL_1"),
		},
		{
			Name:    "Bougou",
			Address: os.Getenv("TO_MAIL_2"),
		},
	}

	cc := []mail.Address{
		{
			Name:    "",
			Address: os.Getenv("TO_MAIL_3"),
		},
	}

	bcc := []mail.Address{
		{
			Name:    "",
			Address: os.Getenv("TO_MAIL_4"),
		},
	}

	froms := []struct {
		addr string
		user string
		pass string
	}{
		{
			addr: os.Getenv("FROM_126MAIL_ADDR"),
			user: os.Getenv("FROM_126MAIL_USER"),
			pass: os.Getenv("FROM_126MAIL_PASS"),
		},
		// {
		// 	addr: os.Getenv("FROM_GMAIL_ADDR"),
		// 	user: os.Getenv("FROM_GMAIL_USER"),
		// 	pass: os.Getenv("FROM_GMAIL_PASS"),
		// },
	}

	subject := "Hello world 🚨🔥✅❗"

	for _, f := range froms {
		from := mail.Address{
			Name:    "发报局",
			Address: f.user,
		}

		sm := mailx.NewSendMail(fmt.Sprintf("%s:%d", f.addr, 25), false)

		sm.WithAuth(true, f.user, f.pass).
			WithFrom(from).
			WithTo(to).
			WithCc(cc).
			WithBcc(bcc).
			WithSubject(subject).
			WithContent(contentPlain, contentHTML).
			WithDebug(true)

		if err := sm.AddAttachments("1.cfg", "2.png"); err != nil {
			t.Errorf("add attachments failed, err: %s", err)
		}

		if err := sm.Send(); err != nil {
			t.Error(err)
		}

		time.Sleep(5 * time.Second)
	}

}
