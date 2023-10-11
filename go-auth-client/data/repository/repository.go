package repository

import (
	"learn.oauth.client/data/model"
)

type Repository interface {
	Insert(token model.TokenResponseData) (int, error)
	GetByAccessToken(accessToken string) (*model.TokenResponseData, error)
}
