# Foreword
First of all, I have never written anything in `go` up until yesterday. I spent 4 hours writing this assignment, and majority of this time went into researching of how to do primitive things in `go`. I basically run out of time, and that's why:
- there are no tests
- decks are stored in-memory instead of db, I am not even sure if it is safe to do so concurrency-wise. But again, I had no time left since I never used `go` before


# To run API:

`docker build -t cards . && docker run --rm -it -p 8080:8080 cards`

After that you'll have an API running at `127.0.0.1:8080`


# API endoints:

**POST /decks** - create a new deck
Optional query paras: 
- `shuffled` - tells if cards should be shuffled, default - false.
- `cards` - tells exact cards that should be used in the new deck, comma-separated, by default - use full deck.

**Sample:** `POST /decks?shuffled=true&cards=2C,4H,QH` - *create a new deck based on 3 cards in shuffled order*
 
 ***

**GET /decks/:deck_id:** - get deck details by id

**Sample:** `GET /decks/696aefad-2eb4-4041-9a2c-ef15660c2084` - *gets deck with id 696aefad-2eb4-4041-9a2c-ef15660c2084*

***

**POST /decks/:deck_id:/draw** - draw card(s) from the deck and return it
Optional query paras: 
- `count` - tells amount of cards to draw, default - 1.

**Sample:** `POST /decks/696aefad-2eb4-4041-9a2c-ef15660c2084/draw?count=5` - *draw 5 cards from deck with id 696aefad-2eb4-4041-9a2c-ef15660c2084*