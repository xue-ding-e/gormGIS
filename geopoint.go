// Based on https://github.com/jinzhu/gorm/issues/142
package gormGIS

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
)

type GeoPoint struct {
	Lng  float64 `json:"lng"`
	Lat  float64 `json:"lat"`
	SRID int     `json:"srid,omitempty"` // default 4326
}

func (p *GeoPoint) SetSRID(srid int) *GeoPoint {
	p.SRID = srid
	return p
}

func NewGeoPoint(lng, lat float64) *GeoPoint {
	return &GeoPoint{
		Lng:  lng,
		Lat:  lat,
		SRID: 4326,
	}
}

func NewGeoPointWithSRID(lng, lat float64, srid int) *GeoPoint {
	return &GeoPoint{
		Lng:  lng,
		Lat:  lat,
		SRID: srid,
	}
}

func (p *GeoPoint) String() string {
	srid := p.SRID
	if srid == 0 {
		srid = 4326
	}
	return fmt.Sprintf("SRID=%d;POINT(%v %v)", srid, p.Lng, p.Lat)
}

func (p *GeoPoint) Scan(val interface{}) error {
	if val == nil {
		return nil
	}

	var b []byte
	var err error

	// 处理不同类型的输入
	switch v := val.(type) {
	case []uint8:
		b, err = hex.DecodeString(string(v))
	case string:
		b, err = hex.DecodeString(v)
	default:
		return fmt.Errorf("cannot scan %T into GeoPoint", val)
	}

	if err != nil {
		return err
	}

	r := bytes.NewReader(b)
	var wkbByteOrder uint8
	if err := binary.Read(r, binary.LittleEndian, &wkbByteOrder); err != nil {
		return err
	}

	var byteOrder binary.ByteOrder
	switch wkbByteOrder {
	case 0:
		byteOrder = binary.BigEndian
	case 1:
		byteOrder = binary.LittleEndian
	default:
		return fmt.Errorf("Invalid byte order %d", wkbByteOrder)
	}

	var wkbGeometryType uint64
	if err := binary.Read(r, byteOrder, &wkbGeometryType); err != nil {
		return err
	}

	if err := binary.Read(r, byteOrder, p); err != nil {
		return err
	}

	return nil
}

func (p GeoPoint) Value() (driver.Value, error) {
	return p.String(), nil
}
