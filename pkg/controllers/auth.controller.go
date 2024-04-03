package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"series/pkg/models"
	"series/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

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

		generatedToken, _ := utils.GenerateVerificationToken()

		_, err = collection.UpdateOne(context.Background(), bson.M{"_id": newUser.ID}, bson.M{"$set": bson.M{"verificationToken": generatedToken}})
		if err != nil {
			log.Println("Error saving verification token:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = utils.SendEmail(newUser.Email, generatedToken)
		if err != nil {
			log.Println("Error sending verification email:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func VerifyEmail(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			log.Println("No token provided")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		collection := client.Database("test").Collection("users")
		var user models.User
		err := collection.FindOneAndUpdate(context.Background(), bson.M{"verificationToken": token}, bson.M{"$set": bson.M{"emailVerified": true}}).Decode(&user)
		if err != nil {
			log.Println("Invalid token:", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}


func Login(client *mongo.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("Logged in successfully"))
	}
}

func Logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out successfully"))
	}
}
