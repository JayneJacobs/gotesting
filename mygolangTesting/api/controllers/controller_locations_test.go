package controllers

import (
	"encoding/json"
	"fmt"
	"gotrain/GoTestingInteg/mygolangTesting/api/domain/locations"
	"gotrain/GoTestingInteg/mygolangTesting/api/services"
	"gotrain/GoTestingInteg/mygolangTesting/api/utils/errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/mercadolibre/golang-restclient/rest"
	"github.com/stretchr/testify/assert"
)

var getCountryFunc func(countryId string) (*locations.Country, *errors.APIerror)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

type locationServiceMock struct{}

func (*locationServiceMock) GetCountry(countryId string) (*locations.Country, *errors.APIerror) {
	return getCountryFunc(countryId)
}
func TestGetCountryNotFound(t *testing.T) {
	getCountryFunc = func(countryId string) (*locations.Country, *errors.APIerror) {
		return nil, &errors.APIerror{Status: http.StatusNotFound, Message: "Country not found"}
	}

	services.LocationsService = &locationServiceMock{}

	// equivalent to below

	// rest.FlushMockups()
	// rest.AddMockups(&rest.Mock{
	// 	URL:          "https://api.mercadolibre.com/countries/AR",
	// 	HTTPMethod:   http.MethodGet,
	// 	RespHTTPCode: http.StatusNotFound,
	// 	RespBody: `{
	// 		 "message": "Country not found",
	// 		 "error": "not_found",
	// 		 "status": 404,
	// 		 "cause": []}`,
	// })

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{Key: "country_id", Value: "AR"},
	}
	GetCountry(c)
	assert.EqualValues(t, http.StatusNotFound, response.Code)

	var apiErr errors.APIerror
	err := json.Unmarshal(response.Body.Bytes(), &apiErr)
	assert.EqualValues(t, http.StatusNotFound, apiErr.Status)
	assert.Nil(t, err)
	assert.EqualValues(t, "Country not found", apiErr.Message)
}

func TestGetCountryNoError(t *testing.T) {
	getCountryFunc = func(countryId string) (*locations.Country, *errors.APIerror) {
		return &locations.Country{Id: "BR", Name: "Brasil", TimeZone: "3.00+GMT"}, nil
	}

	services.LocationsService = &locationServiceMock{}

	response := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(response)
	c.Request, _ = http.NewRequest(http.MethodGet, "", nil)
	c.Params = gin.Params{
		{Key: "country_id", Value: "BR"},
	}
	GetCountry(c)
	assert.EqualValues(t, http.StatusOK, response.Code)

	var country locations.Country
	err := json.Unmarshal(response.Body.Bytes(), &country)
	fmt.Printf("this is the result %v", country)
	assert.EqualValues(t, "Brasil", country.Name)
	assert.Nil(t, err)
	assert.EqualValues(t, "BR", country.Id)
}