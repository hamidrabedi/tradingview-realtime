package tradingview

import "time"

// SocketInterface ...
type SocketInterface interface {
	AddSymbol(symbol string) error
	RemoveSymbol(symbol string) error
	Init() error
	Close() error
	GetSessionId() string
}

// SocketMessage ...
type SocketMessage struct {
	Message string      `json:"m"`
	Payload interface{} `json:"p"`
}

// QuoteMessage ...
type QuoteMessage struct {
	Symbol string     `mapstructure:"n"`
	Status string     `mapstructure:"s"`
	Data   *QuoteData `mapstructure:"v"`
}

// QuoteData ...
type QuoteData struct {
	Price  *float64 `mapstructure:"lp"`
	Volume *float64 `mapstructure:"volume"`
	Bid    *float64 `mapstructure:"bid"`
	Ask    *float64 `mapstructure:"ask"`
}

// Flags ...
type Flags struct {
	Flags []string `json:"flags"`
}

// OnReceiveDataCallback ...
type OnReceiveDataCallback func(symbol string, data *QuoteData)

// OnErrorCallback ...
type OnErrorCallback func(err error, context string)

type SymbolData struct {
	Price       float64
	UpdatedTime time.Time
	LastRequest time.Time
}
