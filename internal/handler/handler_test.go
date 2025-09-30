package handler

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	taskEntity "black-pearl/backend-hackathon/internal/domain/task/entity"
	"black-pearl/backend-hackathon/internal/handler/dto"
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock services with correct context types
type MockPetService struct {
	mock.Mock
}

func (m *MockPetService) UpdateXP(ctx context.Context, xp int, userID int) error {
	args := m.Called(ctx, xp, userID)
	return args.Error(0)
}

func (m *MockPetService) GetPetByUserID(ctx context.Context, userID int) (*entity.Pet, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Pet), args.Error(1)
}

func (m *MockPetService) SetName(ctx context.Context, name string, userID int) error {
	args := m.Called(ctx, name, userID)
	return args.Error(0)
}

type MockPrizeService struct {
	mock.Mock
}

func (m *MockPrizeService) AvailablePrizes(ctx context.Context, userID int) (*[]prizeEntity.Prize, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]prizeEntity.Prize), args.Error(1)
}

func (m *MockPrizeService) MyPrizes(ctx context.Context, user_id int) (*[]prizeEntity.Prize, error) {
	args := m.Called(ctx, user_id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]prizeEntity.Prize), args.Error(1)
}

type MockTaskService struct {
	mock.Mock
}

func (m *MockTaskService) Task(ctx context.Context, taskID int) (*taskEntity.Task, error) {
	args := m.Called(ctx, taskID)
	return args.Get(0).(*taskEntity.Task), args.Error(1)
}

func TestHandler_GetPet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockPet        *entity.Pet
		mockError      error
		expectedStatus int
	}{
		{
			name:   "successful get pet",
			userID: "1",
			mockPet: &entity.Pet{
				ID:   1,
				Name: "Buddy",
				Age:  2,
				Exp:  100,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetSvc := new(MockPetService)
			mockPrizeSvc := new(MockPrizeService)
			mockTaskSvc := new(MockTaskService)

			handler := NewHandler(mockTaskSvc, mockPetSvc, mockPrizeSvc)

			mockPetSvc.On("GetPetByUserID", mock.Anything, 1).Return(tt.mockPet, tt.mockError)

			router := gin.Default()
			router.GET("/pet/:id", handler.GetPet)

			req, _ := http.NewRequest("GET", "/pet/"+tt.userID, nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusOK {
				var response dto.GetPetResp
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, tt.mockPet.ID, response.ID)
				assert.Equal(t, tt.mockPet.Name, response.Name)
				assert.Equal(t, tt.mockPet.Age, response.Age)
				assert.Equal(t, tt.mockPet.Exp, response.Exp)
			}

			mockPetSvc.AssertExpectations(t)
		})
	}
}

func TestHandler_PostName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        dto.SetPetNameReq
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful set pet name",
			request: dto.SetPetNameReq{
				Name:   "Max",
				UserID: 1,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetSvc := new(MockPetService)
			mockPrizeSvc := new(MockPrizeService)
			mockTaskSvc := new(MockTaskService)

			handler := NewHandler(mockTaskSvc, mockPetSvc, mockPrizeSvc)

			mockPetSvc.On("SetName", mock.Anything, tt.request.Name, tt.request.UserID).Return(tt.mockError)

			router := gin.Default()
			router.POST("/pet/name", handler.PostName)

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/pet/name", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockPetSvc.AssertExpectations(t)
		})
	}
}

func TestHandler_PostXP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		request        dto.SendXPReq
		mockError      error
		expectedStatus int
	}{
		{
			name: "successful update XP",
			request: dto.SendXPReq{
				UserID: 1,
				Exp:    50,
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetSvc := new(MockPetService)
			mockPrizeSvc := new(MockPrizeService)
			mockTaskSvc := new(MockTaskService)

			handler := NewHandler(mockTaskSvc, mockPetSvc, mockPrizeSvc)

			mockPetSvc.On("UpdateXP", mock.Anything, tt.request.Exp, tt.request.UserID).Return(tt.mockError)

			router := gin.Default()
			router.POST("/pet/xp", handler.PostXP)

			requestBody, _ := json.Marshal(tt.request)
			req, _ := http.NewRequest("POST", "/pet/xp", bytes.NewBuffer(requestBody))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
			mockPetSvc.AssertExpectations(t)
		})
	}
}

func TestHandler_GetMyPrizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockPrizes     *[]prizeEntity.Prize
		mockError      error
		expectedStatus int
	}{
		{
			name:   "successful get my prizes",
			userID: "1",
			mockPrizes: &[]prizeEntity.Prize{
				{
					Title:       "Prize 1",
					Description: "Description 1",
					ImageURL:    "http://example.com/prize1.jpg",
				},
				{
					Title:       "Prize 2",
					Description: "Description 2",
					ImageURL:    "http://example.com/prize2.jpg",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetSvc := new(MockPetService)
			mockPrizeSvc := new(MockPrizeService)
			mockTaskSvc := new(MockTaskService)

			handler := NewHandler(mockTaskSvc, mockPetSvc, mockPrizeSvc)

			mockPrizeSvc.On("MyPrizes", mock.Anything, 1).Return(tt.mockPrizes, tt.mockError)

			router := gin.Default()
			router.GET("/prizes/:id/my", handler.GetMyPrizes)

			req, _ := http.NewRequest("GET", "/prizes/"+tt.userID+"/my", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusOK {
				var response dto.GetPrizesResp
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Prizes, len(*tt.mockPrizes))
				assert.Equal(t, (*tt.mockPrizes)[0].Title, response.Prizes[0].Title)
			}

			mockPrizeSvc.AssertExpectations(t)
		})
	}
}

func TestHandler_GetAvailablePrizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		userID         string
		mockPrizes     *[]prizeEntity.Prize
		mockError      error
		expectedStatus int
	}{
		{
			name:   "successful get available prizes",
			userID: "1",
			mockPrizes: &[]prizeEntity.Prize{
				{
					Title:       "Available Prize 1",
					Description: "Available Description 1",
					ImageURL:    "http://example.com/available1.jpg",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockPetSvc := new(MockPetService)
			mockPrizeSvc := new(MockPrizeService)
			mockTaskSvc := new(MockTaskService)

			handler := NewHandler(mockTaskSvc, mockPetSvc, mockPrizeSvc)

			mockPrizeSvc.On("AvailablePrizes", mock.Anything, 1).Return(tt.mockPrizes, tt.mockError)

			router := gin.Default()
			router.POST("/prizes/:id/available", handler.GetAvailablePrizes)

			req, _ := http.NewRequest("POST", "/prizes/"+tt.userID+"/available", nil)
			resp := httptest.NewRecorder()

			router.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)

			if tt.expectedStatus == http.StatusOK {
				var response dto.GetPrizesResp
				err := json.Unmarshal(resp.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Len(t, response.Prizes, len(*tt.mockPrizes))
				assert.Equal(t, (*tt.mockPrizes)[0].Title, response.Prizes[0].Title)
			}

			mockPrizeSvc.AssertExpectations(t)
		})
	}
}
