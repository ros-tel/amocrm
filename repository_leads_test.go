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
		"bec5cd07-909b-45e6-9a77-6215c8ca00d8",
		"IL6E4xldF1i80k6DPhIpSDnMtNvmJs18bVvpUlVq1W7gOKkdwVm0y4Qtbi7MIatr",
		"https://social-casino.com/",
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
		token, err = almostValidClient.TokenByCode("def50200010bd7ea2b1ca6452070d2ce476cd021e1f265f5c5122029e18a6b7927c64235d32d2a4753d5f8e50f5d2584eab45da62feafb8c64bbffcc1b4ff3efafd144350d758ff5cda9022427e7274abd560d9b9765e8d49ce8a0a147243b2fe65f75b14b1f59774e7880cf81ea843a40680c2829acef9e0ea9836545d38b0116fbdbbc18e6bc24457fce03d1b1ff452940530845cbde2a34cd1c6c28ac7cc2a0d0fb79367829431635cb101c732423f4932c38f706fd739dd29a1c34f42f3a99450dc7ae9ff98087430c7ed62b03c21e9047ca4d65f8c12bff989811a251573f0f602e99a6d04ec527a64b89c2326191da699c586bbdf38a7e6d5b2e8bac9ee505a2a14b9757054ec8222fbf56b443bfb31d4e5bfbb2efba7573ec44f60b1e1213375f8272d1b739fdb9b7a5cb635157526ad58b59ba543562e413cf8e779217f2fce21806ddf71b67a6582830b093bbee0eb0005103acdf56bbe79dfc0ced88fdbc94c7e21eddc69be7409e1ff67427deaf1355f8dde8e518f7c9b33c4636ab61bde54cdce92b0262640935076f1251e1f7bdffa7c0b80cd10eff09071b01cd79a02a08263ba1cea38a64549b2254ab2398ed3c35cc4543ac5e743dc8097add452a730f")
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
