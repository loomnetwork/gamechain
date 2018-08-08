package battleground

import (
	"fmt"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

func getHeroById(heroList []*zb.Hero, heroId int64) *zb.Hero {
	for _, hero := range heroList {
		if hero.HeroId == heroId {
			return hero
		}
	}
	return nil
}

func validateDeckHero(heroList []*zb.Hero, deckHero int64) error {
	for _, hero := range heroList {
		if hero.HeroId == deckHero {
			return nil
		}
	}

	return fmt.Errorf("hero: %d cannot be part of deck, since it is not owned by User.", deckHero)

}
