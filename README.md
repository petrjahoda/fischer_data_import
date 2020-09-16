# Fischer Data Import Service

## Description
Go service that downloads user and product data from Fischer database and updates/creates users, products and product
 groupsin Zapsi

* Periodocity of download: 10 minutes
* Import only users with hvw_zapsi_zam.Delnik == 1
* Delete user.Rfid when hvw_zapsi_zam_cip.Primarni == 0
   
### User mapping:

|Fischer Name|Zapsi Name|
|------------------|------------------|
|hvw_zapsi_zam.Jmeno|user.FirstName|
|hvw_zapsi_zam.Prijmeni|user.Name|
|PAIR hvw_zapsi_zam.Alias|PAIR user.Login|
|hvw_zapsi_zam_cip.CC|user.Rfid|
|nothing|user.Barcode|
|nothing|user.Pin|
|always insert 1|user.UserTypeID|
|nothing|user.Email|
|nothing|user.Phone|
|always insert 2|UserRoleID|

### Product mapping:
    
|CSV Name|Zapsi Name, product table|
|------------------|------------------|
|hvw_zapsi_artikl.Nazev1|product.Name|
|PAIR hvw_zapsi_artikl.RegCis|PAIR product.Barcode|
|nothing|product.Cycle|
|nothing|product.IdleFromTime|
|always insert 1|product.ProductStatusID|
|nothing|product.Deleted|
|proper productGroupId|product.ProductGroupID|
|nothing|product.Cavity|


© 2020 Petr Jahoda
