package mapper

import (
	ucerrs "backend/internal/app/errs"
	pkgerrs "backend/pkg/errs"
	"errors"
	"net/http"
)

func HttpError(err error) *pkgerrs.OutErr {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, pkgerrs.ErrInvalidJSON):
		return pkgerrs.NewOutError(http.StatusBadRequest, "invalid json", err)

	case errors.Is(err, pkgerrs.ErrInvalidIdentifier):
		return pkgerrs.NewOutError(http.StatusBadRequest, "invalid identifier format", err)
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
			errors.Is(err, ucerrs.ErrGetItemDB),
			errors.Is(err, ucerrs.ErrUpdateItemDB),
			errors.Is(err, ucerrs.ErrDeleteItemDB),
			errors.Is(err, ucerrs.ErrListAllItemsDB),
			errors.Is(err, ucerrs.ErrListItemsByIDsDB),
			errors.Is(err, ucerrs.ErrCreateLocationItemDB),
			errors.Is(err, ucerrs.ErrGetLocationItemByLocationAndItemDB),
			errors.Is(err, ucerrs.ErrUpdateLocationItemDB),
			errors.Is(err, ucerrs.ErrDeleteLocationItemsByItemIDDB),
			errors.Is(err, ucerrs.ErrDeleteLocationItemByLocationIDDB),
			errors.Is(err, ucerrs.ErrListLocationItemsDB),
			errors.Is(err, ucerrs.ErrCreateOrderDB),
			errors.Is(err, ucerrs.ErrCreateOrderItemsDB),
			errors.Is(err, ucerrs.ErrGenerateToken),
			errors.Is(err, ucerrs.ErrGenerateQRCode),
			errors.Is(err, ucerrs.ErrCreatePayment),
			errors.Is(err, ucerrs.ErrCreateTransactionDB):
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

	switch {
	case errors.Is(err, ucerrs.ErrInvalidCredentials):
		return pkgerrs.NewOutError(
			http.StatusUnauthorized,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrLocationNotFound),
		errors.Is(err, ucerrs.ErrItemNotFound),
		errors.Is(err, ucerrs.ErrLocationItemNotFound):
		return pkgerrs.NewOutError(
			http.StatusNotFound,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrCannotActivateLocation),
		errors.Is(err, ucerrs.ErrCannotDeactivateLocation),
		errors.Is(err, ucerrs.ErrLocationAlreadyExists),
		errors.Is(err, ucerrs.ErrCannotSellItem),
		errors.Is(err, ucerrs.ErrOrderAlreadyExists),
		errors.Is(err, ucerrs.ErrTransactionAlreadyExists):
		return pkgerrs.NewOutError(
			http.StatusConflict,
			err.Error(),
			nil,
		)

	case errors.Is(err, ucerrs.ErrCannotGetLocationQRCode),
		errors.Is(err, ucerrs.ErrCannotCreateOrder):
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
