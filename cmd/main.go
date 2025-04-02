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

type Inventory struct {
	UserID   int `json:"user_id"`
	ItemID   int `json:"item_id"`
	Quantity int `json:"quantity"`
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
	var id int
	err = db.QueryRow("SELECT id, password_hash FROM users WHERE username = $1", input.Username).Scan(&id, &hashedPassword) // если пользователя нет, то ID будет 0 - это нужно поправить
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

	exp := time.Now().Add(time.Minute * 5).Unix()

	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = exp
	claims["username"] = input.Username
	claims["sub"] = id

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
			if claims, ok := parsedToken.Claims.(jwt.MapClaims); ok && parsedToken.Valid {
				c.Set("userID", claims["sub"])
				c.Set("username", claims["username"])
				c.Next()
			} else {
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			}

		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "not authorized"})
			return
		}
	}
}

func GetItems(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	fmt.Printf("User %v is accessing items\n", userID)

	var ItemList []Inventory = []Inventory{}

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM inventory WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for rows.Next() {
		i := Inventory{}
		err = rows.Scan(&i.UserID, &i.ItemID, &i.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ItemList = append(ItemList, i)
	}

	c.JSON(http.StatusOK, gin.H{"items": ItemList})
}
