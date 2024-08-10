package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Player struct {
	C                          string `json:"@c"`
	PlayerID                   int    `json:"playerID"`
	TeamID                     int    `json:"teamID"`
	Nationality                int    `json:"nationality"`
	CapitalID                  int    `json:"capitalID"`
	Title                      string `json:"title"`
	FullTitle                  string `json:"fullTitle"`
	NationName                 string `json:"nationName"`
	NationDescription          string `json:"nationDescription"`
	NationDifficulty           string `json:"nationDifficulty"`
	NationAdjective            string `json:"nationAdjective"`
	HighlightedUnits           []any  `json:"highlightedUnits"`
	PrimaryColor               string `json:"primaryColor"`
	SecondaryColor             string `json:"secondaryColor"`
	Faction                    int    `json:"faction"`
	ComputerPlayer             bool   `json:"computerPlayer"`
	LastLogin                  int    `json:"lastLogin"`
	IsGameCreator              bool   `json:"isGameCreator"`
	NativeComputer             bool   `json:"nativeComputer"`
	SiteUserID                 int    `json:"siteUserID"`
	PlayerImageID              int    `json:"playerImageID"`
	LastDayOfPlayerImageUpdate int    `json:"lastDayOfPlayerImageUpdate"`
	FlagImageID                int    `json:"flagImageID"`
	LastDayOfFlagImageUpdate   int    `json:"lastDayOfFlagImageUpdate"`
	Banned                     bool   `json:"banned"`
	Defeated                   bool   `json:"defeated"`
	Retired                    bool   `json:"retired"`
	PremiumUser                bool   `json:"premiumUser"`
	NoobBonus                  int    `json:"noobBonus"`
	ExpansionMoralePenalty     int    `json:"expansionMoralePenalty"`
	AchievementTitleID         int    `json:"achievementTitleID"`
	PassiveAI                  bool   `json:"passiveAI"`
	PremiumBuildSlot           bool   `json:"premiumBuildSlot"`
	PremiumProductionSlot      bool   `json:"premiumProductionSlot"`
	IsTutorialGameFlag         bool   `json:"isTutorialGameFlag"`
	LastLeftTeam               int    `json:"lastLeftTeam"`
	LastKickedFromTeam         int    `json:"lastKickedFromTeam"`
	TransportShipUnitTypeID    int    `json:"transportShipUnitTypeId"`
	PoiTimer                   int    `json:"poiTimer"`
	FactionWasSet              bool   `json:"factionWasSet"`
	Name                       string `json:"name"`
	IsBot                      bool
}
type PlayerStates struct {
	Players map[string]Player `json:"players"`
}

type MarketOrder struct {
	C            string  `json:"@c"`
	Buy          bool    `json:"buy"`
	Amount       int     `json:"amount"`
	Limit        float64 `json:"limit"`
	PlayerID     int     `json:"playerID"`
	ResourceType int     `json:"resourceType"`
	OrderID      int     `json:"orderID"`
	NationName   string
	IsBot        bool
}

type MarketData struct {
	Orders []MarketOrder
}

func (m *MarketData) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	return json.Unmarshal(raw[1], &m.Orders)
}

type MarketStates struct {
	Asks []MarketData `json:"asks"`
	Bids []MarketData `json:"bids"`
}

type States struct {
	PState PlayerStates `json:"1"`
	MState MarketStates `json:"4"`
}

type Data struct {
	States States `json:"states"`
}

type ResponsePayload struct {
	Result Data `json:"result"`
}

func main() {
	requestPayload := []byte(`{"@c":"ultshared.action.UltUpdateGameStateAction","version":"218","client":"s1914-client-ultimate","adminLevel":0,"gameID":8930992,"playerID":0,"rights":"chat"}`)
	resp, err := http.Post("https://xgs-as-fwnq.c.bytro.com", "application/x-www-form-urlencoded; charset=UTF-8", bytes.NewBuffer(requestPayload))
	if err != nil {
		fmt.Println("Post failed: ", err.Error())
		return
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body, ", err.Error())
	}

	var res ResponsePayload
	err = json.Unmarshal(body, &res)
	if err != nil {
		fmt.Println("Failed to decode json, ", err.Error())
	}

	players := res.Result.States.PState.Players
	for _, p := range players {
		p.IsBot = p.LastLogin == 0
		players[strconv.Itoa(p.PlayerID)] = p
	}
	for _, b := range res.Result.States.MState.Bids {
		for i, o := range b.Orders {
			player := players[strconv.Itoa(o.PlayerID)]
			b.Orders[i].NationName = player.NationName
			b.Orders[i].IsBot = player.IsBot
		}
	}

	fmt.Printf("%#v\n", res.Result.States.MState.Bids)
}
