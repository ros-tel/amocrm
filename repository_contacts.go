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
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type ContactsEmbedded struct {
	Tags []FieldValues `json:"tags,omitempty"`
}

type Contact struct {
	Id                 int               `json:"id,omitempty"`
	Name               string            `json:"name,omitempty"`                 //Название контакта
	FirstName          string            `json:"first_name,omitempty"`           //Имя контакта
	LastName           string            `json:"last_name,omitempty"`            //Фамилия контакта
	ResponsibleUserId  int               `json:"responsible_user_id,omitempty"`  //ID пользователя, ответственного за контакт
	CreatedBy          int               `json:"created_by,omitempty"`           //ID пользователя, создавший контакт
	UpdatedBy          int               `json:"updated_by,omitempty"`           //ID пользователя, изменивший контакт
	CreatedAt          int               `json:"created_at,omitempty"`           //Дата создания контакта, передается в Unix Timestamp
	UpdatedAt          int               `json:"updated_at,omitempty"`           //Дата изменения контакта, передается в Unix Timestamp
	CustomFieldsValues []FieldValues     `json:"custom_fields_values,omitempty"` //Массив, содержащий информацию по дополнительным полям, заданным для данной сделки. Поле не является обязательным. Примеры заполнения полей
	Embedded           *ContactsEmbedded `json:"_embedded,omitempty"`            //Данные вложенных сущностей, при создании и редактировании можно передать только теги. Поле не является обязательным
}

// Contacts describes methods available for Contacts entity.
type Contacts interface {
	Contacts(values url.Values) ([]Contact, error)
	Create(contacts []Contact) ([]Contact, error)
}

// Verify interface compliance.
var _ Contacts = contacts{}

type contacts struct {
	api *api
}

func newContacts(api *api) Contacts {
	return contacts{api: api}
}

func (a contacts) Contacts(values url.Values) ([]Contact, error) {
	r, err := a.api.do(contactsEndpoint, http.MethodGet, values, nil, nil)
	if err != nil {
		return nil, err
	}

	var res []Contact
	if err := json.NewDecoder(r.Body).Decode(&res); err != nil {
		return nil, err
	}

	return res, err
}

// Current returns an Contacts entity for current authorized user.
func (a contacts) Create(contacts []Contact) ([]Contact, error) {
	resp, rErr := a.api.do(contactsEndpoint, http.MethodPost, nil, nil, contacts)
	if rErr != nil {
		return nil, fmt.Errorf("get contacts: %w", rErr)
	}

	var res struct {
		Embedded struct {
			Contacts []Contact `json:"contacts"`
		} `json:"_embedded"`
	}
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return res.Embedded.Contacts, nil
}
