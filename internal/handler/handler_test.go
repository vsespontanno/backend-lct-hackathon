package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	petEntity "black-pearl/backend-hackathon/internal/domain/pet/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type MockPetServiceInterface struct {
	mock.Mock
}

func (m *MockPetServiceInterface) GetPetByUserID(ctx context.Context, userID int) (*petEntity.Pet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*petEntity.Pet), args.Error(1)
}

func (m *MockPetServiceInterface) SetName(ctx context.Context, name string, userID int) error {
	args := m.Called(ctx, name, userID)
	return args.Error(0)
}

type MockTaskServiceInterface struct {
	mock.Mock
}

func (m *MockTaskServiceInterface) Task(ctx context.Context, taskID int64) (*taskEntity.Task, error) {
	args := m.Called(ctx, taskID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*taskEntity.Task), args.Error(1)
}

func TestHandler_GetPet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPet := &petEntity.Pet{
		ID:    1,
		Name:  "Fluffy",
		Age:   2,
		Exp:   100,
		Level: 3,
	}

	tests := []struct {
		name           string
		userID         string
		setupMock      func(*MockPetServiceInterface)
		expectedStatus int
		expectedBody   dto.GetPetReq
		expectedError  string
	}{
		{
			name:   "successful get pet",
			userID: "123",
			setupMock: func(m *MockPetServiceInterface) {
				m.On("GetPetByUserID", mock.Anything, 123).Return(mockPet, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: dto.GetPetReq{
				ID:    1,
				Name:  "Fluffy",
				Age:   2,
				Exp:   100,
				Level: 3,
			},
		},
		{
			name:   "pet not found",
			userID: "456",
			setupMock: func(m *MockPetServiceInterface) {
				m.On("GetPetByUserID", mock.Anything, 456).Return(nil, errors.New("pet not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "pet not found",
		},
		{
			name:   "service returns database error",
			userID: "789",
			setupMock: func(m *MockPetServiceInterface) {
				m.On("GetPetByUserID", mock.Anything, 789).Return(nil, errors.New("database connection failed"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database connection failed",
		},
		{
			name:   "zero user ID",
			userID: "0",
			setupMock: func(m *MockPetServiceInterface) {
				m.On("GetPetByUserID", mock.Anything, 0).Return(nil, errors.New("user not found"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetService := new(MockPetServiceInterface)
			mockTaskService := new(MockTaskServiceInterface)
			tt.setupMock(mockPetService)

			handler := NewHandler(mockTaskService, mockPetService)

			router := gin.New()
			router.GET("/pet/:userID", handler.GetPet)

			req, err := http.NewRequest("GET", "/pet/"+tt.userID, nil)
			require.NoError(t, err)

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedError != "" {
				var response gin.H
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			} else {
				var response dto.GetPetReq
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedBody, response)
			}

			mockPetService.AssertExpectations(t)
		})
	}
}

func TestHandler_PostName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		requestBody    dto.SetPetNameReq
		setupMock      func(*MockPetServiceInterface)
		expectedStatus int
		expectedError  string
	}{
		{
			name: "successful set name",
			requestBody: dto.SetPetNameReq{
				Name:   "Buddy",
				UserID: 123,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "Buddy", 123).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "successful set name with empty name",
			requestBody: dto.SetPetNameReq{
				Name:   "",
				UserID: 456,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "", 456).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "zero user ID",
			requestBody: dto.SetPetNameReq{
				Name:   "ZeroPet",
				UserID: 0,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "ZeroPet", 0).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "negative user ID",
			requestBody: dto.SetPetNameReq{
				Name:   "NegativePet",
				UserID: -1,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "NegativePet", -1).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "service returns error",
			requestBody: dto.SetPetNameReq{
				Name:   "Max",
				UserID: 999,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "Max", 999).Return(errors.New("database error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedError:  "database error",
		},
		{
			name: "very long name",
			requestBody: dto.SetPetNameReq{
				Name:   "VeryLongPetNameThatExceedsNormalLengthButShouldStillWork",
				UserID: 123,
			},
			setupMock: func(m *MockPetServiceInterface) {
				m.On("SetName", mock.Anything, "VeryLongPetNameThatExceedsNormalLengthButShouldStillWork", 123).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetService := new(MockPetServiceInterface)
			mockTaskService := new(MockTaskServiceInterface)
			tt.setupMock(mockPetService)

			handler := NewHandler(mockTaskService, mockPetService)

			router := gin.New()
			router.POST("/pet/:userID", handler.PostName)

			reqBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/pet/123", bytes.NewBuffer(reqBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedError != "" {
				var response gin.H
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedError, response["error"])
			}

			mockPetService.AssertExpectations(t)
		})
	}
}

func TestHandler_EdgeCases(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("context propagation to service", func(t *testing.T) {
		mockPetService := new(MockPetServiceInterface)
		mockTaskService := new(MockTaskServiceInterface)

		mockPet := &petEntity.Pet{
			ID:    1,
			Name:  "TestPet",
			Age:   1,
			Exp:   0,
			Level: 1,
		}

		mockPetService.On("GetPetByUserID", mock.Anything, 123).Return(mockPet, nil)

		handler := NewHandler(mockTaskService, mockPetService)

		router := gin.New()
		router.GET("/pet/:userID", handler.GetPet)

		req, _ := http.NewRequest("GET", "/pet/123", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		assert.Equal(t, http.StatusOK, resp.Code)
		mockPetService.AssertExpectations(t)
	})

	t.Run("concurrent requests", func(t *testing.T) {
		mockPetService := new(MockPetServiceInterface)
		mockTaskService := new(MockTaskServiceInterface)

		mockPet := &petEntity.Pet{
			ID:    1,
			Name:  "ConcurrentPet",
			Age:   1,
			Exp:   0,
			Level: 1,
		}

		mockPetService.On("GetPetByUserID", mock.Anything, 123).Return(mockPet, nil).Times(3)

		handler := NewHandler(mockTaskService, mockPetService)

		router := gin.New()
		router.GET("/pet/:userID", handler.GetPet)

		done := make(chan bool, 3)

		for i := 0; i < 3; i++ {
			go func() {
				req, _ := http.NewRequest("GET", "/pet/123", nil)
				resp := httptest.NewRecorder()
				router.ServeHTTP(resp, req)
				assert.Equal(t, http.StatusOK, resp.Code)
				done <- true
			}()
		}

		for i := 0; i < 3; i++ {
			<-done
		}

		mockPetService.AssertExpectations(t)
	})
}
