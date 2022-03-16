package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                string             `bson:"_id"`
	Name              string             `bson:"name"`
	Password          string             `bson:"password"`
	Email             string             `bson:"email"`
	ProxyNumber       string             `bson:"proxyNumber"`
	PhoneNumber       string             `bson:"phoneNumber"`
	VerificationLevel int                `bson:"verificationLevel"`
	DateJoined        primitive.DateTime `bson:"dateJoined"`
}

type Contact struct {
	ID            primitive.ObjectID `bson:"_id"`
	UserID        string             `bson:"userId"`
	Name          string             `bson:"name"`
	Number        string             `bson:"number"`
	IsBlacklisted bool               `bson:"isBlacklisted"`
	IsWhitelisted bool               `bson:"isWhitelisted"`
	CreatedAt     primitive.DateTime `bson:"createdAt"`
	UpdatedAt     primitive.DateTime `bson:"updatedAt"`
}

type Payment struct {
	DateStart     primitive.DateTime `bson:"dateStart"`
	DateEnd       primitive.DateTime `bson:"dateEnd"`
	Amount        float64            `bson:"amount"`
	TransactionId string             `bson:"transactionId"`
	Plan          string             `bson:"plan"`
}

type Call struct {
	ID       primitive.ObjectID `bson:"_id"`
	DateTime time.Time          `bson:"dateTime"`
	CallSid  string             `bson:"callSid"`
	From     string             `bson:"from"`
	ToUserId string             `bson:"toUserId"`
	Action   string             `bson:"action"`
}

type Notification struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   string             `bson:"userId"`
	DateTime time.Time          `bson:"dateTime"`
	Content  string             `bson:"content"`
	Read     bool               `bson:"read"`
	URL      string             `bson:"url"`
}

/**
Get MongoDB context associated
*/
func getMongoCollection(collectionName string) (*mongo.Collection, *mongo.Client, context.Context, context.CancelFunc) {
	clientOptions := options.Client().ApplyURI(os.Getenv("DB_CONN_STRING"))
	clientOptions = clientOptions.SetConnectTimeout(1 * time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	return client.Database("callcaptcha").Collection(collectionName), client, ctx, cancel
}

/**
Get pointer to user object, based on the masked number of a recipient
*/
func getUserFromMaskedNumber(maskedRecipient string) *User {
	fmt.Println("Masked Recipient:", maskedRecipient)

	usersCollection, client, ctx, cancel := getMongoCollection("users")
	defer func() {
		cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	users, err := usersCollection.Find(ctx, bson.M{"maskedNumber": maskedRecipient})
	if err != nil {
		panic(err)
	}

	if users.Next(ctx) {
		var user User
		err = users.Decode(&user)
		if err != nil {
			panic(err)
		}
		return &user
	}

	return nil
}

/*
Gets the contact of a given user ID
*/
func getContactIfExists(recipientUserId string, callerNumber string) *Contact {
	fmt.Println("Getting contact information user ID", recipientUserId, "and", callerNumber)

	contactsCollection, client, ctx, cancel := getMongoCollection("contacts")
	defer func() {
		cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()
	query := bson.M{"userId": recipientUserId, "number": callerNumber}

	matchedContact, err := contactsCollection.Find(ctx, query)

	if err != nil {
		panic(err)
	}

	if matchedContact.Next(ctx) {
		var contact Contact
		err = matchedContact.Decode(&contact)
		if err != nil {
			panic(err)
		}
		return &contact
	}

	return nil
}

/*
Insert a call object to the database
*/
func insertCall(callSid string, fromNumber string, toUserId string) {

	callStruct := &Call{
		ID:       primitive.NewObjectID(),
		DateTime: time.Now(),
		CallSid:  callSid,
		From:     fromNumber,
		ToUserId: toUserId,
		Action:   "in-progress",
	}

	callsCollection, client, ctx, cancel := getMongoCollection("calls")
	defer func() {
		cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	_, err := callsCollection.InsertOne(ctx, callStruct)
	if err != nil {
		panic(err)
	}
}

/*
Update the call object in the database with the following action
*/
func updateCall(callSid string, action string) {

	callsCollection, client, ctx, cancel := getMongoCollection("calls")
	defer func() {
		cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	updateCriteria, updateAction := bson.M{"callSid": callSid}, bson.M{"$set": bson.M{"action": action}}
	_, err := callsCollection.UpdateOne(ctx, updateCriteria, updateAction)
	if err != nil {
		panic(err)
	}
}

/*
Insert new notification
*/
func insertNotification(content string, userId string) primitive.ObjectID {

	notificationStruct := &Notification{
		ID:       primitive.NewObjectID(),
		UserID:   userId,
		DateTime: time.Now(),
		Content:  content,
		Read:     false,
		URL:      "https://google.com",
	}

	notificationsCollection, client, ctx, cancel := getMongoCollection("notifications")
	defer func() {
		cancel()
		err := client.Disconnect(ctx)
		if err != nil {
			panic(err)
		}
	}()

	_, err := notificationsCollection.InsertOne(ctx, notificationStruct)
	if err != nil {
		panic(err)
	}

	return notificationStruct.ID
}
