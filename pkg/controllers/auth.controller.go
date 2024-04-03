package controllers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	// "go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"

	// "series/pkg/mail"
	"series/pkg/models"
	"series/pkg/utils"
)

// EmailSender interface for sending emails
type EmailSender interface {
	SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachFiles []string) error
}

// Register handles user registration
func Register(client *mongo.Client, emailSender EmailSender) http.HandlerFunc {
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

		// Send verification email
		subject := "Welcome to Our Platform! Please Verify Your Email"
		content := "Hi " + newUser.Username + ",\n\n" +
			"Thank you for registering with us!\n\n" +
			"Please click the following link to verify your email address:\n" +
			"http://yourdomain.com/verify?token=" + generatedToken + "\n\n" +
			"Regards,\n" +
			"Your Platform Team"

		err = emailSender.SendEmail(subject, content, []string{newUser.Email}, nil, nil, nil)
		if err != nil {
			log.Println("Error sending verification email:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

// func VerifyEmail(client *mongo.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		token := r.URL.Query().Get("token")
// 		if token == "" {
// 			log.Println("No token provided")
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		collection := client.Database("test").Collection("users")
// 		var user models.User
// 		err := collection.FindOneAndUpdate(context.Background(), bson.M{"verificationToken": token}, bson.M{"$set": bson.M{"emailVerified": true}}).Decode(&user)
// 		if err != nil {
// 			log.Println("Invalid token:", err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		http.Redirect(w, r, "/login", http.StatusSeeOther)
// 	}
// }

//// / login
// func Login(client *mongo.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		defer r.Body.Close()

// 		w.WriteHeader(http.StatusCreated)
// 		w.Write([]byte("Logged in successfully"))
// 	}
// }
// /// logout
// func Logout() http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		w.WriteHeader(http.StatusOK)
// 		w.Write([]byte("Logged out successfully"))
// 	}
// }
// // Forget Password

// func ForgetPassword(client *mongo.Client) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		var user models.User
// 		err := json.NewDecoder(r.Body).Decode(&user)
// 		if err != nil {
// 			log.Println("Error decoding user data:", err)
// 			w.WriteHeader(http.StatusBadRequest)
// 			return
// 		}

// 		collection := client.Database("test").Collection("users")
// 		var foundUser models.User
// 		err = collection.FindOne(context.Background(), bson.M{"username": user.Username}).Decode(&foundUser)
// 		if err != nil {
// 			log.Println("User not found:", err)
// 			w.WriteHeader(http.StatusNotFound)
// 			return
// 		}

// 		// Generate a new password
// 		newPassword := "newPassword" // This should be a randomly generated password

// 		hashedPass, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
// 		if err != nil {
// 			log.Println("Error hashing password:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}
// 		foundUser.Password = string(hashedPass)

// 		// Update the user's password in the database
// 		_, err = collection.UpdateOne(context.Background(), bson.M{"username": user.Username}, bson.M{"$set": bson.M{"password": foundUser.Password}})
// 		if err != nil {
// 			log.Println("Error updating password in database:", err)
// 			w.WriteHeader(http.StatusInternalServerError)
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 	}
// }
