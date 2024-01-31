package main
	
import (
	"net/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/middleware"

	"github.com/heptagons/lenses/symm/server"
)

// go run ./symm/server/cmd/.
func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	server.New(r)
	http.ListenAndServe(":8080", r)
}