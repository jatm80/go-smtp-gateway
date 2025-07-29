package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	send "github.com/jatm80/go-smtp-gateway/send"

	mail "github.com/emersion/go-message/mail"
	sasl "github.com/emersion/go-sasl"
	smtp "github.com/emersion/go-smtp"
)

var (
	smtpAddr     = getEnv("SMTP_ADDR", "127.0.0.1:2525")
	smtpDomain   = getEnv("SMTP_DOMAIN","localhost")
	smtpUser     = getEnv("SMTP_USER","user")
	smtpPass     = getEnv("SMTP_PASS","empty")
	telegramBot  = os.Getenv("TELEGRAM_TOKEN")
	telegramChat = os.Getenv("TELEGRAM_CHAT_ID")
)

const (
	loginNotStarted loginState = iota
	loginWaitingUsername
	loginWaitingPassword
)

type LoginAuthenticator func(username, password string) error
type loginState int

type Backend struct{}

type Session struct {
	auth bool
	from string
	to string
}
type OutboundMsg struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

type loginAuthServer struct {
	state              loginState
	username, password string
	authenticate       LoginAuthenticator
}

func (bkd *Backend) NewSession(c *smtp.Conn) (smtp.Session, error) {
	return &Session{}, nil
}

func (s *Session) AuthMechanisms() []string {
	return []string{sasl.Plain,sasl.Login}
}

func (s *Session) Auth(mech string) (sasl.Server, error) {
	switch mech {
	case sasl.Plain: 	
		return sasl.NewPlainServer(func(identity, username, password string) error {
			hasher := md5.New()
			fmt.Fprintf(hasher, "%s%s", username, password)
			hashInBytes := hasher.Sum(nil)
			md5Hash := hex.EncodeToString(hashInBytes)
			fmt.Printf("Auth MD5 Hash: %s\n", md5Hash)

			if username != smtpUser || password != smtpPass {
				return errors.New("invalid username or password")
			}
			s.auth = true
			return nil
	}), nil

	case sasl.Login: 
		return NewLoginServer(func(username, password string) error {
			hasher := md5.New()
			fmt.Fprintf(hasher, "%s%s", username, password)
			hashInBytes := hasher.Sum(nil)
			md5Hash := hex.EncodeToString(hashInBytes)
			fmt.Printf("Auth MD5 Hash: %s\n", md5Hash)
			if username != smtpUser || password != smtpPass {
				return errors.New("invalid username or password")
			}
		s.auth = true
		return nil
		}), nil
	default:
		return nil, fmt.Errorf("unsupported auth mechanism: %s", mech)
	}
}

func NewLoginServer(authenticator LoginAuthenticator) sasl.Server {
	return &loginAuthServer{authenticate: authenticator}
}

func (a *loginAuthServer) Next(response []byte) (challenge []byte, done bool, err error) {
	switch a.state {
	case loginNotStarted:
		// Check for initial response field, as per RFC4422 section 3
		if response == nil {
			challenge = []byte("Username:")
			break
		}
		a.state++
		fallthrough
	case loginWaitingUsername:
		a.username = string(response)
		challenge = []byte("Password:")
	case loginWaitingPassword:
		a.password = string(response)
		err = a.authenticate(a.username, a.password)
		done = true
	default:
		err = sasl.ErrUnexpectedClientResponse
	}

	a.state++
	return
}


func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	log.Println("Mail from:", from)
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	log.Println("Rcpt to:", to)
	s.to = to
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.auth {
		return smtp.ErrAuthRequired
	}
	log.Println("[SMTP] Email received")
	if err := processEmail(r,s); err != nil {
		log.Println("[ERROR] Processing email:", err)
	}
	return nil
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	return nil
}

func processEmail(r io.Reader, s *Session) error {
	
	mr, err := mail.CreateReader(r)
	if err != nil {
		return err
	}

	var subject, textBody string
	header := mr.Header

	if subj, err := header.Subject(); err == nil {
		subject = subj
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			body, _ := io.ReadAll(part.Body)
			textBody += string(body)
		case *mail.AttachmentHeader:
			filename, _ := h.Filename()
			log.Println("[ATTACHMENT] Saving:", filename)
			if err := saveAndSendAttachment(filename, part.Body); err != nil {
				log.Println("[ERROR] Sending attachment:", err)
			}
		}
	}

	message := fmt.Sprintf("📧 %s\n\n*%s*\n%s",s.from, subject, textBody)
	return sendToTelegram(message)
}

func sendToTelegram(text string) error {

	telegramChatID, err := strconv.ParseInt(telegramChat, 10, 64)
	if err != nil {
		return err
	}

	payload, err := json.Marshal(&OutboundMsg{
		ChatID: telegramChatID,
		Text:   text,
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", telegramBot)
	c := &send.Request{
			Path:   url,
			Method: "POST",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: payload,
		}
		_, err = c.Send()
		if err != nil {
			return err
		}
	return nil
}

func saveAndSendAttachment(filename string, r io.Reader) error {
	tmpPath := "/tmp/" + filename
	file, err := os.Create(tmpPath)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := io.Copy(file, r); err != nil {
		return err
	}

	defer os.Remove(tmpPath)
	return sendFileToTelegram(tmpPath)
}

func sendFileToTelegram(filePath string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", telegramBot)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("document", filePath)
	_, _ = io.Copy(part, file)
	err = writer.WriteField("chat_id", telegramChat)
		if err != nil {
		return err
	}
	writer.Close()

	req, _ := http.NewRequest("POST", url, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func main() {
	if telegramBot == "" || telegramChat == "" {
		log.Fatal("[ERROR] Missing required environment variables.")
	}

	be := &Backend{}

	s := smtp.NewServer(be)

	s.Addr = smtpAddr
	s.Domain = smtpDomain
	s.WriteTimeout = 10 * time.Second
	s.ReadTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
