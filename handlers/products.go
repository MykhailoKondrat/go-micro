package handlers

import (
	"context"
	"github.com/MykhailoKondrat/go-micro/data"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//if r.Method == http.MethodGet {
	//	p.getProducts(w, r)
	//	return
	//}
	//if r.Method == http.MethodPost {
	//	p.addProduct(w, r)
	//	return
	//}
	//if r.Method == http.MethodPut {
	//	re := regexp.MustCompile(`/([0-9]+)`)
	//	g := re.FindAllStringSubmatch(r.URL.Path, -1)
	//	if len(g) != 1 {
	//		http.Error(w, "Invalid URI", http.StatusBadRequest)
	//		return
	//	}
	//	if len(g[0]) != 2 {
	//		http.Error(w, "Invalid URI", http.StatusBadRequest)
	//		return
	//	}
	//
	//	idString := g[0][1]
	//	id, err := strconv.Atoi(idString)
	//	if err != nil {
	//		http.Error(w, "Invalid URI", http.StatusBadRequest)
	//		return
	//	}
	//	p.l.Println("got id", id)
	//	p.updateProduct(id, w, r)
	//	return
	//}
	//w.WriteHeader(http.StatusNotImplemented)
}
func (p *Products) GetProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle post product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)
	data.AddProduct(&prod)
	p.l.Printf("Prod:%#v", prod)

}
func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
	}
	p.l.Println("handle PUT product")
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	//prod := &data.Product{}
	//if err := prod.FromJSON(r.Body); err != nil {
	//	http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
	//}
	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
}

type KeyProduct struct{}

func (p *Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		prod := data.Product{}
		if err := prod.FromJSON(r.Body); err != nil {
			http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), KeyProduct{}, prod)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}
