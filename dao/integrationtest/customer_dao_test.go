package integrationtest

import (
	"context"
	"fmt"
	"log"
	"testing"

	"database/sql"

	"example.com/dao"
	"example.com/model"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go/modules/mysql"
)

type CustomerDaoTestSuite struct {
	suite.Suite
	dao            *dao.CustomerDao
	mysqlContainer *mysql.MySQLContainer
	db             *sql.DB
}

func (suite *CustomerDaoTestSuite) SetupSuite() {
	log.Println("setting up test")
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
	suite.dao = dao.NewCustomerDao(suite.db)

}

func (suite *CustomerDaoTestSuite) TearDownTest() {
	_, err := suite.db.Exec("DELETE FROM customer")
	if err != nil {
		log.Panic("fail to clean up data")
	}
}

func (suite *CustomerDaoTestSuite) TearDownSuite() {
	cleanUp(suite.db, suite.mysqlContainer)
}

func (suite *CustomerDaoTestSuite) TestInsert() {
	customer := model.Customer{ID: "1f648720-3bd3-4c8e-8d00-294516f64bf7", Name: "Customer Name", Tier: 1}

	count, err := suite.dao.Insert(&customer)
	if err != nil {
		fmt.Println("fail to insert", err)
	}

	assert.Equal(suite.T(), int64(1), count)

}

func (suite *CustomerDaoTestSuite) TestFindByID() {
	customer := model.Customer{ID: "1f648720-3bd3-4c8e-8d00-294516f64bf7", Name: "Customer Name", Tier: 1}

	err := insertCustomer(suite.db, &customer)
	if err != nil {
		fmt.Println("fail to insert", err)
	}

	var actual *model.Customer
	actual, err = suite.dao.FindByID(customer.ID)
	if err != nil {
		suite.T().Error("fail to retrieve record", err)
	}
	if actual == nil {
		suite.T().Error("No record found")
	}

	assert.Equal(suite.T(), actual, &customer)

}

func TestCustomerDaoTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerDaoTestSuite))
}

func insertCustomer(db *sql.DB, customer *model.Customer) error {
	_, err := db.Exec("INSERT INTO customer(id, name, tier) VALUES (?, ?, ?)",
		customer.ID, customer.Name, customer.Tier)

	if err != nil {
		log.Fatalf("insert customer: %v", err)
		return err
	}

	return nil
}
