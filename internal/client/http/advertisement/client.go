package advertisement

import (
	"net/http"
	"strconv"
	"testing"

	"api-tests-template/internal/constants/path"
	apiRunner "api-tests-template/internal/helpers/api-runner"
)

// HttpPostAdvertisement создаёт объявление с фотографией.
func HttpPostAdvertisement(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}

// HttpGetAdvertisement получает объявление по ID.
func HttpGetAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpDeleteAdvertisement удаляет объявление по ID.
func HttpDeleteAdvertisement(t *testing.T, token string, id string) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Delete(path.AdvertisementPath).
		Query("id", id).
		Expect(t).
		End().Response
}

// HttpGetAdvertisements ищет объявления по названию.
func HttpGetAdvertisements(
	t *testing.T,
	token string,
	search string,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get("/advertisements").
		Query("search", search).
		Expect(t).
		End().Response
}

// HttpGetAdvertisementPhotos получает фотографии объявления по ID.
func HttpGetAdvertisementPhotos(
	t *testing.T,
	token string,
	id string,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Get("/advertisements/" + id + "/photos").
		Expect(t).
		End().Response
}

// HttpPostAdvertisementWithoutAuth пытается создать объявление без токена.
func HttpPostAdvertisementWithoutAuth(
	t *testing.T,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
) *http.Response {
	return apiRunner.GetRunner().Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}

// HttpPostAdvertisementWithInvalidAuth пытается создать объявление с неверным токеном.
func HttpPostAdvertisementWithInvalidAuth(
	t *testing.T,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
) *http.Response {
	return apiRunner.GetRunner().Auth("invalid-token").Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}

// HttpPostAdvertisementWithoutPhoto пытается создать объявление без фотографии.
func HttpPostAdvertisementWithoutPhoto(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("description", description).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		Expect(t).
		End().Response
}

// HttpPostAdvertisementWithoutDescription пытается создать объявление без описания.
func HttpPostAdvertisementWithoutDescription(
	t *testing.T,
	token string,
	title string,
	price int,
	quantity int,
	photoPath string,
) *http.Response {
	return apiRunner.GetRunner().Auth(token).Create().
		Post(path.AdvertisementPath).
		MultipartFormData("title", title).
		MultipartFormData("price", strconv.Itoa(price)).
		MultipartFormData("quantity", strconv.Itoa(quantity)).
		MultipartFile("photos", photoPath).
		Expect(t).
		End().Response
}
