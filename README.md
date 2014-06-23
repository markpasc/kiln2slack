kiln2slack
==========

This Go application receives commit webhook calls from [Kiln](http://www.fogcreek.com/kiln/) and sends equivalent messages to a [Slack](https://slack.com) chat room incoming webhook.


Install
-------

Web apps can be tricky to set up. Try this:

1. Install Go on your web server.
2. Clone this repository to a `kiln2slack/` directory.
3. `cd kiln2slack/`
4. `mkdir .env/`
5. `GOPATH=.env go build`
6. `./kiln2slack`

This will start running the service on port 10100. At this point you'll need to edit the `kiln2slack.go` source code to change the port number (look at the end of the file), then do another `GOPATH=.env go build`.

I use supervisor to keep the kiln2slack app running, and nginx to plumb requests from the internet to the kiln2slack app.

Once you have requests from the internet to the kiln2slack app, you should be able to go to `https://yoursite.com/kiln/` and see `404 page not found` message which weirdly is the right thing!


Setting up
----------

Once your kiln2slack app can run, edit the `slackNameToUrl.json` file.

```json
{
    "projectname": "https://domainname.slack.com/services/hooks/incoming-webhook?token=token",
    "anotherproject": "https://domainname.slack.com/services/hooks/incoming-webhook?token=anothertoken"
}
```

On the left are **project names** and the right are **Slack webhook URLs**. Go to the to set up an “incoming webhook” integration to get the Slack URL for the right side. The left side is what you make up to put in Kiln. For instance if you called your project `projectname`, when you add the webhook to Kiln it'll be: `https://yoursite.com/kiln/projectname`

Then restart the kiln2slack app (with ctrl-C then running it again, or restarting it in supervisor or what-have-you).

Whew! That's a lot, but then it should work.
