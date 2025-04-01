package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"

	_ "github.com/lib/pq"
)

var mySignKey = []byte("secret")

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Items struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

func main() {
	r := gin.Default()

	r.POST("/auth", Auth)
	protected := r.Group("/api")
	protected.Use(JWTAuthMiddleware())
	{
		protected.GET("/items", GetItems)
	}

	r.Run(":8080")
}

func Auth(c *gin.Context) {
	var input User

	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()
	inputhashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)

	var hashedPassword string
	err = db.QueryRow("SELECT password_hash FROM users WHERE username = $1", input.Username).Scan(&hashedPassword)
	if err == sql.ErrNoRows {
		_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", input.Username, string(inputhashedPassword))
		hashedPassword = string(inputhashedPassword)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Minute * 5).Unix()

	tokenString, err := token.SignedString(mySignKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Header["Token"] != nil {
			parsedToken, err := jwt.Parse(c.Request.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Error")
				}

				return mySignKey, nil
			})
			if err != nil {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
				return
			}

			if parsedToken.Valid {
				fmt.Println("Yes")
			}
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
			return
		}
	}
}

func GetItems(c *gin.Context) {
	var ItemList []Items = []Items{}

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM shop_items")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for rows.Next() {
		i := Items{}
		err = rows.Scan(&i.Id, &i.Name, &i.Price)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ItemList = append(ItemList, i)
	}

	c.JSON(http.StatusOK, gin.H{"items": ItemList})
}
