package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// XML Generation for Call Ending

// Web service health check
func healthCheck(ctx *gin.Context) {
	ctx.XML(http.StatusOK, structHelloWorld())
}

/**
 *  Incoming call handler
 */
func httpIncoming(ctx *gin.Context) {

	callSid := ctx.PostForm("CallSid")
	numberFrom := ctx.Request.FormValue("From")
	numberTo := ctx.Request.FormValue("Called")

	// Get number and user details
	if numberTo == "" {
		numberTo = os.Getenv("DEFAULT_NUMBER_TO")
	}
	userDialed := getUserFromMaskedNumber(numberTo)
	if userDialed == nil {
		ctx.XML(http.StatusOK, structUserNotFound())
		return
	}
	userDialedId := (*getUserFromMaskedNumber(numberTo)).ID

	// Get the number of times the call has happened
	times, err := strconv.Atoi(ctx.Param("times"))
	if err != nil {
		times = 0
	}

	// First instance, record call into database
	if times == 0 {
		insertCall(callSid, numberFrom, userDialedId)
	}

	// Get blacklist/whitelist information
	contactInfo := getContactIfExists(userDialedId, numberFrom)
	if contactInfo != nil {
		if contactInfo.IsBlacklisted {
			ctx.XML(http.StatusOK, structBlacklisted())
			oid := insertNotification("Blocked call from "+numberFrom, userDialed.ID)
			sqsSendNotification(oid.Hex())
			updateCall(callSid, "blacklisted")
			return
		}
		if contactInfo.IsWhitelisted {
			ctx.XML(http.StatusOK, structWhitelisted(userDialed.PhoneNumber, numberFrom))
			oid := insertNotification("Successful call from "+numberFrom, userDialed.ID)
			sqsSendNotification(oid.Hex())
			updateCall(callSid, "whitelisted")
			return
		}
	}

	// Block if called too many times
	if times > 2 {
		updateCall(callSid, "timeout")
		oid := insertNotification("Call timed out from number "+numberFrom, userDialedId)
		fmt.Println("Object ID: + ", oid.String())
		sqsSendNotification(oid.Hex())
		ctx.XML(http.StatusOK, structEndCall())
		return
	}

	// Generate random numbers and words for voice recognition
	randomWord := generateRandomWord()
	randomNumber := generateRandomNumber()
	randomNumberSpaced := addSpaces(strconv.Itoa(randomNumber))
	actionUrl := fmt.Sprintf("/verify/%d/%s/%d", randomNumber, randomWord, times)
	voicePrompt := fmt.Sprintf(
		`Your call is being screened by robo captcha. Please enter the numbers %s or say the word %s.`,
		randomNumberSpaced,
		randomWord,
	)

	twimlStruct := structVerifyCall(actionUrl, voicePrompt, times)
	ctx.XML(http.StatusOK, twimlStruct)
}

/**
 *  Perform verification on incoming calls
 */
func httpVerify(ctx *gin.Context) {

	correctNum := ctx.Param("verifyNum")
	correctWord := ctx.Param("verifyWord")
	callSid := ctx.Request.FormValue("CallSid")
	digitsEntered := ctx.Request.FormValue("Digits")
	numberTo := ctx.Request.FormValue("Called")
	numberFrom := ctx.Request.FormValue("From")
	times, _ := strconv.Atoi(ctx.Param("times"))

	// Get user details
	var userDialed User
	if numberTo == "" {
		numberTo = os.Getenv("DEFAULT_NUMBER_TO")
	}
	userDialed = *getUserFromMaskedNumber(numberTo)

	// Block call if timeout
	if times > 2 {
		ctx.XML(http.StatusOK, structEndCall())
		oid := insertNotification("Call timed out from number "+numberFrom, userDialed.ID)
		sqsSendNotification(oid.Hex())
		updateCall(callSid, "timeout")
		return
	}

	// Get speech from call
	speechResult := strings.ToLower(ctx.Request.FormValue("SpeechResult"))
	fmt.Println(speechResult, digitsEntered)

	speechOk := strings.Contains(speechResult, correctWord)
	digitsOk := digitsEntered == correctNum

	if speechOk || digitsOk {
		forwardTo := userDialed.PhoneNumber

		// Correct, allow the call to pass through
		twimlStruct := structForwardingCall(forwardTo, numberFrom)
		oid := insertNotification("Successful call from "+numberFrom, userDialed.ID)
		sqsSendNotification(oid.Hex())
		updateCall(callSid, "success")
		ctx.XML(http.StatusOK, twimlStruct)
		return
	}

	// Incorrect, try again and increment count
	ctx.XML(http.StatusOK, structIncorrectEntered(times))
}
