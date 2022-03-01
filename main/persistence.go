package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type User struct {
	ID                  primitive.ObjectID   `bson:"_id"`
	Name                string               `bson:"name"`
	Password            string               `bson:"password"`
	Email               string               `bson:"email"`
	ProxyNumber         string               `bson:"proxyNumber"`
	PhoneNumber         string               `bson:"phoneNumber"`
	VerificationLevel   int                  `bson:"verificationLevel"`
	DateJoined          primitive.DateTime   `bson:"dateJoined"`
	Whitelist           []string             `bson:"whitelist"`
	Blacklist           []string             `bson:"blacklist"`
	SubscriptionHistory []Payment            `bson:"subscriptionHistory"`
	CallHistory         []primitive.ObjectID `bson:"callHistory"`
	NotificationHistory []primitive.ObjectID `bson:"notificationHistory"`
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
	DateTime primitive.DateTime `bson:"dateTime"`
	CallSid  string             `bson:"callSid"`
	From     string             `bson:"from"`
	ToUserId primitive.ObjectID `bson:"toUserId"`
	Action   string             `bson:"action"`
}

type Notification struct {
	ID       primitive.ObjectID `bson:"_id"`
	UserID   primitive.ObjectID `bson:"userId"`
	DateTime primitive.DateTime `bson:"dateTime"`
	Content  string             `bson:"content"`
	Read     bool               `bson:"read"`
	URL      string             `bson:"url"`
}

func getUserFromMaskedNumber(maskedRecipient string) *User {

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
	fmt.Println("Connecting to MongoDB:", os.Getenv("DB_CONN_STRING"))
	fmt.Println("Masked Recipient:", maskedRecipient)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_CONN_STRING")))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer client.Disconnect(ctx)
	defer ctxCancel()

	callcaptchaDb := client.Database("callcaptcha")
	matchedUser, err := callcaptchaDb.Collection("users").Find(ctx, bson.M{"maskedNumber": maskedRecipient})
	if err != nil {
		fmt.Println("A")
		fmt.Println(err)
		panic(err)
	}

	if matchedUser.Next(ctx) {
		var user User
		err = matchedUser.Decode(&user)

		if err != nil {
			fmt.Println("B")
			fmt.Println(err)
			panic(err)
		}

		fmt.Println(user)
		return &user
	}

	return nil
}

func getUsersCollection(w http.ResponseWriter, r *http.Request) {

	ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_CONN_STRING")))
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer client.Disconnect(ctx)
	defer ctxCancel()

	callcaptchaDb := client.Database("callcaptcha")
	callcaptchaUsers, err := callcaptchaDb.Collection("users").Find(ctx, bson.D{})

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	var users []User
	err = callcaptchaUsers.All(ctx, &users)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println(len(users))

	for _, user := range users {
		fmt.Println(user)
		w.Write([]byte(user.Name))
		w.Write([]byte(user.Email))
	}

}
