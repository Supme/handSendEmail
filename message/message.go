// https://tools.ietf.org/html/rfc2076
package message

import (
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	boundaryMixed            = "===============1_MIXED========"
	boundaryMixedBegin       = "--" + boundaryMixed + "\r\n"
	boundaryMixedEnd         = "--" + boundaryMixed + "--\r\n"
	boundaryRelated          = "===============2_RELATED======"
	boundaryRelatedBegin     = "--" + boundaryRelated + "\r\n"
	boundaryRelatedEnd       = "--" + boundaryRelated + "--\r\n"
	boundaryAlternative      = "===============3_ALTERNATIVE=="
	boundaryAlternativeBegin = "--" + boundaryAlternative + "\r\n"
	boundaryAlternativeEnd   = "--" + boundaryAlternative + "--\r\n"
)

type Mail struct {
	name  string
	email string
}

func NewMail(name, email string) Mail {
	return Mail{
		name:  name,
		email: email,
	}
}

func (m Mail) String() string {
	if m.name == "" {
		return m.email
	}
	return mime.BEncoding.Encode("utf-8", m.name) + "<" + m.email + ">"
}

func JoinMails(ms []Mail) string {
	msStr := make([]string, len(ms))
	for i := range ms {
		msStr[i] = ms[i].String()
	}
	return strings.Join(msStr, ", ")
}

type Message struct {
	from           Mail
	to             []Mail
	cc             []Mail
	bcc            []Mail
	headers        map[string]string
	subject        string
	textHTML       string
	textPlain      string
	relatedFile    []*os.File
	attachmentFile []*os.File
}

func NewMessage() *Message {
	return new(Message)
}

func (m *Message) From(name, email string) {
	m.from = NewMail(name, email)
}

func (m *Message) To(name, email string) {
	m.to = append(m.to, NewMail(name, email))
}

func (m *Message) Cc(name, email string) {
	m.cc = append(m.cc, NewMail(name, email))
}

func (m *Message) Bcc(name, email string) {
	m.bcc = append(m.bcc, NewMail(name, email))
}

func (m *Message) Subject(subject string) {
	m.subject = subject
}

func (m *Message) AddHeaders(headers map[string]string) {
	for k, v := range headers {
		m.headers[k] = v
	}
}

func (m *Message) TextHTML(textHTML string) {
	m.textHTML = textHTML
}

func (m *Message) TextPlain(textPlain string) {
	m.textPlain = textPlain
}

func (m *Message) AddRelatedFile(file *os.File) {
	m.relatedFile = append(m.relatedFile, file)
}

func (m *Message) AddAttachmentFile(file *os.File) {
	m.attachmentFile = append(m.attachmentFile, file)
}

func (m Message) Write(w io.Writer) {
	m.HeaderWrite(w)
	m.BodyWrite(w)
}

func (m Message) HeaderWrite(w io.Writer) {
	w.Write([]byte("Date: " + time.Now().Format(time.RFC1123Z) + "\r\n"))
	w.Write([]byte("From: " + m.from.String() + "\r\n"))
	w.Write([]byte("To: " + JoinMails(m.to) + "\r\n"))
	if len(m.cc) > 0 {
		w.Write([]byte("Cc: " + JoinMails(m.cc) + "\r\n"))
	}
	if len(m.cc) > 0 {
		w.Write([]byte("Bcc: " + JoinMails(m.bcc) + "\r\n"))
	}
	w.Write([]byte("Subject: " + mime.BEncoding.Encode("utf-8", m.subject) + "\r\n"))
}

func (m Message) BodyWrite(w io.Writer) {
	// Начинаем наше multipart/mixed письмо
	w.Write([]byte("Content-Type: multipart/mixed; boundary=" + boundaryMixed + "\"\r\n"))
	w.Write([]byte("MIME-Version: 1.0\r\n"))
	w.Write([]byte("\r\n"))
	w.Write([]byte("This is a multi-part message in MIME format.\r\n"))
	w.Write([]byte("\r\n"))

	// У нас будут зависящие друг от друга блоки с mixed разделителем вверху
	{
		w.Write([]byte(boundaryMixedBegin))
		w.Write([]byte("Content-Type: multipart/related; boundary=\"" + boundaryRelated + "\"\r\n"))
		w.Write([]byte("MIME-Version: 1.0\r\n"))
		w.Write([]byte("\r\n"))

		// Первым зависящим блоком будут альтернативные версии с related разделителем вверху
		{
			w.Write([]byte(boundaryRelatedBegin))

			{
				w.Write([]byte("Content-Type: multipart/alternative; boundary=\"" + boundaryAlternative + "\"\r\n"))
				w.Write([]byte("MIME-Version: 1.0\r\n"))
				w.Write([]byte("\r\n"))

				// Если textHTML не пуст добавляем альтернативный блок text/html с alternative разделителем вверху
				if m.textHTML != "" {
					w.Write([]byte(boundaryAlternativeBegin))
					w.Write([]byte("Content-Type: text/html; charset=\"utf-8\"\r\n"))
					w.Write([]byte("MIME-Version: 1.0\r\n"))
					w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
					w.Write([]byte("\r\n"))
					// Пишем textHTML кодируя его в base64 с переводом строки и возвратом каретки каждые 76 символов
					base64TextWriter(w, m.textHTML)
					w.Write([]byte("\r\n"))
					w.Write([]byte("\r\n"))
				}

				// Если textPlain не пуст добавляем блок text/plain с alternative разделителем вверху
				if m.textPlain != "" {
					w.Write([]byte(boundaryAlternativeBegin))
					w.Write([]byte("Content-Type: text/plain; charset=\"utf-8\"\r\n"))
					w.Write([]byte("MIME-Version: 1.0\r\n"))
					w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
					w.Write([]byte("\r\n"))
					// Пишем textPlain кодируя аналогично textHTML
					base64TextWriter(w, m.textPlain)
					w.Write([]byte("\r\n"))
					w.Write([]byte("\r\n"))
				}

				// Закрываем блок альтернатив
				w.Write([]byte(boundaryAlternativeEnd))
				w.Write([]byte("\r\n"))

			}

			// Если есть зависящие файлы
			if len(m.relatedFile) > 0 {
				// Будем все отправлять
				for i := range m.relatedFile {
					// Сперва соберём необходимую информацию о файле
					var (
						// нам нужно имя файла
						fileName string
						// его размер
						fileSize string
						// и mime тип
						fileMime string
					)
					fileName = filepath.Base(m.relatedFile[i].Name())
					info, _ := m.relatedFile[i].Stat()
					fileSize = strconv.FormatInt(info.Size(), 10)
					buf := make([]byte, 512)
					m.relatedFile[i].Read(buf)
					fileMime = http.DetectContentType(buf)
					// Вернём курсор чтения файла в начало
					m.relatedFile[i].Seek(0, 0)
					// Пишем заголовок для файла с related разделителем вверху
					w.Write([]byte(boundaryRelatedBegin))
					w.Write([]byte("Content-Type: " + fileMime + "; name=\"" + fileName + "\"\r\n"))
					w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
					w.Write([]byte("Content-ID: <" + fileName + ">\r\n"))
					w.Write([]byte("Content-Disposition: inline; filename=\"" + fileName + "\"; size=" + fileSize + ";\r\n"))
					w.Write([]byte("\r\n"))
					// Пишем файл кодируя в base64 с переносами строк через каждые 76 символов
					base64FileWriter(w, m.relatedFile[i])
					w.Write([]byte("\r\n"))
				}
			}

			// Закрываем блок зависящих
			w.Write([]byte(boundaryRelatedEnd))
			w.Write([]byte("\r\n"))
		}

		// Если есть файлы для вложения
		if len(m.attachmentFile) > 0 {
			// Будем все отправлять
			for i := range m.attachmentFile {
				// Сперва соберём необходимую информацию о файле
				var (
					// нам нужно имя файла
					fileName string
					// его размер
					fileSize string
					// и mime тип
					fileMime string
				)
				fileName = filepath.Base(m.attachmentFile[i].Name())
				info, _ := m.attachmentFile[i].Stat()
				fileSize = strconv.FormatInt(info.Size(), 10)
				buf := make([]byte, 512)
				m.attachmentFile[i].Read(buf)
				fileMime = http.DetectContentType(buf)
				// Вернём курсор чтения файла в начало
				m.attachmentFile[i].Seek(0, 0)
				// Пишем заголовок для файла с mixed разделителем вверху
				w.Write([]byte(boundaryMixedBegin))
				w.Write([]byte("Content-Type: " + fileMime + "; name=\"" + fileName + "\"\r\n"))
				w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
				w.Write([]byte("Content-Disposition: attachment; filename=\"" + fileName + "\"; size=" + fileSize + ";\r\n"))
				w.Write([]byte("\r\n"))
				// Пишем файл кодируя в base64 с переносами строк через каждые 76 символов
				base64FileWriter(w, m.attachmentFile[i])
				w.Write([]byte("\r\n"))
			}
		}

	}

	// И закрываем наше сообщение
	w.Write([]byte(boundaryMixedEnd))
}
