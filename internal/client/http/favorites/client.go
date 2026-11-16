package favorites

import (
	"net/http"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

type AddFavoriteRequest struct {
	AdvertisementID string `json:"advertisement_id"`
}

// HttpPostFavorite добавляет объявление в избранное.
func HttpPostFavorite(t *testing.T, token, advertisementID string) *http.Response {
	request := AddFavoriteRequest{
		AdvertisementID: advertisementID,
	}

	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.FavoritesPath).
		ContentType("application/json").
		JSON(request).
		Expect(t).
		End().Response
}

// HttpGetFavorites получает список избранного.
func HttpGetFavorites(t *testing.T, token string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.FavoritesPath).
		Expect(t).
		End().Response
}

// HttpDeleteFavorite удаляет объявление из избранного.
func HttpDeleteFavorite(t *testing.T, token, advertisementID string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.FavoriteByIDPath + advertisementID).
		Expect(t).
		End().Response
}
