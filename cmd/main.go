package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/supme/handSendEmail/message"
	"github.com/supme/handSendEmail/send"
	"html/template"
	"log"
	"net/http"
	"os"
)

const addr = ":8080"

func main() {
	var (
		err   error
		iface *send.Iface
	)

	ifaces, err := send.GetInterfaces(map[string]string{"192.168.0.10": "1.2.3.4"})
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
		TextHTML("<h1>Съешь ещё этих мягких французских булок да выпей чаю</h1><br><img src='cid:me.gif' alt='Super me'>").
		TextPlain("Съешь ещё этих мягких французских булок да выпей чаю")

	fRelated, err := os.Open("../testdata/me.gif")
	if err != nil {
		log.Fatal(err)
	}
	e.AddRelatedFile(fRelated)

	//fAttachment, err := os.Open("../message.txt")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//e.AddAttachmentFile(fAttachment)

	buf := &bytes.Buffer{}
	e.Write(buf)

	for _, to := range e.GetRecipientEmails() {
		mail := send.NewSmtp(iface)
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
		//		time.Sleep(time.Second*10)

		fmt.Println("QUIT")
		err = mail.CommandQuit()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Ok")
		//time.Sleep(time.Second)

		//fmt.Println("CLOSE")
		//err = mail.CommandClose()
		//if err != nil {
		//	log.Println(err)
		//	return
		//}
		//fmt.Println("Ok")
		//time.Sleep(time.Second)
	}

	return

	mux := http.NewServeMux()
	mux.HandleFunc("/", root)

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		ico, _ := base64.StdEncoding.DecodeString("AAABAAEAEBAAAAEAIABoBAAAFgAAACgAAAAQAAAAIAAAAAEAIAAAAAAAAAQAABILAAASCwAAAAAAAAAAAAByGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/8q2uP9yGSL/yra4/3IZIv/Ktrj/yra4/3IZIv9yGSL/yra4/8q2uP9yGSL/yra4/3IZIv/Ktrj/chki/3IZIv/Ktrj/chki/+je3/9yGSL/yra4/3IZIv/Ktrj/chki/8q2uP9yGSL/chki/8q2uP9yGSL/yra4/3IZIv9yGSL/yra4/+je3//Ktrj/chki/8q2uP9yGSL/yra4/3IZIv/Ktrj/yra4/3IZIv/Ktrj/yra4/3IZIv9yGSL/chki/+je3/9yGSL/yra4/3IZIv/Ktrj/chki/8q2uP9yGSL/yra4/3IZIv9yGSL/yra4/3IZIv/Ktrj/chki/3IZIv/Ktrj/chki/8q2uP9yGSL/yra4/8q2uP9yGSL/chki/8q2uP/Ktrj/chki/8q2uP/Ktrj/yra4/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/+je3//o3t//6N7f/+je3//o3t//6N7f/+je3/9yGSL/6N7f/+je3//o3t//6N7f/+je3//o3t//chki/3IZIv/o3t//yra4/8q2uP/Ktrj/yra4/8q2uP/Ktrj/chki/8q2uP/Ktrj/yra4/8q2uP/Ktrj/6N7f/3IZIv9yGSL/6N7f/8q2uP9yGSL/chki/3IZIv/Ktrj/6N7f/3IZIv/Ktrj/yra4/3IZIv9yGSL/yra4/+je3/9yGSL/chki/+je3//Ktrj/chki/8q2uP/o3t//6N7f/8q2uP9yGSL/6N7f/+je3/9yGSL/chki/8q2uP/o3t//chki/3IZIv/o3t//yra4/3IZIv/o3t//yra4/8q2uP/Ktrj/chki/8q2uP/Ktrj/chki/3IZIv/Ktrj/6N7f/3IZIv9yGSL/6N7f/+je3/9yGSL/chki/3IZIv9yGSL/chki/3IZIv/Ktrj/yra4/8q2uP/Ktrj/6N7f/+je3/9yGSL/chki/+je3//o3t//6N7f/+je3//o3t//6N7f/+je3/9yGSL/6N7f/+je3//o3t//6N7f/+je3//o3t//chki/3IZIv/Ktrj/yra4/8q2uP/Ktrj/yra4/8q2uP/Ktrj/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/chki/3IZIv9yGSL/AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA==")
		w.Write(ico)
	})

	fmt.Println("Listen on", addr)
	panic(http.ListenAndServe(addr, mux))
}

func root(w http.ResponseWriter, _ *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusOK)
		log.Print(err)
	}

	data := map[string]string{
		"_Title": "Index page",
	}

	w.Header().Set("Content-Type", "text/html")
	err = tmpl.ExecuteTemplate(w, "index", data)
	if err != nil {
		log.Print(err)
	}
}
