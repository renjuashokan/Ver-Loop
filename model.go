package main

import (
	"database/sql"
	"net/http"
)

const (
	MAX_WORDS_IN_SENTENCE = 15
	MAX_SENTENCE_IN_PARA  = 10
	MAX_PARA_IN_STORY     = 7
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
	Id              int64
	Title           string
	CurrentSentance string
	DBPtr           datastore
	Body            StoryBody
	WordCount       int16
	SentenceCount   int16
	ParagraphsCount int16
}

type Sentence struct {
	Words string
}

type Paragraphs struct {
	Sentences []string `json:"sentences,omitempty"`
}

type StoryBody struct {
	Paragraphs []Paragraphs `json:"paragraphs,omitempty"`
}

type RespSingleStory struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	CreatedTime string `json:"created_at"`
	UpdateTime  string `json:"updated_at"`
	Paragraphs  string `json:"paragraphs"`
}
type incomingMsg struct {
	Word string `json:"word"`
}

type datastore struct {
	HTTPClient http.Client
	config     readConfig
	psqlDB     *sql.DB
}

type readConfig struct {
	ConfigFile       string
	LogLevel         string
	DeploymentMode   string
	AllowCrossOrigin bool
	CharSet          string
}
