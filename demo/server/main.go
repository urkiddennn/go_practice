package main

import (
	"context"
	"encoding/json"

	//	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Note struct represents the note model
type Note struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title" bson:"title"`
	Content   string             `json:"content" bson:"content"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// Global variables
var (
	client     *mongo.Client
	collection *mongo.Collection
)

// Initialize MongoDB connection
func initDB() {
	// Set Stable API version
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	uri := "mongodb+srv://richardbanguiz:DJdyJxTLBE8nfoX9@cluster0.68wut.mongodb.net/note-deb?retryWrites=true&w=majority&appName=Cluster0"
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI).SetServerSelectionTimeout(10 * time.Second)

	var err error
	client, err = mongo.Connect(context.TODO(), opts)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Ping the server to verify connection
	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	collection = client.Database("notesdb").Collection("notes")
	log.Println("Successfully connected to MongoDB Atlas with Stable API!")
}

// CreateNote handler
func CreateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var note Note
	err := json.NewDecoder(r.Body).Decode(&note) // Fixed ¬e to &note
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.CreatedAt = time.Now()
	note.UpdatedAt = time.Now()

	result, err := collection.InsertOne(context.TODO(), note)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	note.ID = result.InsertedID.(primitive.ObjectID)
	json.NewEncoder(w).Encode(note)
}

// GetNotes handler
func GetNotes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var notes []Note
	ctx := context.TODO()
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var note Note
		if err := cursor.Decode(&note); err != nil { // Fixed ¬e to &note
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		notes = append(notes, note)
	}

	json.NewEncoder(w).Encode(notes)
}

// GetNote handler
func GetNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var note Note
	ctx := context.TODO()
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&note) // Fixed ¬e to &note
	if err != nil {
		http.Error(w, "Note not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(note)
}

// UpdateNote handler
func UpdateNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var note Note
	err = json.NewDecoder(r.Body).Decode(&note) // Fixed ¬e to &note
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	note.UpdatedAt = time.Now()

	update := bson.M{
		"$set": bson.M{
			"title":      note.Title,
			"content":    note.Content,
			"updated_at": note.UpdatedAt,
		},
	}

	ctx := context.TODO()
	err = collection.FindOneAndUpdate(ctx, bson.M{"_id": id}, update).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	note.ID = id
	json.NewEncoder(w).Encode(note)
}

// DeleteNote handler
func DeleteNote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	ctx := context.TODO()
	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Note deleted successfully"})
}

func main() {
	// Initialize database
	initDB()

	// Initialize router
	router := mux.NewRouter()

	// Define endpoints
	router.HandleFunc("/notes", CreateNote).Methods("POST")
	router.HandleFunc("/notes", GetNotes).Methods("GET")
	router.HandleFunc("/notes/{id}", GetNote).Methods("GET")
	router.HandleFunc("/notes/{id}", UpdateNote).Methods("PUT")
	router.HandleFunc("/notes/{id}", DeleteNote).Methods("DELETE")

	// Create HTTP server with timeout settings
	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// Channel to listen for OS signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Server starting on port 8080...")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-stop
	log.Println("Shutting down server...")

	// Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	// Disconnect MongoDB client
	if err := client.Disconnect(ctx); err != nil {
		log.Fatalf("MongoDB disconnect failed: %v", err)
	}

	log.Println("Server stopped gracefully")
}
