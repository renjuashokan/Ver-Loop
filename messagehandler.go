package main

import (
	"net/http"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func (instance *story) addWord(word string) (int, addRespose) {
	words := strings.Fields(word)
	if len(words) > 1 {
		return http.StatusBadRequest, addRespose{}
	}

	// set title if it is not set
	if instance.Title == "" {
		log.Debug("Creating title")
		return instance.setTitle(word)
	} else if len(strings.Fields(instance.Title)) <= 1 {
		instance.Title = instance.Title + " " + word
		instance.DBPtr.UpdateTitle(*instance)
		return http.StatusOK, addRespose{Id: instance.Id,
			Title: instance.Title}
	}

	return instance.UpdateBody(word)
}

func (instance *story) UpdateBody(word string) (int, addRespose) {

	log.Debug("Updating body")

	if len(instance.CurrentSentance) == 0 {
		instance.WordCount = 1

		tmp := make([]string, 1)
		tmp[0] = word
		tmp2 := make([]Paragraphs, 1)
		tmp2[0].Sentences = tmp
		instance.Body = StoryBody{tmp2}
		instance.CurrentSentance = word

		//instance.CurrentSentance = word
	} else if instance.WordCount < MAX_WORDS_IN_SENTENCE {
		instance.Body.Paragraphs[instance.ParagraphsCount].Sentences[instance.SentenceCount] += " " + word
		instance.CurrentSentance += " " + word
		instance.WordCount++
	} else if instance.WordCount == MAX_WORDS_IN_SENTENCE {
		// current sentence is full,
		// move to new sentence

		if instance.SentenceCount < (MAX_SENTENCE_IN_PARA - 1) {
			// create new sentence,
			tmp := make([]string, 1)
			instance.Body.Paragraphs[instance.ParagraphsCount].Sentences = append(instance.Body.Paragraphs[instance.ParagraphsCount].Sentences, tmp...)

			// update sentence count
			instance.SentenceCount++
			instance.Body.Paragraphs[instance.ParagraphsCount].Sentences[instance.SentenceCount] = word
			instance.WordCount = 1
			instance.CurrentSentance = word

		} else if instance.SentenceCount == (MAX_SENTENCE_IN_PARA - 1) {
			// current paragraph is full
			// move to new paragraph

			if instance.ParagraphsCount < (MAX_PARA_IN_STORY - 1) {
				log.Info("create new paragraph")

				tmp2 := make([]Paragraphs, 1)
				tmp := make([]string, 1)
				tmp2[0].Sentences = tmp
				instance.Body.Paragraphs = append(instance.Body.Paragraphs, tmp2...)

				instance.ParagraphsCount++

				instance.SentenceCount = 0
				instance.Body.Paragraphs[instance.ParagraphsCount].Sentences[instance.SentenceCount] = word
				instance.WordCount = 1
				instance.CurrentSentance = word
			} else if instance.ParagraphsCount == (MAX_PARA_IN_STORY - 1) {
				log.Info("Creating new story!!")
				return instance.setTitle(word)
			}
		}
	}

	instance.DBPtr.UpdateBody(*instance)
	log.WithFields(log.Fields{
		"id":                               instance.Id,
		"current sentace word count":       instance.WordCount,
		"current paragraph sentence count": instance.SentenceCount + 1,
		"total number of paragraphs":       instance.ParagraphsCount + 1,
	}).Debug("server status")
	return http.StatusOK, addRespose{Id: instance.Id,
		Title:           instance.Title,
		CurrentSentance: instance.CurrentSentance}
}

func (instance *story) setTitle(word string) (int, addRespose) {
	nxtID, err := instance.DBPtr.GetNextId()
	if err != nil {
		panic(errors.Wrap(err, "something failed"))
	}
	instance.Title = word
	instance.Id = nxtID
	instance.CurrentSentance = ""
	instance.WordCount, instance.SentenceCount, instance.ParagraphsCount = 0, 0, 0
	var body StoryBody
	instance.Body = body

	log.Info("Creating new story!")
	instance.DBPtr.CreateStory(*instance)
	return http.StatusCreated, addRespose{Id: instance.Id,
		Title: instance.Title}
}
