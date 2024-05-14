package storage

import (
	"context"
	"errors"
	"lamoda/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

var (
	getStockAvailability = `
		SELECT is_available
		FROM stocks
		WHERE stock_id = $1`

	getAvailableQty = `
		SELECT available_qty
		FROM product_on_stock
		WHERE stock_id = $1
		AND product_id = $2;`
	reserveProduct = `
		UPDATE product_on_stock
		SET available_qty = available_qty - $3,
			reserved_qty = reserved_qty + $3
		WHERE stock_id  = $1
		AND product_id = $2;`

	getReservedQty = `
		SELECT reserved_qty
		FROM product_on_stock
		WHERE stock_id = $1
		AND product_id = $2;`
	deleteReservation = `
		UPDATE product_on_stock
		SET available_qty = available_qty + $3,
			reserved_qty = reserved_qty - $3
		WHERE stock_id  = $1
		AND product_id = $2;`

	// добавить наименование товара
	selectStockProducts = `
		SELECT product_id, available_qty
		 FROM product_on_stock
		 WHERE stock_id = $1`
)

type storer interface {
	ReserveProducts(ctx context.Context, ReservedProducts model.ReservedProducts) error
	DeleteReservation(ctx context.Context, ReservedProducts model.ReservedProducts) error
	GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error)
	CheckStockAvailability(ctx context.Context, stockID int) error
	Close()
}

func NewStorage(connString string, log *logrus.Logger) (storer, error) {

	log.Info("Запускаем инициализацию хранилища")
	pool, err := pgxpool.New(context.Background(), connString)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return &StoragePostgres{
		pgxPool: pool,
		log:     log,
	}, nil
}

type StoragePostgres struct {
	pgxPool *pgxpool.Pool
	log     *logrus.Logger
}

func (s *StoragePostgres) ReserveProducts(ctx context.Context, ReservedProducts model.ReservedProducts) error {

	s.log.Debug("Резервирование товаров на складе для доставки")

	for _, product := range ReservedProducts.Products {
		tx, err := s.pgxPool.Begin(ctx)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}
		row := tx.QueryRow(ctx, getAvailableQty, ReservedProducts.StockID, product.ProductID)
		var availableQty int64
		err = row.Scan(&availableQty)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrProductNotFound
			}
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
		if availableQty < product.Quantity {
			return model.ErrNotEnoughAvailableQty
		}
		_, err = tx.Exec(ctx, reserveProduct, ReservedProducts.StockID, product.ProductID, product.Quantity)
		if err != nil {
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
	}
	return nil
}

func (s *StoragePostgres) DeleteReservation(ctx context.Context, ReservedProducts model.ReservedProducts) error {
	for _, product := range ReservedProducts.Products {
		tx, err := s.pgxPool.Begin(ctx)
		if err != nil {
			s.log.Error(err.Error())
			return err
		}
		row := tx.QueryRow(ctx, getReservedQty, ReservedProducts.StockID, product.ProductID)
		var reservedQty int64
		err = row.Scan(&reservedQty)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return model.ErrProductNotFound
			}
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
		if reservedQty < product.Quantity {
			product.Quantity = reservedQty
		}
		_, err = tx.Exec(ctx, deleteReservation, ReservedProducts.StockID, product.ProductID, product.Quantity)
		if err != nil {
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
		err = tx.Commit(ctx)
		if err != nil {
			s.log.Error(err.Error())
			tx.Rollback(ctx)
			return err
		}
	}
	return nil
}

func (s *StoragePostgres) GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error) {

	rows, err := s.pgxPool.Query(ctx, selectStockProducts, stockID)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	var product model.Products
	var availableProducts []model.Products
	for rows.Next() {
		err = rows.Scan(&product.ProductID, &product.Quantity)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		availableProducts = append(availableProducts, product)
	}
	return availableProducts, nil
}

func (s *StoragePostgres) CheckStockAvailability(ctx context.Context, stockID int) error {

	s.log.Debug("Проверка доступности склада")
	var isAvailable bool
	row := s.pgxPool.QueryRow(ctx, getStockAvailability, stockID)
	err := row.Scan(&isAvailable)
	if err != nil {
		s.log.Error(err.Error())
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrStockNotFound
		}
		return err
	}
	if !isAvailable {
		return model.ErrStockIsNotAvailable
	}
	return nil
}

func (s *StoragePostgres) Close() {
	s.pgxPool.Close()
}
