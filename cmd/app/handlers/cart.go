package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/igorbelousov/shop-backend/foundation/web"
	"github.com/igorbelousov/shop-backend/internal/data/product"
)

type cartGroup struct {
	product product.Product
}

func (bc cartGroup) query(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

	type CartRequest struct {
		ID  string `json:"id"`
		Qty int    `json:"qty"`
	}

	type CartItem struct {
		Product product.Info
		Qty     int
	}

	ci := []CartItem{}
	cr := []CartRequest{}

	v, ok := ctx.Value(web.KeyValues).(*web.Values)
	if !ok {
		return web.NewShutdownError("web value missing from context")
	}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)

	}
	err = json.Unmarshal(b, &cr)
	if err != nil {
		http.Error(w, err.Error(), 500)
	}

	for i, _ := range cr {
		br, err := bc.product.QueryByID(ctx, v.TraceID, cr[i].ID)
		if err != nil {
			return err
		}
		c := CartItem{
			Product: br,
			Qty:     cr[i].Qty,
		}
		ci = append(ci, c)
	}

	return web.Respond(ctx, w, ci, http.StatusOK)
}
