package battleground

type AbilityActivityType int32

const (
        AbilityActivityType_PASSIVE AbilityActivityType = iota
        AbilityActivityType_ACTIVE AbilityActivityType = iota
)

type AbilityCallType int32

const (
        AbilityCallType_TURN AbilityCallType = iota
        AbilityCallType_ENTRY AbilityCallType = iota
        AbilityCallType_END AbilityCallType = iota
        AbilityCallType_ATTACK AbilityCallType = iota
        AbilityCallType_DEATH AbilityCallType = iota
        AbilityCallType_PERMANENT AbilityCallType = iota
        AbilityCallType_GOT_DAMAGE AbilityCallType = iota
        AbilityCallType_AT_DEFENCE AbilityCallType = iota
        AbilityCallType_IN_HAND AbilityCallType = iota
)

type PlayerAction int32

const (
        PlayerAction_StartGame PlayerAction = iota
        PlayerAction_EndTurn PlayerAction = iota
        PlayerAction_DrawCardPlayer PlayerAction = iota
        PlayerAction_PlayCard PlayerAction = iota
        PlayerAction_CardAttack PlayerAction = iota
        PlayerAction_UseCardAbility PlayerAction = iota
)

type GameStateChangeAction int32

const (
        GameStateChangeAction_None GameStateChangeAction = iota
        GameStateChangeAction_SetPlayerDefense GameStateChangeAction = iota
        GameStateChangeAction_SetPlayerGoo GameStateChangeAction = iota
)

type CustomUiElement int32

const (
        CustomUiElement_None CustomUiElement = iota
        CustomUiElement_Label CustomUiElement = iota
        CustomUiElement_Button CustomUiElement = iota
)

 

