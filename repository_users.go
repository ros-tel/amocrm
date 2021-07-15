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

type User struct {
	Id    int    `json:"id,omitempty"`
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Lang  string `json:"lang,omitempty"`
}

// Leads describes methods available for Leads entity.
type Users interface {
	Get(id string) (*User, error)
}

// Verify interface compliance.
var _ Users = users{}

type users struct {
	api *api
}

func newUsers(api *api) Users {
	return users{api: api}
}

func (a users) Get(id string) (*User, error) {
	resp, rErr := a.api.do(endpoint(fmt.Sprintf("users/%s", id)), http.MethodGet, nil, nil, nil)
	if rErr != nil {
		return nil, fmt.Errorf("get users: %w", rErr)
	}

	var res User
	if err := a.api.read(resp, &res); err != nil {
		return nil, err
	}

	return &res, nil
}
