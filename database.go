package main

import "database/sql"

type user struct {
	OID        int           `gorm:"primary_key;column:OID"`
	Login      string        `gorm:"column:Login"`
	Password   string        `gorm:"column:Password"`
	Name       string        `gorm:"column:Name"`
	FirstName  string        `gorm:"column:FirstName"`
	Rfid       string        `gorm:"column:Rfid"`
	Barcode    string        `gorm:"column:Barcode"`
	Pin        string        `gorm:"column:Pin"`
	Function   string        `gorm:"column:Function"`
	UserTypeID sql.NullInt32 `gorm:"column:UserTypeID"`
	UserRoleID sql.NullInt32 `gorm:"column:UserRoleID"`
	Email      string        `gorm:"column:Email"`
	Phone      string        `gorm:"column:Phone"`
}

func (user) TableName() string {
	return "user"
}

type product struct {
	OID             int     `gorm:"primary_key;column:OID"`
	Name            string  `gorm:"column:Name"`
	Barcode         string  `gorm:"column:Barcode"`
	Cycle           float64 `gorm:"column:Cycle"`
	IdleFromTime    int     `gorm:"column:IdleFromTime"`
	ProductStatusID int     `gorm:"column:ProductStatusID"`
	Deleted         int     `gorm:"column:Deleted"`
	ProductGroupID  int     `gorm:"column:ProductGroupID"`
	Cavity          int     `gorm:"column:Cavity"`
}

func (product) TableName() string {
	return "product"
}

type productGroup struct {
	OID  int    `gorm:"primary_key;column:OID"`
	Name string `gorm:"column:Name"`
}

func (productGroup) TableName() string {
	return "product_group"
}

type hvwZapsiZam struct {
	Alias    string `gorm:"column:Alias"`
	Jmeno    string `gorm:"column:Jmeno"`
	Prijmeni string `gorm:"column:Prijmeni"`
	Delnik   int    `gorm:"column:Delnik"`
}

func (hvwZapsiZam) TableName() string {
	return "hvw_zapsi_zam"
}

type hvwZapsiZamCip struct {
	ID       string `gorm:"column:ID"`
	Alias    string `gorm:"column:Alias"`
	CC       string `gorm:"column:CC"`
	CCH      string `gorm:"column:CCH"`
	Primarni int    `gorm:"column:Primarni"`
}

func (hvwZapsiZamCip) TableName() string {
	return "hvw_zapsi_zam_cip"
}

type hvwZapsiArtikl struct {
	ID      int    `gorm:"column:ID"`
	SkupZbo string `gorm:"column:SkupZbo"`
	RegCis  string `gorm:"column:RegCis"`
	Nazev1  string `gorm:"column:Nazev1"`
}

func (hvwZapsiArtikl) TableName() string {
	return "hvw_zapsi_artikl"
}
