package gormGIS_test

import (
	"fmt"

	_ "github.com/lib/pq"
	"github.com/xue-ding-e/gormGIS"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

var (
	DB *gorm.DB
)

func init() {
	var err error
	fmt.Println("testing postgres...")
	DB, err = gorm.Open(postgres.Open("user=gorm dbname=gormGIS sslmode=disable"), &gorm.Config{})
	DB.Debug()
	if err != nil {
		panic(fmt.Sprintf("No error should happen when connect database, but got %+v", err))
	}

	DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis")
	DB.Exec("CREATE EXTENSION IF NOT EXISTS postgis_topology")
}

type TestPoint struct {
	ID       uint                `gorm:"primarykey"`
	Location gormGIS.GeoPoint   `gorm:"type:geometry(Point,4326)"`
}

func TestGeoPoint(t *testing.T) {
	// 自动迁移表结构
	err := DB.AutoMigrate(&TestPoint{})
	if err != nil {
		t.Errorf("Can't create table: %v", err)
		return
	}

	p := TestPoint{
		Location: gormGIS.GeoPoint{
			GeoPointDB: gormGIS.GeoPointDB{
				Lat: 43.76857094631136,
				Lng: 11.292383687705296,
			},
			SRID: 4326,
		},
	}

	err = DB.Create(&p).Error
	if err != nil {
		t.Errorf("Can't create row: %v", err)
		return
	}

	var res TestPoint
	err = DB.First(&res).Error
	if err != nil {
		t.Errorf("Can't query row: %v", err)
		return
	}

	if res.Location.Lat != 43.76857094631136 {
		t.Errorf("Latitude not correct, expected %f, got %f", 43.76857094631136, res.Location.Lat)
	}

	if res.Location.Lng != 11.292383687705296 {
		t.Errorf("Longitude not correct, expected %f, got %f", 11.292383687705296, res.Location.Lng)
	}
}