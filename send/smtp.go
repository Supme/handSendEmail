package send

import (
	"fmt"
	"net"
	"net/smtp"
	"strings"
	"time"
)

type SMTP struct {
	addr      net.Addr
	LocalName string
	conn      net.Conn
	client    *smtp.Client
}

type Interface struct {
	Addr     net.Addr
	Hostname string
}

func NewSmtp(iface Interface) (*SMTP, error) {
	var s SMTP
	s.LocalName = iface.Hostname
	return &s, nil
}

func (s *SMTP) CommandConnectAndHello(emailTo string) error {
	var err error
	splitEmail := strings.Split(emailTo, "@")
	if len(splitEmail) != 2 {
		return fmt.Errorf("bad email format")
	}
	if err = s.connect(splitEmail[1]); err != nil {
		return err
	}
	return nil
}

func (s *SMTP) CommandVerify(email string) error {
	fmt.Printf("%#v\n", s.client)
	if err := s.client.Verify(email); err != nil {
		return err
	}
	return nil
}

func (s *SMTP) CommandFrom(email string) error {
	if err := s.client.Mail(email); err != nil {
		return err
	}
	return nil
}

func (s *SMTP) CommandRcpt(email string) error {
	if err := s.client.Rcpt(email); err != nil {
		return err
	}
	return nil
}

func (s *SMTP) CommandData(data []byte) error {
	w, err := s.client.Data()
	if err != nil {
		return err
	}
	if _, err = w.Write(data); err != nil {
		return err
	}
	return w.Close()
}

func (s *SMTP) CommandQuit() error {
	return s.client.Quit()
}

func (s *SMTP) CommandClose() error {
	return s.client.Close()
}

func (s *SMTP) connect(host string) error {
	var (
		err  error
		errs []string
		mxs  []*net.MX
	)
	if mxs, err = net.LookupMX(host); err != nil {
		return err
	}
	dialer := net.Dialer{
		Timeout: time.Second * 15,
	}
	for i := range mxs {
		if s.conn, err = dialer.Dial("tcp", mxs[i].Host+":25"); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if s.client, err = smtp.NewClient(s.conn, host); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		if err = s.client.Hello(s.LocalName); err != nil {
			errs = append(errs, err.Error())
			continue
		}
		return nil
	}
	return fmt.Errorf("can not connect, errors: %s", strings.Join(errs, "; "))
}
