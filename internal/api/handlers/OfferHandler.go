package handlers

type IOfferHandler interface{}

type OfferHandler struct{}

func NewOfferHandler() IOfferHandler {
	return &OfferHandler{}
}
