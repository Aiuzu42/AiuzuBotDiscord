package app

import (
	"net/url"
	"strconv"
)

const (
	DEFAULT_COLOR = "#ff0000"
)

func parseFormToEmbed(form url.Values) EmbedMessage {
	em := EmbedMessage{}
	if t := form.Get("title"); t != "" {
		em.Title = t
	}
	if d := form.Get("description"); d != "" {
		em.Content = d
	}
	em.Color = ColorToInt(form.Get("color"))
	for i := 0; i < 25; i++ {
		t := form.Get("field" + strconv.Itoa(i) + "title")
		c := form.Get("field" + strconv.Itoa(i) + "content")
		i := form.Get("field" + strconv.Itoa(i) + "inline")
		if t != "" || c != "" {
			f := Field{Name: t, Value: c}
			if i == "true" {
				f.Inline = true
			}
			em.Fields = append(em.Fields, f)
		}
	}
	return em
}

func ColorToInt(hex string) int {
	hex = hex[1:]
	i64, err := strconv.ParseInt(hex, 16, 64)
	if err != nil {
		return 0
	}
	return int(i64)
}

func ColorToHex(i int) string {
	return "#" + strconv.FormatInt(int64(i), 16)
}

func newEmbedMessage() EmbedMessage {
	em := EmbedMessage{}
	for i := 0; i < 25; i++ {
		em.Fields = append(em.Fields, Field{})
	}
	return em
}

func setFieldsSize(e []Field) []Field {
	for i := len(e); i < 25; i++ {
		e = append(e, Field{})
	}
	return e
}
