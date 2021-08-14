module ttn-to-aprs-bridge

go 1.16

require (
	github.com/ebarkie/aprs v1.0.3
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/prometheus/client_golang v1.11.0
	github.com/streadway/amqp v1.0.0
	github.com/tkanos/gonfig v0.0.0-20210106201359-53e13348de2f
	gorm.io/driver/postgres v1.1.0
	gorm.io/gorm v1.21.13
)

replace github.com/ebarkie/aprs => github.com/ttnmapper/aprs v1.0.4-0.20210814150543-b944bc851c7d
