package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search    string // передаваемое значение search
		wantCount int    // ожидаемое количество кафе в ответе
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}
	// Проходим по всем тестовым случаям
	for _, v := range requests {

		// Создаем рекодер для записи ответа
		response := httptest.NewRecorder()

		//Формируем URL с параметрами
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		// Проверяем статус ответа
		require.Equal(t, http.StatusOK, response.Code)
		// Получаем тело ответа и разбиваем на слайс
		body := response.Body.String()
		lower := strings.ToLower(body)
		cafes := strings.Split(lower, ",")
		// Проверяем наличие нужной подстроки
		for _, cafe := range cafes {
			if cafe == "" {
				continue
			}
			if !strings.Contains(cafe, v.search) {
				t.Errorf("название [%s] не содержет подстроку [%s]", cafe, v.search)
			}
		}
		count := len(cafes)
		if cafes[0] == "" {
			count = 0
		}
		assert.Equal(t, v.wantCount, count)
	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int // передаваемое значение count
		want  int // ожидаемое количество кафе в ответе
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, len(cafeList["moscow"])},
	}
	// Проходим по всем тестовым случаям
	for _, v := range requests {

		// Создаем recorder для записи ответа
		response := httptest.NewRecorder()

		// Формируем URL с параметрами
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", v.count)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		// Проверяем статус ответа
		require.Equal(t, http.StatusOK, response.Code)

		// Получаем тело ответа и разбиваем на элементы
		body := response.Body.String()
		cafes := strings.Split(body, ",")
		count := len(cafes)
		if cafes[0] == "" {
			count = 0
		}
		// Проверяем количество кафе
		assert.Equal(t, v.want, count)

	}
}

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
