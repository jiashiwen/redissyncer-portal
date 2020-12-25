package httpserver

type server interface {
	ListenAndServe() error
}
