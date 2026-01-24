package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Produk struct {
	ID    int    `json:"id" validate:"min:10"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

var produk = []Produk{
	{ID: 1, Nama: "Indomie Rebus", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Martabak Kanji", Harga: 3000, Stok: 40},
}

func main() {
	// POST
	// GET
	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// Krn mainnya masih di In Memory, maka
			// Baca data dari request
			// Masukin data ke dalam variabel produk
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
				return
			}

			//masukkin data ke dalam variabel produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			getProdukByID(w, r)

		} else if r.Method == "PUT" { // //Update Produk
			updateProduk(w, r)
		} else if r.Method == "DELETE" {
			deleteProduk(w, r)
		}

	})

	http.HandleFunc("/api/produks", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			w.Header().Set("Content-Type", "application/json") // -> Set biar jadi JSON
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == "POST" {
			// Baca data dari request
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid Request", http.StatusBadRequest)
			}

			// masukan data ke dalam variabel produk
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated) // 201
			json.NewEncoder(w).Encode(produkBaru)
		}
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
		w.Write([]byte("OK"))
	})

	fmt.Println("Server running di localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("gagal running server:", err)
	}
}

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	if idStr == "" {
		http.Error(w, "ID Tidak boleh kosong", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr) //-> diubah ke Int
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json") // -> Set biar jadi JSON
			json.NewEncoder(w).Encode(p)
			return
		}
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	//GET id dulu
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	//ganti ke int
	if idStr == "" {
		http.Error(w, "ID Tidak boleh kosong", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr) //-> diubah ke Int
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// get data dari request
	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	//loop produk , cari id , ganti sesuai data dari request
	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)

}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	//GET id dulu
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	//ganti ke int
	if idStr == "" {
		http.Error(w, "ID Tidak boleh kosong", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(idStr) //-> diubah ke Int
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	
	//loop produk , cari id yang mau dihapus
	for i := range produk {
		if produk[i].ID == id {
			//bikin slice baru dengan data sebelum dan sesudah index
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message":"sukses delete",
			})
			return
		}
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}
