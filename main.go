package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type card struct {
	Suit  string `json:"suit"`
	Value string `json:"value"`
	Code  string `json:"code"`
}

type deck struct {
	id       string
	shuffled bool
	cards    []string
}

func generate_full_cards() ([]string, map[string]card) {
	suits := [...]string{"SPADES", "DIAMONDS", "CLUBS", "HEARTS"}
	values := [...]string{"ACE", "2", "3", "4", "5", "6", "7", "8", "9", "10", "JACK", "QUEEN", "KING"}
	var cards []string
	detailed_cards := make(map[string]card)

	for i := 0; i < len(suits); i++ {
		for j := 0; j < len(values); j++ {
			card_key := string(values[j][0]) + string(suits[i][0])
			cards = append(cards, card_key)
			detailed_cards[card_key] = card{Suit: suits[i], Value: values[j], Code: card_key}
		}
	}
	return cards, detailed_cards
}

func generate_unique_deck_id(decks map[string]deck) string {
	for {
		id := uuid.New()
		_, is_already_exists := decks[id.String()]
		if !is_already_exists {
			return id.String()
		}
	}
}

func create_new_deck(decks map[string]deck, is_shuffled bool, wanted_cards []string, all_cards []string) deck {
	var cards []string

	if len(wanted_cards) > 0 {
		cards = wanted_cards
	} else {
		cards = make([]string, len(all_cards))
		copy(cards, all_cards)
	}

	if is_shuffled {
		rand.Shuffle(len(cards), func(i, j int) {
			cards[i], cards[j] = cards[j], cards[i]
		})
	}

	deck := deck{id: generate_unique_deck_id(decks), cards: cards, shuffled: is_shuffled}
	decks[deck.id] = deck
	return deck
}

func parse_cards_from_query(raw_cards string, detailed_cards map[string]card) ([]string, error) {
	already_parsed_cards := make(map[string]bool)

	split_fn := func(c rune) bool {
		return c == ','
	}
	card_names := strings.FieldsFunc(raw_cards, split_fn)

	for i := 0; i < len(card_names); i++ {
		_, is_present := detailed_cards[card_names[i]]
		if !is_present {
			return nil, errors.New("There is no such card as " + card_names[i])
		}
		_, is_already_used := already_parsed_cards[card_names[i]]
		if is_already_used {
			return nil, errors.New("This card occured twice in a desired deck: " + card_names[i])
		}
		already_parsed_cards[card_names[i]] = true
	}

	return card_names, nil
}

func main() {
	decks := map[string]deck{} // A storage variable. Not sure about concurrency though.
	// all_cards is used to quickly create new full deck by copying this variable
	// detailed_cards is used to actually build full-text card objects for client
	all_cards, detailed_cards := generate_full_cards()

	router := gin.Default()

	router.POST("/decks", func(c *gin.Context) {
		is_shuffled := false
		if c.DefaultQuery("shuffled", "false") == "true" {
			is_shuffled = true
		}

		wanted_cards, err := parse_cards_from_query(c.DefaultQuery("cards", ""), detailed_cards)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		deck := create_new_deck(decks, is_shuffled, wanted_cards, all_cards)
		c.JSON(http.StatusOK, gin.H{"deck_id": deck.id, "shuffled": deck.shuffled, "remaining": len(deck.cards)})
	})

	router.GET("/decks/:deck_id", func(c *gin.Context) {
		deck_id := c.Param("deck_id")
		deck, is_exists := decks[deck_id]
		if !is_exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "There is no deck with id " + deck_id})
			return
		}

		cards := []card{}
		for i := 0; i < len(deck.cards); i++ {
			cards = append(cards, detailed_cards[deck.cards[i]])
		}

		c.JSON(http.StatusOK, gin.H{"deck_id": deck.id, "shuffled": deck.shuffled, "remaining": len(deck.cards),
			"cards": cards})
	})

	router.POST("/decks/:deck_id/draw", func(c *gin.Context) {
		draw_count, err := strconv.Atoi(c.DefaultQuery("count", "1"))
		if err != nil || draw_count < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "'Count' parameter should be a positive integer"})
			return
		}

		deck_id := c.Param("deck_id")
		deck, is_exists := decks[deck_id]
		if !is_exists {
			c.JSON(http.StatusNotFound, gin.H{"error": "There is no deck with id " + deck_id})
			return
		}

		if len(deck.cards) < draw_count {
			c.JSON(http.StatusBadRequest, gin.H{"error": "There are not enough cards to do the draw in deck" + deck_id})
			return
		}

		cards_to_return := deck.cards[:draw_count]

		var cards []card
		for i := 0; i < len(cards_to_return); i++ {
			cards = append(cards, detailed_cards[deck.cards[i]])
		}

		deck.cards = deck.cards[draw_count:]
		decks[deck_id] = deck

		c.JSON(http.StatusOK, gin.H{"cards": cards})
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
