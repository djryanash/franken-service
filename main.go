package main

import (
	"net/http"
	"os"
	"github.com/gin-gonic/gin"
)

// album represents data about a record album.
type album struct {
	ID     	string  `json:"id"`
	Title  	string  `json:"title"`
	Artist 	string  `json:"artist"`
	Price  	float64 `json:"price"`
}

// albums slice to see record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
  {ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
  {ID: "4", Title: "Joshua Tree", Artist: "U2", Price: 25.99},
}

func main() {
	router := gin.Default()
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)

	// Check if we are in a development or production environment
	env := os.Getenv("ENV")
	if env == "production" {
			port := os.Getenv("PORT")
			if port == "" {
					port = "5000" // Default to 5000 on Elastic Beanstalk
			}
			router.Run(":" + port)
	} else {
			// Local environment
			router.Run("localhost:8080")
	}
}

// getAlbums responds with the list of all albums.
func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON to newAlbum
func postAlbums(c *gin.Context) {
	
	var newAlbum album

	// Call BindJSON to bind the received JSON to a newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return 
	}

	// Add the new album to the slice 
	albums = append(albums, newAlbum)
	c.IndentedJSON(http.StatusCreated, newAlbum)

}

// getAlbumByID locates the album whose ID value matches the ID
// parameter sent by the client, then returns that album as a response.
func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	// Loop over the list of albums looking for
	// an album whos ID value matches the parameter.

	for _, a := range albums {
		if a.ID == id {
			c.IndentedJSON(http.StatusOK, a)
			return 
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})

}