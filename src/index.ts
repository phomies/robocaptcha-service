import twilio from "twilio";
import express from 'express';
import bodyParser from "body-parser";

const VoiceResponse = twilio.twiml.VoiceResponse;

const app = express()
app.use(bodyParser.urlencoded( {extended: true} ));

app.get("/", (req, res) => {
    res.send("It works!!!!");
})


// Initial Incoming Call Handler
app.post("/incoming", (req, res) => {
    
    console.log(req.body);

    const twiml = new VoiceResponse();
    var randomNumber = Math.floor(1000 + Math.random() * 9000);
    var words = ["comprehensive", "popular", "interesting", "singular", "triangular", "absolutely"];
    var randomWord = words[Math.floor(Math.random() * words.length)];

    twiml.gather({
        
        input: ["dtmf", "speech"],
        timeout: 3,
        numDigits: 4,
        action: `/incoming/verify/${randomNumber}/${randomWord}`,
        language: "en-SG"

    }).say(

        {voice: "Polly.Brian"},
        `Your call is being screened for human verification,
        please enter ${randomNumber.toString().split("").join(" ")}.
        Alternatively, you may say the word ${randomWord}`
    );

    twiml.redirect('/incoming');

    res.send(twiml.toString());
    console.log(twiml.toString());
})


// Incoming Word/Digits Verifier
app.post("/incoming/verify/:correctDigits/:correctWord", (req, res) => {

    console.log(req.body);

    const twiml = new VoiceResponse();

    if (req.body.Digits && req.params.correctDigits == req.body.Digits) {

        twiml.say(
            {voice: "Polly.Brian"},
            "You did it, redirecting your call now.");

    } else if (req.body.SpeechResult.toLowerCase().includes(req.params.correctWord)) {
        twiml.say(
            {voice: "Polly.Brian"},
            "You did it, redirecting your call now.");
            
    } else {
        twiml.say(
            {voice: "Polly.Brian"},
            "Incorrect, please try again.");
        twiml.redirect("/incoming");
    }

    res.send(twiml.toString());
    console.log(twiml.toString());

})

app.listen(5000, () => console.log("Running on port 5000"));