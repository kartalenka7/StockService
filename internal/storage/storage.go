package storage

import (
	"context"
	"errors"
	"lamoda/internal/model"

	"github.com/jackc/pgx/v5"
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

func NewStorage(conn *pgx.Conn, log *logrus.Logger) (*StoragePostgres, error) {

	return &StoragePostgres{
		pgxConn: conn,
		log:     log,
	}, nil
}

type StoragePostgres struct {
	pgxConn *pgx.Conn
	log     *logrus.Logger
}

func (s *StoragePostgres) ReserveProduct(ctx context.Context,
	stockId int, product model.Products) error {

	tx, err := s.pgxConn.Begin(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	row := tx.QueryRow(ctx, getAvailableQty, stockId, product.ProductID)
	var availableQty int64
	err = row.Scan(&availableQty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrProductNotFound
		}
		s.log.Error(err.Error())
		return err
	}
	if availableQty < product.Quantity {
		return model.ErrNotEnoughAvailableQty
	}
	_, err = tx.Exec(ctx, reserveProduct, stockId, product.ProductID, product.Quantity)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *StoragePostgres) DeleteReservation(ctx context.Context,
	stockId int, product model.Products) error {

	tx, err := s.pgxConn.Begin(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()
	row := tx.QueryRow(ctx, getReservedQty, stockId, product.ProductID)
	var reservedQty int64
	err = row.Scan(&reservedQty)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ErrProductNotFound
		}
		s.log.Error(err.Error())
		return err
	}
	if reservedQty < product.Quantity {
		product.Quantity = reservedQty
	}
	_, err = tx.Exec(ctx, deleteReservation, stockId, product.ProductID, product.Quantity)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}
	err = tx.Commit(ctx)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	return nil
}

func (s *StoragePostgres) GetAvailableQty(ctx context.Context, stockID int) ([]model.Products, error) {

	rows, err := s.pgxConn.Query(ctx, selectStockProducts, stockID)
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

	var isAvailable bool
	row := s.pgxConn.QueryRow(ctx, getStockAvailability, stockID)
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
