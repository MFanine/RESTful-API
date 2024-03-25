package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"series/pkg/models"
	"time"
	// "os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"github.com/dgrijalva/jwt-go"
	// "github.com/joho/godotenv"
)


// var jwtKey = []byte(os.Getenv("SECRET_KEY")) 

var jwtKey = []byte("SECRET_KEY") 

func Register(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var newUser models.User
		err := json.NewDecoder(r.Body).Decode(&newUser)
		if err != nil {
			log.Println("Error decoding user data:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		newUser.Password = string(hashedPass)

		collection := client.Database("test").Collection("users")
		var existingUser models.User
		err = collection.FindOne(context.Background(), bson.M{"$or": []bson.M{{"username": newUser.Username}, {"email": newUser.Email}}}).Decode(&existingUser)
		if err == nil {
			log.Println("User with the same username or email already exists")
			w.WriteHeader(http.StatusConflict)
			return
		}

		_, err = collection.InsertOne(context.Background(), newUser)
		if err != nil {
			log.Println("Error inserting user into database:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func Login(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			log.Println("Error decoding user data:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		var foundUser models.User
		err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&foundUser)
		if err != nil {
			log.Println("User not found:", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(user.Password))
		if err != nil {
			log.Println("Invalid password:", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &jwt.StandardClaims{
			Subject:   foundUser.Username,
			ExpiresAt: expirationTime.Unix(),
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			log.Println("Error generating token:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
		})

		w.WriteHeader(http.StatusOK)
	}
}
