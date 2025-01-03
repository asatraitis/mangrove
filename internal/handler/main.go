package handler

import "net/http"

type MainHandler interface{}
type mainHandler struct {
	*BaseHandler

	mux *http.ServeMux
}

func NewMainHandler(baseHandler *BaseHandler, mux *http.ServeMux) MainHandler {
	h := &mainHandler{
		BaseHandler: baseHandler,
		mux:         mux,
	}
	h.logger = h.logger.With().Str("subcomponent", "MainHandler").Logger()
	h.register()
	return h
}

func (h *mainHandler) register() {
	h.mux.Handle("GET /", http.FileServer(http.Dir("./dist/main")))
}
