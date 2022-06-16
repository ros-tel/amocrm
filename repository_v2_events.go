// Copyright (c) 2021 Alexey Khan
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package amocrm

import (
	"fmt"
	"net/http"
)

const (
	eventsV2endpoint endpoint = "events"
)

type (
	eventAdd struct {
		Add []Event `json:"add"`
	}

	Event struct {
		Type        string `json:"type"`         // Тип уведомления – phone_call
		PhoneNumber string `json:"phone_number"` // Номер телефона на который поступает звонок. Можно передавать в любом формате
		Users       []int  `json:"users"`        // Пользователи для которых будет отправлено уведомление. Если не передавать этот параметр, то уведомление будет отправлено для всех пользователей
	}

	EventEmbeddedItem struct {
		ElementID   int    `json:"element_id"`
		ElementType int    `json:"element_type"`
		UID         string `json:"uid"`
		PhoneNumber string `json:"phone_number"`
	}
)

// Events describes methods available for Events entity
type EventsV2 interface {
	Add(events []Event) ([]EventEmbeddedItem, error)
}

// Verify interface compliance.
var _ EventsV2 = eventsV2{}

type eventsV2 struct {
	api *api
}

func newEventsV2(api *api) EventsV2 {
	return eventsV2{api: api}
}

// Create returns an Contacts entity for successfully added Calls
func (a eventsV2) Add(events []Event) ([]EventEmbeddedItem, error) {

	resp, rErr := a.api.do(eventsV2endpoint, http.MethodPost, nil, nil, eventAdd{Add: events})
	if rErr != nil {
		return nil, fmt.Errorf("get calls: %w", rErr)
	}

	var res struct {
		Embedded struct {
			Items []EventEmbeddedItem `json:"items"`
		} `json:"_embedded"`
	}
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return res.Embedded.Items, nil
}
