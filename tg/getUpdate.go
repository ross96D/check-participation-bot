package tg

import (
	"encoding/json"
	"errors"

	"github.com/ross96D/battle-log-parser/assert"
)

type AllowedUpdate string

func (au AllowedUpdate) String() string {
	switch au {
	case AUMessage, AUEditedChannelPost, AUCallbackQuery:
		return string(au)
	default:
		panic("unknown allowUpdate " + string(au))
	}
}

const (
	AUMessage           AllowedUpdate = "message"
	AUEditedChannelPost AllowedUpdate = "edited_channel_post"
	AUCallbackQuery     AllowedUpdate = "callback_query"
)

type GetUpdatesParams struct {
	Offset         int64           `json:"offset"`
	Limit          uint8           `json:"limit"`
	Timeout        uint            `json:"timeout"`
	AllowedUpdates []AllowedUpdate `json:"allowed_updates"`
}

func GetUpdates(apiToken string, params GetUpdatesParams) ([]Update, error) {
	if params.Limit == 0 {
		params.Limit = 100
	}
	assert.Assert(1 <= params.Limit && params.Limit <= 100)

	b, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}

	body, err := sendGetBody(b, "getUpdates", apiToken)
	if err != nil {
		return nil, err
	}

	type GetUpdatesResponse struct {
		Ok     bool     `json:"ok"`
		Result []Update `json:"result,omitempty"`
	}

	var updates GetUpdatesResponse
	dec := json.NewDecoder(body)
	err = dec.Decode(&updates)

	if !updates.Ok {
		return nil, errors.New("updates false")
	}

	return updates.Result, err
}
