package main

import (
	"fmt"
	"math/rand"
	"sync"
)

// Player represents a player in the game
type Player struct {
	Name  string
	Score int
	Ball  chan Ball
}

// Ball represents the ball in the game
type Ball struct {
	Score int
}

// Game represents the game
type Game struct {
	Players []*Player
	Mutex   sync.Mutex
}

// NewPlayer creates a new player
func NewPlayer(name string) *Player {
	return &Player{
		Name:  name,
		Score: 0,
		Ball:  make(chan Ball),
	}
}

// NewGame creates a new game
func NewGame(players []*Player) *Game {
	return &Game{
		Players: players,
		Mutex:   sync.Mutex{},
	}
}

// Play plays the game
func (g *Game) Play() {
	// Create the ball
	ball := Ball{Score: 0}

	// Start passing the ball between players
	for {
		// Select a random player to receive the ball
		player := g.Players[rand.Intn(len(g.Players))]

		// Send the ball to the player
		player.Ball <- ball

		// Update the player's score
		g.updatePlayerScore(player)

		// Check if the game is over
		if g.isGameOver() {
			break
		}
	}
}

// updatePlayerScore updates a player's score
func (g *Game) updatePlayerScore(player *Player) {
	g.Mutex.Lock()
	player.Score++
	g.Mutex.Unlock()
}

// isGameOver checks if the game is over
func (g *Game) isGameOver() bool {
	g.Mutex.Lock()
	for _, player := range g.Players {
		if player.Score >= 10 {
			g.Mutex.Unlock()
			return true
		}
	}
	g.Mutex.Unlock()
	return false
}

// ReceiveBall receives the ball and updates the player's score
func (p *Player) ReceiveBall() {
	for {
		ball := <-p.Ball
		fmt.Printf("%s received the ball with score %d\n", p.Name, ball.Score)
	}
}

func main() {
	// Create players
	player1 := NewPlayer("Player 1")
	player2 := NewPlayer("Player 2")
	player3 := NewPlayer("Player 3")

	// Create the game
	game := NewGame([]*Player{player1, player2, player3})

	// Start receiving the ball for each player
	go player1.ReceiveBall()
	go player2.ReceiveBall()
	go player3.ReceiveBall()

	// Start the game
	game.Play()
}
