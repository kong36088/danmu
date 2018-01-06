package danmu


import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)


func StaticHandler(res http.ResponseWriter, req *http.Request) {
	t, err := template.ParseFiles("I:\\ubuntu14.04\\share\\docker\\go\\www\\src\\github.com\\kong36088\\danmu\\client\\index.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	err = WriteTemplateToHttpResponse(res, t)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
	if t == nil || res == nil {
		return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
	}
	var buf bytes.Buffer
	err := t.Execute(&buf, nil)
	if err != nil {
		return err
	}
	res.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, err = res.Write(buf.Bytes())
	return err
}