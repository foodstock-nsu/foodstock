package errs

import "errors"

/*
================ Validation failures ================
*/
var (
	ErrInvalidCredentials       = errors.New("invalid login or password")
	ErrCannotGetLocationQRCode  = errors.New("location is not operational")
	ErrCannotActivateLocation   = errors.New("location is already activated")
	ErrCannotDeactivateLocation = errors.New("location is already deactivated")
	ErrCannotCancelBooking      = errors.New("booking can be cancelled only by owner")

	ErrInvalidInput = errors.New("invalid input") // for rich models
)

/*
================ Infrastructure failures ================
*/
var (
	ErrHashPassword   = errors.New("failed to hash password")
	ErrGenerateToken  = errors.New("failed to generate token")
	ErrGenerateQRCode = errors.New("failed to generate qr code")
)

/*
================ Database failures ================
*/
var (
	ErrCreateLocationDB                 = errors.New("failed to create location using db")
	ErrGetAdminByLoginDB                = errors.New("failed to get admin by login using db")
	ErrDeleteLocationDB                 = errors.New("failed to delete location using db")
	ErrDeleteLocationItemByLocationIDDB = errors.New("failed to delete location items by location id using db")
	ErrCreateItemDB                     = errors.New("failed to create item using db")
	ErrGetLocationByIDDB                = errors.New("failed to get location by id using db")
	ErrListRoomsDB                      = errors.New("failed to get a list of rooms using db")
	ErrDeleteItemDB                     = errors.New("failed to delete item using db")
	ErrDeleteLocationItemsByItemIDDB    = errors.New("failed to delete location items by item id using db")
	ErrGetScheduleDB                    = errors.New("failed to get schedule using db")
	ErrCreateSlotsDB                    = errors.New("failed to create slots using db")
	ErrGetSlotDB                        = errors.New("failed to get slot using db")
	ErrListLocationItemsDB              = errors.New("failed to get a list of location items using db")
	ErrCreateBookingDB                  = errors.New("failed to create booking using db")
	ErrGetBookingDB                     = errors.New("failed to get booking using db")
	ErrUpdateLocationDB                 = errors.New("failed to update location using db")
	ErrListItemsByIDsDB                 = errors.New("failed to get a list of items by ids using db")
	ErrListLocationsDB                  = errors.New("failed to get a list of locations using db")

	ErrLocationNotFound = errors.New("location not found")
	ErrItemNotFound     = errors.New("item not found")
	ErrSlotNotFound     = errors.New("slot not found")
	ErrBookingNotFound  = errors.New("booking not found")

	ErrLocationAlreadyExists = errors.New("location with given slug already exists")
	ErrScheduleAlreadyExists = errors.New("schedule for this room already exists")
)
