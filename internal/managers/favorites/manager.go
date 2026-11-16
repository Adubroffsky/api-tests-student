package favorites

import (
	"io"
	"testing"

	"github.com/stretchr/testify/require"

	favoritesClient "api-tests-template/internal/client/http/favorites"
)

// AddFavorite добавляет объявление в избранное.
func AddFavorite(t *testing.T, token, advertisementID string, expectedStatus int) string {
	response := favoritesClient.HttpPostFavorite(t, token, advertisementID)
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// GetFavorites получает список избранного.
func GetFavorites(t *testing.T, token string, expectedStatus int) string {
	response := favoritesClient.HttpGetFavorites(t, token)
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}

// DeleteFavorite удаляет объявление из избранного.
func DeleteFavorite(t *testing.T, token, advertisementID string, expectedStatus int) string {
	response := favoritesClient.HttpDeleteFavorite(t, token, advertisementID)
	defer response.Body.Close()

	require.Equal(t, expectedStatus, response.StatusCode)

	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)

	return string(body)
}
