package game

import (
	"github.com/ratel-online/server/uno/card"
	"github.com/ratel-online/server/uno/card/color"
)

type playerController struct {
	player Player
	hand   *Hand
}

func newPlayerController(player Player) *playerController {
	return &playerController{
		player: player,
		hand:   NewHand(),
	}
}

func (c *playerController) AddCards(cards []card.Card) {
	c.hand.AddCards(cards)
	c.player.NotifyCardsDrawn(cards)
}

func (c *playerController) Hand() []card.Card {
	return c.hand.Cards()
}

func (c *playerController) Name() string {
	return c.player.NickName()
}

func (c *playerController) ID() int64 {
	return c.player.PlayerID()
}
func (c *playerController) NoCards() bool {
	return c.hand.Empty()
}

func (c *playerController) PickColor(gameState State) color.Color {
	return c.player.PickColor(gameState)
}

func (c *playerController) Play(gameState State, deck *Deck) (selectedCard card.Card, noMatch bool) {
	playableCards := c.hand.PlayableCards(gameState.LastPlayedCard)
	if len(playableCards) == 0 {
		c.player.NotifyNoMatchingCardsInHand(gameState.LastPlayedCard, gameState.CurrentPlayerHand)
		selectedCard = c.tryTopDecking(gameState, deck)
		return selectedCard, true
	}

	for {
		selectedCard = c.player.Play(playableCards, gameState)
		c.hand.RemoveCard(selectedCard)
		return selectedCard, false
	}
}

func (c *playerController) tryTopDecking(gameState State, deck *Deck) card.Card {
	extraCard := deck.DrawOne()
	c.AddCards([]card.Card{extraCard})
	if Playable(extraCard, gameState.LastPlayedCard) {
		c.hand.RemoveCard(extraCard)
		return extraCard
	}
	return nil
}
