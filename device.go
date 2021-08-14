package main

import (
	"sync"
	"time"
	"ttn-to-aprs-bridge/types"
)

var minimumInterval = 30*time.Second
var lastReportCache = sync.Map{}

func GetLastReportCacheKey(message types.TtnMapperUplinkMessage) string {
	return message.NetworkId + "/" + message.AppID + "/" + message.DevID
}

// Returns true if we may report this location.
// Returns false if we have reported less than minimumInterval ago
func MinimumReportDurationPassed(currentMessage types.TtnMapperUplinkMessage) bool {
	deviceKey := GetLastReportCacheKey(currentMessage)

	lastReport, ok := lastReportCache.Load(deviceKey)
	if !ok {
		return true
	}

	lastReportedMessage := lastReport.(types.TtnMapperUplinkMessage)

	if currentMessage.Time - lastReportedMessage.Time > minimumInterval.Nanoseconds() {
		return true
	} else {
		return false
	}
}

func CacheReportedMessage(currentMessage types.TtnMapperUplinkMessage) {
	deviceKey := GetLastReportCacheKey(currentMessage)
	lastReportCache.Store(deviceKey, currentMessage)
}