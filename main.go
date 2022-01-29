package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"syscall"
	"time"

	"github.com/satori/uuid"

	"github.com/gin-gonic/gin"
)

type Venue struct {
	Name     string `json:"name"`
	Location string `json:"location"`
}

type Event struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Venue       Venue     `json:"venue"`
	Description string    `json:"description"`
	Date        string    `json:"date"`
}

var EventList []Event

func listEvents(c *gin.Context) {
	EventList = append(EventList, Event{
		Id:   uuid.Must(uuid.FromString("b44ac834-be8c-4457-bd71-72bc98ebfa4d")),
		Name: "ANDREA BOCELLI 2022",
		Venue: Venue{
			Name:     "Papp László Budapest Sportaréna",
			Location: "1143 Budapest, Stefánia út 2.",
		},
		Description: "Napjaink egyik legnépszerűbb tenorja, Andrea Bocelli 2022. október 15-én a Papp László Budapest Sportarénában ad koncertet!",
		Date:        "2022-10-15T18:00:00Z",
	})

	EventList = append(EventList, Event{
		Id:   uuid.Must(uuid.FromString("38896b7c-e221-43a1-977a-f01f7239f66b")),
		Name: "JAMES BLUNT THE STARS BENEATH MY FEET TOUR",
		Venue: Venue{
			Name:     "VeszprémFest, Veszprém Aréna",
			Location: "8200 Veszprém, Külső-kádártai u. 5.",
		},
		Description: "A VeszprémFest zárónapján a brit szupersztár, JAMES BLUNT ad koncertet. Az énekes 2022 februárjában induló, The Stars Beneath My Feet című Európa turnéjának keretében lép fel Veszprémben.",
		Date:        "2022-07-16T20:00:00Z",
	})
	c.AbortWithStatusJSON(http.StatusOK, EventList)

	//c.JSON(http.StatusOK, EventList)
}

func saveEvent(c *gin.Context) {
	//Create new event and errors variables
	var newEvent Event
	var errors []string

	//Binding payload to variable
	if err := c.ShouldBindJSON(&newEvent); err != nil {
		//Invalid payload error handling
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "[Invalid payload!]"})
	}
	//Validation of payload
	if len(newEvent.Name) < 10 {
		errors = append(errors, "Name is too short, min. 10 characters!")
	}
	if len(newEvent.Description) < 30 {
		errors = append(errors, "Description is too short, min. 30 characters!")
	}
	if len(newEvent.Venue.Name) == 0 {
		errors = append(errors, "Venue name is empty")
	}
	if len(newEvent.Venue.Location) == 0 {
		errors = append(errors, "Venue location name is empty")
	}
	valid, _ := regexp.MatchString("^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$", newEvent.Date)
	if !valid {
		errors = append(errors, "Date is invalid")
	}
	//Check for validation errors
	if len(errors) > 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": errors})
	} else {
		c.AbortWithStatusJSON(http.StatusCreated, newEvent)
	}
	//Generating UUID for the event
	uid := uuid.NewV4()
	newEvent.Id = uid
	//Add event to eventlist
	EventList = append(EventList, newEvent)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router := gin.Default()

	// Creating /api route prefix group
	ag := router.Group("/api")
	{
		// Registering list events handler
		ag.GET("/events", listEvents)

		// Registering event save handler
		ag.PUT("/events", saveEvent)
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			log.Printf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)

	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be caught, so don't need to add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
