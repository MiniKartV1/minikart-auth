package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MiniKartV1/minikart-auth/pkg/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Adapter struct {
	db              *mongo.Database
	usersCollection *mongo.Collection
}

// constructor for our db adapter
func NewAdapter(uri string) *Adapter {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to connect to database %v", err)
	}
	fmt.Println("Connected Succesfully and pinged MongoDB!")
	db := client.Database("naresh-apps")
	usersCollection := db.Collection("users")
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{"email", 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := usersCollection.Indexes().CreateOne(context.Background(), indexModel); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\n Unique index on 'email' created succesfully.")
	return &Adapter{
		db:              db,
		usersCollection: usersCollection,
	}
}

func (dbClient Adapter) CloseDBConnection() {
	// implements db close connection
	log.Output(1, "Closing the db connection")
}

func (dbClient Adapter) AddUser(user *models.User) error {
	// register

	result, err := dbClient.usersCollection.InsertOne(context.TODO(), user)

	if mongoErr, ok := err.(mongo.WriteError); ok { // Type assertion to get the WriteError
		if mongoErr.Code == 11000 {
			// Extract and return the custom error message for duplicate email
			return fmt.Errorf("duplicate email error: the email already exists")
		}
	} else if mongoErr, ok := err.(mongo.WriteException); ok {
		for _, we := range mongoErr.WriteErrors {
			if we.Code == 11000 {
				return fmt.Errorf("EMAIL_EXISTS:the email already exists")
			}
		}
	}
	fmt.Println("Added user to the system", result)
	return nil
}
func (dbClient Adapter) UpdatePassword() {
	// changepassword
	fmt.Println("Updating the password")
}
func (dbClient Adapter) FindUserByEmail(email *string) (*models.User, error) {
	var user models.User
	err := dbClient.usersCollection.FindOne(context.TODO(), bson.M{"email": email}).Decode(&user)

	if err != nil {
		log.Printf("Error in finding user: %v", err)
		return &models.User{}, err
	}
	return &user, nil
}
func (dbClient Adapter) UpdateLastSignedIn(email *string) (*models.User, error) {
	if email == nil {
		return nil, errors.New("NIL_EMAIL: bad parameters passed")
	}
	var updatedUser models.User
	update := bson.M{
		"$set": bson.M{"lastSignedIn": time.Now()},
	}
	filter := bson.M{"email": email}
	res, err := dbClient.usersCollection.UpdateOne(context.TODO(), filter, update)

	if res.MatchedCount == 0 {
		return nil, errors.New("NO_USER_FOUND:no users found for email")
	}
	err = dbClient.usersCollection.FindOne(context.TODO(), filter).Decode(&updatedUser)
	if err != nil {
		return nil, err
	}

	return &updatedUser, nil
}
func (dbClient Adapter) SignOut() {
	fmt.Println("Signing out the user")
}
func (dbClient Adapter) SaveCode() {
	// creates code for the user in the database
	fmt.Println("Saving code.")
}
