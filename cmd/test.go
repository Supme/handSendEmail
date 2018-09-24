package main

import (
	"bytes"
	"fmt"
	"github.com/supme/handSendEmail/email"
	"github.com/supme/handSendEmail/message"
	"log"
	"os"
)

func main() {
	var (
		err   error
		iface *email.Iface
	)

	ifaces, err := email.GetInterfaces(map[string]string{"192.168.0.10": "1.2.3.4"})
	for i := range ifaces {
		fmt.Printf("- %d (Addr: '%s', Hostname: '%s'\n", i, ifaces[i].IP, ifaces[i].Hostname)
	}
	var n int
	fmt.Print("Выберите интерфейс через который будем слать: ")
	fmt.Scanln(&n)
	iface = ifaces[n]
	fmt.Printf("Выбран интерфейс %s ('%s')\n", iface.IP, iface.Hostname)

	e := message.NewMessage().
		From(message.NewMail("Алексей", "alexey@domain.tld")).
		To(message.NewMail("Василий", "vasiliy@domain.tld")).
		To(message.NewMail("Фёдор", "fedor@domain.tld")).
		Cc(message.NewMail("Василий 1", "vasiliy_1@domain.tld")).
		Cc(message.NewMail("Фёдор 1", "fedor_1@domain.tld")).
		Bcc(message.NewMail("Василий 2", "vasiliy_2@domain.tld")).
		Bcc(message.NewMail("Фёдор 2", "fedor_2@domain.tld")).
		Subject("Тестовый email").
		TextHTML("<h1>Привет! Это я.</h1><br><img src=\"cid:me.gif\"/><br><h2>Съешь ещё этих мягких французских булок да выпей чаю</h2>").
		TextPlain("Привет! Это я.\n[картинка меня]\nСъешь ещё этих мягких французских булок да выпей чаю")

	fRelated, err := os.Open("../testdata/me.gif")
	if err != nil {
		log.Fatal(err)
	}
	e.AddRelatedFile(fRelated)

	fAttachment, err := os.Open("../testdata/the_little_go_book.pdf")
	if err != nil {
		log.Fatal(err)
	}
	e.AddAttachmentFile(fAttachment)

	buf := &bytes.Buffer{}
	e.Write(buf)

	for _, to := range e.GetRecipientEmails() {
		mail := email.NewSmtp(iface)
		fmt.Println("Connect...\nHELO", iface.Hostname)
		err = mail.CommandConnectAndHello(to)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")
		//		time.Sleep(time.Second)

		fmt.Println("FROM: ", e.GetFromEmail())
		err = mail.CommandFrom(e.GetFromEmail())
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")
		//		time.Sleep(time.Second)

		fmt.Println("RCPT: ", to)
		err = mail.CommandRcpt(to)
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")
		//		time.Sleep(time.Second)

		fmt.Println("DATA ...you message data...")
		//fmt.Printf("DATA\n%s\n.\n\n", buf.String())
		err = mail.CommandData(buf.Bytes())
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")

		fmt.Println("QUIT")
		err = mail.CommandQuit()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")
	}
}
