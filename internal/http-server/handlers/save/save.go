package save

import (
	"errors"
	"net/http"
	resp "zipinit/internal/lib/api/response"
	sl "zipinit/internal/lib/logger"

	"zipinit/internal/storage"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"golang.org/x/exp/slog"
)

type Request struct {
	URL   string `json:"url" validate:"required,url"`
	Alias string `json:"alias,omitempty"`
}

type Response struct {
	resp.Response
	ShortLink string `json:"shortLink,omitempty"`
}

type URLSaver interface {
	SaveUrl(urlToSave string, alias string) (string, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.url.save.New"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("Failed to decode request", sl.Err(err))
			render.JSON(w, r, resp.Error("Failed to decode request"))
			return
		}
		log.Info("Request decoded", slog.Any("request", req))

		if err := validator.New().Struct(req); err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("invalid request", sl.Err(err))
			render.JSON(w, r, resp.ValidationError(validateErr))
			return
		}

		alias := req.Alias

		shortLink, err := urlSaver.SaveUrl(req.URL, alias)
		if err != nil {
			if errors.Is(err, storage.ErrUrlExists) {
				log.Info("url already exists", slog.String("url", req.URL))
				render.JSON(w, r, resp.Error("url already exists"))
				return
			}
			log.Error("failed to save url", sl.Err(err))
			render.JSON(w, r, resp.Error("failed to save url"))
			return
		}

		log.Info("url saved")

		responseOK(w, r, shortLink)
	}
}

func responseOK(w http.ResponseWriter, r *http.Request, shortLink string) {
	render.JSON(w, r, Response{
		Response:  resp.OK(),
		ShortLink: shortLink,
	})
}
