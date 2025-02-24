package trader

import (
	"fmt"

	"github.com/stellar/go/clients/horizon"
	"github.com/stellar/kelp/support/toml"
	"github.com/stellar/kelp/support/utils"
)

// XLM is a constant for XLM
const XLM = "XLM"

// FeeConfig represents input data for how to deal with network fees
type FeeConfig struct {
	CapacityTrigger float64 `valid:"-" toml:"CAPACITY_TRIGGER" json:"capacity_trigger"`     // trigger when "ledger_capacity_usage" in /fee_stats is >= this value
	Percentile      uint8   `valid:"-" toml:"PERCENTILE" json:"percentile"`                 // percentile computation to use from /fee_stats (10, 20, ..., 90, 95, 99)
	MaxOpFeeStroops uint64  `valid:"-" toml:"MAX_OP_FEE_STROOPS" json:"max_op_fee_stroops"` // max fee in stroops per operation to use
}

// BotConfig represents the configuration params for the bot
type BotConfig struct {
	SourceSecretSeed                   string     `valid:"-" toml:"SOURCE_SECRET_SEED" json:"source_secret_seed"`
	TradingSecretSeed                  string     `valid:"-" toml:"TRADING_SECRET_SEED" json:"trading_secret_seed"`
	AssetCodeA                         string     `valid:"-" toml:"ASSET_CODE_A" json:"asset_code_a"`
	IssuerA                            string     `valid:"-" toml:"ISSUER_A" json:"issuer_a"`
	AssetCodeB                         string     `valid:"-" toml:"ASSET_CODE_B" json:"asset_code_b"`
	IssuerB                            string     `valid:"-" toml:"ISSUER_B" json:"issuer_b"`
	TickIntervalSeconds                int32      `valid:"-" toml:"TICK_INTERVAL_SECONDS" json:"tick_interval_seconds"`
	MaxTickDelayMillis                 int64      `valid:"-" toml:"MAX_TICK_DELAY_MILLIS" json:"max_tick_delay_millis"`
	DeleteCyclesThreshold              int64      `valid:"-" toml:"DELETE_CYCLES_THRESHOLD" json:"delete_cycles_threshold"`
	SubmitMode                         string     `valid:"-" toml:"SUBMIT_MODE" json:"submit_mode"`
	FillTrackerSleepMillis             uint32     `valid:"-" toml:"FILL_TRACKER_SLEEP_MILLIS" json:"fill_tracker_sleep_millis"`
	FillTrackerDeleteCyclesThreshold   int64      `valid:"-" toml:"FILL_TRACKER_DELETE_CYCLES_THRESHOLD" json:"fill_tracker_delete_cycles_threshold"`
	HorizonURL                         string     `valid:"-" toml:"HORIZON_URL" json:"horizon_url"`
	CcxtRestURL                        *string    `valid:"-" toml:"CCXT_REST_URL" json:"ccxt_rest_url"`
	Fee                                *FeeConfig `valid:"-" toml:"FEE" json:"fee"`
	CentralizedPricePrecisionOverride  *int8      `valid:"-" toml:"CENTRALIZED_PRICE_PRECISION_OVERRIDE" json:"centralized_price_precision_override"`
	CentralizedVolumePrecisionOverride *int8      `valid:"-" toml:"CENTRALIZED_VOLUME_PRECISION_OVERRIDE" json:"centralized_volume_precision_override"`
	// Deprecated: use CENTRALIZED_MIN_BASE_VOLUME_OVERRIDE instead
	MinCentralizedBaseVolumeDeprecated *float64                 `valid:"-" toml:"MIN_CENTRALIZED_BASE_VOLUME" deprecated:"true" json:"min_centralized_base_volume"`
	CentralizedMinBaseVolumeOverride   *float64                 `valid:"-" toml:"CENTRALIZED_MIN_BASE_VOLUME_OVERRIDE" json:"centralized_min_base_volume_override"`
	CentralizedMinQuoteVolumeOverride  *float64                 `valid:"-" toml:"CENTRALIZED_MIN_QUOTE_VOLUME_OVERRIDE" json:"centralized_min_quote_volume_override"`
	AlertType                          string                   `valid:"-" toml:"ALERT_TYPE" json:"alert_type"`
	AlertAPIKey                        string                   `valid:"-" toml:"ALERT_API_KEY" json:"alert_api_key"`
	MonitoringPort                     uint16                   `valid:"-" toml:"MONITORING_PORT" json:"monitoring_port"`
	MonitoringTLSCert                  string                   `valid:"-" toml:"MONITORING_TLS_CERT" json:"monitoring_tls_cert"`
	MonitoringTLSKey                   string                   `valid:"-" toml:"MONITORING_TLS_KEY" json:"monitoring_tls_key"`
	GoogleClientID                     string                   `valid:"-" toml:"GOOGLE_CLIENT_ID" json:"google_client_id"`
	GoogleClientSecret                 string                   `valid:"-" toml:"GOOGLE_CLIENT_SECRET" json:"google_client_secret"`
	AcceptableEmails                   string                   `valid:"-" toml:"ACCEPTABLE_GOOGLE_EMAILS" json:"acceptable_google_emails"`
	TradingExchange                    string                   `valid:"-" toml:"TRADING_EXCHANGE" json:"trading_exchange"`
	ExchangeAPIKeys                    toml.ExchangeAPIKeysToml `valid:"-" toml:"EXCHANGE_API_KEYS" json:"exchange_api_keys"`
	ExchangeParams                     toml.ExchangeParamsToml  `valid:"-" toml:"EXCHANGE_PARAMS" json:"exchange_params"`
	ExchangeHeaders                    toml.ExchangeHeadersToml `valid:"-" toml:"EXCHANGE_HEADERS" json:"exchange_headers"`

	// initialized later
	tradingAccount *string
	sourceAccount  *string // can be nil
	assetBase      horizon.Asset
	assetQuote     horizon.Asset
	isTradingSdex  bool
}

// MakeBotConfig factory method for BotConfig
func MakeBotConfig(
	sourceSecretSeed string,
	tradingSecretSeed string,
	assetCodeA string,
	issuerA string,
	assetCodeB string,
	issuerB string,
	tickIntervalSeconds int32,
	maxTickDelayMillis int64,
	deleteCyclesThreshold int64,
	submitMode string,
	fillTrackerSleepMillis uint32,
	fillTrackerDeleteCyclesThreshold int64,
	horizonURL string,
	ccxtRestURL *string,
	fee *FeeConfig,
	centralizedPricePrecisionOverride *int8,
	centralizedVolumePrecisionOverride *int8,
	centralizedMinBaseVolumeOverride *float64,
	centralizedMinQuoteVolumeOverride *float64,
) *BotConfig {
	return &BotConfig{
		SourceSecretSeed:                 sourceSecretSeed,
		TradingSecretSeed:                tradingSecretSeed,
		AssetCodeA:                       assetCodeA,
		IssuerA:                          issuerA,
		AssetCodeB:                       assetCodeB,
		IssuerB:                          issuerB,
		TickIntervalSeconds:              tickIntervalSeconds,
		MaxTickDelayMillis:               maxTickDelayMillis,
		DeleteCyclesThreshold:            deleteCyclesThreshold,
		SubmitMode:                       submitMode,
		FillTrackerSleepMillis:           fillTrackerSleepMillis,
		FillTrackerDeleteCyclesThreshold: fillTrackerDeleteCyclesThreshold,
		HorizonURL:                       horizonURL,
		CcxtRestURL:                      ccxtRestURL,
		Fee:                              fee,
		CentralizedPricePrecisionOverride:  centralizedPricePrecisionOverride,
		CentralizedVolumePrecisionOverride: centralizedVolumePrecisionOverride,
		CentralizedMinBaseVolumeOverride:   centralizedMinBaseVolumeOverride,
		CentralizedMinQuoteVolumeOverride:  centralizedMinQuoteVolumeOverride,
	}
}

// String impl.
func (b BotConfig) String() string {
	return utils.StructString(b, map[string]func(interface{}) interface{}{
		"EXCHANGE_API_KEYS":                     utils.Hide,
		"EXCHANGE_PARAMS":                       utils.Hide,
		"EXCHANGE_HEADERS":                      utils.Hide,
		"SOURCE_SECRET_SEED":                    utils.SecretKey2PublicKey,
		"TRADING_SECRET_SEED":                   utils.SecretKey2PublicKey,
		"ALERT_API_KEY":                         utils.Hide,
		"GOOGLE_CLIENT_ID":                      utils.Hide,
		"GOOGLE_CLIENT_SECRET":                  utils.Hide,
		"ACCEPTABLE_GOOGLE_EMAILS":              utils.Hide,
		"CENTRALIZED_PRICE_PRECISION_OVERRIDE":  utils.UnwrapInt8Pointer,
		"CENTRALIZED_VOLUME_PRECISION_OVERRIDE": utils.UnwrapInt8Pointer,
		"MIN_CENTRALIZED_BASE_VOLUME":           utils.UnwrapFloat64Pointer,
		"CENTRALIZED_MIN_BASE_VOLUME_OVERRIDE":  utils.UnwrapFloat64Pointer,
		"CENTRALIZED_MIN_QUOTE_VOLUME_OVERRIDE": utils.UnwrapFloat64Pointer,
	})
}

// TradingAccount returns the config's trading account
func (b *BotConfig) TradingAccount() string {
	return *b.tradingAccount
}

// SourceAccount returns the config's source account
func (b *BotConfig) SourceAccount() string {
	if b.sourceAccount == nil {
		return ""
	}
	return *b.sourceAccount
}

// AssetBase returns the config's assetBase
func (b *BotConfig) AssetBase() horizon.Asset {
	return b.assetBase
}

// AssetQuote returns the config's assetQuote
func (b *BotConfig) AssetQuote() horizon.Asset {
	return b.assetQuote
}

// IsTradingSdex returns whether the config is set to trade on SDEX
func (b *BotConfig) IsTradingSdex() bool {
	return b.isTradingSdex
}

// Init initializes this config
func (b *BotConfig) Init() error {
	b.isTradingSdex = b.TradingExchange == "" || b.TradingExchange == "sdex"

	if b.AssetCodeA == b.AssetCodeB && b.IssuerA == b.IssuerB {
		return fmt.Errorf("error: both assets cannot be the same '%s:%s'", b.AssetCodeA, b.IssuerA)
	}

	asset, e := utils.ParseAsset(b.AssetCodeA, b.IssuerA)
	if e != nil {
		return fmt.Errorf("Error while parsing Asset A: %s", e)
	}
	b.assetBase = *asset

	asset, e = utils.ParseAsset(b.AssetCodeB, b.IssuerB)
	if e != nil {
		return fmt.Errorf("Error while parsing Asset B: %s", e)
	}
	b.assetQuote = *asset

	b.tradingAccount, e = utils.ParseSecret(b.TradingSecretSeed)
	if e != nil {
		return e
	}
	if b.tradingAccount == nil {
		return fmt.Errorf("no trading account specified")
	}

	b.sourceAccount, e = utils.ParseSecret(b.SourceSecretSeed)
	return e
}
