package errs

import "errors"

/*
================ Validation failures ================
*/
var (
	ErrInvalidCredentials  = errors.New("invalid login or password")
	ErrCannotCreateBooking = errors.New("specified slot is in the past")
	ErrCannotCancelBooking = errors.New("booking can be cancelled only by owner")

	ErrInvalidInput = errors.New("invalid input") // for rich models
)

/*
================ Infrastructure failures ================
*/
var (
	ErrHashPassword  = errors.New("failed to hash password")
	ErrGenerateToken = errors.New("failed to generate token")
	ErrCreateMeeting = errors.New("failed to create conference link")
)

/*
================ Database failures ================
*/
var (
	ErrCreateLocationDB                 = errors.New("failed to create location using db")
	ErrGetAdminByLoginDB                = errors.New("failed to get admin by login using db")
	ErrDeleteLocationDB                 = errors.New("failed to delete location using db")
	ErrDeleteLocationItemByLocationIDDB = errors.New("failed to delete location items by location id using db")
	ErrCreateRoomDB                     = errors.New("failed to create room using db")
	ErrGetRoomDB                        = errors.New("failed to get room using db")
	ErrListRoomsDB                      = errors.New("failed to get a list of rooms using db")
	ErrCreateScheduleDB                 = errors.New("failed to create schedule using db")
	ErrGetScheduleDB                    = errors.New("failed to get schedule using db")
	ErrCreateSlotsDB                    = errors.New("failed to create slots using db")
	ErrGetSlotDB                        = errors.New("failed to get slot using db")
	ErrListLocationItemsDB              = errors.New("failed to get a list of location items using db")
	ErrCreateBookingDB                  = errors.New("failed to create booking using db")
	ErrGetBookingDB                     = errors.New("failed to get booking using db")
	ErrUpdateBookingStatusDB            = errors.New("failed to update booking status using db")
	ErrListItemsByIDsDB                 = errors.New("failed to get a list of items by ids using db")
	ErrListMyBookingsDB                 = errors.New("failed to get a list of bookings by user id using db")

	ErrLocationNotFound = errors.New("location not found")
	ErrScheduleNotFound = errors.New("schedule not found")
	ErrSlotNotFound     = errors.New("slot not found")
	ErrBookingNotFound  = errors.New("booking not found")

	ErrLocationAlreadyExists = errors.New("location with given slug already exists")
	ErrScheduleAlreadyExists = errors.New("schedule for this room already exists")
)
