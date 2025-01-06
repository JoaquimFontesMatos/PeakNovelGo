package controllers

import (
	"archive/zip"
	"backend/internal/models"
	"backend/internal/repositories"
	"backend/internal/utils"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NovelController struct {
	novelRepository *repositories.NovelRepository
}

func NewNovelController(novelRepository *repositories.NovelRepository) *NovelController {
	return &NovelController{novelRepository: novelRepository}
}

type Session struct {
	UserInput        string `json:"user_input,omitempty"`
	OutputPath       string `json:"output_path,omitempty"`
	Completed        bool   `json:"completed,omitempty"`
	DownloadChapters []int  `json:"download_chapters,omitempty"`
}

type Metadata struct {
	Novel   models.ImportedNovel `json:"novel,omitempty"`
	Session Session              `json:"session,omitempty"`
}

// HandleGetNovels handles POST /novel
func (n *NovelController) HandleImportNovel(ctx *gin.Context) {

	metadata := Metadata{}
	bodyCopy := new(bytes.Buffer)

	// Read the whole body
	_, err := io.Copy(bodyCopy, ctx.Request.Body)
	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Error reading API token"})
		return
	}
	bodyData := bodyCopy.Bytes()

	// Replace the body with a reader that reads from the buffer
	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyData))

	err = ctx.Bind(&metadata)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Print the parsed metadata
	fmt.Println("Novel Title:", metadata.Novel.Title)
	fmt.Println("Author(s):", metadata.Novel.Authors)
	fmt.Println("Cover URL:", metadata.Novel.CoverUrl)
	fmt.Println("Synopsis:", metadata.Novel.Synopsis)
	fmt.Println("Language:", metadata.Novel.Language)
	fmt.Println("Status:", metadata.Novel.Status)
	fmt.Println("Genres:", metadata.Novel.Genres)
	fmt.Println("NovelUpdatesURL:", metadata.Novel.NovelUpdatesUrl)
	fmt.Println("Tags:", metadata.Novel.Tags)

	// Print chapters and volumes
	for _, volume := range metadata.Novel.Volumes {
		fmt.Printf("Volume %d: %s (Chapters: %d-%d)\n", volume.ID, volume.Title, volume.StartChapter, volume.EndChapter)
	}

	metadata.Novel.Synopsis = utils.StripHTML(metadata.Novel.Synopsis)

	// Convert the imported novel to a Novel struct
	novel := models.ConvertToNovel(metadata.Novel)

	// Save the novel to the database
	createdNovel, err := n.novelRepository.CreateNovel(&novel)

	if err != nil {
		log.Println(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to save novel to database"})
		return
	}

	log.Println("Novel saved successfully")

	ctx.JSON(200, createdNovel)

	ctx.JSON(200, novel)
}

// HandleUploadNovelZip handles POST /novel/upload
func (c *NovelController) HandleImportChaptersZip(ctx *gin.Context) {
	idParam := ctx.Param("novel_id")
	id, err := strconv.Atoi(idParam)

	if err != nil || id <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	// Convert the id from int to uint (assuming id can be positive)
	uid := uint(id)

	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File not provided"})
		return
	}

	// Save the uploaded file to a temporary location
	tempFile := "./temp_upload.zip"
	if err := ctx.SaveUploadedFile(file, tempFile); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Process the ZIP file
	err = processChaptersZip(tempFile, uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return the extracted metadata
	ctx.JSON(http.StatusOK, gin.H{"message": "Chapters extracted successfully"})
}

func processChaptersZip(filePath string, uid uint) error {
	// Open the ZIP file
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %v", err)
	}
	defer reader.Close()

	// Map to store volumes by name
	var chapters []models.Chapter

	// Iterate through the files in the ZIP
	for _, file := range reader.File {
		if file.FileInfo().IsDir() {
			continue
		}

		// Open the file within the ZIP
		f, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %v", err)
		}
		defer f.Close()

		// Read the file content
		var chapterData struct {
			Title      string `json:"title"`
			Body       string `json:"body"`
			ChapterUrl string `json:"url"`
			VolumeID   uint   `json:"volume"`
		}
		if err := json.NewDecoder(f).Decode(&chapterData); err != nil {
			return fmt.Errorf("failed to decode JSON file: %v", err)
		}

		// Clean up the chapter body
		chapterData.Body = utils.StripHTML(chapterData.Body)

		// Create the chapter model
		chapter := models.Chapter{
			VolumeID:   &chapterData.VolumeID,
			NovelID:    &uid,
			Title:      chapterData.Title,
			Body:       chapterData.Body,
			ChapterUrl: chapterData.ChapterUrl,
		}

		// Append the chapter to the list
		chapters = append(chapters, chapter)
	}

	return nil
}
