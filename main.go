package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type event struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
}

// type eventList []event

// var allEvents = eventList{
// 	{
// 		ID:          "0",
// 		Title:       "mapple leave",
// 		Description: "I don't even know what this item is",
// 	},
// }
var allEvents = []event{
	{
		ID:          "0",
		Title:       "mapple leave",
		Description: "I don't even know what this item is",
	},
}

func rootRoute(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my first server")
}

func createEvent(w http.ResponseWriter, r *http.Request) {
	var newEvent event
	reqBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Fprintf(w, "Enter in the 'Title' and 'Description' of the new event alone to create a new event")
	} else {
		json.Unmarshal(reqBody, &newEvent)
		allEvents = append(allEvents, newEvent)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(allEvents)
	}
}

func getOneEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if len(id) == 0 {
		fmt.Fprintf(w, "You must enter in a valid id for this endpoint")
	} else {
		for _, fetchEvent := range allEvents {
			if fetchEvent.ID == id {
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(fetchEvent)
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Error! No event exist with this id")
	}
}

func getAllEvents(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(allEvents)
}

func updateEvent(w http.ResponseWriter, r *http.Request) {
	var newUpdate event
	id := mux.Vars(r)["id"]
	reqBody, error := ioutil.ReadAll(r.Body)
	if error != nil {
		fmt.Fprintf(w, "An Error has occured with passing request body")
	} else {
		for i, eventInst := range allEvents {
			if eventInst.ID == id {
				json.Unmarshal(reqBody, &newUpdate)
				if newUpdate.Title != "" {
					eventInst.Title = newUpdate.Title
				}
				if newUpdate.Description != "" {
					eventInst.Description = newUpdate.Description
				}

				allEvents = append(allEvents[:i], eventInst)
				w.WriteHeader(http.StatusOK)
				json.NewEncoder(w).Encode(eventInst)
				return
			}
		}
	}
}

func deleteEvent(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	for i, eventInst := range allEvents {
		if eventInst.ID == id {
			allEvents = append(allEvents[:i], allEvents[i+1:]...)
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "Event successfully deleted")
		}
	}
}

func deleteAllEvent(w http.ResponseWriter, r *http.Request) {
	allEvents = []event{}
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Successfully deleted all events")
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", rootRoute).Methods(http.MethodGet)
	router.HandleFunc("/create", createEvent).Methods(http.MethodPost)
	router.HandleFunc("/update/{id}", updateEvent).Methods(http.MethodPatch)
	router.HandleFunc("/fetch/all", getAllEvents).Methods(http.MethodGet)
	router.HandleFunc("/fetch/{id}", getOneEvent).Methods(http.MethodGet)
	router.HandleFunc("/delete/all", deleteAllEvent).Methods(http.MethodDelete)
	router.HandleFunc("/delete/{id}", deleteEvent).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8001", router))
}
