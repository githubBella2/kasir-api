package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	//inisialilasi subtotal -> jumlah total transaksi keseluruhan
	totalAmount := 0
	//inisialisasi modeling transactiondetails -> nanti kita insert ke db
	details := make([]models.TransactionDetails, 0)
	//loop setiap item
	//get product dapet pricing
	//hitung current total = qty * price
	//ditambahin ke dlm subtotal
	//kurangi jumlah stok
	//item nya dimasukkin ke transction details
	//insert trascation
	//insert trasanction details

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetails{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	// err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	// if err != nil {
	// 	return nil, err
	// }

	err = tx.QueryRow("INSERT INTO transactions (total_amount, created_at) VALUES ($1, NOW()) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// for i := range details {
	// 	details[i].TransactionID = transactionID
	// 	_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
	// 		transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// }

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	} 

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (r *TransactionRepository) GetTodayReport() (*models.TodayReport, error) {
    today := time.Now().Truncate(24 * time.Hour)

    // Total revenue
    var totalRevenue sql.NullInt64
    err := r.db.QueryRow("SELECT COALESCE(SUM(total_amount), 0) FROM transactions WHERE created_at >= $1", today).Scan(&totalRevenue)
    if err != nil {
        return nil, err
    }

    // Total transaksi
    var totalTransaksi int
    err = r.db.QueryRow("SELECT COUNT(*) FROM transactions WHERE created_at >= $1", today).Scan(&totalTransaksi)
    if err != nil {
        return nil, err
    }

    // Produk terlaris - pakai struct local dulu
    type topProductStruct struct {
        name     string
        totalQty int
    }
    var topProduct topProductStruct
    
    err = r.db.QueryRow(`
        SELECT p.name, COALESCE(SUM(td.quantity), 0) 
        FROM transactions t 
        JOIN transaction_details td ON t.id = td.transaction_id 
        JOIN products p ON td.product_id = p.id 
        WHERE t.created_at >= $1 
        GROUP BY p.id, p.name 
        ORDER BY SUM(td.quantity) DESC 
        LIMIT 1
    `, today).Scan(&topProduct.name, &topProduct.totalQty)

    // Default value kalau belum ada data
    report := &models.TodayReport{
        TotalRevenue:   int(totalRevenue.Int64),
        TotalTransaksi: totalTransaksi,
        ProdukTerlaris: models.ProdukTerlaris{Nama: "", QtyTerjual: 0},
    }

    // Kalau ada produk terlaris, update
    if err == nil {
        report.ProdukTerlaris = models.ProdukTerlaris{
            Nama:       topProduct.name,
            QtyTerjual: topProduct.totalQty,
        }
    }

    return report, nil
}
