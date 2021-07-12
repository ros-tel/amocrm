package amocrm

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"time"
)

type TokenStorage interface {
	SetToken(Token) error
	GetToken() (Token, error)
}

type JSONFileTokenStorage struct {
	File string
}

type JSONToken struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	TokenType    string    `json:"token_type"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (self JSONFileTokenStorage) SetToken(token Token) error {
	jt := JSONToken{
		AccessToken:  token.AccessToken(),
		RefreshToken: token.RefreshToken(),
		TokenType:    token.TokenType(),
		ExpiresAt:    token.ExpiresAt(),
	}

	data, err := json.Marshal(jt)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(self.File, data, os.ModePerm)
}

func (self JSONFileTokenStorage) GetToken() (Token, error) {
	data, err := ioutil.ReadFile(self.File)
	if err != nil {
		return nil, nil
	}

	var jt JSONToken
	if err := json.Unmarshal(data, &jt); err != nil {
		return nil, err
	}

	return NewToken(jt.AccessToken, jt.RefreshToken, jt.TokenType, jt.ExpiresAt), nil
}
