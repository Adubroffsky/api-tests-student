package advertisement

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"

	advertisementManager "api-tests-template/internal/managers/advertisement"
	"api-tests-template/internal/managers/auth"
	"api-tests-template/internal/managers/auth/models"

	base "api-tests-template/tests"
)

const (
	testPhotoPath   = "testdata/test_photo.webp"
	testDescription = "Test description"
	testPrice       = 100
	testQuantity    = 1
)

type TestSuite struct {
	suite.Suite
	loginData models.LoginOkResponse
}

func TestSuiteRun(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (s *TestSuite) SetupSuite() {
	base.SetupSuite()

	base.Precondition("Авторизация пользователя")
	s.loginData = auth.Login(
		s.T(),
		os.Getenv("TEST_LOGIN"),
		os.Getenv("TEST_PASSWORD"),
	)
}

func (s *TestSuite) TearDownSuite() {
	base.TearDownSuite()
}

// TestCreateAdvertisementPositive проверяет создание объявления.
func (s *TestSuite) TestCreateAdvertisementPositive() {
	title := fmt.Sprintf("API test advertisement %d", time.Now().UnixNano())
	description := "Test advertisement description"
	price := testPrice
	quantity := testQuantity
	photoPath := testPhotoPath

	// Проверяем наличие тестовой фотографии.
	_, err := os.Stat(photoPath)
	require.NoError(s.T(), err, "Тестовая фотография не найдена")

	// Создаём объявление.
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

	id := gjson.Get(createBody, "id").String()
	require.NotEmpty(s.T(), id, "ID созданного объявления пустой")

	// Удаляем объявление после завершения проверок,
	// в том числе если одна из проверок завершится ошибкой.
	defer advertisementManager.DeleteAdvertisement(
		s.T(),
		s.loginData.Token,
		id,
		http.StatusNoContent,
	)

	// Проверяем получение объявления по ID.
	getBody := advertisementManager.GetAdvertisement(
		s.T(),
		s.loginData.Token,
		id,
		http.StatusOK,
	)

	require.Equal(s.T(), id, gjson.Get(getBody, "id").String())
	require.Equal(s.T(), title, gjson.Get(getBody, "title").String())
	require.Equal(s.T(), description, gjson.Get(getBody, "description").String())
	require.Equal(s.T(), int64(price), gjson.Get(getBody, "price").Int())
	require.Equal(s.T(), int64(quantity), gjson.Get(getBody, "quantity").Int())
	require.Equal(s.T(), s.loginData.User.Id, gjson.Get(getBody, "user_id").String())

	// Проверяем наличие фотографии.
	photos := gjson.Get(getBody, "photos").Array()
	require.NotEmpty(s.T(), photos, "У объявления отсутствуют фотографии")

	// Проверяем поиск объявления по полному названию.
	searchBody := advertisementManager.GetAdvertisements(
		s.T(),
		s.loginData.Token,
		title,
		http.StatusOK,
	)

	searchItems := gjson.Get(searchBody, "items").Array()
	require.NotEmpty(s.T(), searchItems, "Поиск не вернул объявления")

	found := false

	for _, item := range searchItems {
		if item.Get("title").String() == title &&
			item.Get("id").String() == id {
			found = true
			break
		}
	}

	require.True(s.T(), found, "Созданное объявление не найдено по полному названию")

	// Проверяем получение фотографии отдельной ручкой.
	photosBody := advertisementManager.GetAdvertisementPhotos(
		s.T(),
		s.loginData.Token,
		id,
		http.StatusOK,
	)

	photoItems := gjson.Parse(photosBody).Array()
	require.NotEmpty(s.T(), photoItems, "Ручка фотографий вернула пустой список")

	require.NotEmpty(
		s.T(),
		photoItems[0].Get("id").String(),
		"У фотографии отсутствует ID",
	)

	require.NotEmpty(
		s.T(),
		photoItems[0].Get("url").String(),
		"У фотографии отсутствует URL",
	)
}

// TestCreateAdvertisementWithLongTitle проверяет ограничение длины названия.
func (s *TestSuite) TestCreateAdvertisementWithLongTitle() {
	title := strings.Repeat("A", 51)

	body := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		title,
		testDescription,
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusBadRequest,
	)

	require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String())
	require.Equal(
		s.T(),
		"title must not exceed 50 characters provided 51",
		gjson.Get(body, "message").String(),
	)
}

// TestCreateAdvertisementWithEmptyTitle проверяет запрет пустого названия.
func (s *TestSuite) TestCreateAdvertisementWithEmptyTitle() {
	body := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		"",
		testDescription,
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusBadRequest,
	)

	require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String())
	require.NotEmpty(s.T(), gjson.Get(body, "message").String())
}

// TestCreateAdvertisementWithNegativePrice проверяет запрет отрицательной цены.
func (s *TestSuite) TestCreateAdvertisementWithNegativePrice() {
	body := advertisementManager.CreateAdvertisement(
		s.T(),
		s.loginData.Token,
		"Advertisement with negative price",
		testDescription,
		-1,
		testQuantity,
		testPhotoPath,
		http.StatusBadRequest,
	)

	require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String())
	require.NotEmpty(s.T(), gjson.Get(body, "message").String())
}

// TestCreateAdvertisementWithoutAuth проверяет создание объявления без авторизации.
func (s *TestSuite) TestCreateAdvertisementWithoutAuth() {
	body := advertisementManager.CreateAdvertisementWithoutAuth(
		s.T(),
		"Advertisement without auth",
		testDescription,
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusUnauthorized,
	)

	require.Equal(
		s.T(),
		"unauthorized",
		gjson.Get(body, "error").String(),
	)

	require.Equal(
		s.T(),
		"Authorization header required",
		gjson.Get(body, "message").String(),
	)
}

// TestCreateAdvertisementWithInvalidToken проверяет создание объявления с неверным токеном.
func (s *TestSuite) TestCreateAdvertisementWithInvalidToken() {
	body := advertisementManager.CreateAdvertisementWithInvalidAuth(
		s.T(),
		"Advertisement with invalid token",
		testDescription,
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusUnauthorized,
	)

	require.Equal(s.T(), "unauthorized", gjson.Get(body, "error").String())
	require.NotEmpty(s.T(), gjson.Get(body, "message").String())
}

// TestCreateAdvertisementWithoutPhoto проверяет создание объявления без фотографии.
func (s *TestSuite) TestCreateAdvertisementWithoutPhoto() {
	body := advertisementManager.CreateAdvertisementWithoutPhoto(
		s.T(),
		s.loginData.Token,
		"Advertisement without photo",
		testDescription,
		testPrice,
		testQuantity,
		http.StatusBadRequest,
	)

	require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String())
	require.NotEmpty(s.T(), gjson.Get(body, "message").String())
}

// TestCreateAdvertisementWithoutDescription проверяет создание объявления без описания.
func (s *TestSuite) TestCreateAdvertisementWithoutDescription() {
	body := advertisementManager.CreateAdvertisementWithoutDescription(
		s.T(),
		s.loginData.Token,
		"Advertisement without description",
		testPrice,
		testQuantity,
		testPhotoPath,
		http.StatusBadRequest,
	)

	require.Equal(s.T(), "invalid_request", gjson.Get(body, "error").String())
	require.NotEmpty(s.T(), gjson.Get(body, "message").String())
}

func (s *TestSuite) TestCreateAdvertisementWithStringPrice() {
	advertisementManager.CreateAdvertisementWithRawData(
		s.T(),
		s.loginData.Token,
		"Test advertisement with string price",
		testDescription,
		"invalid-price",
		"1",
		http.StatusBadRequest,
		testPhotoPath,
	)
}

func (s *TestSuite) TestCreateAdvertisementWithStringQuantity() {
	advertisementManager.CreateAdvertisementWithRawData(
		s.T(),
		s.loginData.Token,
		"Test advertisement with string quantity",
		testDescription,
		"100",
		"invalid-quantity",
		http.StatusBadRequest,
		testPhotoPath,
	)
}

func (s *TestSuite) TestCreateAdvertisementWithMoreThanThreePhotos() {
	advertisementManager.CreateAdvertisementWithRawData(
		s.T(),
		s.loginData.Token,
		"Test advertisement with more than three photos",
		testDescription,
		"100",
		"1",
		http.StatusBadRequest,
		testPhotoPath,
		testPhotoPath,
		testPhotoPath,
		testPhotoPath,
	)
}

func (s *TestSuite) TestCreateAdvertisementWithUnsupportedPhotoFormat() {
	advertisementManager.CreateAdvertisementWithRawData(
		s.T(),
		s.loginData.Token,
		"Test advertisement with unsupported photo format",
		testDescription,
		"100",
		"1",
		http.StatusBadRequest,
		"testdata/test_photo.txt",
	)
}
