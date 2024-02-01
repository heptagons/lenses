package main
	
import (
	"fmt"
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
	port := ":8080"
	fmt.Printf("symm service running at %s\n", port)
	http.ListenAndServe(port, r)
}