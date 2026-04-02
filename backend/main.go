package main

import (
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Song struct {
	ID           int     `json:"id"`
	Title        string  `json:"title"`
	Artist       string  `json:"artist"`
	Tag          string  `json:"tag"`
	YoutubeURL   string  `json:"youtubeUrl,omitempty"`
	VideoID      string  `json:"videoId,omitempty"`
	ThumbnailURL string  `json:"thumbnailUrl,omitempty"`
	EmbedURL     string  `json:"embedUrl,omitempty"`
	ViewCount    int64   `json:"viewCount,omitempty"`
	LikeCount    int64   `json:"likeCount,omitempty"`
	CommentCount int64   `json:"commentCount,omitempty"`
	Score        float64 `json:"score,omitempty"`
}

var songs []Song
var nextID = 1

func main() {
	_ = loadEnvFile(".env")

	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType},
	}))

	e.POST("/songs", postSongHandler)
	e.GET("/songs", getSongHandler)
	e.GET("/recommend", recommendHandler)

	e.Logger.Fatal(e.Start(":8080"))
}

func postSongHandler(c echo.Context) error {
	var song Song

	contentType := c.Request().Header.Get(echo.HeaderContentType)
	if !strings.HasPrefix(contentType, echo.MIMEApplicationJSON) {
		return echo.NewHTTPError(http.StatusUnsupportedMediaType, "Content-Type must be application/json")
	}

	decoder := json.NewDecoder(c.Request().Body)
	if err := decoder.Decode(&song); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid JSON body")
	}

	if strings.TrimSpace(song.YoutubeURL) != "" {
		apiKey := os.Getenv("YOUTUBE_API_KEY")
		if strings.TrimSpace(apiKey) == "" {
			return echo.NewHTTPError(http.StatusInternalServerError, "YOUTUBE_API_KEY is not configured")
		}

		videoID, err := extractVideoID(song.YoutubeURL)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		meta, err := fetchYouTubeVideoInfo(videoID, apiKey)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		song.VideoID = videoID
		song.EmbedURL = "https://www.youtube.com/embed/" + videoID
		song.ThumbnailURL = meta.ThumbnailURL
		song.Title = meta.Title
		song.Artist = meta.ChannelTitle
		song.ViewCount = meta.ViewCount
		song.LikeCount = meta.LikeCount
		song.CommentCount = meta.CommentCount
		song.Score = calculateScore(meta.ViewCount, meta.LikeCount, meta.CommentCount)

		if strings.TrimSpace(song.Tag) == "" {
			song.Tag = "youtube"
		}
	} else {
		if strings.TrimSpace(song.Title) == "" || strings.TrimSpace(song.Artist) == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "title and artist are required")
		}
	}

	song.ID = nextID
	nextID++
	songs = append(songs, song)

	return c.JSON(http.StatusCreated, song)
}

func getSongHandler(c echo.Context) error {
	if songs == nil {
		return c.JSON(http.StatusOK, []Song{})
	}

	return c.JSON(http.StatusOK, songs)
}

func recommendHandler(c echo.Context) error {
	targetTag := "夜"

	var result []Song

	for _, s := range songs {
		if s.Tag == targetTag {
			result = append(result, s)
		}
	}

	if result == nil {
		return c.JSON(http.StatusOK, []Song{})
	}

	return c.JSON(http.StatusOK, result)
}

type youtubeVideoInfo struct {
	Title        string
	ChannelTitle string
	ThumbnailURL string
	ViewCount    int64
	LikeCount    int64
	CommentCount int64
}

func extractVideoID(rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return "", fmt.Errorf("invalid youtube url")
	}

	host := strings.TrimPrefix(parsed.Host, "www.")
	switch host {
	case "youtu.be":
		videoID := strings.Trim(parsed.Path, "/")
		if videoID == "" {
			return "", fmt.Errorf("invalid youtube url")
		}
		return videoID, nil
	case "youtube.com", "m.youtube.com":
		if parsed.Path == "/watch" {
			videoID := parsed.Query().Get("v")
			if videoID == "" {
				return "", fmt.Errorf("invalid youtube url")
			}
			return videoID, nil
		}

		if strings.HasPrefix(parsed.Path, "/shorts/") || strings.HasPrefix(parsed.Path, "/embed/") {
			parts := strings.Split(strings.Trim(parsed.Path, "/"), "/")
			if len(parts) >= 2 && parts[1] != "" {
				return parts[1], nil
			}
		}
	}

	return "", fmt.Errorf("unsupported youtube url")
}

func fetchYouTubeVideoInfo(videoID, apiKey string) (youtubeVideoInfo, error) {
	endpoint := fmt.Sprintf("https://www.googleapis.com/youtube/v3/videos?part=snippet,statistics&id=%s&key=%s", url.QueryEscape(videoID), url.QueryEscape(apiKey))
	resp, err := http.Get(endpoint)
	if err != nil {
		return youtubeVideoInfo{}, fmt.Errorf("failed to call youtube api")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return youtubeVideoInfo{}, fmt.Errorf("youtube api error: status %d", resp.StatusCode)
	}

	var payload struct {
		Items []struct {
			Snippet struct {
				Title        string `json:"title"`
				ChannelTitle string `json:"channelTitle"`
				Thumbnails   struct {
					High struct {
						URL string `json:"url"`
					} `json:"high"`
					Medium struct {
						URL string `json:"url"`
					} `json:"medium"`
					Default struct {
						URL string `json:"url"`
					} `json:"default"`
				} `json:"thumbnails"`
			} `json:"snippet"`
			Statistics struct {
				ViewCount    string `json:"viewCount"`
				LikeCount    string `json:"likeCount"`
				CommentCount string `json:"commentCount"`
			} `json:"statistics"`
		} `json:"items"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return youtubeVideoInfo{}, fmt.Errorf("failed to decode youtube api response")
	}

	if len(payload.Items) == 0 {
		return youtubeVideoInfo{}, fmt.Errorf("video not found")
	}

	item := payload.Items[0]
	thumbnailURL := item.Snippet.Thumbnails.High.URL
	if thumbnailURL == "" {
		thumbnailURL = item.Snippet.Thumbnails.Medium.URL
	}
	if thumbnailURL == "" {
		thumbnailURL = item.Snippet.Thumbnails.Default.URL
	}

	return youtubeVideoInfo{
		Title:        item.Snippet.Title,
		ChannelTitle: item.Snippet.ChannelTitle,
		ThumbnailURL: thumbnailURL,
		ViewCount:    parseCount(item.Statistics.ViewCount),
		LikeCount:    parseCount(item.Statistics.LikeCount),
		CommentCount: parseCount(item.Statistics.CommentCount),
	}, nil
}

func parseCount(raw string) int64 {
	if raw == "" {
		return 0
	}

	v, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return 0
	}

	return v
}

func calculateScore(viewCount, likeCount, commentCount int64) float64 {
	if viewCount <= 0 && likeCount <= 0 && commentCount <= 0 {
		return 0
	}

	views := float64(viewCount)
	likes := float64(likeCount)
	comments := float64(commentCount)

	engagementRate := (likes + comments*3.0) / (views + 300.0)
	volumeBonus := math.Log1p(views) * 0.15

	return engagementRate*1000.0 + volumeBonus
}

func loadEnvFile(path string) error {
	content, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if key == "" {
			continue
		}

		if os.Getenv(key) == "" {
			_ = os.Setenv(key, value)
		}
	}

	return nil
}
