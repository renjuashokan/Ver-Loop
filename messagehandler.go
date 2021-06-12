package main

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

func (instance *story) addWord(word string) (int, addRespose) {
	words := strings.Fields(word)
	if len(words) > 1 {
		return http.StatusBadRequest, addRespose{}
	}

	log.SetFormatter(&log.TextFormatter{TimestampFormat: "2006-01-02 15:04:05", FullTimestamp: true})

	// set title if it is not set
	if instance.Title == "" {
		fmt.Println("title is empty")
		return instance.setTitle(word)
	} else if len(strings.Fields(instance.Title)) <= 1 {
		instance.Title = instance.Title + " " + word
		return http.StatusOK, addRespose{Id: instance.Id,
			Title: instance.Title}
	} else {
		// title is set, check for body

		if len(instance.CurrentSentance) == 0 {
			instance.SentanceWordCount = 1
			instance.CurrentSentance = word
		} else if instance.SentanceWordCount < 15 {
			instance.CurrentSentance = instance.CurrentSentance + " " + word
			instance.SentanceWordCount++
		} else if instance.SentanceWordCount == 15 {

			// max length reached for current sentance.
			// check we can add to current paragraph
			if len(instance.CurrentParagraph) < 10 {
				instance.CurrentParagraph = append(instance.CurrentParagraph, instance.CurrentSentance)
				log.Info("adding to current paragraph")

			} else if len(instance.CurrentParagraph) == 10 {
				// check we can add
				if len(instance.Paragraphs) < 7 {
					instance.Paragraphs = append(instance.Paragraphs, instance.CurrentParagraph)
					instance.CurrentParagraph = nil
					log.Debug("Adding to paragraphs")
				} else if len(instance.Paragraphs) == 7 {
					log.Info("Creating new story!!")
					return instance.setTitle(word)
				}
			}

			instance.CurrentSentance = word
			instance.SentanceWordCount = 1
		}
		log.WithFields(log.Fields{
			"id":                               instance.Id,
			"current sentace word count":       instance.SentanceWordCount,
			"current paragraph sentence count": len(instance.CurrentParagraph),
			"total number of paragraphs":       len(instance.Paragraphs),
		}).Info("server status")

		return http.StatusOK, addRespose{Id: instance.Id,
			Title: instance.Title, CurrentSentance: instance.CurrentSentance}
	}

}

func (instance *story) setTitle(word string) (int, addRespose) {
	instance.Title = word
	instance.Id++
	instance.CurrentParagraph = nil
	instance.CurrentSentance = ""
	instance.Paragraphs = nil
	log.Info("Creating new story!")
	return http.StatusCreated, addRespose{Id: instance.Id,
		Title: instance.Title}
}
