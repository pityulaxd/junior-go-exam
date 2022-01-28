# InterTicket junior Go developer test

## Prerequisites

- [https://github.com](github.com) account
- go (1.16+)
- git client

## Tasks

Create a REST endpoint which will save and lists out events.
The exam has two parts: just listing out the existing events and the second part is to save the event via a rest call, and list out the saved events too.
The application should produce `application/json` content type and should process the `application/json` body content with `PUT` method to validate and save the event payload.

> the prepared go code available in main.go file. You should pass the tests defined in main_test.go file. The test can be runnable with `go test ./` command

⚠️ Before you can start the task please create a fork from this repository open a feature branch and start your work there! When you finished with your work please open a pull request and I will check your code as soon as I can!

### Task 1: The event list handler

Implement the **GET** handler which returns a list of event objects at `/api/events`

For a `GET` method you should list the stored events in JSON format. The data should be generated from code, (with go struct marshalled by json package) and should return the following: [event-list.json](test/event-list.json) obviously in the same format.

### Task 2: Event creation

Implement a **PUT** handler where you can save the events: `/api/events`

This `PUT` handler is much more complex then the previous task, the endpoint wants a JSON request body like this:

```json
{
    "name": "Manon Lescaut",
    "venue": {
        "name": "Erkel Színház",
        "location": "1087 Budapest, II. János Pál pápa tér 30."
    },
    "description": "Giacomo Puccini: Manon Lescaut Opera két részben, négy felvonásban, olasz nyelven, magyar és angol felirattal",
    "date": "2022-02-06T18:00:00Z"
}
```

The endpoint should validate the sent payload (JSON Request body) as follows:

- The name should be minimum 10 characters
- The description field should have at least 30 characters
- the venue name should not be blank/empty
- the venue location field should not be blank/empty
- the date must be in UTC format (ex.: 	2022-01-28T11:41:57Z)

When error occurs the response status should be 400 (HTTP Bad Request) with an error message in the following structure:

```json
{
    "message": "[error message]"
}
```

After the request was ok, you have to generate an unique id for the event and store the event somewhere (put it in a slice/array). The endpoint should return a 201 code (HTTP Created) and should return the saved **event** object in JSON with the generated ID field.

## Bonus task (Optional)

Extend the application main function to use OS environment variable called: `PORT` to set the listening address of the HTTP server.

For example if you run the code with `PORT=5000 go run main.go` the server should listen on port `5000` for HTTP requests.

---

If you have any questions about the tasks feel free to contact me: <gyula.paal@interticket.hu>

Good luck, and have a lot of fun!
