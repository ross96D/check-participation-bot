package tg

import (
	"encoding/json"

	"sync"
	"time"
)

type chat struct {
	ID       int64
	mut      *sync.Mutex
	apiToken string
}

func (c chat) sendMessage(sm SendMessage) error {
	c.mut.Lock()
	timer := time.NewTimer(time.Second)
	defer func() {
		<-timer.C
		c.mut.Unlock()
	}()

	b, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	return send(b, "sendMessage", c.apiToken)
}

func (c chat) editMessage(sm EditMessage) error {
	c.mut.Lock()
	timer := time.NewTimer(time.Second)
	defer func() {
		<-timer.C
		c.mut.Unlock()
	}()

	b, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	return send(b, "editMessageText", c.apiToken)
}

func (c chat) deleteMessage(sm DeleteMessage) error {
	c.mut.Lock()
	timer := time.NewTimer(time.Second)
	defer func() {
		<-timer.C
		c.mut.Unlock()
	}()

	b, err := json.Marshal(sm)
	if err != nil {
		return err
	}

	return send(b, "deleteMessage", c.apiToken)
}

type TelegramSender interface {
	SendMessage(SendMessage) error
	EditMessageText(EditMessage) error
	DeleteMessage(DeleteMessage) error
}

type Tg struct {
	chats    *sync.Map
	wg       *sync.WaitGroup
	apiToken string
}

func New(apiToken string) TelegramSender {
	return Tg{
		chats:    &sync.Map{},
		wg:       &sync.WaitGroup{},
		apiToken: apiToken,
	}
}

func (tg Tg) SendMessage(sm SendMessage) error {
	if sm.ParseMode == "MarkdownV2" {
		return tg.sendMessageMarkdown(sm)
	} else {
		return tg.sendMessageText(sm)
	}
}

// TODO make this more general
func (tg Tg) sendMessageMarkdown(sm SendMessage) error {
	v, ok := tg.chats.Load(sm.ChatID)
	if !ok {
		v = chat{ID: sm.ChatID, mut: &sync.Mutex{}, apiToken: tg.apiToken}
		tg.chats.Store(sm.ChatID, v)
	}

	err := call(v.(chat).sendMessage, sm)
	if err != nil {
		return err
	}

	// for _, rune := range sm.Text {
	// 	if len(toSend) > 4080 {
	// 		newSm := sm
	// 		newSm.Text = string(toSend) + "```"
	// 		if count > 0 {
	// 			newSm.Text = "```" + string(toSend)
	// 		}
	// 		err := call(v.(chat).sendMessage, newSm)
	// 		if err != nil {
	// 			return err
	// 		}

	// 		toSend = []byte{}
	// 		count++
	// 	}
	// 	toSend = utf8.AppendRune(toSend, rune)
	// }
	// for i := 0; i < len(sm.Text); {
	// 	newSm := sm
	// 	start := i
	// 	end := i + 4080
	// 	if len(sm.Text) < end {
	// 		end = len(sm.Text)
	// 	}
	// 	newSm.Text = sm.Text[start:end] + "```"
	// 	if i != 0 {
	// 		newSm.Text = "```" + sm.Text
	// 	}
	// 	lenss := len(newSm.Text)
	// 	_ = lenss
	// 	err := call(v.(chat).sendMessage, newSm)
	// 	if err != nil {
	// 		return fmt.Errorf("telegramSendMessage %s %w", newSm.Text, err)
	// 	}
	// 	i += 4080
	// }
	return nil
}

func (tg Tg) sendMessageText(sm SendMessage) error {
	v, ok := tg.chats.Load(sm.ChatID)
	if !ok {
		v = chat{ID: sm.ChatID, mut: &sync.Mutex{}, apiToken: tg.apiToken}
		tg.chats.Store(sm.ChatID, v)
	}
	for i := 0; i < len(sm.Text); {
		newSm := sm
		start := i
		end := i + 4096
		if len(sm.Text) < end {
			end = len(sm.Text)
		}
		newSm.Text = sm.Text[start:end]
		err := call(v.(chat).sendMessage, newSm)
		if err != nil {
			return err
		}
		i += 4096
	}
	return nil
}

func (tg Tg) EditMessageText(sm EditMessage) error {
	v, ok := tg.chats.Load(sm.ChatID)
	if !ok {
		v = chat{ID: sm.ChatID, mut: &sync.Mutex{}, apiToken: tg.apiToken}
		tg.chats.Store(sm.ChatID, v)
	}
	return call(v.(chat).editMessage, sm)
}

func (tg Tg) DeleteMessage(sm DeleteMessage) error {
	v, ok := tg.chats.Load(sm.ChatID)
	if !ok {
		v = chat{ID: sm.ChatID, mut: &sync.Mutex{}, apiToken: tg.apiToken}
		tg.chats.Store(sm.ChatID, v)
	}

	return call(v.(chat).deleteMessage, sm)
}

type t interface {
	SendMessage | EditMessage | DeleteMessage
}

func call[T t](fn func(T) error, p T) error {
	for {
		err := fn(p)
		if err == nil {
			return nil
		}
		if err, ok := err.(TelegramError); ok {
			if err.ErrorCode != 429 {
				return err
			}
			if val, ok := err.Parameters["retry_after"]; ok {
				if val, ok := val.(int); ok {
					time.Sleep(time.Duration(val) * time.Millisecond)
					continue
				}
			}
		}
		return err
	}
}
