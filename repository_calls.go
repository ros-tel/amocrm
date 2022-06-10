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
	callsEndpoint endpoint = "calls"
)

type (
	Call struct {
		Direction         string `json:"direction"`                     // Направление звонка. inbound – входящий, outbound – исходящий. Обязательный параметр
		Uniq              string `json:"uniq,omitempty"`                // Уникальный идентификатор звонка. Необязательный параметр
		Duration          int    `json:"duration"`                      // Длительность звонка в секундах. Обязательный параметр
		Source            string `json:"source"`                        // Источник звонка. Обязательный параметр
		Link              string `json:"link,omitempty"`                // Ссылка на запись звонка. Необязательный параметр
		Phone             string `json:"phone"`                         // Номер телефона, по которому будет произведен поиск. Обязательный параметр
		CallResult        string `json:"call_result,omitempty"`         // Результат звонка. Необязательный параметр
		CallStatus        int    `json:"call_status,omitempty"`         // Статус звонка. Доступные варианты: 1 – оставил сообщение, 2 – перезвонить позже, 3 – нет на месте, 4 – разговор состоялся, 5 – неверный номер, 6 – Не дозвонился, 7 – номер занят. Необязательный параметр
		ResponsibleUserID int    `json:"responsible_user_id,omitempty"` // ID пользователя, ответственного за звонок
		CreatedBy         int    `json:"created_by,omitempty"`          // ID пользователя, создавший звонок
		UpdatedBy         int    `json:"updated_by,omitempty"`          // ID пользователя, изменивший звонок
		CreatedAt         int    `json:"created_at,omitempty"`          // Дата создания звонка, передается в Unix Timestamp
		UpdatedAt         int    `json:"updated_at,omitempty"`          // Дата изменения звонка, передается в Unix Timestamp
		RequestID         string `json:"request_id,omitempty"`          // Поле, которое вернется вам в ответе без изменений и не будет сохранено. Необязательный параметр
	}
)

// Calls describes methods available for Calls entity
type Calls interface {
	Create(calls []Call) ([]Contact, error)
}

// Verify interface compliance.
var _ Calls = calls{}

type calls struct {
	api *api
}

func newCalls(api *api) Calls {
	return calls{api: api}
}

// Create returns an Contacts entity for successfully added Calls
func (a calls) Create(calls []Call) ([]Contact, error) {
	resp, rErr := a.api.do(callsEndpoint, http.MethodPost, nil, nil, calls)
	if rErr != nil {
		return nil, fmt.Errorf("get calls: %w", rErr)
	}

	var res struct {
		Embedded struct {
			Contacts []Contact `json:"calls"`
		} `json:"_embedded"`
	}
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return res.Embedded.Contacts, nil
}
