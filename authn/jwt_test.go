package authn

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type stubKeyProvider struct {
	key string
}

func (s *stubKeyProvider) GetKey(token *jwt.Token) (interface{}, error) {
	return []byte(s.key), nil
}

type AuthnTestSuite struct {
	suite.Suite
	handler    http.Handler
	nextCalled int
}

func (suite *AuthnTestSuite) SetupTest() {
	suite.nextCalled = 0
	suite.handler = JwtAuthn(
		&stubKeyProvider{"secret"},
		WithRespondEmptyHeader(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "NoAuthHeader", http.StatusUnauthorized)
		}),
		WithRespondInvalidJwt(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "InvalidToken", http.StatusUnauthorized)
		}),
	)(func(w http.ResponseWriter, r *http.Request) {
		sub, _ := r.Context().Value("token").(*jwt.Token).Claims.GetSubject()
		suite.nextCalled++
		_, _ = w.Write([]byte(sub))
	})
}

func (suite *AuthnTestSuite) TestAuthenticated() {

	validTk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-id",
	})
	validTkStr, err := validTk.SignedString([]byte("secret"))
	suite.NoError(err)

	// Prepare request stack
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer "+validTkStr)

	rr := httptest.NewRecorder()

	suite.handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusOK, rr.Code)
	suite.Equal("test-id", rr.Body.String())
	suite.Equal(1, suite.nextCalled)
}

func (suite *AuthnTestSuite) TestAuthWrongSig() {
	wrongTk := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-id",
	})
	wrongTkStr, err := wrongTk.SignedString([]byte("wrong-secret"))
	suite.NoError(err)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer "+wrongTkStr)

	rr := httptest.NewRecorder()

	suite.handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.True(strings.HasPrefix(rr.Body.String(), "InvalidToken"))
	suite.Equal(0, suite.nextCalled)
}

func (suite *AuthnTestSuite) TestNoToken() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	suite.handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.True(strings.HasPrefix(rr.Body.String(), "NoAuthHeader"))
	suite.Equal(0, suite.nextCalled)
}

func (suite *AuthnTestSuite) TestInvalidToken() {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Add("Authorization", "Bearer invalid-token")
	rr := httptest.NewRecorder()

	suite.handler.ServeHTTP(rr, req)
	suite.Equal(http.StatusUnauthorized, rr.Code)
	suite.True(strings.HasPrefix(rr.Body.String(), "InvalidToken"))
	suite.Equal(0, suite.nextCalled)
}

func TestAuthn(t *testing.T) {
	suite.Run(t, new(AuthnTestSuite))
}
