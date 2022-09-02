package handlers

import (
	"github.com/MykhailoKondrat/go-micro/data"
	"log"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		p.getProducts(w, r)
		return
	}
	if r.Method == http.MethodPost {
		p.addProduct(w, r)
		return
	}
	if r.Method == http.MethodPut {
		re := regexp.MustCompile(`/([0-9]+)`)
		g := re.FindAllStringSubmatch(r.URL.Path, -1)
		if len(g) != 1 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		if len(g[0]) != 2 {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := g[0][1]
		id, err := strconv.Atoi(idString)
		if err != nil {
			http.Error(w, "Invalid URI", http.StatusBadRequest)
			return
		}
		p.l.Println("got id", id)
		p.updateProduct(id, w, r)
		return
	}
	w.WriteHeader(http.StatusNotImplemented)
}
func (p *Products) getProducts(w http.ResponseWriter, h *http.Request) {
	lp := data.GetProducts()
	err := lp.ToJSON(w)
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
	}
}
func (p *Products) addProduct(w http.ResponseWriter, r *http.Request) {
	p.l.Println("handle post product")
	prod := &data.Product{}
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
	}
	data.AddProduct(prod)
	p.l.Printf("Prod:%#v", prod)

}
func (p *Products) updateProduct(id int, w http.ResponseWriter, r *http.Request) {

	p.l.Println("handle post product")
	prod := &data.Product{}
	if err := prod.FromJSON(r.Body); err != nil {
		http.Error(w, "Unable to unmarshall JSON", http.StatusBadRequest)
	}
	err := data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "Product not Found", http.StatusNotFound)
		return
	}
}
