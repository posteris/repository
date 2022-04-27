package repository

import (
	"math/rand"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/posteris/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	conndata "github.com/posteris/ci-database-starter/conn-data"
)

// Product is a domain entity
type product struct {
	ID          uint8
	Name        string
	IsAvailable bool
}

// ProductGorm is DTO used to map Product entity to database
type productModel struct {
	ID          uint8  `gorm:"primarykey"`
	Name        string `gorm:"column:name"`
	IsAvailable bool   `gorm:"column:is_available"`
}

func (g productModel) ToEntity() product {
	return product(g)
}

func (g productModel) FromEntity(product product) interface{} {
	return productModel(product)
}

func createAleatoryData() product {
	return product{
		Name:        uuid.NewString(),
		IsAvailable: rand.Intn(2) == 1,
	}

}

func getDatabaseinstance(t *testing.T, tt conndata.Test) *gorm.DB {
	os.Setenv(database.DatabaseTypeLabel, tt.Type)
	os.Setenv(database.DatabaseDsnLabel, tt.Args.DSN)

	db, err := database.NewFromEnv(nil)
	if err != nil {
		t.Error("unable to connect to database")
	}

	return db
}

func Test_Create(t *testing.T) {
	//obtains the connection database parameters from the ci-database-starter
	tests := conndata.GetTestData()

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			db := getDatabaseinstance(t, tt)
			db.AutoMigrate(&productModel{})

			//create new repository instance
			repository := New[productModel, product](db)

			prd := createAleatoryData()

			errCreate := repository.Create(&prd)
			if errCreate != nil {
				t.Errorf("unable to create register: %v", errCreate)
			}
		})
	}
}

func Test_Updates(t *testing.T) {
	//obtains the connection database parameters from the ci-database-starter
	tests := conndata.GetTestData()

	//remove clickhouse - it's trash at gorm update
	tests = conndata.RemoveFromTestData(conndata.Clickhouse, tests)

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {

			db := getDatabaseinstance(t, tt)
			db.AutoMigrate(&productModel{})

			//create new repository instance
			repository := New[productModel, product](db)

			prd := createAleatoryData()

			errCreate := repository.Create(&prd)
			if errCreate != nil {
				t.Errorf("unable to create register: %v", errCreate)
			}

			t.Log(prd.ID, prd.Name, prd.IsAvailable)

			var oldName string = prd.Name

			prd.Name = uuid.NewString()

			errUpdate := repository.Updates(&prd)
			if errUpdate != nil {
				t.Errorf("unable to update register: %v", errUpdate)
			}

			assert.NotEqual(t, oldName, prd.Name)
		})
	}
}
