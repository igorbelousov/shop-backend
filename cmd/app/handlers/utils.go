package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/igorbelousov/shop-backend/foundation/web"
)

type utilsGroup struct {
}

func (ug utilsGroup) Upload(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	filename := uuid.New().String() + ".png"

	file, _, err := r.FormFile("file")
	if err != nil {
		log.Println("File no in request: ", err)
	}

	web.Upload("./media/", file, filename)

	return web.Respond(ctx, w, nil, http.StatusOK)
}
