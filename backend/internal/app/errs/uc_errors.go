package errs

import "errors"

/*
================ Validation failures ================
*/
var (
	ErrCannotActivateLocation   = errors.New("location is already activated")
	ErrCannotDeactivateLocation = errors.New("location is already deactivated")
	ErrCannotGetLocationQRCode  = errors.New("cannot get qr-code: location is not operational")
	ErrCannotCreateOrder        = errors.New("cannot create order: location is not operational")
	ErrCannotSellItem           = errors.New("cannot sell one of the chosen items")

	ErrInvalidCredentials = errors.New("invalid login or password")
	ErrInvalidInput       = errors.New("invalid input") // for rich models
)

/*
================ Infrastructure failures ================
*/
var (
	ErrGenerateToken  = errors.New("failed to generate token")
	ErrGenerateQRCode = errors.New("failed to generate qr code")
)

/*
================ Database failures ================
*/
var (
	ErrGetAdminByLoginDB = errors.New("failed to get admin by login using db")

	ErrCreateLocationDB  = errors.New("failed to create location using db")
	ErrGetLocationByIDDB = errors.New("failed to get location by id using db")
	ErrUpdateLocationDB  = errors.New("failed to update location using db")
	ErrDeleteLocationDB  = errors.New("failed to delete location using db")
	ErrListLocationsDB   = errors.New("failed to get a list of locations using db")

	ErrCreateItemDB     = errors.New("failed to create item using db")
	ErrGetItemDB        = errors.New("failed to get item using db")
	ErrUpdateItemDB     = errors.New("failed to update item using db")
	ErrDeleteItemDB     = errors.New("failed to delete item using db")
	ErrListAllItemsDB   = errors.New("failed to get a list of items by ids using db")
	ErrListItemsByIDsDB = errors.New("failed to get a list of all items using db")

	ErrCreateLocationItemDB               = errors.New("failed to create location item using db")
	ErrGetLocationItemByLocationAndItemDB = errors.New("failed to get location item by location id and item id")
	ErrUpdateLocationItemDB               = errors.New("failed to update location item using db")
	ErrDeleteLocationItemsByItemIDDB      = errors.New("failed to delete location items by item id using db")
	ErrDeleteLocationItemByLocationIDDB   = errors.New("failed to delete location items by location id using db")
	ErrListLocationItemsDB                = errors.New("failed to get a list of location items using db")

	ErrCreateOrderDB = errors.New("failed to create order using db")
	ErrListRoomsDB   = errors.New("failed to get a list of rooms using db")
	ErrGetScheduleDB = errors.New("failed to get schedule using db")
	ErrGetSlotDB     = errors.New("failed to get slot using db")
	ErrGetBookingDB  = errors.New("failed to get booking using db")

	ErrCreateOrderItemsDB = errors.New("failed to create order items using db")

	ErrCreateTransactionDB = errors.New("failed to create transaction using db")

	ErrLocationNotFound     = errors.New("location not found")
	ErrItemNotFound         = errors.New("item not found")
	ErrLocationItemNotFound = errors.New("location item not found")
	ErrSlotNotFound         = errors.New("slot not found")
	ErrBookingNotFound      = errors.New("booking not found")

	ErrLocationAlreadyExists    = errors.New("location with given slug already exists")
	ErrOrderAlreadyExists       = errors.New("order already exists")
	ErrTransactionAlreadyExists = errors.New("transaction with given slug already exists")
)

/*
================ Payment Gateway failures ================
*/

var (
	ErrCreatePayment = errors.New("failed to create a payment")
)
