package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

)

// album represents data about a record album.
type album struct {
	ID     	string  `json:"id"`
	Title  	string  `json:"title"`
	Artist 	string  `json:"artist"`
	Price  	float64 `json:"price"`
}

type Proxy struct {
	IPPrefix 				string `json:"ip_prefix"`
	Region					string `json:"region"`
	Service					string `json:"service"`
	BorderGroup			string `json:"network_border_group"`
	
}


type ProxyList struct {
	Proxies []Proxy `json:"prefixes"`
}

// albums slice to see record album data.
var albums = []album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99},
	{ID: "4", Title: "Joshua Tree", Artist: "U2", Price: 25.99},
}

func main() {
	
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}
	
	router := gin.Default()
	
	trustedProxies, err := getProxies()
	if err != nil {
		log.Println("Error getting proxies")
	}
	
	setProxyError := router.SetTrustedProxies(trustedProxies)
	if setProxyError != nil {
		log.Fatal("Error setting trusted proxies:", setProxyError)
	} else {
		log.Println("Trusted proxies set: ", len(trustedProxies))
	}
	
	
	// Serve static files.
	//router.StaticFile("/", "./static/index.html")
	
	// Routes.
	router.GET("/", getRoot)
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", postAlbums)
	
	// Check if we are in a development or production environment.
	env := os.Getenv("ENV")
	port := os.Getenv("PORT")
	log.Println("ENV: " + env)
	log.Println("PORT: " + port)
	
	// Use PORT environment variable or default to 5000 if not set (for production)
	if port == "" {
		port = "5000" 
	}
	
	if env == "" {
		env = "production"
	}
	
	router.Run(":" + port)
	
	
}

// getRoot returns the root path.
func getRoot(context *gin.Context) {
	path := "This is the root path"
	
	context.IndentedJSON(http.StatusOK, gin.H{"message": path})
	
	
}

// getAlbums responds with the list of all albums.
func getAlbums(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, albums)
}

// postAlbums adds an album from JSON to newAlbum.
func postAlbums(c *gin.Context) {
	
	var newAlbum album
	
	// Call BindJSON to bind the received JSON to a newAlbum.
	if err := c.BindJSON(&newAlbum); err != nil {
		return 
	}
	
	// Add the new album to the slice.
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

var myClient = &http.Client{Timeout: 10 * time.Second}

func getProxies() ([]string, error) {
	
	var proxyList ProxyList
	var result []string
	
	url := "https://ip-ranges.amazonaws.com/ip-ranges.json"
	
	response, err := myClient.Get(url)
	if err != nil {
		return result,err
	}
	
	defer response.Body.Close()
	
	json.NewDecoder(response.Body).Decode(&proxyList)

	filterRegion := "us-east-1"
	
	for _, proxy := range proxyList.Proxies {
		if proxy.Region == filterRegion {
			result = append(result, proxy.IPPrefix)

		}
	}
	return result, nil
	
}

