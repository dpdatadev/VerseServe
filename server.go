package main

/*
	TODO, add KJV.db (Sqlite)
	TODO, add CLI functionality akin to https://github.com/dtjm/bible
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
)

type ReferenceDTO struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Text    string `json:"text"`
}

func init() {
	log.SetPrefix("[=]VERSE SERVE[=]")
	log.SetFlags(0)
	localLog := GetFileWrite("executions.log")
	log.SetOutput(localLog)
	log.Printf("SERVER START: %s\n", time.Now().Local().String())
}

// Return files for Logging or dumping
func GetFileWrite(fileName string) *os.File {
	if fileName == "" {
		log.Fatalf("errors.New(\"\"): %v\n", errors.New("WRITE FILE ERROR"))
		return nil
	}

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("errors.New(\"\"): %v\n", err)
		return nil
	}

	return file
}

func WriteHttpJson(v any, w http.ResponseWriter) {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println("failed to serialize response:", err)
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.Write(b)
}

func GetReference(r *http.Request) string {
	ref := r.URL.Query().Get("ref")
	if ref == "" {
		missingMsg := "missing reference argument"
		log.Println(missingMsg)
		return errors.New(missingMsg).Error()
	}

	log.Println("parsing..", ref)

	return ref
}

func VerseServeHTTP(w http.ResponseWriter, r *http.Request) {

	ref := GetReference(r)
	if ref == "" {
		http.Error(w, "missing reference argument", http.StatusBadRequest)
		return
	}

	// Initialize an empty GoBible object
	b := NewGoBible()

	// Load a GoBible formatted JSON file
	b.Load("data/KJV.json")

	log.Println("parsing BIBLE REF..", ref)
	verses, err := b.ParseReference(ref)
	if err != nil {
		log.Println("failed to parse reference:", err)
		http.Error(w, "invalid reference format", http.StatusBadRequest)
		return
	}

	verseMap := make(map[string][]ReferenceDTO)

	verseList := make([]ReferenceDTO, len(verses))
	for i, v := range verses {
		verseList[i] = ReferenceDTO{
			Book:    v.Book,
			Chapter: v.Chapter,
			Verse:   v.Verse,
			Text:    v.VerseRef.Text,
		}
	}

	verseMap["Verses"] = verseList

	WriteHttpJson(verseMap, w)
}

func main() {
	// example usage: http://127.0.1:7777/verse?ref=Genesis%201:1-3
	http.HandleFunc("/verse", VerseServeHTTP)
	// health check endpoint http://127.0.1:7777/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	fmt.Println("Starting [/health] & [/verse] endpoints on :7777...")
	if err := http.ListenAndServe(":7777", nil); err != nil {
		log.Fatalf("failed to start server: %v\n", err)
	}
}
