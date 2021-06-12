package main

import (
	"database/sql"
	"net/http"
)

type addRespose struct {
	Id              int64  `json:"id"`
	Title           string `json:"title"`
	CurrentSentance string `json:"current_sentence"`
}

// TODO: use protobuf enum
type titleStatus struct {
}

type story struct {
	Id                int64
	Title             string
	CurrentSentance   string
	CurrentParagraph  []string
	Paragraphs        [][]string
	SentanceWordCount int16
}

type incomingMsg struct {
	Word string `json:"word"`
}

type datastore struct {
	HTTPClient http.Client
	config     readConfig
	psqlDB     *sql.DB
	id         int64
}

type readConfig struct {
	ConfigFile       string
	LogLevel         string
	DeploymentMode   string
	AllowCrossOrigin bool
	CharSet          string
}
