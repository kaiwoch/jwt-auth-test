package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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

type Send struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type Inventory struct {
	ItemName string `json:"name"`
	Quantity int    `json:"quantity"`
}

type Received struct {
	FromUser string `json:"from"`
	Amount   int    `json:"amount"`
}

type Sent struct {
	ToUser string `json:"to"`
	Amount int    `json:"amount"`
}

type CoinHistory struct {
	Received []Received `json:"received"`
	Sent     []Sent     `json:"sent"`
}

type Info struct {
	Balance     int         `json:"balance"`
	Inventory   []Inventory `json:"inventory"`
	CoinHistory CoinHistory `json:"coinHistory"`
}

func main() {
	r := gin.Default()

	r.POST("/auth", Auth)
	protected := r.Group("/api")
	protected.Use(JWTAuthMiddleware())
	{
		protected.GET("/info", GetInfo)
		//protected.GET("/buy/{}", BuyItem)
		protected.POST("/sendCoin", SendCoin)
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
		// решил добавить бонус при регистрации
		db.QueryRow("SELECT id FROM users WHERE username = $1", input.Username).Scan(&id)
		_, err = db.Exec("INSERT INTO user_wallet (user_id, coin_balance) VALUES ($1, $2)", id, 100)
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

func GetInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	fmt.Printf("User %v is accessing info\n", userID)

	var info Info = Info{Inventory: make([]Inventory, 0), CoinHistory: CoinHistory{Received: make([]Received, 0), Sent: make([]Sent, 0)}}

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	db.QueryRow("SELECT coin_balance FROM user_wallet WHERE user_id = $1", userID).Scan(&info.Balance)

	rows, err := db.Query("SELECT shop_items.name, inventory.quantity FROM inventory JOIN shop_items ON shop_items.id = inventory.item_id WHERE user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for rows.Next() {
		i := Inventory{}
		err = rows.Scan(&i.ItemName, &i.Quantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		info.Inventory = append(info.Inventory, i)
	}

	rows, err = db.Query("SELECT users.username, transactions.amount FROM transactions INNER JOIN users ON users.id = transactions.from_user_id WHERE transactions.to_user_id = $1", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for rows.Next() {
		h := Received{}
		err = rows.Scan(&h.FromUser, &h.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		info.CoinHistory.Received = append(info.CoinHistory.Received, h)
	}

	rows, err = db.Query("SELECT users.username, transactions.amount FROM transactions JOIN users ON users.id = transactions.to_user_id WHERE transactions.from_user_id = $1;", userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	for rows.Next() {
		h := Sent{}
		err = rows.Scan(&h.ToUser, &h.Amount)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		info.CoinHistory.Sent = append(info.CoinHistory.Sent, h)
	}

	c.JSON(http.StatusOK, gin.H{"response": info})
}

func SendCoin(c *gin.Context) {
	var input Send
	var id int
	var balance, toBalance int
	err := json.NewDecoder(c.Request.Body).Decode(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User information not found"})
		return
	}

	fmt.Printf("User %v is accessing info\n", userID)

	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost/coins?sslmode=disable")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	defer db.Close()

	db.QueryRow("SELECT coin_balance FROM user_wallet WHERE user_id = $1", userID).Scan(&balance)
	if balance < input.Amount {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not enough money"})
		return
	} else {
		db.QueryRow("SELECT users.id, user_wallet.coin_balance FROM users INNER JOIN user_wallet ON users.id = user_wallet.user_id WHERE username = $1", input.ToUser).Scan(&id, &toBalance)
		fmt.Println(id)
		if id != 0 && id != userID {
			_, err := db.Exec("UPDATE user_wallet SET coin_balance = $1 WHERE user_id = $2", balance-input.Amount, userID)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = db.Exec("UPDATE user_wallet SET coin_balance = $1 WHERE user_id = $2", toBalance+input.Amount, id)
			if err != nil {
				log.Println(err)
				return
			}
			_, err = db.Exec("INSERT INTO transactions (from_user_id, to_user_id, amount) VALUES ($1, $2, $3)", userID, id, input.Amount)
			if err != nil {
				log.Println(err)
				return
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User with this id does not exist"})
			return
		}
	}

}
