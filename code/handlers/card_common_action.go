package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	larkcard "github.com/larksuite/oapi-sdk-go/v3/card"
)

type CardHandlerMeta func(cardMsg CardMsg, m MessageHandler) CardHandlerFunc

type CardHandlerFunc func(ctx context.Context, cardAction *larkcard.CardAction) (
	interface{}, error)

var ErrNextHandler = fmt.Errorf("next handler")

func NewCardHandler(m MessageHandler) CardHandlerFunc {
	handlers := []CardHandlerMeta{
		NewClearCardHandler,
		NewPicResolutionHandler,
		NewVisionResolutionHandler,
		NewPicTextMoreHandler,
		NewPicModeChangeHandler,
		NewRoleTagCardHandler,
		NewRoleCardHandler,
		NewAIModeCardHandler,
		NewVisionModeChangeHandler,
	}

	return func(ctx context.Context, cardAction *larkcard.CardAction) (interface{}, error) {
		var cardMsg CardMsg
		actionValue := cardAction.Action.Value
		actionValueJson, _ := json.Marshal(actionValue)
		if err := json.Unmarshal(actionValueJson, &cardMsg); err != nil {
			return nil, err
		}
		//pp.Println(cardMsg)
		//logger.Debug("cardMsg ", cardMsg)
		for _, handler := range handlers {
			h := handler(cardMsg, m)
			i, err := h(ctx, cardAction)
			if err == ErrNextHandler {
				continue
			}
			return i, err
		}
		return nil, nil
	}
}
