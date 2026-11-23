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

const (
	testPhotoPath               = "testdata/test_photo.webp"
	testPrice                   = 100
	testQuantity                = 1
	nonExistentAdvertisementID  = "00000000-0000-4000-8000-000000000000"
	unauthorizedAdvertisementID = "non-existent-advertisement-id"
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

	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)

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
	title := fmt.Sprintf("API favorite test %d", time.Now().UnixNano())
	description := "Test advertisement for favorites"

	_, err := os.Stat(testPhotoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		description,
		testPrice,
		testQuantity,
		testPhotoPath,
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

	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

	addBody := favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	require.NotEmpty(s.T(), gjson.Get(addBody, "favorite_id").String())

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

	favoritesManager.DeleteFavorite(
		s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
	)
	favoriteAdded = false

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
	title := fmt.Sprintf("API deleted favorite test %d", time.Now().UnixNano())
	description := "Test advertisement for favorites deletion"

	_, err := os.Stat(testPhotoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		description,
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID, "ID созданного объявления пустой")

	adDeleted := false
	defer func() {
		if !adDeleted {
			advertisementManager.DeleteAdvertisement(
				s.T(), s.loginData.Token, adID, http.StatusNoContent,
			)
		}
	}()

	favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	favoriteAdded := true
	defer func() {
		if favoriteAdded {
			favoritesManager.DeleteFavorite(
				s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
			)
		}
	}()

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

	advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)
	adDeleted = true
	favoriteAdded = false

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
		unauthorizedAdvertisementID,
		http.StatusUnauthorized,
	)

	require.NotEmpty(s.T(), responseBody, "Ответ с ошибкой пустой")
}

func (s *TestSuite) TestDeleteFavoriteWithoutAuthorization() {
	favoritesManager.DeleteFavorite(
		s.T(),
		"",
		unauthorizedAdvertisementID,
		http.StatusUnauthorized,
	)
}

func (s *TestSuite) TestAddNonExistentFavorite() {
	responseBody := favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		nonExistentAdvertisementID,
		http.StatusNotFound,
	)

	require.Equal(s.T(), "not_found", gjson.Get(responseBody, "error").String())
	require.Equal(s.T(), "Advertisement not found", gjson.Get(responseBody, "message").String())
}

func (s *TestSuite) TestDeleteNonExistentFavorite() {
	responseBody := favoritesManager.DeleteFavorite(
		s.T(),
		s.secondLoginData.Token,
		nonExistentAdvertisementID,
		http.StatusNotFound,
	)

	require.Equal(s.T(), "not_found", gjson.Get(responseBody, "error").String())
	require.Equal(s.T(), "Favorite entry not found", gjson.Get(responseBody, "message").String())
}

func (s *TestSuite) TestAddFavoriteTwice() {
	title := fmt.Sprintf("API duplicate favorite test %d", time.Now().UnixNano())

	_, err := os.Stat(testPhotoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		"Test advertisement for duplicate favorite",
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID)

	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

	favoritesManager.AddFavorite(
		s.T(),
		s.secondLoginData.Token,
		adID,
		http.StatusCreated,
	)

	defer favoritesManager.DeleteFavorite(
		s.T(), s.secondLoginData.Token, adID, http.StatusNoContent,
	)

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
	title := fmt.Sprintf("API own favorite test %d", time.Now().UnixNano())

	_, err := os.Stat(testPhotoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	createBody := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		"Test own advertisement",
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusCreated,
	)

	adID := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), adID)

	defer advertisementManager.DeleteAdvertisement(
		s.T(), s.loginData.Token, adID, http.StatusNoContent,
	)

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
