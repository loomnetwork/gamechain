package battleground

import (
	"fmt"

	"github.com/loomnetwork/zombie_battleground/types/zb"
)

func getHeroInfoById(heroInfoList []*zb.HeroInfo, heroId int64) *zb.HeroInfo {
	for _, heroInfo := range heroInfoList {
		if heroInfo.HeroId == heroId {
			return heroInfo
		}
	}
	return nil
}

func validateDeckHero(heroInfoList []*zb.HeroInfo, deckHero int64) error {
	for _, heroInfo := range heroInfoList {
		if heroInfo.HeroId == deckHero {
			return nil
		}
	}

	return fmt.Errorf("hero: %d cannot be part of deck, since it is not owned by User.", deckHero)

}
