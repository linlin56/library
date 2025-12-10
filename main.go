package main

import (
	"bibliotheque/constants"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PageData struct {
	Title   string
	Content string
}

func readHTMLContent(filename string) (string, error) {
	// Lire le fichier HTML
	content, err := ioutil.ReadFile(constants.ContentDir + filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
func main() {
	// Create a Gin router with default middleware (logger and recovery)
	r := gin.Default()

	// Set custom template function
	r.SetFuncMap(template.FuncMap{
		"safe": func(s string) template.HTML {
			return template.HTML(s)
		},
	})

	// Load HTML templates from the templates directory
	r.LoadHTMLGlob("templates/**/*.html")

	r.Static("/static", "./static")

	// Accueil route
	r.GET("/", func(c *gin.Context) {
		content, err := readHTMLContent("accueil.html")

		if err != nil {
			log.Printf("Couldn't read markdown : %v", err)
			c.String(http.StatusInternalServerError, "Couldn't read content")
			return
		}

		data := PageData{
			Title:   "Accueil - Bibliothèque",
			Content: content,
		}
		log.Printf("DATA: Title='%s', Content='%s'", data.Title, data.Content)
		c.HTML(http.StatusOK, "accueil", data)
	})

	// Start server on port 8080 (default)
	// Server will listen on 0.0.0.0:8080 (localhost:8080 on Windows)
	r.Run()
}
