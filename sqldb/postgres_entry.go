package sqldb

type AprsDevice struct {
	ID uint

	NetworkId string `gorm:"index:aprs_network_app_dev,unique"`
	AppId   string `gorm:"index:aprs_network_app_dev,unique"`
	DevId   string `gorm:"index:aprs_network_app_dev,unique"`

	Callsign string
	Ssid int

	SymbolTable string
	Symbol string
}