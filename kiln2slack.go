package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    //"net/mail"
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


func SendToSlack(r *http.Request) {

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
        /*
        addr, err := mail.ParseAddress(update.Commits[i].Author)
        if err != nil {
            // LOL
            log.Println("Error parsing email address", update.Commits[i].Author, ":", err)
            return
        }
        */

        message.Attachments[i].Color = "good"
        message.Attachments[i].Text = fmt.Sprintf("<%s|%d> %s – %s",
            update.Commits[i].Url,
            update.Commits[i].Revision,
            update.Commits[i].Message,
            //addr.Name,
            update.Commits[i].Author,
        )
    }

    jsonMessage, err := json.Marshal(message)
    if err != nil {
        log.Println("Error marshalling message to JSON:", err)
        return
    }
    outputReader := bytes.NewReader(jsonMessage)
    resp, err := http.Post("https://domainname.slack.com/services/hooks/incoming-webhook?token=token",
        "application/json", outputReader)
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

    http.HandleFunc("/kiln/unityproject", func (w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "OK")
        go SendToSlack(r)
    })

    log.Fatal(http.ListenAndServe(":10100", nil))

}
