package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/betom84/kl-logger/api"
	"github.com/betom84/kl-logger/repository/testdata"

	"github.com/stretchr/testify/assert"
)

type assertion func(*testing.T, map[string]interface{})

func TestServer(t *testing.T) {
	tt := []struct {
		url         string
		code        int
		contentType string
		assertions  []assertion
	}{
		{
			url:         "/weather",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			url:         "/weather/0",
			code:        http.StatusOK,
			contentType: "application/json",
			assertions: []assertion{
				assertResponseContains("name", "Description 0"),
				assertResponseContains("temperature", 20.0),
			},
		},
		{
			url:         "/weather/8",
			code:        http.StatusOK,
			contentType: "application/json",
			assertions: []assertion{
				assertResponseContains("name", "Description 8"),
				assertResponseContains("temperature", 28.0),
			},
		},
		{
			url:         "/weather/9",
			code:        http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			url:         "/weather/abc",
			code:        http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			url:         "/config",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			url:         "/config/0",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			url:         "/config/8",
			code:        http.StatusOK,
			contentType: "application/json",
		},
		{
			url:         "/config/9",
			code:        http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
		{
			url:         "/config/abc",
			code:        http.StatusNotFound,
			contentType: "text/plain; charset=utf-8",
		},
	}

	w := testdata.MockWeatherSample(func(m *testdata.WeatherSampleMock) {
		for i := 0; i <= 8; i++ {
			m.On("Temperature", i).Return(float32(20 + i))
		}
	})

	c := testdata.MockConfiguration(func(m *testdata.ConfigurationMock) {
		for i := 0; i <= 8; i++ {
			m.On("Description", i).Return(fmt.Sprintf("Description %d", i))
		}
	})

	s := httptest.NewServer(api.NewServer(testdata.MockRepository(w, c), nil))

	for _, tc := range tt {
		t.Run(fmt.Sprintf("endpoint%s", tc.url), func(t *testing.T) {
			r, err := s.Client().Get(fmt.Sprintf("%s%s", s.URL, tc.url))
			assert.NoError(t, err)

			assert.Equal(t, tc.code, r.StatusCode)
			assert.Equal(t, tc.contentType, r.Header.Get("Content-Type"))

			if tc.assertions == nil {
				return
			}

			var content interface{}
			err = json.NewDecoder(r.Body).Decode(&content)
			assert.NoError(t, err)

			values, ok := content.(map[string]interface{})
			if !ok {
				t.Errorf("invalid json response")
			}

			for _, a := range tc.assertions {
				a(t, values)
			}
		})
	}
}

func assertResponseContains(expectedKey string, expectedValue interface{}) assertion {
	return func(t *testing.T, values map[string]interface{}) {
		v, ok := values[expectedKey]
		if !ok {
			t.Log(values)
			t.Errorf("response does not contain key '%s'", expectedKey)
			t.FailNow()
		}

		assert.Equal(t, expectedValue, v)
	}
}
