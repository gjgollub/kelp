package backend

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/kelp/plugins"

	"github.com/stellar/go/support/config"
	"github.com/stellar/kelp/gui/model"
	"github.com/stellar/kelp/trader"
)

type botConfigResponse struct {
	Name           string                `json:"name"`
	Strategy       string                `json:"strategy"`
	TraderConfig   trader.BotConfig      `json:"trader_config"`
	StrategyConfig plugins.BuySellConfig `json:"strategy_config"`
}

func (s *APIServer) getBotConfig(w http.ResponseWriter, r *http.Request) {
	botName, e := s.parseBotName(r)
	if e != nil {
		s.writeError(w, fmt.Sprintf("error parsing bot name in getBotConfig: %s\n", e))
		return
	}

	filenamePair := model.GetBotFilenames(botName, "buysell")
	traderFilePath := fmt.Sprintf("%s/%s", s.configsDir, filenamePair.Trader)
	var botConfig trader.BotConfig
	e = config.Read(traderFilePath, &botConfig)
	if e != nil {
		s.writeErrorJson(w, fmt.Sprintf("cannot read bot config at path '%s': %s\n", traderFilePath, e))
		return
	}
	strategyFilePath := fmt.Sprintf("%s/%s", s.configsDir, filenamePair.Strategy)
	var buysellConfig plugins.BuySellConfig
	e = config.Read(strategyFilePath, &buysellConfig)
	if e != nil {
		s.writeErrorJson(w, fmt.Sprintf("cannot read strategy config at path '%s': %s\n", strategyFilePath, e))
		return
	}

	response := botConfigResponse{
		Name:           botName,
		Strategy:       "buysell",
		TraderConfig:   botConfig,
		StrategyConfig: buysellConfig,
	}
	jsonBytes, e := json.MarshalIndent(response, "", "  ")
	if e != nil {
		s.writeErrorJson(w, fmt.Sprintf("cannot marshal botConfigResponse: %s\n", e))
		return
	}
	log.Printf("getBotConfig response for botName '%s': %s\n", botName, string(jsonBytes))
	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}
