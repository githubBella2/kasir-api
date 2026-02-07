package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/repositories"
	"kasir-api/services"

	"github.com/spf13/viper"
)

type Produk struct {
	ID    int    `json:"id" validate:"min:10"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var produk = []Produk{
	{ID: 1, Name: "Indomie Rebus", Price: 3500, Stock: 10},
	{ID: 2, Name: "Martabak Kanji", Price: 3000, Stock: 40},
}

func main() {
	// viper.SetConfigFile(".env")
	// viper.SetConfigType("env")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:    viper.GetString("PORT"),
		DB_CONN: viper.GetString("DB_CONN"),
	}

	// Setup database
	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	productRepo := repositories.NewProductRepository(db)
	ProductService := services.NewProductService(productRepo)
	ProductHandler := handlers.NewProductHandler(ProductService)

	//setup route
	http.HandleFunc("/api/products", ProductHandler.HandleProducts)
	http.HandleFunc("/api/products/", ProductHandler.HandleProductByID)

	//Transaction
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", transactionHandler.HandleCheckout) // POST

	// DEBUG WAJIB
	fmt.Println("PORT =", config.Port)
	fmt.Println("DB_CONN =", config.DB_CONN)

	if config.DB_CONN == "" {
		log.Fatal("DB_CONN KOSONG — .env tidak terbaca")
	}

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

	fmt.Println("Server running di localhost:" + config.Port)
	err = http.ListenAndServe(":"+config.Port, nil)
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
				"message": "sukses delete",
			})
			return
		}
	}
	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// ============================
type Config struct {
	Port    string `mapstructure:"PORT"`
	DB_CONN string `mapstructure:"DB_CONN"`
}

// Setup database
