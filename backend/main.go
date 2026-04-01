package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Song struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist string `json:"artist"`
	Tag    string `json:"tag"`
}

var songs []Song
var nextID = 1

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderContentType},
	}))

	e.POST("/songs", postSongHandler)
	e.GET("/songs", getSongHandler)

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

	if strings.TrimSpace(song.Title) == "" || strings.TrimSpace(song.Artist) == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "title and artist are required")
	}

	song.ID = nextID
	nextID++
	songs = append(songs, song)

	return c.JSON(http.StatusCreated, song)
}

func getSongHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, songs)
}
