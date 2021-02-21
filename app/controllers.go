package app

import (
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/index.gtpl")
		if err != nil {
			log.Error("[indexController]Cant parse template: " + err.Error())
		}
		t.Execute(w, nil)
	}
}

func msgController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/msg.gtpl")
		channels, err := AppSession.GetChannels()
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		t.Execute(w, channels)
	} else if r.Method == "POST" {
		r.ParseForm()
		var toChannel string
		if r.Form.Get("channelid") != "" {
			toChannel = r.Form.Get("channelid")
		} else {
			toChannel = r.Form.Get("channel")
		}
		m := Message{ChannelID: toChannel, Content: r.Form.Get("message")}
		err := AppSession.SendMessage(m)
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		retMessage := Response{Message: "Mensaje enviado!"}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		t.Execute(w, retMessage)
	}
}

func editController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/edit.gtpl")
		channels, err := AppSession.GetChannels()
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		t.Execute(w, channels)
	} else if r.Method == "POST" {
		r.ParseForm()
		var toChannel string
		if r.Form.Get("channelid") != "" {
			toChannel = r.Form.Get("channelid")
		} else {
			toChannel = r.Form.Get("channel")
		}
		msgID := r.Form.Get("messageid")
		mw, err := AppSession.GetMessage(toChannel, msgID)
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		var t *template.Template
		if mw.IsEmbed {
			t, _ = template.ParseFiles("web/editmsgembed.gtpl")
		} else {
			t, _ = template.ParseFiles("web/editmsg.gtpl")
		}
		t.Execute(w, mw)
	}
}

func editMsgController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		err := AppSession.EditMessage(r.Form.Get("channelid"), r.Form.Get("messageid"), r.Form.Get("message"))
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		retMessage := Response{Message: "Mensaje editado!"}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		t.Execute(w, retMessage)
	}
}

func editMsgEmbedController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		em := parseFormToEmbed(r.Form)
		err := em.Validate()
		if err != nil {
			errorEditMessageEmbedController(w, em, r.Form.Get("channelid"), r.Form.Get("messageid"), err.Error())
			return
		}
		err = AppSession.EditMessageEmbed(r.Form.Get("channelid"), r.Form.Get("messageid"), em)
		if err != nil {
			errorEditMessageEmbedController(w, em, r.Form.Get("channelid"), r.Form.Get("messageid"), err.Error())
			return
		}
		retMessage := Response{Message: "Mensaje editado!"}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		t.Execute(w, retMessage)
	}
}

func msgEmbedController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, _ := template.ParseFiles("web/msgembed.gtpl")
		channels, err := AppSession.GetChannels()
		channels.HexColor = DEFAULT_COLOR
		e := newEmbedMessage()
		channels.Emb = e
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		t.Execute(w, channels)
	} else if r.Method == "POST" {
		r.ParseForm()
		var toChannel string
		if r.Form.Get("channelid") != "" {
			toChannel = r.Form.Get("channelid")
		} else {
			toChannel = r.Form.Get("channel")
		}
		e := parseFormToEmbed(r.Form)
		err := e.Validate()
		if err != nil {
			errorMessageEmbedController(w, e, err.Error())
			return
		}
		e.ChannelID = toChannel
		err = AppSession.SendMessageEmbed(e)
		if err != nil {
			errorMessageEmbedController(w, e, err.Error())
			return
		}
		retMessage := Response{Message: "Mensaje enviado!"}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		t.Execute(w, retMessage)
	}
}

func errorPageController(w http.ResponseWriter, s ServerError) {
	t, _ := template.ParseFiles("web/error.gtpl")
	t.Execute(w, s)
}

func errorMessageEmbedController(w http.ResponseWriter, e EmbedMessage, errMsg string) {
	t, _ := template.ParseFiles("web/msgembed.gtpl")
	channels, err := AppSession.GetChannels()
	if err != nil {
		errorPageController(w, ServerError{Message: err.Error(), Code: 500})
		return
	}
	channels.Emb = e
	channels.ErrorMsg = errMsg
	e.Fields = setFieldsSize(e.Fields)
	channels.HexColor = ColorToHex(e.Color)
	t.Execute(w, channels)
}

func errorEditMessageEmbedController(w http.ResponseWriter, e EmbedMessage, cID string, mID string, errMsg string) {
	t, _ := template.ParseFiles("web/editmsgembed.gtpl")
	channels, err := AppSession.GetChannels()
	if err != nil {
		errorPageController(w, ServerError{Message: err.Error(), Code: 500})
		return
	}
	channels.Emb = e
	channels.ErrorMsg = errMsg
	channels.Emb.ChannelID = cID
	channels.MessageID = mID
	e.Fields = setFieldsSize(e.Fields)
	channels.HexColor = ColorToHex(e.Color)
	t.Execute(w, channels)
}
