package repositories

import (
	"backend/internal/models"
	"errors"
	"log"

	"gorm.io/gorm"
)

type NovelRepository struct {
	db *gorm.DB
}

func NewNovelRepository(db *gorm.DB) *NovelRepository {
	return &NovelRepository{db: db}
}

func (n *NovelRepository) CreateAuthor(author *models.Author) (*models.Author, error) {
	if IsAuthorCreated, err := n.IsAuthorCreated(author); err != nil || IsAuthorCreated {
		return nil, errors.New("author already exists")
	}

	// Save the author
	if err := n.db.Create(&author).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create author")
	}

	return author, nil
}

func (n *NovelRepository) IsAuthorCreated(author *models.Author) (bool, error) {
	var existingAuthor models.Author
	if err := n.db.Where("name = ?", author.Name).First(&existingAuthor).Error; err != nil {
		log.Println(err)
		return false, errors.New("failed to check if author already exists")
	}
	return existingAuthor.ID != 0, nil
}

func (n *NovelRepository) CreateGenre(genre *models.Genre) (*models.Genre, error) {
	if IsGenreCreated, err := n.IsGenreCreated(genre); err != nil || IsGenreCreated {
		return nil, errors.New("genre already exists")
	}

	// Save the genre
	if err := n.db.Create(&genre).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create genre")
	}

	return genre, nil
}

func (n *NovelRepository) IsGenreCreated(genre *models.Genre) (bool, error) {
	var existingGenre models.Genre
	if err := n.db.Where("name = ?", genre.Name).First(&existingGenre).Error; err != nil {
		log.Println(err)
		return false, errors.New("failed to check if genre already exists")
	}
	return existingGenre.ID != 0, nil
}

func (n *NovelRepository) CreateTag(tag *models.Tag) (*models.Tag, error) {
	if IsTagCreated, err := n.IsTagCreated(tag); err != nil || IsTagCreated {
		return nil, errors.New("tag already exists")
	}

	// Save the tag
	if err := n.db.Create(&tag).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create tag")
	}

	return tag, nil
}

func (n *NovelRepository) IsTagCreated(tag *models.Tag) (bool, error) {
	var existingTag models.Tag
	if err := n.db.Where("name = ?", tag.Name).First(&existingTag).Error; err != nil {
		log.Println(err)
		return false, errors.New("failed to check if tag already exists")
	}
	return existingTag.ID != 0, nil
}

func (n *NovelRepository) CreateVolume(volume *models.Volume) (*models.Volume, error) {
	if IsVolumeCreated, err := n.IsVolumeCreated(volume); err != nil || IsVolumeCreated {
		return nil, errors.New("volume already exists")
	}

	// Save the volume with relationships
	if err := n.db.Create(&volume).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create volume")
	}

	return volume, nil
}

func (n *NovelRepository) IsVolumeCreated(volume *models.Volume) (bool, error) {
	var existingVolume models.Volume
	if err := n.db.Where("title = ?", volume.Title).First(&existingVolume).Error; err != nil {
		log.Println(err)
		return false, errors.New("failed to check if volume already exists")
	}
	return existingVolume.ID != 0, nil
}

func (n *NovelRepository) CreateNovel(novel *models.Novel) (*models.Novel, error) {
	if IsNovelCreated := n.IsNovelCreated(novel); IsNovelCreated {
		return nil, errors.New("novel already exists")

	}

	// Save the novel with relationships
	if err := n.db.Create(&novel).Error; err != nil {
		log.Println(err)
		return nil, errors.New("failed to create novel")
	}

	return novel, nil
}

func (n *NovelRepository) IsNovelCreated(novel *models.Novel) bool {
	var existingNovel models.Novel
	if err := n.db.Where("url = ?", novel.Url).First(&existingNovel).Error; err != nil {
		return false
	}
	return existingNovel.ID != 0
}

func (n *NovelRepository) CreateChapters(chapters []models.Chapter) error {
	// Filter out already existing chapters
	newChapters := []models.Chapter{}
	for _, chapter := range chapters {
		if !n.IsChapterCreated(&chapter) {
			newChapters = append(newChapters, chapter)
		}
	}

	// Return early if no new chapters to save
	if len(newChapters) == 0 {
		return nil
	}

	// Use GORM's Create method for batch insertion
	if err := n.db.Create(&newChapters).Error; err != nil {
		log.Println(err)
		return errors.New("failed to batch save chapters")
	}

	return nil
}

func (n *NovelRepository) IsChapterCreated(chapter *models.Chapter) bool {
	var existingChapter models.Chapter
	if err := n.db.Where("chapter_url = ?", chapter.ChapterUrl).First(&existingChapter).Error; err != nil {
		return false
	}
	return existingChapter.ID != 0
}
