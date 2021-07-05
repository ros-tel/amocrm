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

package amocrm_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/sintanial/amocrm"
)

type TokenStored struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func TestLeads_Create(t *testing.T) {
	var err error

	almostValidClient := amocrm.New(
		clientID,
		clientSecret,
		redirectURL,
	)

	if err = almostValidClient.SetDomain("topslotsapp.amocrm.ru"); err != nil {
		t.Fatal(err)
	}

	var token amocrm.Token
	if data, err := ioutil.ReadFile("amocrm_token.json"); err == nil {
		var st TokenStored
		if err := json.Unmarshal(data, &st); err != nil {
			t.Fatal(err)
		}
		token = amocrm.NewToken(st.AccessToken, st.RefreshToken, st.TokenType, st.ExpiresAt)
	}

	if token == nil {
		token, err = almostValidClient.TokenByCode("")
		if err != nil {
			t.Fatal(err)
		}

		data, err := json.Marshal(TokenStored{
			AccessToken:  token.AccessToken(),
			RefreshToken: token.RefreshToken(),
			TokenType:    token.TokenType(),
			ExpiresAt:    token.ExpiresAt(),
		})

		if err != nil {
			t.Fatal(err)
		}

		if err := ioutil.WriteFile("amocrm_token.json", data, os.ModePerm); err != nil {
			t.Fatal(err)
		}
	}

	if err := almostValidClient.SetToken(token); err != nil {
		t.Fatal(err)
	}

	contacts, err := almostValidClient.Contacts().Create([]amocrm.Contact{
		{
			Name:      "+79185436238",
			FirstName: "Roman",
			LastName:  "Martynov",
			CustomFieldsValues: []amocrm.FieldValues{
				{
					"field_code": "PHONE",
					"values": []amocrm.FieldValues{
						{"value": "+79185436238"},
					},
				},
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(contacts)

	leads, err := almostValidClient.Leads().Create([]amocrm.Lead{
		{
			Name: "+79185436238",
			Embedded: &amocrm.LeadEmbedded{
				Tags: []amocrm.FieldValues{
					{
						"name": "partner_reg",
					},
				},
				Contacts: []amocrm.FieldValues{
					{
						"id":      contacts[0].Id,
						"is_main": true,
					},
				},
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(leads)

	almostValidClient.Leads().Update([]amocrm.Lead{
		{
			Id:       leads[0].Id,
			StatusId: 41138881,
		},
	})
}
