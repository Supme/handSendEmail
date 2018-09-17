package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/supme/handSendEmail/message"
	"html/template"
	"log"
	"net/http"
	"os"
)

const addr = ":8080"

func main() {
	buf := bytes.NewBufferString("Email body:\r\n")
	e := message.NewMessage()

	e.From("Алексей", "alexey@domain.tld")
	e.To("Василий", "vasiliy@domain.tld")
	e.To("Фёдор", "fedor@domain.tld")
	e.Cc("Василий 1", "vasiliy_1@domain.tld")
	e.Cc("Фёдор 1", "fedor_1@domain.tld")
	e.Bcc("Василий 2", "vasiliy_2@domain.tld")
	e.Bcc("Фёдор 2", "fedor_2@domain.tld")
	e.Subject("Тестовый email")
	e.TextHTML("<h1>Съешь ещё этих мягких французских булок да выпей чаю</h1>")
	e.TextPlain("Съешь ещё этих мягких французских булок да выпей чаю")
	fRelated, err := os.Open("template/index.html")
	if err != nil {
		log.Fatal(err)
	}
	e.AddRelatedFile(fRelated)

	fAttachment, err := os.Open("../message.txt")
	if err != nil {
		log.Fatal(err)
	}
	e.AddAttachmentFile(fAttachment)

	e.Write(buf)
	fmt.Println(buf.String())
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
