package main

import (
	"log/slog"

	"errors"
)

type tracker struct {
	companyTracker map[string]*Price // only for most current timestamp
	// companyTracker map[string][]*Price // history
}

type Price struct {
	Price     float64
	Timestamp int
}

func (t *tracker) Update(timestamp int, companyId string, price float64) {
	// var price Price
	t.companyTracker[companyId] = &Price{
		Price:     price,
		Timestamp: timestamp,
	}
}

func (t *tracker) GetCurrentPrice(companyId string) (*Price, error) {
	if price, ok := t.companyTracker[companyId]; ok {
		return price, nil
	}
	// what to do if company does not exist
	return nil, errors.New("company does not exist")
}

func AssetPriceTracker() *tracker {
	return &tracker{
		companyTracker: make(map[string]*Price),
	}
}

func main() {
	tracker := AssetPriceTracker()
	tracker.Update(1627683600, "AAPL", 145.30)
	price, err := tracker.GetCurrentPrice("AAPL")
	if err != nil {
		slog.Error("Error getting price", "error", err.Error())
	}

	slog.Info("APPLE", "current_price", price) // Output: 145.30

	tracker.Update(1627683660, "AAPL", 145.50)
	tracker.Update(1627683720, "GOOG", 2729.89)

	price, err = tracker.GetCurrentPrice("AAPL")
	if err != nil {
		slog.Error("Error getting price", "error", err.Error())
	}

	slog.Info("APPL", "current_price", price.Price) // Output: 145.50

	price, err = tracker.GetCurrentPrice("GOOG")
	if err != nil {
		slog.Error("Error getting price", "error", err.Error())
	}

	slog.Info("GOOG", "current_price", price.Price) // Output: 2729.89

	// print(tracker.get_current_price("AAPL"))  # Output: 145.50
	// print(tracker.get_current_price("GOOG"))  # Output: 2729.89

}

// Coding Challenge: Real-Time Data Processing

// Problem Statement:

// Pillar's platform needs to process real-time streams of financial data from multiple vendors to make quick, accurate hedging decisions. Your task is to implement a system that can efficiently process these data streams and output the most recent price of a given asset in constant time.

// Specifications:

// - You're provided with a continuous stream of price updates for different assets. Each update is a tuple (timestamp, asset_id, price).
// - Implement a class AssetPriceTracker that supports the following operations:
//   - update(timestamp, asset_id, price): Updates the price of the asset at the given timestamp.
//   - get_current_price(asset_id) -> float: Returns the most recent price of the asset.
  
// Constraints:
// - The system should be able to handle a large number of assets and frequent updates efficiently.
// - You should optimize for both time and space complexity.

// Example:
// tracker = AssetPriceTracker()
// tracker.update(1627683600, "AAPL", 145.30)
// tracker.update(1627683660, "AAPL", 145.50)
// tracker.update(1627683720, "GOOG", 2729.89)

// print(tracker.get_current_price("AAPL"))  # Output: 145.50
// print(tracker.get_current_price("GOOG"))  # Output: 2729.89


