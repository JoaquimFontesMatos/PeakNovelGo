package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Chapter struct {
	ID          int    `json:"id"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Volume      int    `json:"volume"`
	VolumeTitle string `json:"volume_title"`
}

type Volume struct {
	ID            int    `json:"id"`
	Title         string `json:"title"`
	StartChapter  int    `json:"start_chapter"`
	FinalChapter  int    `json:"final_chapter"`
	ChapterCount  int    `json:"chapter_count"`
}

type Novel struct {
	URL           string     `json:"url"`
	Title         string     `json:"title"`
	Authors       []string   `json:"authors"`
	CoverURL      string     `json:"cover_url"`
	Chapters      []Chapter  `json:"chapters"`
	Volumes       []Volume   `json:"volumes"`
	IsRTL         bool       `json:"is_rtl"`
	Synopsis      string     `json:"synopsis"`
	Language      string     `json:"language"`
	Tags          []string   `json:"novel_tags"`
	Status        string     `json:"status"`
	Genres        []string   `json:"genres"`
	NovelUpdatesURL string   `json:"novelupdates_url"`
}

type Session struct {
	UserInput   string   `json:"user_input"`
	OutputPath  string   `json:"output_path"`
	Completed   bool     `json:"completed"`
	DownloadChapters []int `json:"download_chapters"`
}

type Metadata struct {
	Novel   Novel   `json:"novel"`
	Session Session `json:"session"`
}

func main() {
	// Assuming this is the JSON response you received
	data := `{
		"novel": {
			"url": "http://novelhall.com/Reverend-Insanity-179",
			"title": "Reverend Insanity",
			"authors": ["Gu Zhen Ren"],
			"cover_url": "https://www.novelhall.com/upload/images/article/20190428/Reverend-Insanity.jpg",
			"chapters": [
				{
					"id": 1,
					"url": "http://novelhall.com/Reverend-Insanity-179/54426.html",
					"title": "1 The Heart Of A Demon Never Has Regret Even In Death",
					"volume": 1,
					"volume_title": "Volume 1"
				},
				{
					"id": 2,
					"url": "http://novelhall.com/Reverend-Insanity-179/54429.html",
					"title": "2 Going Back In Time With 500 Years Of Knowledge",
					"volume": 1,
					"volume_title": "Volume 1"
				}
			],
			"volumes": [
				{
					"id": 1,
					"title": "Volume 1",
					"start_chapter": 1,
					"final_chapter": 100,
					"chapter_count": 100
				},
				{
					"id": 2,
					"title": "Volume 2",
					"start_chapter": 101,
					"final_chapter": 200,
					"chapter_count": 100
				}
			],
			"is_rtl": false,
			"synopsis": "<p>Humans are clever in tens of thousands of ways...</p>",
			"language": "en",
			"novel_tags": ["Xianxia"],
			"status": "Unknown",
			"genres": [],
			"novelupdates_url": null
		},
		"session": {
			"user_input": "reverend insanity",
			"output_path": "C:\\Users\\QuimQuimOTerceiro\\Downloads\\Lightnovels\\novelhall-com\\Reverend Insanity",
			"completed": true,
			"download_chapters": [1, 2],
			"good_file_name": "Reverend Insanity"
		}
	}`

	// Unmarshal JSON into Go struct
	var metadata Metadata
	err := json.Unmarshal([]byte(data), &metadata)
	if err != nil {
		log.Fatal(err)
	}

	// Print the parsed metadata
	fmt.Println("Novel Title:", metadata.Novel.Title)
	fmt.Println("Author(s):", metadata.Novel.Authors)
	fmt.Println("Cover URL:", metadata.Novel.CoverURL)
	fmt.Println("Synopsis:", metadata.Novel.Synopsis)

	// Print chapters and volumes
	for _, volume := range metadata.Novel.Volumes {
		fmt.Printf("Volume %d: %s (Chapters: %d-%d)\n", volume.ID, volume.Title, volume.StartChapter, volume.FinalChapter)
	}
	for _, chapter := range metadata.Novel.Chapters {
		fmt.Printf("Chapter %d: %s (Volume %d)\n", chapter.ID, chapter.Title, chapter.Volume)
	}
}
