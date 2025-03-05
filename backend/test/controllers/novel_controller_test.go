package controllers

import (
	"backend/internal/controllers"
	"backend/internal/models"
	"backend/internal/services"
	"backend/internal/types"
	"backend/test/mocks"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandleImportNovel(t *testing.T) {
	tests := []struct {
		name                       string
		novelUpdatesID             string
		expectedCode               int
		isAuthorized               bool
		hasPermissions             bool
		hasStableNetworkConnection bool
		isSourceWebsiteAvailable   bool
		isDatabaseOnline           bool
		isScriptReady              bool
		isNovelAlreadyImported     bool
	}{
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_USER_HAS_PERMISSIONS-SUCCESS",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusCreated,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_USER_LACKS_PERMISSIONS-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusForbidden,
			isAuthorized:               true,
			hasPermissions:             false,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_USER_UNAUTHORIZED-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusUnauthorized,
			isAuthorized:               false,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NO_NETWORK_CONNECTIVITY-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusServiceUnavailable,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: false,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_SOURCE_WEBSITE_DOWN-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusServiceUnavailable,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   false,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_DATABASE_OFFLINE-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusServiceUnavailable,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           false,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_SCRIPT_MISCONFIGURED-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusServiceUnavailable,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              false,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_AND_SOURCE_NOVEL_DOESNT_EXISTS-FAILURE",
			novelUpdatesID:             "aaa",
			expectedCode:               http.StatusNotFound,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_ALREADY_CREATED-FAILURE",
			novelUpdatesID:             "reverend-insanity",
			expectedCode:               http.StatusConflict,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     true,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_UPPERCASE-SUCCESS",
			novelUpdatesID:             "Reverend-Insanity",
			expectedCode:               http.StatusCreated,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_DIVIDED_BY_SPACES-SUCCESS",
			novelUpdatesID:             "reverend%20insanity",
			expectedCode:               http.StatusCreated,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_SPECIAL_CHARS-FAILURE",
			novelUpdatesID:             "reverend-@insanity",
			expectedCode:               http.StatusBadRequest,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_ONE_WORD-SUCCESS",
			novelUpdatesID:             "turning",
			expectedCode:               http.StatusCreated,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_MINIMUM_VALUE-SUCCESS",
			novelUpdatesID:             "a",
			expectedCode:               http.StatusNotFound,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_MAXIMUM_VALUE-SUCCESS",
			novelUpdatesID:             strings.Repeat("a", 255),
			expectedCode:               http.StatusNotFound,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			novelUpdatesID:             "",
			expectedCode:               http.StatusNotFound,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
		{
			name:                       "#TCO-B3/F1/IMPORT_NOVEL_NOVELUPDATESID_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			novelUpdatesID:             strings.Repeat("a", 256),
			expectedCode:               http.StatusBadRequest,
			isAuthorized:               true,
			hasPermissions:             true,
			hasStableNetworkConnection: true,
			isSourceWebsiteAvailable:   true,
			isDatabaseOnline:           true,
			isScriptReady:              true,
			isNovelAlreadyImported:     false,
		},
	}

	dir, _ := os.Getwd()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})
			err := os.Chdir(dir)
			if err != nil {
				fmt.Println("Error changing directory:", err)
			}

			// Mock user data
			user := models.User{
				Username:           "example",
				Email:              "test@example.com",
				Password:           "12345678",
				EmailVerified:      false,
				VerificationToken:  "",
				ProfilePicture:     "https://example.com/profile2.jpg",
				Bio:                "I am a software engineer",
				Roles:              "user",
				LastLogin:          time.Time{},
				DateOfBirth:        time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
				PreferredLanguage:  "ru",
				ReadingPreferences: "{}",
				IsDeleted:          false,
			}

			if tt.isNovelAlreadyImported {
				novel := models.Novel{
					Title:          "Test",
					Synopsis:       "Test",
					CoverUrl:       "https://example.com/cover.jpg",
					Language:       "en",
					Status:         "ongoing",
					NovelUpdatesID: tt.novelUpdatesID,
				}

				_, err := novelRepo.CreateNovel(novel)
				if err != nil {
					t.Errorf("Failed to create novel: %v", err)
				}
			}

			if tt.isScriptReady {
				err := os.Chdir("../../")
				if err != nil {
					fmt.Println("Error changing directory:", err)
				}
			}

			if !tt.isSourceWebsiteAvailable {
				scriptExecutor = &mocks.MockScriptExecutorSourceWebsiteDown{}
			}

			if !tt.hasStableNetworkConnection {
				scriptExecutor = &mocks.MockScriptExecutorNetworkDown{}
			}

			novelService = services.NewNovelService(novelRepo, scriptExecutor)

			if !tt.isDatabaseOnline {
				// Create a mock repository
				mockRepo := new(mocks.MockNovelRepository)

				// Simulate database offline
				mockRepo.On("IsDown").Return(true)

				novelService = services.NewNovelService(mockRepo, scriptExecutor)
				mockRepo.On("CreateNovel", mock.AnythingOfType("models.Novel")).Return((*models.Novel)(nil), types.WrapError(types.DATABASE_ERROR, "Failed to create novel", nil))
			}

			novelController = *controllers.NewNovelController(novelService)

			if tt.hasPermissions {
				user.Roles = "admin;user"
			}

			err = userRepo.CreateUser(&user)
			if err != nil {
				t.Errorf("Failed to create user: %v", err)
			}
			createdUser, err := userRepo.GetUserByID(1)
			if err != nil {
				t.Errorf("Failed to get user: %v", err)
			}

			accessToken, refreshToken, err := authService.GenerateToken(createdUser)

			if err != nil {
				t.Errorf("Failed to generate token: %v", err)

			}

			// Create a request
			req := httptest.NewRequest(http.MethodPost, "/novels/"+tt.novelUpdatesID, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")
			if tt.isAuthorized {
				req.Header.Set("Authorization", "Bearer "+accessToken)
				req.AddCookie(&http.Cookie{
					Name:  "refreshToken",
					Value: refreshToken,
				})
			}

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.POST("/novels/:novel_updates_id", myMiddleware.AuthMiddleware(), myMiddleware.PermissionMiddleware("novels", "create"), novelController.HandleImportNovelByNovelUpdatesID)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusCreated {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusCreated {
				_, total, err := novelService.GetNovels(1, 10)
				if total <= 0 {
					t.Errorf("Novel not imported")
				}

				if err != nil {
					t.Errorf("Failed to get novel by updates ID: %v", err)
				}
			}
		})
	}
}
func TestHandleGetNovelByNovelUpdatesID(t *testing.T) {
	tests := []struct {
		name             string
		novelUpdatesID   string
		expectedCode     int
		isDatabaseOnline bool
		isNovelCreated   bool
	}{
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BASE_CASE-SUCCESS",
			novelUpdatesID:   "reverend-insanity",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_DATABASE_OFFLINE-FAILURE",
			novelUpdatesID:   "reverend-insanity",
			expectedCode:     http.StatusServiceUnavailable,
			isDatabaseOnline: false,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVEL_NOT_IMPORTED-FAILURE",
			novelUpdatesID:   "reverend-insanity",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isNovelCreated:   false,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_UPPERCASE-SUCCESS",
			novelUpdatesID:   "Reverend-Insanity",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_DIVIDED_BY_SPACES-SUCCESS",
			novelUpdatesID:   "reverend%20insanity",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_SPECIAL_CHARS-FAILURE",
			novelUpdatesID:   "reverend-@insanity",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_ONE_WORD-SUCCESS",
			novelUpdatesID:   "turning",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_MINIMUM_VALUE-SUCCESS",
			novelUpdatesID:   "a",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_MAXIMUM_VALUE-SUCCESS",
			novelUpdatesID:   strings.Repeat("a", 255),
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			novelUpdatesID:   "",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
		{
			name:             "#TCO-B3/F2/GET_NOVEL_BY_NOVELUPDATESID_BUT_NOVELUPDATESID_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			novelUpdatesID:   strings.Repeat("a", 256),
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isNovelCreated:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			if tt.isNovelCreated {
				novel := models.Novel{
					Title:          "Test",
					Synopsis:       "Test",
					CoverUrl:       "https://example.com/cover.jpg",
					Language:       "en",
					Status:         "ongoing",
					NovelUpdatesID: "reverend-insanity",
				}

				novel2 := models.Novel{
					Title:          "Test2",
					Synopsis:       "Test2",
					CoverUrl:       "https://example.com/cover2.jpg",
					Language:       "en",
					Status:         "ongoing",
					NovelUpdatesID: "turning",
				}

				_, err := novelRepo.CreateNovel(novel)
				if err != nil {
					t.Errorf("Failed to create novel: %v", err)
				}

				_, err = novelRepo.CreateNovel(novel2)
				if err != nil {
					t.Errorf("Failed to create novel2: %v", err)
				}
			}

			if !tt.isDatabaseOnline {
				// Create a mock repository
				mockRepo := new(mocks.MockNovelRepository)

				// Simulate database offline
				mockRepo.On("IsDown").Return(true)

				novelService = services.NewNovelService(mockRepo, scriptExecutor)
				mockRepo.On("GetNovelByUpdatesID", mock.AnythingOfType("string")).Return((*models.Novel)(nil), types.WrapError(types.DATABASE_ERROR, "Database offline", nil))
			}

			novelController = *controllers.NewNovelController(novelService)

			// Create a request
			req := httptest.NewRequest(http.MethodGet, "/novels/title/"+tt.novelUpdatesID, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/novels/title/:title", novelController.GetNovelByUpdatesID)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusCreated {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusCreated {
				_, total, err := novelService.GetNovels(1, 10)
				if total <= 0 {
					t.Errorf("Novel not imported")
				}

				if err != nil {
					t.Errorf("Failed to get novel by updates ID: %v", err)
				}
			}
		})
	}
}

func TestHandleGetNovelsByAuthor(t *testing.T) {
	tests := []struct {
		name             string
		author           string
		page             string
		limit            string
		expectedCode     int
		isDatabaseOnline bool
		isFound          bool
		hasTestParam     bool
		isPageNotFound   bool
	}{
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BASE_CASE-SUCCESS",
			author:           "John",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_DATABASE_OFFLINE-FAILURE",
			author:           "John",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusServiceUnavailable,
			isDatabaseOnline: false,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_NO_NOVELS_BY_AUTHOR-FAILURE",
			author:           "John2",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_NO_PARAMS-SUCCESS",
			author:           "John",
			page:             "",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_ONLY_PAGE_PARAM-SUCCESS",
			author:           "John",
			page:             "2",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_ONLY_LIMIT_PARAM-SUCCESS",
			author:           "John",
			page:             "",
			limit:            "11",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_MORE_THAN_ACCEPTED_QUERY_PARAMS-FAILURE",
			author:           "John",
			page:             "1",
			limit:            "1",
			expectedCode:     http.StatusBadRequest,
			hasTestParam:     true,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_AUTHOR_IS_A_SPACE-FAILURE",
			author:           "%20",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_IS_NOT_AN_INT-FAILURE",
			author:           "John",
			page:             "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_DOES_NOT_EXIST-FAILURE",
			author:           "John",
			page:             "2",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
			isPageNotFound:   true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_LIMIT_IS_NOT_AN_INT-FAILURE",
			author:           "John",
			limit:            "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_AUTHOR_IS_MINIMUM_LENGTH-SUCCESS",
			author:           "a",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_AUTHOR_IS_MAXIMUM_LENGTH-SUCCESS",
			author:           strings.Repeat("a", 255),
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_AUTHOR_IS_JUST_BELOW_MINIMUM_LENGTH-FAILURE",
			author:           "",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_AUTHOR_IS_JUST_ABOVE_MAXIMUM_LENGTH-FAILURE",
			author:           strings.Repeat("a", 256),
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_IS_MINIMUM_VALUE-SUCCESS",
			author:           "John",
			page:             "1",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_IS_MAXIMUM_VALUE-SUCCESS",
			author:           "John",
			page:             "1000",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			author:           "John",
			page:             "0",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_PAGE_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			author:           "John",
			page:             "1001",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_LIMIT_IS_MINIMUM_VALUE-SUCCESS",
			author:           "John",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_LIMIT_IS_MAXIMUM_VALUE-SUCCESS",
			author:           "John",
			limit:            "100",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_LIMIT_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			author:           "John",
			limit:            "9",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F3/GET_NOVEL_BY_AUTHOR_BUT_LIMIT_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			author:           "John",
			limit:            "101",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			novel := models.Novel{
				Title:          "t",
				Synopsis:       "t",
				CoverUrl:       "https://example.com/t.jpg",
				Language:       "en",
				Status:         "ongoing",
				NovelUpdatesID: "t",
				Authors:        []models.Author{{Name: "test"}},
			}

			_, err := novelRepo.CreateNovel(novel)
			if err != nil {
				t.Errorf("Failed to create novel2: %v", err)
			}

			// Parse page parameter
			pageStr := tt.page
			if pageStr == "" {
				pageStr = "1" // Default to page 1 if not provided
			}
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil {
				pageInt = 1 // Default to page 1 if parsing fails
			}

			limitStr := tt.limit
			if limitStr == "" {
				limitStr = "10" // Default to page 1 if not provided
			}
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				limitInt = 10 // Default to page 1 if parsing fails
			}

			if tt.isPageNotFound {
				pageInt = 1
			}

			if tt.isFound {
				for i := range pageInt * limitInt {
					novel2 := models.Novel{
						Title:          fmt.Sprintf("Test%d", i),
						Synopsis:       fmt.Sprintf("Test%d", i),
						CoverUrl:       fmt.Sprintf("https://example.com/t%d.jpg", i),
						Language:       "en",
						Status:         "ongoing",
						NovelUpdatesID: fmt.Sprintf("t%d", i),
						Authors:        []models.Author{{Name: tt.author}},
					}

					_, err := novelRepo.CreateNovel(novel2)
					if err != nil {
						t.Errorf("Failed to create novel: %v", err)
					}
				}
			}

			if !tt.isDatabaseOnline {
				mockRepo := new(mocks.MockNovelRepository)
				mockRepo.On("IsDown").Return(true)
				mockRepo.On("GetNovelsByAuthorName", mock.Anything, mock.Anything, mock.Anything).Return(([]models.Novel)(nil), (int64)(0), types.WrapError(types.DATABASE_ERROR, "Database is offline", nil))
				novelService = services.NewNovelService(mockRepo, scriptExecutor)
			}

			novelController = *controllers.NewNovelController(novelService)

			// Build the query parameters
			queryParams := url.Values{}
			if tt.page != "" {
				queryParams.Add("page", tt.page)
			}
			if tt.limit != "" {
				queryParams.Add("limit", tt.limit)
			}
			if tt.hasTestParam {
				queryParams.Add("test", "test")
			}

			// Construct the URL with query parameters
			urlPath := "/novels/authors/" + tt.author
			if len(queryParams) > 0 {
				urlPath += "?" + queryParams.Encode()
			}

			// Create a request
			req := httptest.NewRequest(http.MethodGet, urlPath, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/novels/authors/:author_name", novelController.GetNovelsByAuthorName)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusOK {
				// Fetch novels from the database
				novels, total, err := novelRepo.GetNovelsByAuthorName(tt.author, pageInt, limitInt)
				if err != nil {
					t.Errorf("Failed to fetch novels from the database: %v", err)
					return
				}

				// Assert the number of novels returned
				expectedCount := limitInt
				if len(novels) != expectedCount {
					t.Errorf("Expected %d novels, but got %d", expectedCount, len(novels))
				}

				// Assert the total number of novels for the author
				expectedTotal := pageInt * limitInt
				if total != int64(expectedTotal) {
					t.Errorf("Expected total novels to be %d, but got %d", expectedTotal, total)
				}

				// Optionally, assert the content of the novels
				for i, novel := range novels {
					expectedTitle := fmt.Sprintf("Test%d", (pageInt-1)*limitInt+i)
					if novel.Title != expectedTitle {
						t.Errorf("Expected novel title to be %s, but got %s", expectedTitle, novel.Title)
					}
				}
			}
		})
	}
}

func TestHandleGetNovelsByGenre(t *testing.T) {
	tests := []struct {
		name             string
		genre            string
		page             string
		limit            string
		expectedCode     int
		isDatabaseOnline bool
		isFound          bool
		hasTestParam     bool
		isPageNotFound   bool
	}{
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BASE_CASE-SUCCESS",
			genre:            "Action",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_DATABASE_OFFLINE-FAILURE",
			genre:            "Action",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusServiceUnavailable,
			isDatabaseOnline: false,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_NO_NOVELS_BY_GENRE-FAILURE",
			genre:            "Action2",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_NO_PARAMS-SUCCESS",
			genre:            "Action",
			page:             "",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_ONLY_PAGE_PARAM-SUCCESS",
			genre:            "Action",
			page:             "2",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_ONLY_LIMIT_PARAM-SUCCESS",
			genre:            "Action",
			page:             "",
			limit:            "11",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_MORE_THAN_ACCEPTED_QUERY_PARAMS-FAILURE",
			genre:            "Action",
			page:             "1",
			limit:            "1",
			expectedCode:     http.StatusBadRequest,
			hasTestParam:     true,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_GENRE_IS_A_SPACE-FAILURE",
			genre:            "%20",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_IS_NOT_AN_INT-FAILURE",
			genre:            "Action",
			page:             "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_DOES_NOT_EXIST-FAILURE",
			genre:            "Action",
			page:             "2",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
			isPageNotFound:   true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_LIMIT_IS_NOT_AN_INT-FAILURE",
			genre:            "Action",
			limit:            "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_GENRE_IS_MINIMUM_LENGTH-SUCCESS",
			genre:            "a",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_GENRE_IS_MAXIMUM_LENGTH-SUCCESS",
			genre:            strings.Repeat("a", 255),
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_GENRE_IS_JUST_BELOW_MINIMUM_LENGTH-FAILURE",
			genre:            "",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_GENRE_IS_JUST_ABOVE_MAXIMUM_LENGTH-FAILURE",
			genre:            strings.Repeat("a", 256),
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_IS_MINIMUM_VALUE-SUCCESS",
			genre:            "Action",
			page:             "1",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_IS_MAXIMUM_VALUE-SUCCESS",
			genre:            "Action",
			page:             "1000",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			genre:            "Action",
			page:             "0",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_PAGE_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			genre:            "Action",
			page:             "1001",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_LIMIT_IS_MINIMUM_VALUE-SUCCESS",
			genre:            "Action",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_LIMIT_IS_MAXIMUM_VALUE-SUCCESS",
			genre:            "Action",
			limit:            "100",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_LIMIT_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			genre:            "Action",
			limit:            "9",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F4/GET_NOVEL_BY_GENRE_BUT_LIMIT_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			genre:            "Action",
			limit:            "101",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			novel := models.Novel{
				Title:          "t",
				Synopsis:       "t",
				CoverUrl:       "https://example.com/t.jpg",
				Language:       "en",
				Status:         "ongoing",
				NovelUpdatesID: "t",
				Genres:         []models.Genre{{Name: "test"}},
			}

			_, err := novelRepo.CreateNovel(novel)
			if err != nil {
				t.Errorf("Failed to create novel2: %v", err)
			}

			// Parse page parameter
			pageStr := tt.page
			if pageStr == "" {
				pageStr = "1" // Default to page 1 if not provided
			}
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil {
				pageInt = 1 // Default to page 1 if parsing fails
			}

			limitStr := tt.limit
			if limitStr == "" {
				limitStr = "10" // Default to page 1 if not provided
			}
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				limitInt = 10 // Default to page 1 if parsing fails
			}

			if tt.isPageNotFound {
				pageInt = 1
			}

			if tt.isFound {
				for i := range pageInt * limitInt {
					novel2 := models.Novel{
						Title:          fmt.Sprintf("Test%d", i),
						Synopsis:       fmt.Sprintf("Test%d", i),
						CoverUrl:       fmt.Sprintf("https://example.com/t%d.jpg", i),
						Language:       "en",
						Status:         "ongoing",
						NovelUpdatesID: fmt.Sprintf("t%d", i),
						Genres:         []models.Genre{{Name: tt.genre}},
					}

					_, err := novelRepo.CreateNovel(novel2)
					if err != nil {
						t.Errorf("Failed to create novel: %v", err)
					}
				}
			}

			if !tt.isDatabaseOnline {
				mockRepo := new(mocks.MockNovelRepository)
				mockRepo.On("IsDown").Return(true)
				mockRepo.On("GetNovelsByGenreName", mock.Anything, mock.Anything, mock.Anything).Return(([]models.Novel)(nil), (int64)(0), types.WrapError(types.DATABASE_ERROR, "Database is offline", nil))
				novelService = services.NewNovelService(mockRepo, scriptExecutor)
			}

			novelController = *controllers.NewNovelController(novelService)

			// Build the query parameters
			queryParams := url.Values{}
			if tt.page != "" {
				queryParams.Add("page", tt.page)
			}
			if tt.limit != "" {
				queryParams.Add("limit", tt.limit)
			}
			if tt.hasTestParam {
				queryParams.Add("test", "test")
			}

			// Construct the URL with query parameters
			urlPath := "/novels/genres/" + tt.genre
			if len(queryParams) > 0 {
				urlPath += "?" + queryParams.Encode()
			}

			// Create a request
			req := httptest.NewRequest(http.MethodGet, urlPath, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/novels/genres/:genre_name", novelController.GetNovelsByGenreName)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusOK {
				// Fetch novels from the database
				novels, total, err := novelRepo.GetNovelsByGenreName(tt.genre, pageInt, limitInt)
				if err != nil {
					t.Errorf("Failed to fetch novels from the database: %v", err)
					return
				}

				// Assert the number of novels returned
				expectedCount := limitInt
				if len(novels) != expectedCount {
					t.Errorf("Expected %d novels, but got %d", expectedCount, len(novels))
				}

				// Assert the total number of novels for the genre
				expectedTotal := pageInt * limitInt
				if total != int64(expectedTotal) {
					t.Errorf("Expected total novels to be %d, but got %d", expectedTotal, total)
				}

				// Optionally, assert the content of the novels
				for i, novel := range novels {
					expectedTitle := fmt.Sprintf("Test%d", (pageInt-1)*limitInt+i)
					if novel.Title != expectedTitle {
						t.Errorf("Expected novel title to be %s, but got %s", expectedTitle, novel.Title)
					}
				}
			}
		})
	}
}

func TestHandleGetNovelsByTag(t *testing.T) {
	tests := []struct {
		name             string
		tag              string
		page             string
		limit            string
		expectedCode     int
		isDatabaseOnline bool
		isFound          bool
		hasTestParam     bool
		isPageNotFound   bool
	}{
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BASE_CASE-SUCCESS",
			tag:              "Action",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_DATABASE_OFFLINE-FAILURE",
			tag:              "Action",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusServiceUnavailable,
			isDatabaseOnline: false,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_NO_NOVELS_BY_TAG-FAILURE",
			tag:              "Action2",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_NO_PARAMS-SUCCESS",
			tag:              "Action",
			page:             "",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_ONLY_PAGE_PARAM-SUCCESS",
			tag:              "Action",
			page:             "2",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_ONLY_LIMIT_PARAM-SUCCESS",
			tag:              "Action",
			page:             "",
			limit:            "11",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_MORE_THAN_ACCEPTED_QUERY_PARAMS-FAILURE",
			tag:              "Action",
			page:             "1",
			limit:            "1",
			expectedCode:     http.StatusBadRequest,
			hasTestParam:     true,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_TAG_IS_A_SPACE-FAILURE",
			tag:              "%20",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_IS_NOT_AN_INT-FAILURE",
			tag:              "Action",
			page:             "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_DOES_NOT_EXIST-FAILURE",
			tag:              "Action",
			page:             "2",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
			isPageNotFound:   true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_LIMIT_IS_NOT_AN_INT-FAILURE",
			tag:              "Action",
			limit:            "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_TAG_IS_MINIMUM_LENGTH-SUCCESS",
			tag:              "a",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_TAG_IS_MAXIMUM_LENGTH-SUCCESS",
			tag:              strings.Repeat("a", 255),
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_TAG_IS_JUST_BELOW_MINIMUM_LENGTH-FAILURE",
			tag:              "",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_TAG_IS_JUST_ABOVE_MAXIMUM_LENGTH-FAILURE",
			tag:              strings.Repeat("a", 256),
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_IS_MINIMUM_VALUE-SUCCESS",
			tag:              "Action",
			page:             "1",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_IS_MAXIMUM_VALUE-SUCCESS",
			tag:              "Action",
			page:             "1000",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			tag:              "Action",
			page:             "0",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_PAGE_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			tag:              "Action",
			page:             "1001",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_LIMIT_IS_MINIMUM_VALUE-SUCCESS",
			tag:              "Action",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_LIMIT_IS_MAXIMUM_VALUE-SUCCESS",
			tag:              "Action",
			limit:            "100",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_LIMIT_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			tag:              "Action",
			limit:            "9",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F5/GET_NOVEL_BY_TAG_BUT_LIMIT_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			tag:              "Action",
			limit:            "101",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			novel := models.Novel{
				Title:          "t",
				Synopsis:       "t",
				CoverUrl:       "https://example.com/t.jpg",
				Language:       "en",
				Status:         "ongoing",
				NovelUpdatesID: "t",
				Tags:           []models.Tag{{Name: "test"}},
			}

			_, err := novelRepo.CreateNovel(novel)
			if err != nil {
				t.Errorf("Failed to create novel2: %v", err)
			}

			// Parse page parameter
			pageStr := tt.page
			if pageStr == "" {
				pageStr = "1" // Default to page 1 if not provided
			}
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil {
				pageInt = 1 // Default to page 1 if parsing fails
			}

			limitStr := tt.limit
			if limitStr == "" {
				limitStr = "10" // Default to page 1 if not provided
			}
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				limitInt = 10 // Default to page 1 if parsing fails
			}

			if tt.isPageNotFound {
				pageInt = 1
			}

			if tt.isFound {
				for i := range pageInt * limitInt {
					novel2 := models.Novel{
						Title:          fmt.Sprintf("Test%d", i),
						Synopsis:       fmt.Sprintf("Test%d", i),
						CoverUrl:       fmt.Sprintf("https://example.com/t%d.jpg", i),
						Language:       "en",
						Status:         "ongoing",
						NovelUpdatesID: fmt.Sprintf("t%d", i),
						Tags:           []models.Tag{{Name: tt.tag}},
					}

					_, err := novelRepo.CreateNovel(novel2)
					if err != nil {
						t.Errorf("Failed to create novel: %v", err)
					}
				}
			}

			if !tt.isDatabaseOnline {
				mockRepo := new(mocks.MockNovelRepository)
				mockRepo.On("IsDown").Return(true)
				mockRepo.On("GetNovelsByTagName", mock.Anything, mock.Anything, mock.Anything).Return(([]models.Novel)(nil), (int64)(0), types.WrapError(types.DATABASE_ERROR, "Database is offline", nil))
				novelService = services.NewNovelService(mockRepo, scriptExecutor)
			}

			novelController = *controllers.NewNovelController(novelService)

			// Build the query parameters
			queryParams := url.Values{}
			if tt.page != "" {
				queryParams.Add("page", tt.page)
			}
			if tt.limit != "" {
				queryParams.Add("limit", tt.limit)
			}
			if tt.hasTestParam {
				queryParams.Add("test", "test")
			}

			// Construct the URL with query parameters
			urlPath := "/novels/tags/" + tt.tag
			if len(queryParams) > 0 {
				urlPath += "?" + queryParams.Encode()
			}

			// Create a request
			req := httptest.NewRequest(http.MethodGet, urlPath, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/novels/tags/:tag_name", novelController.GetNovelsByTagName)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusOK {
				// Fetch novels from the database
				novels, total, err := novelRepo.GetNovelsByTagName(tt.tag, pageInt, limitInt)
				if err != nil {
					t.Errorf("Failed to fetch novels from the database: %v", err)
					return
				}

				// Assert the number of novels returned
				expectedCount := limitInt
				if len(novels) != expectedCount {
					t.Errorf("Expected %d novels, but got %d", expectedCount, len(novels))
				}

				// Assert the total number of novels for the tag
				expectedTotal := pageInt * limitInt
				if total != int64(expectedTotal) {
					t.Errorf("Expected total novels to be %d, but got %d", expectedTotal, total)
				}

				// Optionally, assert the content of the novels
				for i, novel := range novels {
					expectedTitle := fmt.Sprintf("Test%d", (pageInt-1)*limitInt+i)
					if novel.Title != expectedTitle {
						t.Errorf("Expected novel title to be %s, but got %s", expectedTitle, novel.Title)
					}
				}
			}
		})
	}
}

func TestHandleGetNovels(t *testing.T) {
	tests := []struct {
		name             string
		page             string
		limit            string
		expectedCode     int
		isDatabaseOnline bool
		isFound          bool
		hasTestParam     bool
		isPageNotFound   bool
	}{
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BASE_CASE-SUCCESS",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_DATABASE_OFFLINE-FAILURE",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusServiceUnavailable,
			isDatabaseOnline: false,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_NO_NOVELS-FAILURE",
			page:             "1",
			limit:            "10",
			expectedCode:     http.StatusNotFound,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_NO_PARAMS-SUCCESS",
			page:             "",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_ONLY_PAGE_PARAM-SUCCESS",
			page:             "2",
			limit:            "",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_ONLY_LIMIT_PARAM-SUCCESS",
			page:             "",
			limit:            "11",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_MORE_THAN_ACCEPTED_QUERY_PARAMS-FAILURE",
			page:             "1",
			limit:            "1",
			expectedCode:     http.StatusBadRequest,
			hasTestParam:     true,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_IS_NOT_AN_INT-FAILURE",
			page:             "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_DOES_NOT_EXIST-FAILURE",
			page:             "2",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
			isPageNotFound:   true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_LIMIT_IS_NOT_AN_INT-FAILURE",
			limit:            "test",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          false,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_IS_MINIMUM_VALUE-SUCCESS",
			page:             "1",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_IS_MAXIMUM_VALUE-SUCCESS",
			page:             "1000",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			page:             "0",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_PAGE_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			page:             "1001",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_LIMIT_IS_MINIMUM_VALUE-SUCCESS",
			limit:            "10",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_LIMIT_IS_MAXIMUM_VALUE-SUCCESS",
			limit:            "100",
			expectedCode:     http.StatusOK,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_LIMIT_IS_JUST_BELOW_MINIMUM_VALUE-FAILURE",
			limit:            "9",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
		{
			name:             "#TCO-B3/F6/GET_NOVEL_BUT_LIMIT_IS_JUST_ABOVE_MAXIMUM_VALUE-FAILURE",
			limit:            "101",
			expectedCode:     http.StatusBadRequest,
			isDatabaseOnline: true,
			isFound:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(func() {
				cleanDB()
			})

			// Parse page parameter
			pageStr := tt.page
			if pageStr == "" {
				pageStr = "1" // Default to page 1 if not provided
			}
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil {
				pageInt = 1 // Default to page 1 if parsing fails
			}

			limitStr := tt.limit
			if limitStr == "" {
				limitStr = "10" // Default to page 1 if not provided
			}
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				limitInt = 10 // Default to page 1 if parsing fails
			}

			if tt.isPageNotFound {
				pageInt = 1
			}

			if tt.isFound {
				for i := range pageInt * limitInt {
					novel := models.Novel{
						Title:          fmt.Sprintf("Test%d", i),
						Synopsis:       fmt.Sprintf("Test%d", i),
						CoverUrl:       fmt.Sprintf("https://example.com/t%d.jpg", i),
						Language:       "en",
						Status:         "ongoing",
						NovelUpdatesID: fmt.Sprintf("t%d", i),
					}

					_, err := novelRepo.CreateNovel(novel)
					if err != nil {
						t.Errorf("Failed to create novel: %v", err)
					}
				}
			}

			if !tt.isDatabaseOnline {
				mockRepo := new(mocks.MockNovelRepository)
				mockRepo.On("IsDown").Return(true)
				mockRepo.On("GetNovels", mock.Anything, mock.Anything).Return(([]models.Novel)(nil), (int64)(0), types.WrapError(types.DATABASE_ERROR, "Database is offline", nil))
				novelService = services.NewNovelService(mockRepo, scriptExecutor)
			}

			novelController = *controllers.NewNovelController(novelService)

			// Build the query parameters
			queryParams := url.Values{}
			if tt.page != "" {
				queryParams.Add("page", tt.page)
			}
			if tt.limit != "" {
				queryParams.Add("limit", tt.limit)
			}
			if tt.hasTestParam {
				queryParams.Add("test", "test")
			}

			// Construct the URL with query parameters
			urlPath := "/novels/"
			if len(queryParams) > 0 {
				urlPath += "?" + queryParams.Encode()
			}

			// Create a request
			req := httptest.NewRequest(http.MethodGet, urlPath, strings.NewReader(""))
			req.Header.Set("Content-Type", "application/json")

			// Mock context and recorder
			w := httptest.NewRecorder()
			router := gin.Default()
			router.GET("/novels/", novelController.GetNovels)
			router.ServeHTTP(w, req)

			// Assertions
			assert.Equal(t, tt.expectedCode, w.Code)

			if w.Code != http.StatusOK {
				log.Println(w.Body.String())
			}

			// Check if the database reflects changes
			if tt.expectedCode == http.StatusOK {
				// Fetch novels from the database
				novels, total, err := novelRepo.GetNovels(pageInt, limitInt)
				if err != nil {
					t.Errorf("Failed to fetch novels from the database: %v", err)
					return
				}

				// Assert the number of novels returned
				expectedCount := limitInt
				if len(novels) != expectedCount {
					t.Errorf("Expected %d novels, but got %d", expectedCount, len(novels))
				}

				// Assert the total number of novels for the tag
				expectedTotal := pageInt * limitInt
				if total != int64(expectedTotal) {
					t.Errorf("Expected total novels to be %d, but got %d", expectedTotal, total)
				}

				// Optionally, assert the content of the novels
				for i, novel := range novels {
					expectedTitle := fmt.Sprintf("Test%d", (pageInt-1)*limitInt+i)
					if novel.Title != expectedTitle {
						t.Errorf("Expected novel title to be %s, but got %s", expectedTitle, novel.Title)
					}
				}
			}
		})
	}
}
