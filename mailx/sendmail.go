package mailx

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/mail"
	"path/filepath"
	"strings"
)

type SendMail struct {
	addr string
	ssl  bool

	auth     bool
	authUser string
	authPass string

	debug bool

	// 发件人
	// 发件人：腾讯视频 <film@tencent.com>
	// 发件人：RedHat <email@redhat.com>
	// 发件人：Google <no-reply@accounts.google.com>
	// see: https://askleo.com/why_are_email_addresses_sometimes_in_anglebrackets/
	// the valid from format is:
	// "some@domain.com"
	// "Display Name <some@domain.com>"
	from mail.Address
	// 收件人 The identity of the primary recipients of the message.
	to []mail.Address
	// 抄送 The identity of the secondary (informational) recipients of the message.
	cc []mail.Address
	// 密送 The identity of additional recipients of the message.
	// The contents of this field are not included in copies of the message sent to the primary and secondary recipients.
	bcc []mail.Address

	// 主题
	subject string

	// 正文
	contentPlain string
	contentHTML  string

	// 附件
	// Use slice to keep order
	attachments []Attachment
}

type Attachment struct {
	Name    string
	Content []byte
}

// NewSendMail creates SendMail, the smtpServerAddr must contain tcp port.
func NewSendMail(smtpServerAddr string, ssl bool) *SendMail {
	return &SendMail{
		addr:        smtpServerAddr,
		ssl:         ssl,
		attachments: make([]Attachment, 0),
	}
}

// WithAuth set auth info. If auth is false, the user and pass are ignored.
func (s *SendMail) WithAuth(auth bool, user string, pass string) *SendMail {
	s.auth = auth
	s.authUser = user
	s.authPass = pass
	return s
}

func (s *SendMail) WithDebug(debug bool) *SendMail {
	s.debug = debug
	return s
}

func (s *SendMail) WithFrom(from mail.Address) *SendMail {
	s.from = from
	return s
}

func (s *SendMail) WithTo(to []mail.Address) *SendMail {
	s.to = to
	return s
}

func (s *SendMail) WithCc(cc []mail.Address) *SendMail {
	s.cc = cc
	return s
}

func (s *SendMail) WithBcc(bcc []mail.Address) *SendMail {
	s.bcc = bcc
	return s
}

// WithSubject sets the subject for the mail.
// Accepts non-ascii characters in the subject.
// SendMail will automatically do the conversion of encoding.
func (s *SendMail) WithSubject(subject string) *SendMail {
	s.subject = subject
	return s
}

func (s *SendMail) WithContent(plainContent string, htmlContent string) *SendMail {
	s.contentPlain = plainContent
	s.contentHTML = htmlContent
	return s
}

func (s *SendMail) AddAttachments(files ...string) error {
	fmt.Println("Add att", files)
	for _, file := range files {
		b, err := ioutil.ReadFile(file)
		if err != nil {
			return fmt.Errorf("read file from (%s) failed, err: %s", file, err)
		}
		_, filename := filepath.Split(file)
		s.attachments = append(s.attachments, Attachment{
			Name:    filename,
			Content: b,
		})
	}
	return nil
}

// Recipients return all recipients including to, cc, and bcc
func (s *SendMail) Recipients() []mail.Address {
	t := make([]mail.Address, 0, len(s.to)+len(s.cc)+len(s.bcc))
	t = append(t, s.to...)
	t = append(t, s.cc...)
	t = append(t, s.bcc...)

	return t
}

func (s *SendMail) FromDisplayName() string {
	return s.from.String()
}

func (s *SendMail) ToDisplayName() string {
	return strings.Join(plain(s.to), ",")
}

func (s *SendMail) CcDisplayName() string {
	return strings.Join(plain(s.cc), ",")
}

// SubjectEncoded encodes the subject to pure ascii characters.
// Thus, the subject can contain non-ascci-characters, like emoji...
// see: https://www.telemessage.com/developer/faq/how-do-i-encode-non-ascii-characters-in-an-email-subject-line/
func (s *SendMail) SubjectEncoded() string {
	return fmt.Sprintf("=?utf-8?B?%s?=", base64.StdEncoding.EncodeToString([]byte(s.subject)))
}

func (s *SendMail) Validate() error {
	var errs []error
	for _, rcpt := range s.Recipients() {
		if rcpt.Address == "" {
			errs = append(errs, fmt.Errorf("the address (%v) is not valid email address", rcpt))
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("errs: %s", errs)
	}
	return nil
}
