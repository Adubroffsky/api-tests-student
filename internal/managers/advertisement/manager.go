package advertisement

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	advertisementClient "api-tests-template/internal/client/http/advertisement"
)

// CreateAdvertisement создаёт объявление и возвращает тело ответа.
func CreateAdvertisement(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpPostAdvertisement(
		t, token, title, description, price, quantity, photoPath,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// GetAdvertisement получает объявление и возвращает тело ответа.
func GetAdvertisement(
	t *testing.T,
	token string,
	id string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpGetAdvertisement(t, token, id)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// DeleteAdvertisement удаляет объявление.
func DeleteAdvertisement(
	t *testing.T,
	token string,
	id string,
	expectedStatusCode int,
) {
	response := advertisementClient.HttpDeleteAdvertisement(t, token, id)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)
}

// GetAdvertisements ищет объявления по названию.
func GetAdvertisements(
	t *testing.T,
	token string,
	search string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpGetAdvertisements(t, token, search)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// GetAdvertisementPhotos получает фотографии объявления.
func GetAdvertisementPhotos(
	t *testing.T,
	token string,
	id string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpGetAdvertisementPhotos(t, token, id)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// CreateAdvertisementWithoutAuth проверяет создание объявления без авторизации.
func CreateAdvertisementWithoutAuth(
	t *testing.T,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpPostAdvertisementWithoutAuth(
		t,
		title,
		description,
		price,
		quantity,
		photoPath,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// CreateAdvertisementWithInvalidAuth проверяет создание объявления с неверным токеном.
func CreateAdvertisementWithInvalidAuth(
	t *testing.T,
	title string,
	description string,
	price int,
	quantity int,
	photoPath string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpPostAdvertisementWithInvalidAuth(
		t,
		title,
		description,
		price,
		quantity,
		photoPath,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// CreateAdvertisementWithoutPhoto проверяет создание объявления без фотографии.
func CreateAdvertisementWithoutPhoto(
	t *testing.T,
	token string,
	title string,
	description string,
	price int,
	quantity int,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpPostAdvertisementWithoutPhoto(
		t,
		token,
		title,
		description,
		price,
		quantity,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// CreateAdvertisementWithoutDescription проверяет создание объявления без описания.
func CreateAdvertisementWithoutDescription(
	t *testing.T,
	token string,
	title string,
	price int,
	quantity int,
	photoPath string,
	expectedStatusCode int,
) string {
	response := advertisementClient.HttpPostAdvertisementWithoutDescription(
		t,
		token,
		title,
		price,
		quantity,
		photoPath,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

func CreateAdvertisementWithRawData(
	t *testing.T,
	token string,
	title string,
	description string,
	price string,
	quantity string,
	expectedStatusCode int,
	photoPaths ...string,
) string {
	response := advertisementClient.HttpPostAdvertisementWithRawData(
		t,
		token,
		title,
		description,
		price,
		quantity,
		photoPaths...,
	)
	defer response.Body.Close()

	require.Equal(t, expectedStatusCode, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}
