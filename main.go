package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/AnotherCoolDude/golang-angular/handlers"
	"github.com/auth0-community/go-auth0"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"
)

var (
	audience string
	domain   string
	token    string
)

func main() {
	setAuth0Variables()
	r := gin.Default()

	r.NoRoute(func(c *gin.Context) {
		dir, file := path.Split(c.Request.RequestURI)
		ext := path.Ext(file)
		if file == "" || ext == "" {
			c.File("./ui/dist/ui/index.html")
		} else {
			c.File("./ui/dist/ui/" + path.Join(dir, file))
		}

	})

	authorized := r.Group("/")
	authorized.Use(authRequired())

	authorized.GET("/todo", handlers.GetTodoListHandler)
	authorized.POST("/todo", handlers.AddTodoHandler)
	authorized.DELETE("/todo/:id", handlers.DeleteTodoHandler)
	authorized.PUT("/todo", handlers.CompleteTodoHandler)

	err := r.Run(":3000")
	if err != nil {
		panic(err)
	}

	todos, err := getTodos(token)
	if err != nil {
		log.Println("error getting todos")
		log.Println(err)
	}
	fmt.Println(todos)
}

func setAuth0Variables() {
	audience = "https://my-golang-api"
	domain = "dev-3tt1ae45.auth0.com"
	t, err := receiveToken()
	if err != nil {
		fmt.Println(err)
	}
	token = t
}

func authRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Header.Set("Content-Type", "application/json")
		var auth0Domain = "https://" + domain + "/"
		client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: auth0Domain + ".well-known/jwks.json"}, nil)
		config := auth0.NewConfiguration(client, []string{audience}, auth0Domain, jose.RS256)
		validator := auth0.NewValidator(config, nil)

		_, err := validator.ValidateRequest(c.Request)

		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}

		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}

func receiveToken() (string, error) {
	url := "https://dev-3tt1ae45.eu.auth0.com/oauth/token"

	payload := strings.NewReader("{\"client_id\":\"Vi72HiE8p3clPpdji4L64i0axSba3q1u\",\"client_secret\":\"ONmN5BSndMCyPCHqnswN9d5ycoenQ72-XHmmgqr4WfqzqVAPbwCpaAU-ilTKNWq5\",\"audience\":\"https://my-golang-api\",\"grant_type\":\"client_credentials\"}")

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		log.Println("error creating request")
		return "", err
	}

	req.Header.Add("content-type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request")
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading response")
		return "", err
	}

	fmt.Println(res)
	fmt.Println()
	fmt.Println(string(body))
	return string(body), nil
}

func getTodos(bearerToken string) (string, error) {
	url := "localhost:3000/todo"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("error creating request")
		return "", err
	}

	req.Header.Add("authorization", "Bearer "+bearerToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("error sending request")
		return "", err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error reading response")
		return "", err
	}

	fmt.Println(res)
	fmt.Println(string(body))
	return string(body), nil
}
