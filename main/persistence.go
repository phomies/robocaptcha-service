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
    DateTime time.Time          `bson:"dateTime"`
    CallSid  string             `bson:"callSid"`
    From     string             `bson:"from"`
    ToUserId primitive.ObjectID `bson:"toUserId"`
    Action   string             `bson:"action"`
}

type Notification struct {
    ID       primitive.ObjectID `bson:"_id"`
    UserID   primitive.ObjectID `bson:"userId"`
    DateTime time.Time			`bson:"dateTime"`
    Content  string             `bson:"content"`
    Read     bool               `bson:"read"`
    URL      string             `bson:"url"`
}

func getMongoContext() (context.Context, *mongo.Client, context.CancelFunc) {
    ctx, ctxCancel := context.WithTimeout(context.Background(), 5*time.Second)
    fmt.Println("Connecting to MongoDB:", os.Getenv("DB_CONN_STRING"))

    client, err := mongo.Connect(ctx, options.Client().ApplyURI(os.Getenv("DB_CONN_STRING")))
    if err != nil {
        fmt.Println(err)
        panic(err)
    }

    return ctx, client, ctxCancel
}

func getUserFromMaskedNumber(maskedRecipient string) *User {
    fmt.Println("Masked Recipient:", maskedRecipient)
    ctx, client, ctxCancel := getMongoContext()

    defer client.Disconnect(ctx)
    defer ctxCancel()

    callcaptchaDb := client.Database("callcaptcha")
    matchedUser, err := callcaptchaDb.Collection("users").Find(ctx, bson.M{"maskedNumber": maskedRecipient})
    if err != nil {
        panic(err)
    }

    if matchedUser.Next(ctx) {
        var user User

        err = matchedUser.Decode(&user)
        if err != nil {
            panic(err)
        }

        //fmt.Println(user)
        return &user
    }
    return nil
}

func insertCall(callSid string, fromNumber string, toUserId primitive.ObjectID) {
    ctx, client, ctxCancel := getMongoContext()
    defer client.Disconnect(ctx)
    defer ctxCancel()

    callStruct := &Call{
        ID:       primitive.NewObjectID(),
        DateTime: time.Now(),
        CallSid:  callSid,
        From:     fromNumber,
        ToUserId: toUserId,
        Action:   "in-progress",
    }

    callcaptchaDb := client.Database("callcaptcha")
    callcaptchaDb.Collection("calls").InsertOne(ctx, callStruct)
}

func updateCall(callSid string, action string) {
    ctx, client, ctxCancel := getMongoContext()
    defer client.Disconnect(ctx)
    defer ctxCancel()

    callcaptchaDb := client.Database("callcaptcha")
    callcaptchaDb.Collection("calls").UpdateOne(ctx, bson.M{"callSid": callSid}, bson.M{"$set": bson.M{"action": action}})
}

func insertNotification(content string, userId primitive.ObjectID) primitive.ObjectID {
    ctx, client, ctxCancel := getMongoContext()
    defer client.Disconnect(ctx)
    defer ctxCancel()

    notificationStruct := &Notification{
        ID: primitive.NewObjectID(),
        UserID: userId,
        DateTime: time.Now(),
        Content: content,
        Read: false,
        URL: "http://google.com",
    }

    callcaptchaDb := client.Database("callcaptcha")
    callcaptchaDb.Collection("notifications").InsertOne(ctx, notificationStruct)
    return notificationStruct.ID
}
