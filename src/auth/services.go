package auth

import (
	"context"
	"crypto/md5"
	"errors"
	"fmt"
	database "systems-management-api/core/database"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

/* AUTHENTICATION */

// AuthenticationService common interface for all authentication service providers
type AuthenticationService interface {
	Authenticate(email string, password string) bool
}

// databaseAuthenticationService database specific authentication service
type databaseAuthenticationService struct{}

// Authenticate authenticates the user with the provided email and password
func (dbAuth *databaseAuthenticationService) Authenticate(email string, password string) bool {
	var user User
	db := database.DB()
	collection := db.D.Collection("user")

	md5Password := []byte(password)

	err := collection.FindOne(
		context.TODO(),
		bson.M{"email": email, "password": fmt.Sprintf("%x", md5.Sum(md5Password))},
	).Decode(&user)

	if err != nil {
		zap.S().Debugw("AuthenticationService, user not found: ", "error", err)
		return false
	} else {
		zap.S().Debugw("AuthenticationService, found user", "email", user.Email)
		return true
	}
}

// NewDatabaseAuthenticationService constructor for the databaseAuthenticationService
func NewDatabaseAuthenticationService() AuthenticationService {
	return &databaseAuthenticationService{}
}

/* JWT */

type JWTService interface {
	GenerateToken(email string) string
	ValidateToken(token string) (*JwtClaim, error)
}
type JwtClaim struct {
	Email string `json:"email"`
	jwt.StandardClaims
}

type jwtService struct {
	secretKey string
	issuer    string
}

// JWTAuthService a service which provides generate key and validate key methods
func JWTAuthService() JWTService {
	return &jwtService{
		secretKey: getSecretKey(),
		issuer:    "Otto",
	}
}

func getSecretKey() string {
	secret := viper.GetString("jwt.secret")
	if secret == "" {
		secret = "secret"
	}
	return secret
}

func (service *jwtService) GenerateToken(email string) string {
	claims := &JwtClaim{
		email,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 48).Unix(),
			Issuer:    service.issuer,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// encoded string
	t, err := token.SignedString([]byte(service.secretKey))
	if err != nil {
		panic(err)
	}
	return t
}

// https://betterprogramming.pub/hands-on-with-jwt-in-golang-8c986d1bb4c0
func (service *jwtService) ValidateToken(signedToken string) (*JwtClaim, error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JwtClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(service.secretKey), nil
		},
	)

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JwtClaim)
	if !ok {
		err = errors.New("Couldn't parse claims")
		return nil, err
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("JWT is expired")
		return nil, err
	}

	return claims, nil

}

// UserService service which provides methos to access and modify database data
type UserService struct{}

func (service *UserService) all() (*[]User, error) {
	db := database.DB()
	collection := db.D.Collection("user")
	cursor, err := collection.Find(context.TODO(), bson.D{{}})

	if err != nil {
		return nil, err
	} else {
		users := []User{}
		for cursor.Next(context.TODO()) {
			var user User
			cursor.Decode(&user)
			users = append(users, user)
		}
		return &users, nil
	}
}

func (service *UserService) GetById(id string) (*User, error) {
	db := database.DB()
	user := User{}

	if err := db.GetById("user", id, &user); err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// Returns an user instance given an email
func (service *UserService) GetByEmail(email string) (*User, error) {
	var user User
	db := database.DB()
	collection := db.D.Collection("user")
	err := collection.FindOne(
		context.TODO(),
		bson.M{"email": email},
	).Decode(&user)

	if err != nil {
		return nil, err
	} else {
		return &user, nil
	}
}

// Saves the user model to database
// Returns boolean result and error
func (service *UserService) Save(user *User) (bool, error) {
	db := database.DB()
	collection := db.D.Collection("user")

	if user.ID.IsZero() {
		// insert
		res, err := collection.InsertOne(context.TODO(), user)

		if err != nil {
			zap.S().Error("Error inserting user: ", err)
			return false, err
		} else {
			zap.S().Info(fmt.Sprintf("User %s inserted succesfully", user.Email))
			// update user ID
			user.ID = res.InsertedID.(primitive.ObjectID)
			return true, nil
		}
	} else {
		// update
		filter := bson.M{"_id": user.ID}
		_, err := collection.ReplaceOne(context.TODO(), filter, user)

		if err != nil {
			zap.S().Error("Error inserting user: ", err)
			return false, err
		} else {
			zap.S().Info(fmt.Sprintf("User %s updated succesfully", user.Email))
			// update user ID
			return true, nil
		}
	}
}

// Deleted the user model from databse
// Returns boolean result and error
func (service *UserService) Delete(user *User) (bool, error) {
	db := database.DB()
	collection := db.D.Collection("user")

	_, err := collection.DeleteOne(context.TODO(), bson.M{"_id": user.ID})

	if err != nil {
		zap.S().Error("Error deleting user: ", err)
		return false, err
	} else {
		zap.S().Info(fmt.Sprintf("User %s deleted succesfully", user.Email))
		// update user ID
		return true, nil
	}
}
