package mapper

import (
	ucerrs "backend/internal/app/errs"
	pkgerrs "backend/pkg/errs"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func HttpError(err error) *pkgerrs.OutErr {
	if err == nil {
		return nil
	}

	var w *ucerrs.WrappedError
	if errors.As(err, &w) {
		switch {
		case errors.Is(err, ucerrs.ErrGetAdminByLoginDB),
			errors.Is(err, ucerrs.ErrCreateLocationDB),
			errors.Is(err, ucerrs.ErrGetLocationByIDDB),
			errors.Is(err, ucerrs.ErrUpdateLocationDB),
			errors.Is(err, ucerrs.ErrDeleteLocationDB),
			errors.Is(err, ucerrs.ErrListLocationsDB),
			errors.Is(err, ucerrs.ErrCreateItemDB),
			errors.Is(err, ucerrs.ErrDeleteItemDB),
			errors.Is(err, ucerrs.ErrListAllItemsDB),
			errors.Is(err, ucerrs.ErrListItemsByIDsDB),
			errors.Is(err, ucerrs.ErrCreateLocationItemDB),
			errors.Is(err, ucerrs.ErrDeleteLocationItemsByItemIDDB),
			errors.Is(err, ucerrs.ErrDeleteLocationItemByLocationIDDB),
			errors.Is(err, ucerrs.ErrListLocationItemsDB),
			errors.Is(err, ucerrs.ErrGenerateToken),
			errors.Is(err, ucerrs.ErrGenerateQRCode):
			return pkgerrs.NewOutError(
				http.StatusInternalServerError,
				"internal error",
				w.Reason,
			)

		case errors.Is(err, ucerrs.ErrInvalidInput):
			return pkgerrs.NewOutError(
				http.StatusBadRequest,
				w.Public.Error(),
				w.Reason,
			)

		default:
			return pkgerrs.NewOutError(
				http.StatusInternalServerError,
				"internal error",
				w.Reason,
			)
		}
	}

	var echoErr *echo.HTTPError
	if errors.As(err, &echoErr) {
		msg := fmt.Sprintf("%v", echoErr.Message)

		if echoErr.Code == http.StatusBadRequest {
			msg = "invalid json"
		}

		return pkgerrs.NewOutError(
			echoErr.Code,
			msg,
			err,
		)
	}

	switch {
	case errors.Is(err, ucerrs.ErrInvalidCredentials):
		return pkgerrs.NewOutError(
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrLocationNotFound),
		errors.Is(err, ucerrs.ErrItemNotFound):
		return pkgerrs.NewOutError(
			http.StatusNotFound,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrCannotActivateLocation),
		errors.Is(err, ucerrs.ErrCannotDeactivateLocation),
		errors.Is(err, ucerrs.ErrLocationAlreadyExists):
		return pkgerrs.NewOutError(
			http.StatusConflict,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrCannotGetLocationQRCode):
		return pkgerrs.NewOutError(
			http.StatusUnprocessableEntity,
			err.Error(),
			nil,
		)
	}

	return pkgerrs.NewOutError(
		http.StatusInternalServerError,
		"internal error",
		nil,
	)
}
