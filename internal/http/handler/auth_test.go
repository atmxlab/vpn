package handler_test

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/atmxlab/vpn/internal/domain/dto/usecase"
	"github.com/atmxlab/vpn/internal/http/handler"
	"github.com/atmxlab/vpn/internal/http/handler/mocks"
	mockhttp "github.com/atmxlab/vpn/internal/http/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

type HTTPTester struct {
	t                  *testing.T
	ctrl               *gomock.Controller
	responseWriterMock *mockhttp.MockResponseWriter
}

func NewHTTPTester(
	t *testing.T,
) *HTTPTester {
	ctrl := gomock.NewController(t)

	responseWriterMock := mockhttp.NewMockResponseWriter(ctrl)
	return &HTTPTester{
		t:                  t,
		ctrl:               ctrl,
		responseWriterMock: responseWriterMock,
	}
}

func (h *HTTPTester) ResponseWriterMock() *mockhttp.MockResponseWriter {
	return h.responseWriterMock
}

func (h *HTTPTester) Ctrl() *gomock.Controller {
	return h.ctrl
}

func (h *HTTPTester) Request(method string, target string, body any) *http.Request {
	bodyJson, err := json.Marshal(body)
	require.NoError(h.t, err)

	buf := bytes.NewBuffer(bodyJson)

	return httptest.NewRequest(method, target, buf)
}

func TestHandle(t *testing.T) {
	tester := NewHTTPTester(t)
	defer tester.Ctrl().Finish()

	usecaseMock := mocks.NewMockUsecase(tester.Ctrl())

	type body struct {
		Key string `json:"key"`
		IP  string `json:"ip"`
	}

	bd := body{
		Key: "123",
		IP:  "1.1.1.1",
	}

	usecaseMock.EXPECT().Auth(gomock.Any(), usecase.AuthOptions{
		IP:  net.IPv4(1, 1, 1, 1),
		Key: []byte("123"),
	}).Return(&usecase.AuthResult{}, nil)

	tester.ResponseWriterMock().EXPECT().WriteHeader(200)
	tester.ResponseWriterMock().EXPECT().Write(gomock.Any())

	authHandler := handler.NewAuth(usecaseMock)
	authHandler.Handle(
		tester.ResponseWriterMock(),
		tester.Request("GET", "http://localhost:8080", bd),
	)
}
