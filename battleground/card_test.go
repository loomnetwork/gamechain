package battleground

import (
	"encoding/csv"
	"io"
	"strconv"
	"strings"
	"testing"

	"github.com/loomnetwork/zombie_battleground/types/zb"
	"github.com/stretchr/testify/assert"
)

func TestCardToGenesisInit(t *testing.T) {
	cardlist, err := loadcards(strings.NewReader(cardCSV))
	assert.Nil(t, err)
	for _, card := range cardlist.Cards {
		assert.NotEmpty(t, card.Name)
	}
}

func loadcards(reader io.Reader) (*zb.CardList, error) {
	r := csv.NewReader(reader)
	var cards []*zb.Card
	var isHead = true
	var current int
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		// skip the first record
		if isHead {
			isHead = false
			continue
		}
		// if name is missing, we are still at the previous card
		if record[0] != "" {
			card, err := cardFrom(record)
			if err != nil {
				return nil, err
			}
			card.Id = int64(current + 1)
			cards = append(cards, card)
		}
		// convert effect
		effect := effectFrom(record[8:])
		if effect.Trigger != "" {
			cards[current].Effects = append(cards[current].Effects, effect)
		}

		current++
	}

	return &zb.CardList{
		Cards: cards,
	}, nil
}

func cardFrom(record []string) (*zb.Card, error) {
	var card zb.Card
	card.Name = record[0]
	card.Set = record[1]
	card.Rank = record[2]
	card.Type = record[3]
	i, _ := strconv.ParseInt(record[4], 10, 32)
	card.Damage = int32(i)
	i, _ = strconv.ParseInt(record[5], 10, 32)
	card.Health = int32(i)
	i, _ = strconv.ParseInt(record[6], 10, 32)
	card.Cost = int32(i)
	card.Ability = record[7]
	return &card, nil
}

func effectFrom(record []string) *zb.Effect {
	effect := zb.Effect{}
	effect.Trigger = record[0]
	effect.Effect = record[1]
	effect.Duration = record[2]
	effect.Target = record[3]
	return &effect
}

const cardCSV = `NAME,ELEMENT,RANK,TYPE,ATTACK,DEFENSE,GOO COST,ABILITY,trigger,Effect,duration,target,limit,13,39,5,19,12
Banshee,Air,Minion,Feral,2,1,2,Feral,entry,feral,permanent,self,,entry,feral,permanent,self,
Breezee,Air,Minion,Walker,1,1,1,-,death,attack_strength_buff,,friendly_selectable,,death,attack_strength_buff,,friendly_selectable,card.goo <= 2
Buffer,Air,Minion,Walker,1,1,2,Death: Give a friendly zombie +1 Attack,entry,draw_card(x),,,,post_turn(x),draw_card(x),one_move,,max_target = 3
Soothsayer,Air,Minion,Walker,1,1,2,Entry: Draw a card.,entry,lower_goo_cost,,friendly_random,,,lower_goo_cost,while_alive,friendly_random,target.defense >= 2
Wheezy,Air,Minion,Walker,1,2,2,Enter: Lower the goo cost of a random card in your hand by 1.,entry,return_to_hand,,friendly_selectable,,on_attacked,return_to_hand,one_turn,friendly_hand,target.type != heavy
Whiffer,Air,Minion,Walker,3,2,2,Entry: Return a friendly zombie to your hand.,post_turn(x),attack_strength_buff(x),,self,,pre_entry,attack_strength_buff(x),once,enemy_hand,if target.element == water
Whizpar,Air,Minion,Feral,1,1,2,Feral,,,,,,pre_attack,,,adjacent,if friendly.in_play == 1
Zhocker,Air,Minion,Walker,0,2,1,Becomes 2/2 at the end of your turn.,,,,,,after_entry(x),set_heavy,,all_except_self,if target.type == heavy
Bouncer,Air,Officer,Heavy,2,3,4,Heavy,entry,set_heavy,,self,,before_turn_end,summon_zombie,,enemy_random,if friendly_hero.life <= 10
Dragger,Air,Officer,Walker,1,2,3,Entry: Summon a zombie from your hand that costs 2 Goo or less.,entry,summon_zombie,,friendly_hand,card.goo <= 2,on_use,reduce_goo_cost(x),,enemy_all,if target.element == fire
Guzt,Air,Officer,Walker,3,3,3,"When this takes damage and does not die, return it to your hand.",on_attacked,return_to_hand,,self,,after_entry(3),draw_card,,enemy_selectable,x = unique_type_in_play
Pushhh,Air,Officer,Walker,3,3,3,Entry: Return a zombie to hand,,,,,,before_turn,deal_damage,,friendly_adjacent,current_goo
Two-rnado,Air,Officer,Walker,3,6,5,Entry: Draw 2 cards.,,,,,,on_enemy_killed,summon_enemy_zombie,,enemy_hero,if target.status == frozen
Zephyr,Air,Officer,Walker,3,3,4,Costs 1 Goo less for each Air zombie you have in play.,pre_entry,reduce_goo_cost(x),,self,,after_entry(1),attack_strength_boost,,friendly_all,
Draft,Air,Commander,Walker,4,5,5,Entry: Draw a card from the enemy deck. (This removes the card form his deck.),entry,draw_card,,enemy_hand,,,kill,,friendly_hero,
MonZoon,Air,Commander,Walker,6,6,9,Costs 1 Goo less for each Air zombie in your hand.,,,,,,,defense_boost,,target,
Stormmcaller,Air,Commander,Walker,4,5,6,Entry: damage 1 to 3 adjacent cards,entry,deal_damage,,adjacent,max_target = 3,,remove_heavy,,friendly_all_with_hero,
Zyclone,Air,Commander,Walker,2,2,6,Entry: Return all other zombies in the battleground.,entry,return_to_hand,,all_except_self,,,disable_hero_attack,,attacker,
Mind Flayer,Air,General,Walker,5,5,6,Entry: Select an enemy zombie to join your forces,entry,summon_enemy_zombie,,enemy_random,,,set_feral,,enemy,
ZeuZ,Air,General,Walker,8,7,10,Entry: Damage 3 all enemy zombies,entry,deal_damage,,enemy_all,,,gain_move(x),,enemy_deck,
Blocker,Earth,Minion,Heavy,0,3,1,Heavy,,,,,,,distract,,,
Crumbz,Earth,Minion,Heavy,0,2,2,"Heavy, Death: Draw a card.",death,draw_card,,,,,heal,,,
Hardy,Earth,Minion,Walker,3,3,3,Heavy,,,,,,,add_goo,,,
Pebble,Earth,Minion,Walker,1,1,1,Entry: Deal 1 damage to an enemy zombie with 2 Defense or higher.,entry,deal_damage,,enemy_selectable,target.defense >= 2,,freeze,,,
Protector,Earth,Minion,Heavy,1,2,2,"Heavy, Death: Give a random friendly zombie Heavy  ",death,set_heavy,,friendly_random,target.type != heavy,,lose_defense_on_attack_debuff,,,
Rockky,Earth,Minion,Walker,1,1,1,Attack: 1 additional damage for water zombies,pre_attack,attack_strength_boost,one_move,self,if target.element == water,,defense_boost(x),,,
Shale,Earth,Minion,Feral,2,3,3,Feral,,,,,,,add_swing,,,
Slab,Earth,Minion,Walker,3,4,3,-,,,,,,,summon_special_zombie(1),,,
Defender,Earth,Officer,Walker,3,3,4,Death: Give a random friendly zombie +2 Defense.,,,,,,,summon_duplicate,,,
Golem,Earth,Officer,Heavy,2,6,4,Heavy,,,,,,,summon_special_zombie(2),,,
Groundy,Earth,Officer,Walker,5,5,5,"Entry: If this is your only zombie in the battleground, gain Heavy.",entry,set_heavy,,self,if friendly.in_play == 1,,random_heal(10),,,
Spiker,Earth,Officer,Heavy,2,3,3,"Heavy, Attack: Deal +1 damage to other Heavy Zombies",pre_attack,attack_strength_boost,one_move,enemy_selectable,if target.type == heavy,,summon_special_zombie(3),,,
Tiny,Earth,Officer,Heavy,0,7,4,Heavy,,,,,,,z_virus_ability,,,
Walley,Earth,Officer,Walker,2,2,4,Make adjacent zombies Heavy,entry,set_heavy,while_alive,friendly_adjacent,,,reduce_goo,,,
Bolderr,Earth,Commander,Walker,3,5,5,Entry: attack target for 3 damage,entry,kill,,enemy_selectable,if target.type == heavy,,zeeter_ability,,,
Earthshaker,Earth,Commander,Walker,4,4,5,Entry: Destroy another Heavy zombie,,,,,,,set_goo(0),,,
IgneouZ,Earth,Commander,Heavy,3,3,4,"Entry: If you have 10 life or less, this gets +2 Defense.",entry,defense_boost,,self,if friendly_hero.life <= 10,,return_to_enemy_hand,,,
Pyrite,Earth,Commander,Heavy,0,8,2,"Heavy, Loses Heavy after 2 turns, but gains +1 Attack.",after_entry(x),remove_heavy,,self,,,discard_from_hand,,,
Gaea,Earth,General,Walker,4,5,10,Entry: Give all your Earth zombies in the battleground +2 Attack and +2 Defense.,,,,,,,discard_from_top_of_deck,,,
Mountain,Earth,General,Heavy,8,8,10,"Heavy, Attack: damage 2 adjacent targets for 4 dmg",,,,,,,discard_from_hand(enemy.hand.count - friendly.hand.count),,,
BlaZter,Fire,Minion,Walker,1,1,2,Entry: Deal 1 damage to an enemy zombie.,,,,,,,,,,
BurZt,Fire,Minion,Feral,4,1,3,"Feral, Cannot attack the enemy hero this turn",entry,disable_hero_attack,one_turn,enemy_hero,,,,,,
Ember,Fire,Minion,Walker,3,2,2,-,,,,,,,,,,
Firecaller,Fire,Minion,Walker,1,1,2,Entry: Draw a card if there is a damaged zombie in the battleground.,entry,draw_card,,friendly_hand,,,,,,
Firewall,Fire,Minion,Heavy,2,1,2,Heavy,,,,,,,,,,
Pyromaz,Fire,Minion,Walker,1,1,1,Attack: 1 additional damage to life zombie,,,,,,,,,,
RabieZ,Fire,Minion,Feral,1,1,2,"Feral, Death: Give a random friendly minion Feral.",,,,,,,,,,
Sparky,Fire,Minion,Feral,2,1,1,"Feral, Deal 1 damage to this at the end of your turn.",,,,,,,,,,
Alpha,Fire,Officer,Walker,3,4,5,Entry: All other Fire zombies in the battleground gain Feral.,entry,set_feral,while_alive,friendly_all,if target.element == fire,,,,,
Burrrnn,Fire,Officer,Feral,2,2,3,Feral,,,,,,,,,,
Enrager,Fire,Officer,Feral,3,3,5,"Feral, Entry: If you have 10 life or less this gets +2 Attack.",,,,,,,,,,
Rager,Fire,Officer,Walker,4,4,4,"Entry: If this is your only zombie in the battleground, gain Feral.",,,,,,,,,,
Werezomb,Fire,Officer,Walker,1,1,3,Entry: Give a friendly zombie Feral.,,,,,,,,,,
Zlinger,Fire,Officer,Feral,4,4,3,"Feral, Entry: Deal 2 damage to your hero.",entry,deal_damage,,friendly_hero,,,,,,
Cynderman,Fire,Commander,Walker,4,3,5,Enter: Damage a target for 2 damage,,,,,,,,,,
Fire-Maw,Fire,Commander,Feral,3,3,6,"Feral, Flash (attacks twice) ",entry,gain_move(x),once,self,,,,,,
Volcan,Fire,Commander,Walker,5,5,6,"Cannot Attack. At the end of your turn, deal 5 damage to a random enemy.",before_turn_end,deal_damage,,enemy_random,,,,,,
Zhampion,Fire,Commander,Feral,5,2,5,"Feral, Entry: If this is your only zombie in the battleground,  it gains +2 Attack and +2 Defense.",,,,,,,,,,
Cerberus,Fire,General,Feral,6,6,8,"Feral, Flash x 2",,,,,,,,,,
Gargantua,Fire,General,Heavy,6,6,8,Entry: Damage all enemy zombies (including heros) for 2 Damage,,,,,,,,,,
Bat,Item,Minion,Item,NA,NA,3,Deal 2 damage to a zombie and Distract it.,on_use,distract,,target,,,,,,
Molotov,Item,Minion,Item,3,NA,4,Deal 3 damage to a zombie and its adjacent zombies,,,,,,,,,,
Stapler,Item,Minion,Item,NA,NA,4,Heals a zombie for 4 HP,on_use,heal,,friendly_selectable,,,,,,
Super Goo Serum,Item,Minion,Item,NA,NA,4,Give a friendly zombie +3/+3.,,,,,,,,,,
Supply Drop,Item,Minion,Item,NA,NA,3,Each player puts a random zombie card from their deck to play.,on_use,summon_zombie,,enemy_random,,,,,,
Torch,Item,Minion,Item,NA,NA,2,Entry: Destroy a zombie that costs 3 Goo or less.,,,,,,,,,,
Whistle,Item,Minion,Item,NA,NA,0,Draw a card.,,,,,,,,,,
Beaker,Item,Officer,Item,NA,NA,0,Get +1 Goo for this turn only.,on_use,add_goo,once,,,,,,,
Fire Extinguisher,Item,Officer,Item,NA,NA,4,Entry: Disable (freeze) all enemy zombies.,on_use,freeze,,enemy_all,,,,,,
Junk Spear,Item,Officer,Item,3,1,4,Loses 1 Defense every time you attack. Entry: Gets +1 Defense for each zombie type you have in the battleground. ,on_use,lose_defense_on_attack_debuff,one_turn,friendly_all,,,,,,
Lawnmower,Item,Officer,Item,NA,NA,6,Entry: Deal 2 damage to all enemy zombies.,on_use,defense_boost(x),,,x = unique_type_in_play,,,,,
Nail Bomb,Item,Officer,Item,10,NA,4,Damage target and 2 adjacent zombies for 5 damage.,,,,,,,,,,
Shovel,Item,Officer,Item,NA,NA,3,Entry: Deal 4 damage to an enemy OR restore 5 health to a friendly character.,,,,,,,,,,
Zed Kit,Item,Officer,Item,NA,NA,2,Restore defense to each zombie you have for each zombie element type you have in the battleground.,,,,,,,,,,
Fresh Meat,Item,Commander,Item,NA,NA,5,Entry: Give all enemy zombies -3 attack until the end of your turn.,,,,,,,,,,
Goo Bottles,Item,Commander,Item,NA,NA,0,adds 2 full Goo vials,on_use,add_goo,once,,,,,,,
Leash,Item,Commander,Item,NA,NA,8,Entry: Gain control of an enemy zombie.,,,,,,,,,,
Shopping Cart,Item,Commander,Item,NA,NA,6,"Entry: All your zombies gain ""Swing"" (Swing: when attacking, deal damage to the target, as well as any adjacent zombies.)",on_use,add_swing,,friendly_all,,,,,,
Bulldozer,Item,General,Item,NA,NA,10,Entry: Destroy all zombies in the battleground.,on_use,kill,,enemy_all,,,,,,
Chainsaw,Item,General,Item,5,3,2,Loses 1 HP on attack,,,,,,,,,,
Amber,Life,Minion,Walker,0,3,2,Becomes a 3/3 at the end of your next turn.,,,,,,,,,,
Azuraz,Life,Minion,Walker,1,1,1,Attack: 1 additional damage to life zombies,,,,,,,,,,
Bark,Life,Minion,Walker,2,2,3,-,,,,,,,,,,
Bloomer,Life,Minion,Walker,1,1,2,Entry: Draw a card if you have a life zombie in the battleground.,,,,,,,,,,
Grower,Life,Minion,Walker,1,1,2,Entry: Your hero gains 2 life.,,,,,,,,,,
Medic,Life,Minion,Walker,1,1,2,Entry: Restore 2 life to a zombie.,,,,,,,,,,
WiZp,Life,Minion,Walker,1,1,2,Death: gives 1 goo for next turn,,,,,,,,,,
EverlaZting,Life,Officer,Walker,2,3,4,Death: Put this card back into your deck.,,,,,,,,,,
Healz,Life,Officer,Walker,2,3,4,"Entry: Gain 3 life. If you have 10 life or less, gain 5 life instead.",,,,,,,,,,
Keeper,Life,Officer,Walker,1,3,3,Entry: Summon a 1/1 Feral zombie.,entry,summon_special_zombie(1),,,,,,,,
Puffer,Life,Officer,Walker,2,2,3,Entry: add 1 attack to all life zombies on the Battleground,,,,,,,,,,
Sapper,Life,Officer,Item,1,4,5,Attack: Gain 1 life for each damage this deals.,,,,,,,,,,
Shroom,Life,Officer,Walker,4,2,4,Entry: 2 damage to target,,,,,,,,,,
Weed,Life,Officer,Walker,1,2,3,Entry: Place another copy of this minion in the battleground.,entry,summon_duplicate,,self,,,,,,
Blight,Life,Commander,Walker,0,6,6,"At the end of 3 turns, place 2 4/4 copies of this zombie in the battleground and discard this.",after_entry(3),summon_special_zombie(2),,,,,,,,
Rainz,Life,Commander,Walker,3,4,6,Entry: Restore 10 HP randomly split among all zombies in the battleground and your hero.,entry,random_heal(10),,friendly_all_with_hero,,,,,,
Vindrom,Life,Commander,Walker,5,5,6,Attack: disable enemy for 1 turn (entangle vines),,,,,,,,,,
Zeeder,Life,Commander,Walker,2,5,5,"At the end of your turn, summon a 0/2 zombie with Heavy.",,,,,,,,,,
Shammann,Life,General,Walker,5,6,8,Turn: summons a 1/1 zombie minion start of each turn,before_turn,summon_special_zombie(3),,,,,,,,
Yggdrazil,Life,General,Walker,4,4,11,Entry: Revive all life zombies that have died this game.,entry,,,,,,,,,
Z-Virus,Life,General,Walker,NA,NA,7,"Devour friendly zombies, and replace them with a zombie with their combined attack and defense.",entry,z_virus_ability,,,,,,,,
Ectoplasm,Toxic,Minion,Walker,3,2,1,Entry: Lose 1 Goo.,,,,,,,,,,
Germ,Toxic,Minion,Walker,2,1,1,-,,,,,,,,,,
Poizom,Toxic,Minion,Walker,1,1,1,Attack: 1 additional damage to earth zombies,,,,,,,,,,
Spikez,Toxic,Minion,Heavy,2,2,3,Has +1 Attack on your opponent's turn.,,,,,,,,,,
Wazte,Toxic,Minion,Heavy,3,3,2,"Heavy, Entry: Lose 1 Goo ",entry,reduce_goo,,,,,,,,
Zcavenger,Toxic,Minion,Walker,1,1,2,"If this kills a zombie, draw a card.",on_enemy_killed,draw_card,,,,,,,,
Zlimey,Toxic,Minion,Walker,3,1,1,Entry: Deal 2 damage to your hero.,,,,,,,,,,
Zpitter,Toxic,Minion,Walker,2,2,3,Entry: Deal 1 damage randomly enemy,,,,,,,,,,
Azzazzin,Toxic,Officer,Feral,5,1,3,"Feral, Can only attack zombies.",,,,,,,,,,
Boomer,Toxic,Officer,Walker,3,3,4,Death: Adjacent zombies get +1/+1.,,,,,,,,,,
Ghoul,Toxic,Officer,Walker,3,2,2,Attack: looses 1 DMG (not HP),,,,,,,,,,
Polluter,Toxic,Officer,Walker,3,3,4,Death: Gain 1 Goo.,,,,,,,,,,
RelentleZZ,Toxic,Officer,Walker,3,6,3,Entry: Deal 2 damage to this zombie.,,,,,,,,,,
Zlopper,Toxic,Officer,Walker,1,4,3,Gains +1 Attack whenever you play a Toxic zombie.,,,,,,,,,,
Hazmaz,Toxic,Commander,Feral,3,3,4,Feral,,,,,,,,,,
Zeeter,Toxic,Commander,Walker,2,3,5,Entry: Destroy a friendly zombie and gain its base attack and damage.,entry,zeeter_ability,,,,,,,,
Zludge,Toxic,Commander,Walker,4,4,5,"Rage: Whenever this gets damaged, it gets +2 Attack.",,,,,,,,,,
Zteroid,Toxic,Commander,Walker,5,4,6,Death: Give all of your zombies in the battleground +2 Attack until end of turn.,,,,,,,,,,
Cherno-bill,Toxic,General,Heavy,8,8,7,"Heavy, Death: Damage 2 to ALL zombies and heroes",,,,,,,,,,
GooZilla,Toxic,General,Walker,1,1,0,Entry: Spend all your Goo. This has Attack and Defense equal to the spent Goo.,entry,set_goo(0),,,,,,,,
FroZen,Water,Minion,Walker,1,2,2,Death: Draw a card if there is a frozen zombie in the battleground.,entry,deal_damage,,enemy_hero,current_goo,,,,,
HoZer,Water,Minion,Walker,1,1,2,Entry: Gets +1 Attack if you have a Water Zombie in the battleground.,,,,,,,,,,
Izze,Water,Minion,Walker,1,1,1,Attack: disable enemy for 1 turn (freeze),,,,,,,,,,
Slider,Water,Minion,Walker,1,2,2,Entry: Distract an enemy zombie.,,,,,,,,,,
Zhatterer,Water,Minion,Walker,3,3,5,Entry: Destroy a frozen zombie.,entry,kill,,enemy_selectable,if target.status == frozen,,,,,
Znowman,Water,Minion,Heavy,0,5,4,"Heavy, Zombies that attack this become frozen.",on_attacked,freeze,,attacker,,,,,,
Znowy,Water,Minion,Walker,1,2,1,-,,,,,,,,,,
Ztalagmite,Water,Minion,Feral,3,3,3,Feral; Entry: Deal 1 damage to your hero.,,,,,,,,,,
Igloo,Water,Officer,Heavy,4,5,5,Heavy,,,,,,,,,,
Izicle,Water,Officer,Walker,3,5,5,Entry: Deal 2 damage to an enemy if a Water Zombie on the Battleground.,,,,,,,,,,
Jetter,Water,Officer,Walker,3,3,4,Entry: 1 damage to enemy horde,,,,,,,,,,
Sub-Zero,Water,Officer,Walker,2,2,3,Entry: Freeze an enemy zombie and deal 1 damage to it.,,,,,,,,,,
Vortex,Water,Officer,Walker,3,4,4,Entry: Return an enemy zombie to its owner's deck.,entry,return_to_enemy_hand,,enemy_selectable,,,,,,
Blizzard,Water,Commander,Walker,3,3,5,Entry: Freeze all enemy zombies,,,,,,,,,,
Freezzee,Water,Commander,Walker,4,4,6,"Entry: Freeze target and 2 adjacent zombies, if frozen 3 damage",,,,,,,,,,
Froztbite,Water,Commander,Walker,0,6,4,Becomes 6/6 at the start of your next turn.,,,,,,,,,,
Zplash,Water,Commander,Walker,4,4,6,Entry: Deal 1 damage to a random enemy for each water zombie you have in hand.,,,,,,,,,,
Zpring,Water,Commander,Walker,3,3,5,"As long as this is in the battleground, you have 1 extra Goo every turn.",entry,add_goo,while_alive,,,,,,,
Maelstrom,Water,General,Walker,6,6,8,Entry: Shuffle all zombies in the battleground back into their owners' decks.,,,,,,,,,,
Tsunami,Water,General,Walker,8,7,9,Entry: Creates a title wave and 2 damages all enemy targets,,,,,,,,,,
Feeble,Void,Minion,Walker,1,1,2,Entry: Give an enemy zombie -1 Attack.,,,,,,,,,,
Drainer,Void,Minion,Walker,1,2,3,Entry: Opponent has 1 less Goo next turn.,after_entry(1),reduce_goo,,enemy,,,,,,
Drifter,Void,Minion,Walker,1,2,3,Entry: Opponent discards a card from hand at random.,entry,discard_from_hand,,enemy_hand,,,,,,
Trasher,Void,Minion,Walker,2,2,1,Entry: Both players discard the top card of their decks.,entry,discard_from_top_of_deck,,enemy_deck,,,,,,
Ztalker,Void,Minion,Feral,3,1,1,"Feral, When this dies, you take 3 damage.",,,,,,,,,,
Zinge,Void,Minion,Heavy,2,3,3,"Heavy, Weaken 1 to attacking zombie.",,,,,,,,,,
Doubt,Void,Minion,Walker,2,1,2,Death: enemy zombies have -1 Attack until the end of your next turn.,,,,,,,,,,
Faded,Void,Minion,Feral,1,1,2,"Feral, Attack: Weaken 1.",,,,,,,,,,
Violator,Void,Officer,Walker,3,3,5,Entry: Destroy a random enemy zombie that costs 2 or less.,,,,,,,,,,
Choker,Void,Officer,Walker,2,4,4,Entry: Opponent loses 1 Goo.,,,,,,,,,,
Corrupter,Void,Officer,Walker,3,4,5,Entry: Weaken 2.,,,,,,,,,,
Defiler,Void,Officer,Feral,3,4,5,"After this attacks a hero, that hero loses 1 Goo until the end of his next turn.",,,,,,,,,,
Crippler,Void,Officer,Heavy,3,6,4,"Heavy, Weaken 2 to attacking zombie.",,,,,,,,,,
Zlendermaz,Void,Commander,Walker,4,5,6,"Entry: If you have a Void zombie in hand, destroy an enemy zombie.",,,,,,,,,,
Ztrangler,Void,Commander,Walker,5,5,7,Entry: Opponent loses 2 Goo.,,,,,,,,,,
Delirium,Void,Commander,Walker,4,4,6,Death: Opponent discards the top 2 cards of his deck.,,,,,,,,,,
Fear,Void,Commander,Feral,6,3,6,"Feral, Entry: All opponent zombies' Attack is reduced to 0 until end of turn.",,,,,,,,,,
Envy,Void,Commander,Walker,5,7,6,Death: Opponent discards cards from hand until both players have the same number of cards.,death,discard_from_hand(enemy.hand.count - friendly.hand.count),,,,,,,,
Nightmare,Void,General,Walker,6,6,10,Entry: Weaken 3 to all enemy zombies.,,,,,,,,,,
Oblivion,Void,General,Heavy,8,8,10,"Heavy, Death: Destroy all enemy zombies that cost 5 or less in the battleground.",,,,,,,,,,`
