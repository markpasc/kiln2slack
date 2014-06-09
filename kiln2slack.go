package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "strings"
)

type SlackAttachment struct {
    Color string `json:"color"`
    Text string `json:"text"`
    Fallback string `json:"fallback"`
}

type SlackMessage struct {
    Text string `json:"text"`
    Attachments []SlackAttachment `json:"attachments"`
}

type KilnCommit struct {
    Url string `json:"url"`
    Revision int `json:"revision"`
    Author string `json:"author"`
    Message string `json:"message"`
}

type KilnPusher struct {
    Email string `json:"email"`
    FullName string `json:"fullName"`
}

type KilnRepository struct {
    ID int `json:"id"`
    Name string `json:"name"`
    Url string `json:"url"`
}

type KilnUpdate struct {
    Commits []KilnCommit `json:"commits"`
    Pusher KilnPusher `json:"pusher"`
    Repository KilnRepository `json:"repository"`
}


var SlackUrlForName map[string]string


func SendToSlack(slackUrl string, r *http.Request) {

    err := r.ParseForm()
    if err != nil {
        log.Println("Error parsing request's form data:", err)
        return
    }
    payload := r.Form["payload"]
    if len(payload) < 1 {
        log.Println("Request form did not include a 'payload' form value???")
        return
    }

    var update KilnUpdate
    err = json.Unmarshal([]byte(payload[0]), &update)
    if err != nil {
        log.Println("Error unmarshalling JSON from Kiln:", err)
        return
    }

    message := SlackMessage{}
    message.Text = fmt.Sprintf("%s pushed to <%s|%s>", update.Pusher.FullName, update.Repository.Url, update.Repository.Name)
    message.Attachments = make([]SlackAttachment, len(update.Commits), len(update.Commits))
    for i := 0; i < len(update.Commits); i++ {
        // TODO: use mail.ParseAddress in 1.3 when it exists?
        parts := strings.SplitN(update.Commits[i].Author, " <", 2)
        name := parts[0]

        message.Attachments[i].Color = "good"
        message.Attachments[i].Text = fmt.Sprintf("<%s|%d> %s – %s",
            update.Commits[i].Url,
            update.Commits[i].Revision,
            update.Commits[i].Message,
            name,
        )
    }

    jsonMessage, err := json.Marshal(message)
    if err != nil {
        log.Println("Error marshalling message to JSON:", err)
        return
    }
    outputReader := bytes.NewReader(jsonMessage)
    resp, err := http.Post(slackUrl, "application/json", outputReader)
    if err != nil {
        log.Println("Error posting message to Slack:", err)
        return
    }
    if resp.StatusCode != 200 {
        log.Println("Response from Slack had status code", resp.StatusCode)
        return
    }
}


func main() {

    file, err := os.Open("./slackUrlForName.json")
    if err != nil {
        log.Fatalln("Error opening slack URL map file:", err)
    }
    dec := json.NewDecoder(file)
    err = dec.Decode(&SlackUrlForName)
    if err != nil {
        log.Fatalln("Error reading slack URL map file:", err)
    }
    file.Close()

    http.HandleFunc("/kiln/", func (w http.ResponseWriter, r *http.Request) {

        pathParts := strings.SplitN(r.URL.Path, "/", 3)
        hookName := pathParts[2]
        slackUrl, ok := SlackUrlForName[hookName]
        if !ok {
            http.NotFound(w, r)
            return
        }

        fmt.Fprintf(w, "OK")
        go SendToSlack(slackUrl, r)
    })

    log.Fatal(http.ListenAndServe(":10100", nil))

}
