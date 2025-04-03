package middlewares

import (
	"1/internal/usecase/mocks"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Валдиный токен - 200", func(t *testing.T) {
		authServiceMock := new(mocks.AuthServiceMock)

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": 1,
			"exp": time.Now().Add(time.Second * 5).Unix(),
		})
		token.Valid = true

		fmt.Println(token)

		authServiceMock.On("ValidateToken", "fake_token").Return(token, nil)

		router := gin.New()
		router.Use(JWTAuthMiddleware(authServiceMock))
		router.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer fake_token")
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		authServiceMock.AssertCalled(t, "ValidateToken", "fake_token")
	})

	t.Run("Отсутствует Authorization header - ошибка 401", func(t *testing.T) {
		authServiceMock := new(mocks.AuthServiceMock)

		router := gin.New()
		router.Use(JWTAuthMiddleware(authServiceMock))
		router.GET("/test", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"message": "success"})
		})

		req, _ := http.NewRequest("GET", "/test", nil)
		resp := httptest.NewRecorder()
		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusUnauthorized, resp.Code)
		assert.Contains(t, resp.Body.String(), "Authorization header is required")
	})
}
