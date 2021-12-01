package repository

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pajarwp/auth-fetch-app/config"
	"github.com/pajarwp/auth-fetch-app/entity/fetch"
)

type FetchRepository interface {
	GetClaims(token string) (fetch.UserClaims, error)
}

type fetchRepository struct {
}

func NewAuthAppFetchRepository() FetchRepository {
	return fetchRepository{}
}

func (a fetchRepository) GetClaims(token string) (fetch.UserClaims, error) {
	var claims fetch.UserClaims
	var client = &http.Client{}
	url := config.GetEnvVariable("AUTH_APP_HOST")
	request, err := http.NewRequest("GET", url+"/user", nil)
	if err != nil {
		return claims, err
	}
	request.Header.Set(echo.HeaderAuthorization, token)
	resp, err := client.Do(request)
	if err != nil {
		return claims, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return claims, err
	}
	respMap := make(map[string]interface{})
	err = json.Unmarshal(body, &respMap)
	if resp.StatusCode != 200 {
		return claims, fmt.Errorf(respMap["msg"].(string))
	}
	data := respMap["data"].(map[string]interface{})
	claims.Name = data["name"].(string)
	claims.Phone = data["phone"].(string)
	claims.Role = data["role"].(string)
	claims.CreatedAt = data["created_at"].(string)
	return claims, nil
}
