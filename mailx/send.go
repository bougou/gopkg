package mailx

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
)

func (s *SendMail) Send() error {
	if err := s.Validate(); err != nil {
		return fmt.Errorf("validate mail failed, err: %s", err)
	}

	message, err := s.BuildMessageMIME()
	if err != nil {
		return fmt.Errorf("build msg failed, err: %s", err)
	}

	if s.debug {
		fmt.Println(message)
	}

	host, _, _ := net.SplitHostPort(s.addr)
	auth := smtp.PlainAuth("", s.authUser, s.authPass, host)
	// TLS config
	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	var c *smtp.Client
	var conn *tls.Conn

	if s.ssl {
		// Require SSL, normally for smtp servers running on 465 that require an ssl connection.
		// You need to call tls.Dial instead of smtp.Dial	from the very beginning.
		conn, err = tls.Dial("tcp", s.addr, tlsconfig)
		if err != nil {
			return fmt.Errorf("dial %s failed, err: %s", s.addr, err)
		}

		c, err = smtp.NewClient(conn, host)
		if err != nil {
			return fmt.Errorf("NewClient failed, err: %s", err)
		}
	} else {
		// for 25, 587
		c, err = smtp.Dial(s.addr)
		if err != nil {
			return fmt.Errorf("smtp Dial %s failed, err: %s", s.addr, err)
		}

		// No need to call c.Hello beforehand, it will autotically call it if needed.
		if ok, _ := c.Extension("STARTTLS"); ok {
			if err := c.StartTLS(tlsconfig); err != nil {
				return fmt.Errorf("StartTLS failed, err: %s", err)
			}
		}
	}

	if err := c.Auth(auth); err != nil {
		return fmt.Errorf("auth failed, err: %s", err)
	}

	if err := c.Mail(s.from.Address); err != nil {
		return fmt.Errorf("MAIL command failed, err: %s", err)
	}

	for _, mailAddress := range s.Recipients() {
		if s.debug {
			fmt.Println("rcpt:", mailAddress.Address)
		}
		if err = c.Rcpt(mailAddress.Address); err != nil {
			return fmt.Errorf("RCPT command failed, err: %s", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return fmt.Errorf("DATA command failed, err: %s", err)
	}

	_, err = w.Write([]byte(message))
	if err != nil {
		return fmt.Errorf("write message failed, err: %s", err)
	}

	if err := w.Close(); err != nil {
		return fmt.Errorf("close failed, err: %s", err)
	}

	if err := c.Quit(); err != nil {
		return fmt.Errorf("QUIT command failed, err: %s", err)
	}

	return nil
}
