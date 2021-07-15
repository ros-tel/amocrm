// The MIT License (MIT)
//
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

type LeadEmbedded struct {
	Tags      []FieldValues `json:"tags,omitempty"`
	Contacts  []FieldValues `json:"contacts,omitempty"`
	Companies []FieldValues `json:"companies,omitempty"`
}

type Lead struct {
	Id                 int           `json:"id,omitempty"`
	Name               string        `json:"name,omitempty"`                 //Название сделки. Поле не является обязательным
	Price              int           `json:"price,omitempty"`                //Бюджет сделки. Поле не является обязательным
	StatusId           int           `json:"status_id,omitempty"`            //ID статуса, в который добавляется сделка. Поле не является обязательным, по-умолчанию – первый этап главной воронки
	PipelineId         int           `json:"pipeline_id,omitempty"`          //ID воронки, в которую добавляется сделка. Поле не является обязательным
	CreatedBy          int           `json:"created_by,omitempty"`           //ID пользователя, создающий сделку. При передаче значения 0, сделка будет считаться созданной роботом. Поле не является обязательным
	UpdatedBy          int           `json:"updated_by,omitempty"`           //ID пользователя, изменяющий сделку. При передаче значения 0, сделка будет считаться измененной роботом. Поле не является обязательным
	ClosedAt           int           `json:"closed_at,omitempty"`            //Дата закрытия сделки, передается в Unix Timestamp. Поле не является обязательным
	CreatedAt          int           `json:"created_at,omitempty"`           //Дата создания сделки, передается в Unix Timestamp. Поле не является обязательным
	UpdatedAt          int           `json:"updated_at,omitempty"`           //Дата изменения сделки, передается в Unix Timestamp. Поле не является обязательным
	LossReasonId       int           `json:"loss_reason_id,omitempty"`       //ID причины отказа. Поле не является обязательным
	ResponsibleUserId  int           `json:"responsible_user_id,omitempty"`  //ID пользователя, ответственного за сделку. Поле не является обязательным
	CustomFieldsValues []FieldValues `json:"custom_fields_values,omitempty"` //Массив, содержащий информацию по дополнительным полям, заданным для данной сделки. Поле не является обязательным. Примеры заполнения полей
	Embedded           *LeadEmbedded `json:"_embedded,omitempty"`            //Данные вложенных сущностей, при создании и редактировании можно передать только теги. Поле не является обязательным
}

// Leads describes methods available for Leads entity.
type Leads interface {
	Get(id string) (*Lead, error)
	Create(leads []Lead) ([]Lead, error)
	Update(leads []Lead) ([]Lead, error)
}

// Verify interface compliance.
var _ Leads = leads{}

type leads struct {
	api *api
}

func newLeads(api *api) Leads {
	return leads{api: api}
}

func (a leads) Get(id string) (*Lead, error) {
	resp, rErr := a.api.do(endpoint(fmt.Sprintf("leads/%s", id)), http.MethodGet, nil, nil, nil)
	if rErr != nil {
		return nil, fmt.Errorf("get leads: %w", rErr)
	}

	var res Lead
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

// Current returns an Leads entity for current authorized user.
func (a leads) Create(leads []Lead) ([]Lead, error) {
	resp, rErr := a.api.do(leadsEndpoint, http.MethodPost, nil, nil, leads)
	if rErr != nil {
		return nil, fmt.Errorf("create leads: %w", rErr)
	}

	var res struct {
		Embedded struct {
			Leads []Lead `json:"leads"`
		} `json:"_embedded"`
	}
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return res.Embedded.Leads, nil
}

// Current returns an Leads entity for current authorized user.
func (a leads) Update(leads []Lead) ([]Lead, error) {
	resp, rErr := a.api.do(leadsEndpoint, http.MethodPatch, nil, nil, leads)
	if rErr != nil {
		return nil, fmt.Errorf("update leads: %w", rErr)
	}

	var res struct {
		Embedded struct {
			Leads []Lead `json:"leads"`
		} `json:"_embedded"`
	}
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return res.Embedded.Leads, nil
}
