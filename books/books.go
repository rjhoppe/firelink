package books

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rjhoppe/firelink/utils"
)

type CheckForBookAPI struct {
	Count    int `json:"count"`
	Next     any `json:"next"`
	Previous any `json:"previous"`
	Results  []struct {
		ID      int    `json:"id"`
		Title   string `json:"title"`
		Authors []struct {
			Name      string `json:"name"`
			BirthYear int    `json:"birth_year"`
			DeathYear int    `json:"death_year"`
		} `json:"authors"`
		Translators []any    `json:"translators"`
		Subjects    []string `json:"subjects"`
		Bookshelves []string `json:"bookshelves"`
		Languages   []string `json:"languages"`
		Copyright   bool     `json:"copyright"`
		MediaType   string   `json:"media_type"`
		Formats     struct {
			TextHTML                    string `json:"text/html"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
		DownloadCount int `json:"download_count"`
		Formats0      struct {
			TextHTML                    string `json:"text/html"`
			TextHTMLCharsetUtf8         string `json:"text/html; charset=utf-8"`
			ApplicationEpubZip          string `json:"application/epub+zip"`
			ApplicationXMobipocketEbook string `json:"application/x-mobipocket-ebook"`
			TextPlainCharsetUtf8        string `json:"text/plain; charset=utf-8"`
			ApplicationRdfXML           string `json:"application/rdf+xml"`
			ImageJpeg                   string `json:"image/jpeg"`
			ApplicationOctetStream      string `json:"application/octet-stream"`
			TextPlainCharsetUsASCII     string `json:"text/plain; charset=us-ascii"`
		} `json:"formats,omitempty"`
		Formats1 struct {
			TextHTML                string `json:"text/html"`
			AudioOgg                string `json:"audio/ogg"`
			AudioMp4                string `json:"audio/mp4"`
			AudioMpeg               string `json:"audio/mpeg"`
			TextPlainCharsetUsASCII string `json:"text/plain; charset=us-ascii"`
			ApplicationRdfXML       string `json:"application/rdf+xml"`
		} `json:"formats,omitempty"`
		Formats2 struct {
			TextPlainCharsetUsASCII string `json:"text/plain; charset=us-ascii"`
			TextHTMLCharsetUsASCII  string `json:"text/html; charset=us-ascii"`
			AudioMpeg               string `json:"audio/mpeg"`
			ApplicationRdfXML       string `json:"application/rdf+xml"`
			ApplicationOctetStream  string `json:"application/octet-stream"`
		} `json:"formats,omitempty"`
	} `json:"results"`
}

func CheckForBook(c *gin.Context, title string) {
	var book CheckForBookAPI

	if title == "help" {
		c.JSON(http.StatusOK, gin.H{"body": "To see if an ebook is available at Project Gutenberg, send a book title to the /ebook/find/{title} endpoint"})
	}

	url := "https://gutendex.com/books/?search=" + title
	resp, err := http.Get(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"body": "Error retrieving data from source api"})
	}
	if err := json.NewDecoder(resp.Body).Decode(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"body": "Error decoding api response body"})
	}
	if book.Count > 0 {
		if utils.ContainsString(book.Results[0].Languages, "en") {
			if book.Results[0].Formats.ApplicationEpubZip != "" {
				c.JSON(http.StatusOK, gin.H{"body": "Book found in epub format"})
			} else {
				c.JSON(http.StatusOK, gin.H{"body": "Book found, but not in epub format"})
			}
		}
	} else {
		c.JSON(http.StatusOK, gin.H{"body": "Book not found"})
	}
}

func GetBook(c *gin.Context, title string) {
	fmt.Println(title)
	// gets the book from Project Gutenberg
	return
}
