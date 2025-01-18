package dao

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"testing"
	"time"

	"database/sql"

	"example.com/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type ForexRateDaoTestSuite struct {
	suite.Suite
	dao            *ForexRateDao
	mysqlContainer *mysql.MySQLContainer
	db             *sql.DB
}

func (suite *ForexRateDaoTestSuite) SetupTest() {
	ctx := context.Background()

	// start up MySQL testcontainer
	var err error
	suite.mysqlContainer, err = startMySQLContainer()
	if err != nil {
		log.Panic("fail to start up MySQL container\n", err)
	}
	log.Println("MySQL container started")

	suite.db, err = openDB(ctx, suite.mysqlContainer)
	if err != nil {
		log.Panic("fail to connect to MySQL container\n", err)
	}
	suite.dao = &ForexRateDao{suite.db}

}

func (suite *ForexRateDaoTestSuite) TearDownTest() {

	log.Println("closing DB handler")
	if err := suite.db.Close(); err != nil {
		log.Println("failed to close database", err)
	}

	log.Println("shutting down MySQL container")
	if err := testcontainers.TerminateContainer(suite.mysqlContainer); err != nil {
		log.Printf("failed to terminate container: %s\n", err)
	}
}

func (suite *ForexRateDaoTestSuite) TestInsert() {
	duration, err := time.ParseDuration("10m")
	if err != nil {
		log.Panic("fail to parse duration")
	}
	booking := model.ForexRateBooking{ID: "1f648720-3bd3-4c8e-8d00-294516f64bf7", Timestamp: time.Now(), BaseCurrency: "GBP",
		CounterCurrency: "USD", Rate: 0.25, TradeAction: "BUY", BaseCurrencyAmount: 1000,
		BookingRef: "ABCD100", ExpiryTime: time.Now().Add(duration), CustomerID: "f1440302-01ab-4083-88fd-8864ae83d435"}

	count, err := suite.dao.insert(&booking)
	if err != nil {
		fmt.Println("fail to insert", err)
	}

	assert.Equal(suite.T(), int64(1), count)

}

func TestForexRateDaoTestSuite(t *testing.T) {
	suite.Run(t, new(ForexRateDaoTestSuite))
}

func openDB(ctx context.Context, mysqlContainer *mysql.MySQLContainer) (*sql.DB, error) {

	connStr, err := mysqlContainer.ConnectionString(ctx)
	if err != nil {
		return nil, err
	}
	log.Println("MySQL connection string: ", connStr)

	db, err := sql.Open("mysql", connStr)
	return db, err
}

func startMySQLContainer() (*mysql.MySQLContainer, error) {
	ctx := context.Background()

	mysqlContainer, err := mysql.Run(ctx,
		"mysql:8.0.36",
		mysql.WithDatabase("forex"),
		mysql.WithUsername("root"),
		mysql.WithPassword("password"),
		mysql.WithScripts(filepath.Join("testdata", "schema.sql")),
	)

	if err != nil {
		log.Printf("failed to start container: %s\n", err)
		return nil, err
	}

	return mysqlContainer, nil

}
