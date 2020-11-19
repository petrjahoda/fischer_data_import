package main

import (
	"database/sql"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"strconv"
	"time"
)

func importData() {
	timer := time.Now()
	logInfo("MAIN", "Importing process started")
	zapsiUsers, zapsiProducts, downloadedFromZapsi := downloadDataFromZapsi()
	fischerUsers, fischerProducts, fischerChipsAsMap, downloadedFromFischer := downloadDataFromFischer()
	if downloadedFromZapsi && downloadedFromFischer {
		logInfo("MAIN", "Zapsi Users: "+strconv.Itoa(len(zapsiUsers)))
		logInfo("MAIN", "Zapsi Products: "+strconv.Itoa(len(zapsiProducts)))
		logInfo("MAIN", "Fischer Users: "+strconv.Itoa(len(fischerUsers)))
		logInfo("MAIN", "Fischer Products: "+strconv.Itoa(len(fischerProducts)))
		logInfo("MAIN", "Fischer Chips: "+strconv.Itoa(len(fischerChipsAsMap)))
		updatedUsers, createdUsers := processUsers(zapsiUsers, fischerUsers, fischerChipsAsMap)
		updatedProducts, createdProducts := processProducts(zapsiProducts, fischerProducts)
		logInfo("MAIN", "Updated users: "+strconv.Itoa(updatedUsers))
		logInfo("MAIN", "Created users: "+strconv.Itoa(createdUsers))
		logInfo("MAIN", "Updated products: "+strconv.Itoa(updatedProducts))
		logInfo("MAIN", "Created products: "+strconv.Itoa(createdProducts))
	}
	logInfo("MAIN", "Importing process complete, time elapsed: "+time.Since(timer).String())
}

func processProducts(zapsiProducts map[string]product, fischerProducts []hvwZapsiArtikl) (int, int) {
	timer := time.Now()
	logInfo("MAIN", "Processing products")
	updated := 0
	created := 0
	for _, fischerProduct := range fischerProducts {
		if serviceRunning {
			_, productInZapsi := zapsiProducts[fischerProduct.RegCis]
			if productInZapsi {
				updateProductInZapsi(fischerProduct)
				updated++
			} else {
				createProductInZapsi(fischerProduct)
				created++
			}
		}
	}
	logInfo("MAIN", "Products processed, time elapsed: "+time.Since(timer).String())
	return updated, created
}

func createProductInZapsi(fischerProduct hvwZapsiArtikl) {
	logInfo("MAIN", fischerProduct.Nazev1+": Product does not exist in Zapsi, creating...")
	productGroupId := getProductGroupId(fischerProduct)
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	var product product
	product.Name = fischerProduct.Nazev1
	product.Barcode = fischerProduct.RegCis
	product.ProductGroupID = productGroupId
	product.ProductStatusID = 1
	product.Cycle = float64(fischerProduct.ID)
	db.Save(&product)
}

func updateProductInZapsi(fischerProduct hvwZapsiArtikl) {
	productGroupId := getProductGroupId(fischerProduct)
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}

	db.Model(&product{}).Where(user{Barcode: fischerProduct.RegCis}).Updates(product{
		Name:           fischerProduct.Nazev1,
		ProductGroupID: productGroupId,
	})
}

func getProductGroupId(fischerProduct hvwZapsiArtikl) int {
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return 1
	}
	var existingProductGroup productGroup
	db.Where("Name like ?", fischerProduct.SkupZbo).Find(&existingProductGroup)
	if existingProductGroup.OID > 0 {
		return existingProductGroup.OID
	}
	logInfo("MAIN", "Product group "+fischerProduct.SkupZbo+" does not exist, creating ...")
	var newProductGroup productGroup
	newProductGroup.Name = fischerProduct.SkupZbo
	db.Save(&newProductGroup)
	var brandNewProductGroup productGroup
	db.Where("Name like ?", fischerProduct.SkupZbo).Find(&brandNewProductGroup)
	if brandNewProductGroup.OID > 0 {
		return brandNewProductGroup.OID
	}
	return 1
}

func processUsers(zapsiUsers map[string]user, fischerUsers []hvwZapsiZam, fischerChipsAsMap map[string]hvwZapsiZamCip) (int, int) {
	timer := time.Now()
	logInfo("MAIN", "Processing users")
	updated := 0
	created := 0
	for _, fischerUser := range fischerUsers {
		if serviceRunning {
			_, userInZapsi := zapsiUsers[fischerUser.Alias]
			if userInZapsi {
				updateUserInZapsi(fischerUser, zapsiUsers[fischerUser.Alias], fischerChipsAsMap)
				updated++
			} else {
				createUserInZapsi(fischerUser, fischerChipsAsMap)
				created++
			}
		}
	}
	logInfo("MAIN", "Users processed, time elapsed: "+time.Since(timer).String())
	return updated, created
}

func updateUserInZapsi(fischerUser hvwZapsiZam, zapsiUser user, fischerChipsAsMap map[string]hvwZapsiZamCip) {
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	userChip, exists := fischerChipsAsMap[fischerUser.Alias]
	var rfidToInsert = ""
	if userChip.Primarni == 1 && exists {
		rfidToInsert = userChip.CC
	}
	db.Model(&user{}).Where(user{Login: zapsiUser.Login}).Updates(user{
		FirstName:  fischerUser.Jmeno,
		Name:       fischerUser.Prijmeni,
		Rfid:       rfidToInsert,
		UserTypeID: sql.NullInt32{Int32: 1, Valid: true},
		UserRoleID: sql.NullInt32{Int32: 2, Valid: true},
	})
}

func createUserInZapsi(fischerUser hvwZapsiZam, fischerChipsAsMap map[string]hvwZapsiZamCip) {
	logInfo("MAIN", fischerUser.Jmeno+","+fischerUser.Prijmeni+": User does not exist in Zapsi, creating...")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return
	}
	userChip, exists := fischerChipsAsMap[fischerUser.Alias]
	var rfidToInsert = ""
	if userChip.Primarni == 1 && exists {
		rfidToInsert = userChip.CC
	}
	var user user
	user.Login = fischerUser.Alias
	user.FirstName = fischerUser.Jmeno
	user.Name = fischerUser.Prijmeni
	user.Rfid = rfidToInsert
	user.Barcode = ""
	user.Pin = ""
	user.UserTypeID = sql.NullInt32{Int32: 1, Valid: true}
	user.UserRoleID = sql.NullInt32{Int32: 2, Valid: true}
	db.Save(&user)
}

func downloadDataFromFischer() ([]hvwZapsiZam, []hvwZapsiArtikl, map[string]hvwZapsiZamCip, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Zapsi")
	db, err := gorm.Open(sqlserver.Open(fischerConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return nil, nil, nil, false
	}
	var users []hvwZapsiZam
	var products []hvwZapsiArtikl
	var chips []hvwZapsiZamCip
	db.Where("Delnik = ?", 1).Find(&users)
	db.Find(&products)
	db.Where("Primarni = ?", 1).Find(&chips)
	fischerChipsAsMap := make(map[string]hvwZapsiZamCip, len(chips))
	for _, chip := range chips {
		fischerChipsAsMap[chip.Alias] = chip
	}
	logInfo("MAIN", "Zapsi data downloaded, time elapsed: "+time.Since(timer).String())
	return users, products, fischerChipsAsMap, true
}

func downloadDataFromZapsi() (map[string]user, map[string]product, bool) {
	timer := time.Now()
	logInfo("MAIN", "Downloading data from Zapsi")
	db, err := gorm.Open(mysql.Open(zapsiConfig), &gorm.Config{})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
		logError("MAIN", "Problem opening database: "+err.Error())
		return nil, nil, false
	}
	var users []user
	var products []product
	db.Find(&users)
	db.Find(&products)
	returnProducts := make(map[string]product, len(products))
	returnUsers := make(map[string]user, len(users))
	for _, product := range products {
		returnProducts[product.Barcode] = product
	}
	for _, user := range users {
		returnUsers[user.Login] = user
	}
	logInfo("MAIN", "Zapsi data downloaded, time elapsed: "+time.Since(timer).String())
	return returnUsers, returnProducts, true
}
