package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/mitchellh/mapstructure"
	"github.com/pajarwp/auth-fetch-app/entity/fetch"
)

type AuthAppFetchRepository interface {
	GetClaims(token string) (fetch.UserClaims, error)
}

type authAppFetchRepository struct {
}

func NewAuthAppFetchRepository() AuthAppFetchRepository {
	return authAppFetchRepository{}
}

func (a authAppFetchRepository) GetClaims(token string) (fetch.UserClaims, error) {
	var claims fetch.UserClaims
	resp, err := http.Get(os.Getenv("AUTH_APP_HOST") + "/user")
	resp.Header.Add(echo.HeaderAuthorization, "Bearer "+token)
	if err != nil {
		return claims, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return claims, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, respMap)
	if resp.StatusCode != 200 {
		return claims, fmt.Errorf(respMap["message"].(string))
	}
	mapstructure.Decode(respMap["data"], claims)
	return claims, nil
}
