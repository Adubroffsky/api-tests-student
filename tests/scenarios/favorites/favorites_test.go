package favorites

import (
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"
	favoritesManager "api-tests-template/internal/managers/favorites"

	base "api-tests-template/tests"
)

type TestSuite struct {
	suite.Suite
	loginData       models.LoginOkResponse
	secondLoginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()
	base.Precondition("Авторизация двух пользователей")

	// Авторизация первого пользователя.
	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

	// Авторизация второго пользователя.
	s.secondLoginData = auth.Login(
		s.T(),
		os.Getenv("SECOND_TEST_LOGIN"),
		os.Getenv("SECOND_TEST_PASSWORD"),
	)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

func (s *TestSuite) TestAddAndDeleteFavoritePositive() {
	// 1. Создаём объявление от имени первого пользователя.
	title := fmt.Sprintf("API favorite test %d", time.Now().UnixNano())
	description := "Test advertisement for favorites"
	price := 100
	quantity := 1
	photoPath := "../../../test_photo.webp"

	_, err := os.Stat(photoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		description,
		price,
		quantity,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID, "ID созданного объявления пустой")

	require.NotEqual(
		s.T(),
		s.loginData.User.Id,
		s.secondLoginData.User.Id,
		"Для теста нужны два разных пользователя",
	)

	// Удаляем объявление в конце теста от имени владельца.
	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

	// 2. Второй пользователь добавляет объявление в избранное.
	addBody := favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	require.NotEmpty(s.T(), gjson.Get(addBody, "favorite_id").String())

	// Если тест завершится с ошибкой, удаляем избранное.
	favoriteAdded := true
	defer func() {
		if favoriteAdded {
			favoritesManager.DeleteFavorite(
				s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
			)
		}
	}()

	require.Equal(
		s.T(),
		adID,
		gjson.Get(addBody, "advertisement.id").String(),
	)

	// 3. Второй пользователь проверяет список избранного.
	favoritesBody := favoritesManager.GetFavorites(
		s.T(), s.secondLoginData.Token, http.StatusOK,
	)

	items := gjson.Get(favoritesBody, "items").Array()
	found := false

	for _, item := range items {
		if item.Get("advertisement.id").String() == adID {
			found = true
			break
		}
	}

	require.True(s.T(), found, "Объявление не найдено в избранном")

	// 4. Второй пользователь удаляет объявление из избранного.
	favoritesManager.DeleteFavorite(
		s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
	)
	favoriteAdded = false

	// 5. Проверяем, что объявления больше нет в избранном.
	favoritesBody = favoritesManager.GetFavorites(
		s.T(), s.secondLoginData.Token, http.StatusOK,
	)

	items = gjson.Get(favoritesBody, "items").Array()
	found = false

	for _, item := range items {
		if item.Get("advertisement.id").String() == adID {
			found = true
			break
		}
	}

	require.False(s.T(), found, "Объявление всё ещё находится в избранном")
}

func (s *TestSuite) TestDeletedAdvertisementDisappearsFromFavorites() {
	// 1. Первый пользователь создаёт объявление.
	title := fmt.Sprintf("API deleted favorite test %d", time.Now().UnixNano())
	description := "Test advertisement for favorites deletion"
	price := 100
	quantity := 1
	photoPath := "../../../test_photo.webp"

	_, err := os.Stat(photoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		description,
		price,
		quantity,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID, "ID созданного объявления пустой")

	// Если тест завершится раньше удаления, удаляем объявление.
	adDeleted := false
	defer func() {
		if !adDeleted {
			advertisementManager.DeleteAdvertisement(
				s.T(), s.loginData.Token, adID, http.StatusNoContent,
			)
		}
	}()

	// 2. Второй пользователь добавляет объявление в избранное.
	favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	// Если тест завершится раньше удаления объявления,
	// удаляем его из избранного.
	favoriteAdded := true
	defer func() {
		if favoriteAdded {
			favoritesManager.DeleteFavorite(
				s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
			)
		}
	}()

	// 3. Проверяем, что объявление появилось в избранном.
	favoritesBody := favoritesManager.GetFavorites(
		s.T(), s.secondLoginData.Token, http.StatusOK,
	)

	found := false
	for _, item := range gjson.Get(favoritesBody, "items").Array() {
		if item.Get("advertisement.id").String() == adID {
			found = true
			break
		}
	}
	require.True(s.T(), found, "Объявление не найдено в избранном до удаления")

	// 4. Первый пользователь удаляет собственное объявление.
	advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)
	adDeleted = true
	favoriteAdded = false

	// 5. Проверяем, что удалённого объявления больше нет в избранном.
	favoritesBody = favoritesManager.GetFavorites(
		s.T(), s.secondLoginData.Token, http.StatusOK,
	)

	found = false
	for _, item := range gjson.Get(favoritesBody, "items").Array() {
		if item.Get("advertisement.id").String() == adID {
			found = true
			break
		}
	}

	require.False(
		s.T(),
		found,
		"Удалённое объявление всё ещё находится в избранном",
	)
}

func (s *TestSuite) TestGetFavoritesWithoutAuthorization() {
	favoritesBody := favoritesManager.GetFavorites(
		s.T(),
		"",
		http.StatusUnauthorized,
	)

	require.NotEmpty(s.T(), favoritesBody, "Ответ с ошибкой пустой")
}

func (s *TestSuite) TestAddFavoriteWithoutAuthorization() {
	responseBody := favoritesManager.AddFavorite(
		s.T(),
		"",
		"non-existent-advertisement-id",
		http.StatusUnauthorized,
	)

	require.NotEmpty(s.T(), responseBody, "Ответ с ошибкой пустой")
}

func (s *TestSuite) TestDeleteFavoriteWithoutAuthorization() {
	favoritesManager.DeleteFavorite(
		s.T(),
		"",
		"non-existent-advertisement-id",
		http.StatusUnauthorized,
	)
}

func (s *TestSuite) TestAddNonExistentFavorite() {
	nonExistentAdID := "00000000-0000-4000-8000-000000000000"

	responseBody := favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		nonExistentAdID,
		http.StatusNotFound,
	)

	require.Equal(s.T(), "not_found", gjson.Get(responseBody, "error").String())
	require.Equal(s.T(), "Advertisement not found", gjson.Get(responseBody, "message").String())
}

func (s *TestSuite) TestDeleteNonExistentFavorite() {
	nonExistentAdID := "00000000-0000-4000-8000-000000000000"

	responseBody := favoritesManager.DeleteFavorite(
		s.T(),
		s.secondLoginData.Token,
		nonExistentAdID,
		http.StatusNotFound,
	)

	require.Equal(s.T(), "not_found", gjson.Get(responseBody, "error").String())
	require.Equal(s.T(), "Favorite entry not found", gjson.Get(responseBody, "message").String())
}

func (s *TestSuite) TestAddFavoriteTwice() {
	// Первый пользователь создаёт объявление.
	title := fmt.Sprintf("API duplicate favorite test %d", time.Now().UnixNano())
	photoPath := "../../../test_photo.webp"

	_, err := os.Stat(photoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		"Test advertisement for duplicate favorite",
		100,
		1,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID)

	// Удаляем объявление после теста.
	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

	// Второй пользователь добавляет объявление первый раз.
	favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	// Удаляем избранное после теста.
	defer favoritesManager.DeleteFavorite(
		s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
	)

	// Повторное добавление должно вернуть 409 Conflict.
	responseBody := favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusConflict,
	)

	require.Equal(
		s.T(),
		"already_exists",
		gjson.Get(responseBody, "error").String(),
	)
	require.Equal(
		s.T(),
		"Advertisement is already in favorites",
		gjson.Get(responseBody, "message").String(),
	)
}

func (s *TestSuite) TestAddOwnAdvertisementToFavorites() {
	// Первый пользователь создаёт собственное объявление.
	title := fmt.Sprintf("API own favorite test %d", time.Now().UnixNano())
	photoPath := "../../../test_photo.webp"

	_, err := os.Stat(photoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		"Test own advertisement",
		100,
		1,
		photoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID)

	// Удаляем объявление после теста.
	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

	// Владелец пытается добавить своё объявление в избранное.
	responseBody := favoritesManager.AddFavorite(
		s.T(),
		s.loginData.Token,
		adID,
		http.StatusForbidden,
	)

	require.Equal(
		s.T(),
		"forbidden",
		gjson.Get(responseBody, "error").String(),
	)
	require.Equal(
		s.T(),
		"Cannot add own advertisement to favorites",
		gjson.Get(responseBody, "message").String(),
	)
}
