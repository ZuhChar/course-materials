package wyoassign

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Response struct {
	Assignments []Assignment `json:"assignments"`
}

type Assignment struct {
	Id          string `json:"id"`
	Title       string `json:"title`
	Description string `json:"desc"`
	Points      int    `json:"points"`
	// DueDate string `json:"DueDate"`
}

var Assignments []Assignment

const Valkey string = "FooKey"

func InitAssignments() {
	var assignmnet Assignment
	assignmnet.Id = "Mike1A"
	assignmnet.Title = "Lab 4 "
	assignmnet.Description = "Some lab this guy made yesteday?"
	assignmnet.Points = 20
	Assignments = append(Assignments, assignmnet)
}

func APISTATUS(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "API is up and running")
}

func GetAssignments(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	var response Response

	response.Assignments = Assignments

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	jsonResponse, err := json.Marshal(response)

	if err != nil {
		return
	}

	w.Write(jsonResponse)
}

func GetAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	params := mux.Vars(r)

	response := make(map[string]string)

	found := false

	response["Status"] = "No such Assignment exists"
	for _, assignment := range Assignments {
		if assignment.Id == params["id"] {
			w.WriteHeader(http.StatusFound)
			json.NewEncoder(w).Encode(assignment)
			found = true
			break
		}
	}

	// reason for this style is so the print out can be uninterrupted by other Status info
	if !found {
		jsonResponse, err := json.Marshal(response)

		if err != nil {
			return
		}

		w.Write(jsonResponse)
		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s DELETE end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/txt")
	w.WriteHeader(http.StatusOK)
	params := mux.Vars(r)

	response := make(map[string]string)

	response["status"] = "No Such ID to Delete"
	for index, assignment := range Assignments {
		if assignment.Id == params["id"] {
			Assignments = append(Assignments[:index], Assignments[index+1:]...)
			response["status"] = "Successfully deleted"
			break
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResponse)
}

func UpdateAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")

	// same parsing format used in create, unsure if nessesary
	r.ParseForm()
	params := mux.Vars(r)

	response := make(map[string]string)

	response["status"] = "No such Assignment exists"
	// checks all ID's to find the assignment in need of updating
	for index, assignment := range Assignments {
		if assignment.Id == params["id"] {
			Assignments[index].Title = r.FormValue("title")
			Assignments[index].Description = r.FormValue("desc")
			Assignments[index].Points, _ = strconv.Atoi(r.FormValue("points"))
			w.WriteHeader(http.StatusOK)
			response["status"] = "Assignment updated"
			break
		}
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResponse)

}

func CreateAssignment(w http.ResponseWriter, r *http.Request) {
	log.Printf("Entering %s end point", r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	var assignmnet Assignment
	r.ParseForm()

	alreadyexists := false

	response := make(map[string]string)

	// meant to check if the assignment id already exists
	for _, assignment := range Assignments {
		if assignment.Id == r.FormValue("id") {
			alreadyexists = true
			break
		}
	}

	// Builds assignmnet and checks most values entered by user 
	if r.FormValue("id") != "" && !alreadyexists {
		assignmnet.Id = r.FormValue("id")
		if r.FormValue("title") == "" {
			assignmnet.Title = "Assignment"
		} else {
			assignmnet.Title = r.FormValue("title")
		}
		assignmnet.Description = r.FormValue("desc") // Description can be empty if teach is evil
		var points, _ = strconv.Atoi(r.FormValue("points"))
		if points < 0 {
			points = 0
			assignmnet.Points = points
		} else {
			assignmnet.Points = points
		}
		Assignments = append(Assignments, assignmnet)
		w.WriteHeader(http.StatusCreated)
		response["Status"] = "Assignment created"
	} else if alreadyexists {
		response["Status"] = "Assignment already exists"
	} else if r.FormValue("id") == "" {
		response["Status"] = "Invalid ID"
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return
	}
	w.Write(jsonResponse)

	w.WriteHeader(http.StatusNotFound)
}