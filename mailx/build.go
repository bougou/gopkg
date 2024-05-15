package mailx

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

// BuildMessageMIME constructs the whole data of a mail message.
func (s *SendMail) BuildMessageMIME() (string, error) {
	buf := bytes.NewBuffer(nil)

	// Header
	buf.WriteString(s.buildHeader())

	if len(s.attachments) == 0 {
		// Just Text(plain/html), no attachements
		text, err := s.buildText()
		if err != nil {
			return "", fmt.Errorf("buildText failed, err: %s", err)
		}
		buf.WriteString(text)
	} else {
		// Text(plain/html) + Attachements
		mixed, err := s.buildMixed()
		if err != nil {
			return "", fmt.Errorf("buildMixed failed, err: %s", err)
		}
		buf.WriteString(mixed)
	}

	return buf.String(), nil
}

// builderHeader constructs the header part of a mail message.
func (s *SendMail) buildHeader() string {
	buf := bytes.NewBuffer(nil)
	if s.from != nil {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "From", s.FromDisplayName()))
	}
	if s.sender != nil {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Sender", s.SenderDisplayName()))
	}
	if s.replyTo != nil {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Reply-To", s.ReplyToDisplayName()))
	}

	if len(s.to) != 0 {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "To", s.ToDisplayName()))
	}
	if len(s.cc) != 0 {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Cc", s.CcDisplayName()))
	}
	if len(s.bcc) != 0 {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", s.SubjectEncoded()))
	}
	return buf.String()
}

// buildText constructs the mail content.（邮件正文）
//
// You should use `buildText` when there's no attachement in the mail,
// and use `buildMixed` when there are attachements in the mail.
func (s *SendMail) buildText() (string, error) {
	buf := bytes.NewBuffer(nil)

	if s.contentHTML != "" && s.contentPlain != "" {
		body, err := s.buildAlternative()
		if err != nil {
			return "", fmt.Errorf("buildAlternativeBody failed, err: %s", err)
		}
		buf.WriteString(body)
	} else if s.contentHTML != "" {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", "text/html; charset=UTF-8"))
		buf.WriteString("\r\n")
		buf.Write([]byte(s.contentHTML))
		buf.WriteString("\r\n")
	} else {
		buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", "text/plain"))
		buf.WriteString("\r\n")
		buf.Write([]byte(s.contentPlain))
		buf.WriteString("\r\n")
	}

	return buf.String(), nil
}

// 邮件正文可以同时提供 HTML 形式和纯文本形式。
// 收件方可以根据设备的能力支持选择其中一种进行展示。
func (s *SendMail) buildAlternative() (string, error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)

	contentType := fmt.Sprintf("multipart/alternative; boundary=%s", writer.Boundary())
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", contentType))
	buf.WriteString("\r\n")

	plainHeader := textproto.MIMEHeader{}
	plainHeader.Add("Content-Type", "text/plain")
	plainPart, err := writer.CreatePart(plainHeader)
	if err != nil {
		return "", fmt.Errorf("create html part failed, err: %s", err)
	}
	plainPart.Write([]byte(s.contentPlain))

	htmlHeader := textproto.MIMEHeader{}
	htmlHeader.Add("Content-Type", "text/html; charset=UTF-8")
	htmlPart, err := writer.CreatePart(htmlHeader)
	if err != nil {
		return "", fmt.Errorf("create html part failed, err: %s", err)
	}
	htmlPart.Write([]byte(s.contentHTML))

	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("writer close failed, err: %s", err)
	}
	buf.WriteString("\r\n")

	return buf.String(), nil
}

// buildMixed constructs the content containing the mail text and mail attachements.
func (s *SendMail) buildMixed() (string, error) {
	buf := bytes.NewBuffer(nil)

	writer := multipart.NewWriter(buf)
	boundary := writer.Boundary()

	contentType := fmt.Sprintf("multipart/mixed; boundary=%s", boundary)
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Content-Type", contentType))
	buf.WriteString("\r\n")

	text, err := s.buildText()
	if err != nil {
		return "", fmt.Errorf("buildText failed, err: %s", err)
	}
	buf.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	buf.WriteString(text)

	for _, v := range s.attachments {
		header := textproto.MIMEHeader{}
		header.Add("Content-Type", http.DetectContentType(v.Content))
		header.Add("Content-Transfer-Encoding", "base64")
		header.Add("Content-Disposition", fmt.Sprintf("attachment; filename=%s", v.Name))
		attachPart, err := writer.CreatePart(header)
		if err != nil {
			return "", fmt.Errorf("create attach part failed, err: %s", err)
		}

		b := make([]byte, base64.StdEncoding.EncodedLen(len(v.Content)))
		base64.StdEncoding.Encode(b, v.Content)
		attachPart.Write(b)
	}
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("writer close failed, err: %s", err)
	}

	return buf.String(), nil
}
