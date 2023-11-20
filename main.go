package main

import (
	"fmt"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"sync"
	"time"
	"tv/tradingview"
	socket "tv/tradingview"
)

var SymbolDataMutex sync.Mutex
var LatestSymbolData = make(map[string]tradingview.SymbolData)
var tradingviewsocket *socket.Socket

func handleSymbolUpdates(symbol string, data *socket.QuoteData) {
	prevData, exists := LatestSymbolData[symbol]
	if !exists || data.Price != nil {
		if data.Price != nil {
			prevData.Price = *data.Price
			prevData.UpdatedTime = time.Now()
			SymbolDataMutex.Lock()
			LatestSymbolData[symbol] = prevData
			SymbolDataMutex.Unlock()
		}
	}
}

func establishWebSocketConnection() {
	tradingviewsocket, _ = socket.Connect(
		handleSymbolUpdates,
		func(err error, context string) {
			// Handle errors, if needed
		},
	)
}

func removeInactiveSymbolsAndSocketSymbols() {
	currentTimestamp := time.Now()

	for symbol, symbolData := range LatestSymbolData {
		if currentTimestamp.Sub(symbolData.LastRequest) > 20*time.Second {
			fmt.Printf("Symbol %s inactive for more than 30 seconds. Removing symbol and initiating socket removal\n", symbol)
			go func(sym string) {
				SymbolDataMutex.Lock()
				tradingviewsocket.RemoveSymbol(sym)
				delete(LatestSymbolData, symbol)
				SymbolDataMutex.Unlock()
			}(symbol)
		}
	}
}

// todo search the symbol first on the api call and make sure it exists
func getLatestPrice(c *gin.Context) {
	symbol := c.Query("symbol")
	if symbol == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Symbol not provided"})
		return
	}

	var symbolExists bool

	//tradingview.SymbolDataMutex.Lock()
	_, symbolExists = LatestSymbolData[symbol]
	//tradingview.SymbolDataMutex.Unlock()

	if symbolExists {
		symbolData := LatestSymbolData[symbol]
		symbolData.LastRequest = time.Now()
		LatestSymbolData[symbol] = symbolData
		c.JSON(http.StatusOK, gin.H{
			"symbol":       symbol,
			"price":        symbolData.Price,
			"last_update":  symbolData.UpdatedTime,
			"last_request": symbolData.LastRequest,
		})
		return
	}
	//go func() {
	//	tradingviewsocket, _ := socket.Connect(
	//		func(symbol string, data *socket.QuoteData) {
	//			// Update the latest data for each symbol
	//			SymbolDataMutex.Lock()
	//			prevData, exists := LatestSymbolData[symbol]
	//			if !exists || data.Price != nil {
	//				if data.Price != nil {
	//					prevData.Price = *data.Price
	//					prevData.UpdatedTime = time.Now()
	//					LatestSymbolData[symbol] = prevData
	//				}
	//			}
	//			SymbolDataMutex.Unlock()
	//		},
	//		func(err error, context string) {
	//			// Handle errors, if needed
	//		},
	//	)
	//	tradingviewsocket.AddSymbol(symbol)
	//}()
	SymbolDataMutex.Lock()
	LatestSymbolData[symbol] = tradingview.SymbolData{} // Initialize data for a new symbol
	tradingviewsocket.AddSymbol(symbol)
	SymbolDataMutex.Unlock()

	c.JSON(http.StatusNotFound, gin.H{"error": "Data for symbol does not exist"})
}

func init() {
	godotenv.Load()
}

func main() {
	// Establish the trading view socket connection
	router := gin.Default()
	port := os.Getenv("PORT")
	host := os.Getenv("HOST")

	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(10).Seconds().Do(removeInactiveSymbolsAndSocketSymbols)
	scheduler.StartAsync()
	establishWebSocketConnection()

	router.GET("/latest-price", getLatestPrice)
	pprof.Register(router)

	router.Run(host + ":" + port)
}
