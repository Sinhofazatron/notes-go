package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type appHandler func (w http.ResponseWriter, r *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		err := h(w, r)
		otherError := fmt.Errorf("NoAuthErr")

		if err != nil {
			w.Header().Set("Content-Type", "application/json")

			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					w.Write(ErrNotFound.Marshal())

					return
				} else {
					if errors.Is(err, otherError) {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write(ErrNotFound.Marshal())

						return
					}
				}

				err = err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				w.Write(appErr.Marshal())
			}

			w.WriteHeader(http.StatusTeapot)
			w.Write(systemError(err).Marshal())
		}
	}
}
