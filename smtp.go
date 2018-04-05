package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

type emailStruct struct {
	textHTML       string
	textPlain      string
	relatedFile    []*os.File
	attachmentFile []*os.File
}

func emailBody(w io.Writer, e emailStruct) {
	// Начинаем наше multipart/mixed письмо
	w.Write([]byte("Content-Type: multipart/mixed; boundary=" + boundaryMixed + "\"\r\n"))
	w.Write([]byte("MIME-Version: 1.0\r\n"))
	w.Write([]byte("\r\n"))
	w.Write([]byte("This is a multi-part message in MIME format.\r\n"))
	w.Write([]byte("\r\n"))

	// У нас будут зависящие друг от друга блоки с mixed разделителем вверху
	w.Write([]byte(boundaryMixedBegin))
	w.Write([]byte("Content-Type: multipart/related; boundary=\"" + boundaryRelated + "\"\r\n"))
	w.Write([]byte("MIME-Version: 1.0\r\n"))
	w.Write([]byte("\r\n"))

	// Первым зависящим блоком будут альтернативные версии с related разделителем вверху
	w.Write([]byte(boundaryRelatedBegin))
	w.Write([]byte("Content-Type: multipart/alternative; boundary=\"" + boundaryAlternative + "\"\r\n"))
	w.Write([]byte("MIME-Version: 1.0\r\n"))
	w.Write([]byte("\r\n"))

	// Если textHTML не пуст добавляем альтернативный блок text/html с alternative разделителем вверху
	if e.textHTML != "" {
		w.Write([]byte(boundaryAlternativeBegin))
		w.Write([]byte("Content-Type: text/html; charset=\"utf-8\"\r\n"))
		w.Write([]byte("MIME-Version: 1.0\r\n"))
		w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
		w.Write([]byte("\r\n"))
		// Пишем textHTML кодируя его в base64 с переводом строки и возвратом каретки каждые 76 символов
		base64TextWriter(w, e.textHTML)
		w.Write([]byte("\r\n"))
		w.Write([]byte("\r\n"))
	}

	// Если textPlain не пуст добавляем блок text/plain с alternative разделителем вверху
	if e.textPlain != "" {
		w.Write([]byte("Content-Type: text/plain; charset=\"utf-8\"\r\n"))
		w.Write([]byte("MIME-Version: 1.0\r\n"))
		w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
		w.Write([]byte("\r\n"))
		// Пишем textPlain кодируя аналогично textHTML
		base64TextWriter(w, e.textPlain)
		w.Write([]byte("\r\n"))
		w.Write([]byte("\r\n"))
	}

	// Закрываем блок  альтернатив
	w.Write([]byte(boundaryAlternativeEnd))
	w.Write([]byte("\r\n"))

	// Если есть зависящие файлы
	if len(e.relatedFile) > 0 {
		// Будем все отправлять
		for i := range e.relatedFile {
			// Сперва соберём необходимую информацию о файле
			var (
				// нам нужно имя файла
				name string
				// его размер
				size string
				// и mime тип
				mime string
			)
			name = filepath.Base(e.relatedFile[i].Name())
			info, _ := e.relatedFile[i].Stat()
			size = strconv.FormatInt(info.Size(), 10)
			buf := make([]byte, 512)
			e.relatedFile[i].Read(buf)
			mime = http.DetectContentType(buf)
			// Вернём курсор чтения файла в начало
			e.relatedFile[i].Seek(0, 0)
			// Пишем заголовок для файла с related разделителем вверху
			w.Write([]byte(boundaryRelatedBegin))
			w.Write([]byte("Content-Type: " + mime + "; name=\"" + name + "\"\r\n"))
			w.Write([]byte("Content-Transfer-Encoding: base64\r\n"))
			w.Write([]byte("Content-ID: <" + name + ">\r\n"))
			w.Write([]byte("Content-Disposition: inline; filename=\"" + name + "\"; size=" + size + ";\r\n"))
			w.Write([]byte("\r\n"))
			// Пишем файл кодируя в base64 с переносами строк через каждые 76 символов
			base64FileWriter(w, e.relatedFile[i])
			w.Write([]byte("\r\n"))
			w.Write([]byte("\r\n"))
		}
	}

	// Закрываем блок зависящих
	w.Write([]byte(boundaryRelatedEnd))
	w.Write([]byte("\r\n"))
	// И закрываем наше сообщение
	w.Write([]byte(boundaryMixedEnd))
}
