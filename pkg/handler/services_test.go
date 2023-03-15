package handler

import (
	"bytes"
	"errors"
	"fmt"
	"net/http/httptest"
	"testing"

	avitotask "github.com/ant0nix/avitoTask"
	"github.com/ant0nix/avitoTask/pkg/service"
	mock_service "github.com/ant0nix/avitoTask/pkg/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
)

func TestCreateUser(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStart, user avitotask.User)

	testTable := []struct {
		name               string
		inputBody          string
		inputUser          avitotask.User
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"uname":"testName","balance":100,"reserved":0}`,
			inputUser: avitotask.User{
				UName:    "testName",
				Balance:  100,
				Reserved: 0,
			},
			mockBehavior: func(s *mock_service.MockStart, user avitotask.User) {
				s.EXPECT().CreateUser(user).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"User created successfully"}`,
		},
		{
			name:      "TestBad Fake DB err",
			inputBody: `{"uname":"testName","balance":100,"reserved":0}`,
			inputUser: avitotask.User{
				UName:    "testName",
				Balance:  100,
				Reserved: 0,
			},
			mockBehavior: func(s *mock_service.MockStart, user avitotask.User) {
				s.EXPECT().CreateUser(user).Return(errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"uname":123,"balance":100,"reserved":0}`,
			mockBehavior:       func(s *mock_service.MockStart, user avitotask.User) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
		{
			name:      "TestBad Body(negativ values)",
			inputBody: `{"uname":"testName","balance":-100,"reserved":0}`,
			inputUser: avitotask.User{
				UName:    "testName",
				Balance:  -100,
				Reserved: 0,
			},
			mockBehavior:       func(s *mock_service.MockStart, user avitotask.User) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"negative values"}`,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			start := mock_service.NewMockStart(c)
			testCase.mockBehavior(start, testCase.inputUser)

			services := &service.Services{Start: start}

			h := NewHandler(services)

			r := gin.New()
			r.POST("/user", h.CreateUser)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/user", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)
			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestChangeBalance(t *testing.T) {
	type mockBehavior func(s *mock_service.MockInternalServices, balance avitotask.Balance)

	testTableBalance := []struct {
		name               string
		inputBody          string
		inputBalance       avitotask.Balance
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"uid":1,"change":123}`,
			inputBalance: avitotask.Balance{
				UserID:        1,
				ChangeBalance: 123,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, balance avitotask.Balance) {
				s.EXPECT().ChangeBalance(balance).Return("Balance changed successfull", nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"Balance changed successfull"}`,
		},
		{
			name:      "TestBad Fake db err",
			inputBody: `{"uid":1,"change":123}`,
			inputBalance: avitotask.Balance{
				UserID:        1,
				ChangeBalance: 123,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, balance avitotask.Balance) {
				s.EXPECT().ChangeBalance(balance).Return("", errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"uid":"-","change":10}`,
			mockBehavior:       func(s *mock_service.MockInternalServices, balance avitotask.Balance) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
	}
	for _, testCase := range testTableBalance {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			internal := mock_service.NewMockInternalServices(c)
			testCase.mockBehavior(internal, testCase.inputBalance)

			services := &service.Services{InternalServices: internal}

			h := NewHandler(services)

			r := gin.New()
			r.POST("/change-balance", h.ChangeBalance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/change-balance", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestShowBalance(t *testing.T) {
	type mockBehavior func(s *mock_service.MockInternalServices, balance avitotask.Balance)

	testTableBalance := []struct {
		name               string
		inputBody          string
		inputBalance       avitotask.Balance
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"uid":1}`,
			inputBalance: avitotask.Balance{
				UserID: 1,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, balance avitotask.Balance) {
				s.EXPECT().ShowBalance(balance).Return(0, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":0}`,
		},
		{
			name:      "TestBad Fake db err",
			inputBody: `{"uid":1}`,
			inputBalance: avitotask.Balance{
				UserID: 1,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, balance avitotask.Balance) {
				s.EXPECT().ShowBalance(balance).Return(0, errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"uid":"-"}`,
			mockBehavior:       func(s *mock_service.MockInternalServices, balance avitotask.Balance) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
	}
	for _, testCase := range testTableBalance {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			internal := mock_service.NewMockInternalServices(c)
			testCase.mockBehavior(internal, testCase.inputBalance)

			services := &service.Services{InternalServices: internal}

			h := NewHandler(services)

			r := gin.New()
			r.GET("/show-balance", h.ShowBalance)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/show-balance", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestP2P(t *testing.T) {
	type mockBehavior func(s *mock_service.MockInternalServices, balance avitotask.P2p)

	testTable := []struct {
		name               string
		inputBody          string
		inputP2P           avitotask.P2p
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"sid":1,"did":2,"amount":100}`,
			inputP2P: avitotask.P2p{
				SId:    1,
				DId:    2,
				Amount: 100,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, p2p avitotask.P2p) {
				s.EXPECT().P2p(p2p).Return("transaction have done", nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"transaction have done"}`,
		},
		{
			name:      "TestBad Fake db err",
			inputBody: `{"sid":1,"did":2,"amount":100}`,
			inputP2P: avitotask.P2p{
				SId:    1,
				DId:    2,
				Amount: 100,
			},
			mockBehavior: func(s *mock_service.MockInternalServices, p2p avitotask.P2p) {
				s.EXPECT().P2p(p2p).Return("", errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"sid":"-","did":10,"amount":100}`,
			mockBehavior:       func(s *mock_service.MockInternalServices, p2p avitotask.P2p) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
		{
			name:               "TestBad (negative values)",
			inputBody:          `{"sid":2,"did":10,"amount":-100}`,
			mockBehavior:       func(s *mock_service.MockInternalServices, p2p avitotask.P2p) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"negative values"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			internal := mock_service.NewMockInternalServices(c)
			testCase.mockBehavior(internal, testCase.inputP2P)

			services := &service.Services{InternalServices: internal}

			h := NewHandler(services)

			r := gin.New()
			r.PUT("/p2p", h.P2p)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("PUT", "/p2p", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestCreateServices(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStart, service avitotask.Service)

	testTable := []struct {
		name               string
		inputBody          string
		inputService       avitotask.Service
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"sid":1,"price":100}`,
			inputService: avitotask.Service{
				Id:    1,
				Price: 100,
			},
			mockBehavior: func(s *mock_service.MockStart, service avitotask.Service) {
				s.EXPECT().CreateServices(service).Return(nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"Service created successfully"}`,
		},
		{
			name:      "TestBad Fake db err",
			inputBody: `{"sid":1,"price":100}`,
			inputService: avitotask.Service{
				Id:    1,
				Price: 100,
			},
			mockBehavior: func(s *mock_service.MockStart, service avitotask.Service) {
				s.EXPECT().CreateServices(service).Return(errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"sid":"-","price":100}`,
			mockBehavior:       func(s *mock_service.MockStart, service avitotask.Service) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
		{
			name:               "TestBad (negative values)",
			inputBody:          `{"sid":2,"price":-100}`,
			mockBehavior:       func(s *mock_service.MockStart, service avitotask.Service) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"negative values"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			start := mock_service.NewMockStart(c)
			testCase.mockBehavior(start, testCase.inputService)

			services := &service.Services{Start: start}

			h := NewHandler(services)

			r := gin.New()
			r.POST("/services", h.CreateServices)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/services", bytes.NewBufferString(testCase.inputBody))

			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestMakeOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockService, order avitotask.Order)

	testTable := []struct {
		name               string
		inputBody          string
		inputOrder         avitotask.Order
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:      "TestOK",
			inputBody: `{"sid":1,"uid":13}`,
			inputOrder: avitotask.Order{
				SId: 1,
				UId: 13,
			},
			mockBehavior: func(s *mock_service.MockService, order avitotask.Order) {
				s.EXPECT().MakeOrder(order).Return("Order created", nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":"Order created"}`,
		},
		{
			name:      "TestBad Fake db err",
			inputBody: `{"sid":1,"uid":13}`,
			inputOrder: avitotask.Order{
				SId: 1,
				UId: 13,
			},
			mockBehavior: func(s *mock_service.MockService, order avitotask.Order) {
				s.EXPECT().MakeOrder(order).Return("", errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Body",
			inputBody:          `{"sid":"-","uid":100}`,
			mockBehavior:       func(s *mock_service.MockService, order avitotask.Order) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad JSON"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			serv := mock_service.NewMockService(c)
			testCase.mockBehavior(serv, testCase.inputOrder)

			services := &service.Services{Service: serv}

			h := NewHandler(services)

			r := gin.New()
			r.POST("/new-order", h.MakeOrder)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/new-order", bytes.NewBufferString(testCase.inputBody))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestListServicesr(t *testing.T) {
	type mockBehavior func(s *mock_service.MockStart)

	testTable := []struct {
		name               string
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name: "TestOK",
			mockBehavior: func(s *mock_service.MockStart) {
				var tmp []avitotask.Service
				s.EXPECT().ShowServices().Return(tmp, nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":null}`,
		},
		{
			name: "TestBad Fake db err",
			mockBehavior: func(s *mock_service.MockStart) {
				var tmp []avitotask.Service
				s.EXPECT().ShowServices().Return(tmp, errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			start := mock_service.NewMockStart(c)
			testCase.mockBehavior(start)

			services := &service.Services{Start: start}

			h := NewHandler(services)

			r := gin.New()
			r.GET("/", h.ListServices)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", bytes.NewBufferString(""))
			r.ServeHTTP(w, req)

			assert.Equal(t, testCase.expectedStatusCode, w.Code)
			assert.Equal(t, testCase.expectedBody, w.Body.String())
		})
	}
}

func TestDoOrder(t *testing.T) {
	type mockBehavior func(s *mock_service.MockService, id int)
	testTable := []struct {
		name               string
		inputParam         int
		mockBehavior       mockBehavior
		expectedStatusCode int
		expectedBody       string
	}{
		{
			name:       "TestOK",
			inputParam: 1,
			mockBehavior: func(s *mock_service.MockService, id int) {
				s.EXPECT().DoOrder(id).Return("", nil)
			},
			expectedStatusCode: 200,
			expectedBody:       `{"message":""}`,
		},
		{
			name:       "TestBad Fake db err",
			inputParam: 1,
			mockBehavior: func(s *mock_service.MockService, id int) {
				s.EXPECT().DoOrder(id).Return("", errors.New("db is down"))
			},
			expectedStatusCode: 500,
			expectedBody:       `{"message":"db is down"}`,
		},
		{
			name:               "TestBad Param",
			inputParam:         -1,
			mockBehavior:       func(s *mock_service.MockService, id int) {},
			expectedStatusCode: 400,
			expectedBody:       `{"message":"bad Params"}`,
		},
	}
	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			srv := mock_service.NewMockService(c)
			testCase.mockBehavior(srv, testCase.inputParam)

			services := &service.Services{Service: srv}
			h := NewHandler(services)
			r := gin.New()
			r.PATCH("/do-order/:id", h.DoOrder)
			if testCase.name == "TestBad Param" {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("PATCH", "/do-order/asdhasd", bytes.NewBufferString(""))
				r.ServeHTTP(w, req)
				assert.Equal(t, testCase.expectedStatusCode, w.Code)
				assert.Equal(t, testCase.expectedBody, w.Body.String())
			} else {
				w := httptest.NewRecorder()
				req := httptest.NewRequest("PATCH", fmt.Sprintf("/do-order/%d", testCase.inputParam), bytes.NewBufferString(""))
				r.ServeHTTP(w, req)
				assert.Equal(t, testCase.expectedStatusCode, w.Code)
				assert.Equal(t, testCase.expectedBody, w.Body.String())
			}
		})
	}
}
