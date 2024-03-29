package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/zamachnoi/viewthis/data"
	"github.com/zamachnoi/viewthis/models"
	"github.com/zamachnoi/viewthis/util"
)

func GetAllQueuesHandler(w http.ResponseWriter, r *http.Request) {
	// get page limit and search query
	limit, page := util.ParseLimitAndPage(r)
	search := r.URL.Query().Get("search")

	queues, count, err := data.GetAllQueues(page, limit, search)
	if err != nil {
		log.Println("Error retrieving queues: ", err)
		http.Error(w, "Error retrieving queues", http.StatusInternalServerError)
		return
	}

	totalPages := count / limit
	if count%limit != 0 {
		totalPages++
	}

	response := map[string]interface{}{
		"queues":    queues,
		"page":      page,
		"totalPages": totalPages,
		"count":     count,
		"limit":     limit,
	}
	json.NewEncoder(w).Encode(response)
}
// create a queue
func CreateQueueHandler(w http.ResponseWriter, r *http.Request) {
	// get the user id from the context
	user, ok := r.Context().Value(util.UserKey).(util.SessionJWTWithClaims)
	if !ok {
		http.Error(w, "Error getting user from context", http.StatusInternalServerError)
		return
	}
	var newQueue models.Queue
	err := json.NewDecoder(r.Body).Decode(&newQueue)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	newQueue.DiscordID = user.DiscordID
	newQueue.Avatar = user.Avatar
	newQueue.Username = user.Username
	newQueue.UserID = user.DBID
	
	createdQueue, err := data.CreateQueue(newQueue)
	if err != nil {
		http.Error(w, "Error creating queue", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(createdQueue)
}


func GetQueueByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        http.Error(w, "Invalid queue ID", http.StatusBadRequest)
        return
    }

	var queue *models.Queue

	// if there is query param name then get queue by name
	name := r.URL.Query().Get("name")
	if name != "" {
		queue, err = data.GetQueueByName(name)
	} else {
		queue, err = data.GetQueueByID(uint(id))
	}

	if err != nil {
		http.Error(w, "Error retrieving queue", http.StatusInternalServerError)
		return
	}

	if queue == nil {
		http.Error(w, "Queue not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(queue)
}

func UpdateQueueHandler(w http.ResponseWriter, r *http.Request) {
	var updatedQueue *models.Queue
	err := json.NewDecoder(r.Body).Decode(&updatedQueue)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedQueue, err = data.UpdateQueue(*updatedQueue)
	if err != nil {
		http.Error(w, "Error updating queue", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedQueue)
}

func DeleteQueueHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid queue ID", http.StatusBadRequest)
		return
	}

	err = data.DeleteQueue(uint(id))
	if err != nil {
		http.Error(w, "Error deleting queue", http.StatusInternalServerError)
		return
	}
	
	response := map[string]string{
        "status": "success",
        "message": "Queue deleted",
        "queueId": idStr,
    }

	json.NewEncoder(w).Encode(response)
}

func ClearQueueByIDHandler(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid queue ID", http.StatusBadRequest)
		return
	}

	err = data.ClearQueueByID(uint(id))
	if err != nil {
		http.Error(w, "Error clearing queue", http.StatusInternalServerError)
		return
	}

	response := map[string]string{
        "status": "success",
        "message": "Queue cleared",
        "queueId": idStr,
    }

	json.NewEncoder(w).Encode(response)
}
