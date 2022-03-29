package mailx

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/textproto"
)

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
		// Text + Attachements
		mixed, err := s.buildMixed()
		if err != nil {
			return "", fmt.Errorf("buildMixed failed, err: %s", err)
		}
		buf.WriteString(mixed)
	}

	return buf.String(), nil
}

func (s *SendMail) buildHeader() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "From", s.FromDisplayName()))
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "To", s.ToDisplayName()))
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Cc", s.CcDisplayName()))
	buf.WriteString(fmt.Sprintf("%s: %s\r\n", "Subject", s.SubjectEncoded()))
	return buf.String()
}

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

// 邮件正文可以同时提供 HTML 形式和纯文本形式，相同内容使用不同形式表示。
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
			return "", fmt.Errorf("create a1 part failed, err: %s", err)
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
