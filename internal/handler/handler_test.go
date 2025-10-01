package handler

import (
	"black-pearl/backend-hackathon/internal/domain/pet/entity"
	prizeEntity "black-pearl/backend-hackathon/internal/domain/prize/entity"
	quizEntity "black-pearl/backend-hackathon/internal/domain/quiz/entity"
	sectionItemsEntity "black-pearl/backend-hackathon/internal/domain/sectionItems/entity"
	sectionEntity "black-pearl/backend-hackathon/internal/domain/sections/entity"
	theoryEntity "black-pearl/backend-hackathon/internal/domain/theory/entity"
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
	"go.uber.org/zap"
)

// Mock services
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

type MockQuizService struct {
	mock.Mock
}

func (m *MockQuizService) GetQuiz(ctx context.Context, quizID int64) (*quizEntity.Quiz, error) {
	args := m.Called(ctx, quizID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*quizEntity.Quiz), args.Error(1)
}

type MockSectionService struct {
	mock.Mock
}

func (m *MockSectionService) GetSections(ctx context.Context) (*[]sectionEntity.Sections, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]sectionEntity.Sections), args.Error(1)
}

func (m *MockSectionService) NewSection(ctx context.Context, title string) (*sectionEntity.Sections, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sectionEntity.Sections), args.Error(1)
}

type MockSectionItemsService struct {
	mock.Mock
}

func (m *MockSectionItemsService) GetSectionItemsBySectionID(ctx context.Context, sectionID int64) (*[]sectionItemsEntity.SectionItem, error) {
	args := m.Called(ctx, sectionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*[]sectionItemsEntity.SectionItem), args.Error(1)
}

func (m *MockSectionItemsService) NewSectionItem(ctx context.Context, sectionID int64, title string, isTest bool, itemId int64) (*sectionItemsEntity.SectionItem, error) {
	args := m.Called(ctx, sectionID, title, isTest, itemId)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*sectionItemsEntity.SectionItem), args.Error(1)
}

type MockTheoryService struct {
	mock.Mock
}

func (m *MockTheoryService) GetTheoryByID(ctx context.Context, theoryID int64) (*theoryEntity.Theory, error) {
	args := m.Called(ctx, theoryID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*theoryEntity.Theory), args.Error(1)
}

func (m *MockTheoryService) NewTheory(ctx context.Context, title, content string) (*theoryEntity.Theory, error) {
	args := m.Called(ctx, title, content)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*theoryEntity.Theory), args.Error(1)
}

func TestHandler_GetPet(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockPet := &entity.Pet{
		ID:   1,
		Name: "Buddy",
		Age:  2,
		Exp:  100,
	}
	mockPetSvc.On("GetPetByUserID", mock.Anything, 1).Return(mockPet, nil)

	router := gin.Default()
	router.GET("/pet/:id", handler.GetPet)

	req, _ := http.NewRequest("GET", "/pet/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.GetPetResp
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockPet.ID, response.ID)
	assert.Equal(t, mockPet.Name, response.Name)
	assert.Equal(t, mockPet.Age, response.Age)
	assert.Equal(t, mockPet.Exp, response.Exp)

	mockPetSvc.AssertExpectations(t)
}

func TestHandler_PostName(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	request := dto.SetPetNameReq{
		Name:   "Max",
		UserID: 1,
	}
	mockPetSvc.On("SetName", mock.Anything, request.Name, request.UserID).Return(nil)

	router := gin.Default()
	router.POST("/pet/name", handler.PostName)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/pet/name", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockPetSvc.AssertExpectations(t)
}

func TestHandler_PostXP(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	request := dto.SendXPReq{
		UserID: 1,
		Exp:    50,
	}
	mockPetSvc.On("UpdateXP", mock.Anything, request.Exp, request.UserID).Return(nil)

	router := gin.Default()
	router.POST("/pet/xp", handler.PostXP)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/pet/xp", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	mockPetSvc.AssertExpectations(t)
}

func TestHandler_GetMyPrizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockPrizes := &[]prizeEntity.Prize{
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
	}
	mockPrizeSvc.On("MyPrizes", mock.Anything, 1).Return(mockPrizes, nil)

	router := gin.Default()
	router.GET("/prizes/:id/my", handler.GetMyPrizes)

	req, _ := http.NewRequest("GET", "/prizes/1/my", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.GetPrizesResp
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Prizes, len(*mockPrizes))
	assert.Equal(t, (*mockPrizes)[0].Title, response.Prizes[0].Title)

	mockPrizeSvc.AssertExpectations(t)
}

func TestHandler_GetAvailablePrizes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockPrizes := &[]prizeEntity.Prize{
		{
			Title:       "Available Prize 1",
			Description: "Available Description 1",
			ImageURL:    "http://example.com/available1.jpg",
		},
	}
	mockPrizeSvc.On("AvailablePrizes", mock.Anything, 1).Return(mockPrizes, nil)

	router := gin.Default()
	router.POST("/prizes/:id/available", handler.GetAvailablePrizes)

	req, _ := http.NewRequest("POST", "/prizes/1/available", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response dto.GetPrizesResp
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Prizes, len(*mockPrizes))
	assert.Equal(t, (*mockPrizes)[0].Title, response.Prizes[0].Title)

	mockPrizeSvc.AssertExpectations(t)
}

func TestHandler_GetQuiz(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockQuiz := &quizEntity.Quiz{
		ID:    1,
		Title: "Test Quiz",
	}
	mockQuizSvc.On("GetQuiz", mock.Anything, int64(1)).Return(mockQuiz, nil)

	router := gin.Default()
	router.GET("/quiz/:id", handler.GetQuiz)

	req, _ := http.NewRequest("GET", "/quiz/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response quizEntity.Quiz
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockQuiz.ID, response.ID)
	assert.Equal(t, mockQuiz.Title, response.Title)

	mockQuizSvc.AssertExpectations(t)
}

func TestHandler_GetSectionsWithItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockSections := &[]sectionEntity.Sections{
		{
			ID:    1,
			Title: "Section 1",
		},
		{
			ID:    2,
			Title: "Section 2",
		},
	}
	mockItems := &[]sectionItemsEntity.SectionItem{
		{
			SectionID: 1,
			Title:     "Item 1",
			IsTest:    false,
			ItemID:    1,
		},
	}

	mockSectionSvc.On("GetSections", mock.Anything).Return(mockSections, nil)
	mockSectionItemsSvc.On("GetSectionItemsBySectionID", mock.Anything, int64(1)).Return(mockItems, nil)
	mockSectionItemsSvc.On("GetSectionItemsBySectionID", mock.Anything, int64(2)).Return(mockItems, nil)

	router := gin.Default()
	router.GET("/sections", handler.GetSectionsWithItems)

	req, _ := http.NewRequest("GET", "/sections", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response []dto.SectionWithItemsResp
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, (*mockSections)[0].Title, response[0].Title)
	assert.Len(t, response[0].Items, 1)

	mockSectionSvc.AssertExpectations(t)
	mockSectionItemsSvc.AssertExpectations(t)
}

func TestHandler_NewSection(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	request := dto.NewSectionReq{
		Title: "New Section",
	}
	mockSection := &sectionEntity.Sections{
		ID:    1,
		Title: request.Title,
	}
	mockSectionSvc.On("NewSection", mock.Anything, request.Title).Return(mockSection, nil)

	router := gin.Default()
	router.POST("/sections", handler.NewSection)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/sections", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response sectionEntity.Sections
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockSection.ID, response.ID)
	assert.Equal(t, mockSection.Title, response.Title)

	mockSectionSvc.AssertExpectations(t)
}

func TestHandler_GetSectionItems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockItems := &[]sectionItemsEntity.SectionItem{
		{
			SectionID: 1,
			Title:     "Item 1",
			IsTest:    false,
			ItemID:    1,
		},
		{
			SectionID: 1,
			Title:     "Item 2",
			IsTest:    true,
			ItemID:    2,
		},
	}
	mockSectionItemsSvc.On("GetSectionItemsBySectionID", mock.Anything, int64(1)).Return(mockItems, nil)

	router := gin.Default()
	router.GET("/sections/:id/items", handler.GetSectionItems)

	req, _ := http.NewRequest("GET", "/sections/1/items", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response []sectionItemsEntity.SectionItem
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response, 2)
	assert.Equal(t, (*mockItems)[0].Title, response[0].Title)

	mockSectionItemsSvc.AssertExpectations(t)
}

func TestHandler_NewSectionItem(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	request := dto.NewSectionItemReq{
		Title:  "New Item",
		IsTest: true,
		ItemID: 5,
	}
	mockItem := &sectionItemsEntity.SectionItem{
		SectionID: 1,
		Title:     request.Title,
		IsTest:    request.IsTest,
		ItemID:    request.ItemID,
	}
	mockSectionItemsSvc.On("NewSectionItem", mock.Anything, int64(1), request.Title, request.IsTest, request.ItemID).Return(mockItem, nil)

	router := gin.Default()
	router.POST("/sections/:id/items", handler.NewSectionItem)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/sections/1/items", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response sectionItemsEntity.SectionItem
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockItem.Title, response.Title)
	assert.Equal(t, mockItem.IsTest, response.IsTest)
	assert.Equal(t, mockItem.ItemID, response.ItemID)

	mockSectionItemsSvc.AssertExpectations(t)
}

func TestHandler_GetTheory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	mockTheory := &theoryEntity.Theory{
		ID:      1,
		Title:   "Test Theory",
		Content: "Theory content",
	}
	mockTheorySvc.On("GetTheoryByID", mock.Anything, int64(1)).Return(mockTheory, nil)

	router := gin.Default()
	router.GET("/theory/:id", handler.GetTheory)

	req, _ := http.NewRequest("GET", "/theory/1", nil)
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var response theoryEntity.Theory
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockTheory.ID, response.ID)
	assert.Equal(t, mockTheory.Title, response.Title)
	assert.Equal(t, mockTheory.Content, response.Content)

	mockTheorySvc.AssertExpectations(t)
}

func TestHandler_NewTheory(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mockPetSvc := new(MockPetService)
	mockPrizeSvc := new(MockPrizeService)
	mockQuizSvc := new(MockQuizService)
	mockSectionSvc := new(MockSectionService)
	mockSectionItemsSvc := new(MockSectionItemsService)
	mockTheorySvc := new(MockTheoryService)
	logger := zap.NewNop().Sugar()

	handler := NewHandler(mockQuizSvc, mockPetSvc, mockSectionSvc, mockSectionItemsSvc, mockTheorySvc, mockPrizeSvc, logger)

	request := dto.NewTheoryReq{
		Title:   "New Theory",
		Content: "Theory content",
	}
	mockTheory := &theoryEntity.Theory{
		ID:      1,
		Title:   request.Title,
		Content: request.Content,
	}
	mockTheorySvc.On("NewTheory", mock.Anything, request.Title, request.Content).Return(mockTheory, nil)

	router := gin.Default()
	router.POST("/theory", handler.NewTheory)

	requestBody, _ := json.Marshal(request)
	req, _ := http.NewRequest("POST", "/theory", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	router.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code)

	var response theoryEntity.Theory
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, mockTheory.ID, response.ID)
	assert.Equal(t, mockTheory.Title, response.Title)
	assert.Equal(t, mockTheory.Content, response.Content)

	mockTheorySvc.AssertExpectations(t)
}
