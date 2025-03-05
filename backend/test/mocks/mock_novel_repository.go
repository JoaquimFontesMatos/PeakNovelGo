package mocks

import (
	"backend/internal/models"
	"github.com/stretchr/testify/mock"
)

type MockNovelRepository struct {
	mock.Mock
}

// IsDown checks if the database is offline by pinging it
func (m *MockNovelRepository) IsDown() bool {
	args := m.Called()
	return args.Bool(0)
}

// CreateNovel creates a new novel in the database
func (m *MockNovelRepository) CreateNovel(novel models.Novel) (*models.Novel, error) {
	args := m.Called(novel)
	
	// Avoid panic if nil is returned
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	
	return args.Get(0).(*models.Novel), args.Error(1)
}


// GetNovels gets a list of novels
func (m *MockNovelRepository) GetNovels(page, limit int) ([]models.Novel, int64, error) {
	args := m.Called(page, limit)
	return args.Get(0).([]models.Novel), args.Get(1).(int64), args.Error(2)
}

// GetNovelsByAuthorName gets a list of novels by author name
func (m *MockNovelRepository) GetNovelsByAuthorName(authorName string, page, limit int) ([]models.Novel, int64, error) {
	args := m.Called(authorName, page, limit)
	return args.Get(0).([]models.Novel), args.Get(1).(int64), args.Error(2)
}

// GetNovelsByGenreName gets a list of novels by genre name
func (m *MockNovelRepository) GetNovelsByGenreName(genreName string, page, limit int) ([]models.Novel, int64, error) {
	args := m.Called(genreName, page, limit)
	return args.Get(0).([]models.Novel), args.Get(1).(int64), args.Error(2)
}

// GetNovelsByTagName gets a list of novels by tag name
func (m *MockNovelRepository) GetNovelsByTagName(tagName string, page, limit int) ([]models.Novel, int64, error) {
	args := m.Called(tagName, page, limit)
	return args.Get(0).([]models.Novel), args.Get(1).(int64), args.Error(2)
}

// GetNovelByID gets a novel by ID
func (m *MockNovelRepository) GetNovelByID(id uint) (*models.Novel, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Novel), args.Error(1)
}

// GetNovelByUpdatesID gets a novel by novel updates id
func (m *MockNovelRepository) GetNovelByUpdatesID(title string) (*models.Novel, error) {
	args := m.Called(title)
	return args.Get(0).(*models.Novel), args.Error(1)
}
