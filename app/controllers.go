package app

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/aiuzu42/AiuzuBotDiscord/config"

	log "github.com/sirupsen/logrus"
)

func indexController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		t, err := template.ParseFiles("web/index.gtpl")
		if err != nil {
			log.Error("[indexController]Cant parse template: " + err.Error())
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Error("[indexController]Cant execute template: " + err.Error())
		}
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
		err = t.Execute(w, channels)
		if err != nil {
			log.Error("[msgController]Cant execute template: " + err.Error())
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Error("[msgController]Cant parse form: " + err.Error())
		}
		var toChannel string
		if r.Form.Get("channelid") != "" {
			toChannel = r.Form.Get("channelid")
		} else {
			toChannel = r.Form.Get("channel")
		}
		m := Message{ChannelID: toChannel, Content: r.Form.Get("message")}
		d := r.Form.Get("delay")
		delay, convErr := strconv.Atoi(d)
		rMsg := ""
		if convErr == nil && delay > 0 {
			go sendMessageDelayed(m, delay)
			rMsg = "El mensaje se enviara en " + d + " minutos!"
		} else {
			err := AppSession.SendMessage(m)
			if err != nil {
				errorPageController(w, ServerError{Message: err.Error(), Code: 500})
				return
			}
			rMsg = "Mensaje enviado!"
		}
		retMessage := Response{Message: rMsg}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		err = t.Execute(w, retMessage)
		if err != nil {
			log.Error("[msgController]Cant execute template2: " + err.Error())
		}
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
		err = t.Execute(w, channels)
		if err != nil {
			log.Error("[editController]Cant execute template: " + err.Error())
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Error("[editController]Cant parse form: " + err.Error())
		}
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
		err = t.Execute(w, mw)
		if err != nil {
			log.Error("[editController]Cant execute template2: " + err.Error())
		}
	}
}

func editMsgController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Error("[editMsgController]Cant parse form: " + err.Error())
		}
		err = AppSession.EditMessage(r.Form.Get("channelid"), r.Form.Get("messageid"), r.Form.Get("message"))
		if err != nil {
			errorPageController(w, ServerError{Message: err.Error(), Code: 500})
			return
		}
		retMessage := Response{Message: "Mensaje editado!"}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		err = t.Execute(w, retMessage)
		if err != nil {
			log.Error("[editMsgController]Cant execute template: " + err.Error())
		}
	}
}

func editMsgEmbedController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Error("[editMsgEmbedController]Cant parse form: " + err.Error())
		}
		em := parseFormToEmbed(r.Form)
		err = em.Validate()
		if err != nil {
			errorEditMessageEmbedController(w, em, r.Form.Get("channelid"), r.Form.Get("messageid"), err.Error())
			return
		}
		retMessage := Response{}
		if r.Form.Get("clone") == "yes" {
			em.ChannelID = r.Form.Get("channel")
			err = AppSession.SendMessageEmbed(em)
			if err != nil {
				errorEditMessageEmbedController(w, em, r.Form.Get("channelid"), r.Form.Get("messageid"), err.Error())
				return
			}
			retMessage.Message = "Mensaje enviado!"
		} else {
			err = AppSession.EditMessageEmbed(r.Form.Get("channelid"), r.Form.Get("messageid"), em)
			if err != nil {
				errorEditMessageEmbedController(w, em, r.Form.Get("channelid"), r.Form.Get("messageid"), err.Error())
				return
			}
			retMessage.Message = "Mensaje editado!"
		}

		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		err = t.Execute(w, retMessage)
		if err != nil {
			log.Error("[editMsgEmbedController]Cant execute form: " + err.Error())
		}
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
		err = t.Execute(w, channels)
		if err != nil {
			log.Error("[msgEmbedController]Cant execute template: " + err.Error())
		}
	} else if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			log.Error("[msgEmbedController]Cant parse form: " + err.Error())
		}
		var toChannel string
		if r.Form.Get("channelid") != "" {
			toChannel = r.Form.Get("channelid")
		} else {
			toChannel = r.Form.Get("channel")
		}
		e := parseFormToEmbed(r.Form)
		err = e.Validate()
		if err != nil {
			errorMessageEmbedController(w, e, err.Error())
			return
		}
		e.ChannelID = toChannel
		d := r.Form.Get("delay")
		delay, convErr := strconv.Atoi(d)
		rMsg := ""
		if convErr == nil && delay > 0 {
			go sendMessageEmbedDelayed(e, delay)
			rMsg = "El mensaje se enviara en " + d + " minutos!"
		} else {
			err = AppSession.SendMessageEmbed(e)
			if err != nil {
				errorMessageEmbedController(w, e, err.Error())
				return
			}
			rMsg = "Mensaje enviado!"
		}
		retMessage := Response{Message: rMsg}
		t, _ := template.ParseFiles("web/succesmsg.gtpl")
		err = t.Execute(w, retMessage)
		if err != nil {
			log.Error("[msgEmbedController]Cant execute template2: " + err.Error())
		}
	}
}

func errorPageController(w http.ResponseWriter, s ServerError) {
	t, _ := template.ParseFiles("web/error.gtpl")
	err := t.Execute(w, s)
	if err != nil {
		log.Error("[errorPageController]Cant execute template2: " + err.Error())
	}
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
	err = t.Execute(w, channels)
	if err != nil {
		log.Error("[errorMessageEmbedController]Cant execute template2: " + err.Error())
	}
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
	err = t.Execute(w, channels)
	if err != nil {
		log.Error("[errorEditMessageEmbedController]Cant execute template: " + err.Error())
	}
}

func sendMessageDelayed(m Message, min int) {
	time.Sleep(time.Duration(min) * time.Minute)
	err := AppSession.SendMessage(m)
	if err != nil {
		log.Error("[sendMessageDelayed]Error sending message: " + err.Error())
	}
}

func sendMessageEmbedDelayed(em EmbedMessage, min int) {
	time.Sleep(time.Duration(min) * time.Minute)
	err := AppSession.SendMessageEmbed(em)
	if err != nil {
		log.Error("[sendMessageEmbedDelayed]Error sending message: " + err.Error())
	}
}

func sendMsgDirectController(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		j := struct {
			Msg  string `json:"message"`
			ChId string `json:"channelId"`
		}{}
		u, p, ok := r.BasicAuth()
		if !ok || !checkCredentials(u, p) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&j)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = AppSession.SendMessage(Message{ChannelID: j.ChId, Content: j.Msg})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func checkCredentials(user string, pass string) bool {
	if user != config.ApiUsr || pass != config.ApiPass {
		return false
	}
	return true
}
