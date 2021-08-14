package sqldb

import (
	"fmt"
	"github.com/patrickmn/go-cache"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type Credentials struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string

	DebugLog bool
}

var (
	db *gorm.DB
	aprsDevicesCache *cache.Cache
)

func (credentials *Credentials) Init() {
	var err error

	connectionString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		credentials.Host, credentials.Port, credentials.User, credentials.Password, credentials.Database)

	var gormLogLevel = logger.Silent
	if credentials.DebugLog {
		log.Println("Database debug logging enabled")
		gormLogLevel = logger.Info
	}

	db, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		panic(err.Error())
	}

	// Create tables if they do not exist
	log.Println("Performing auto migrate")
	if err := db.AutoMigrate(
		&AprsDevice{},
	); err != nil {
		log.Println("Unable autoMigrateDB - " + err.Error())
	}

	aprsDevicesCache = cache.New(120*time.Minute, 10*time.Minute)
}

func GetAprsDevice(networkId string, appId string, devId string) (AprsDevice, error) {
	deviceKey := networkId + "/" + appId + "/" + devId

	// Try to load transforms from cache first
	if x, found := aprsDevicesCache.Get(deviceKey); found {
		//log.Println("  [d] Cache hit")
		aprsDevice := x.(AprsDevice)
		return aprsDevice, nil
	}

	device := AprsDevice{
		NetworkId:   networkId,
		AppId:       appId,
		DevId:       devId,
	}
	err := db.Find(&device, &device).Error
	if err != nil {
		return device, err
	}

	// Store in cache
	aprsDevicesCache.Set(deviceKey, device, cache.DefaultExpiration)

	return device, nil
}

func IsAprsDevice(networkId string, appId string, devId string) bool {
	device, err := GetAprsDevice(networkId, appId, devId)

	if err != nil {
		log.Println(err.Error())
		return false
	}

	if device.ID == 0 {
		return false
	}

	return true
}