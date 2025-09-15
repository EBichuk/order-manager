package service

import (
	"errors"
	"log/slog"
	"order-manager/internal/models"
	"order-manager/pkg/errorx"

	"github.com/go-playground/validator/v10"
)

type repository interface {
	GetOrderByUID(string) (*models.Order, error)
	SaveOrder(*models.Order) error
	GetAllOrders(int) ([]models.Order, error)
}

type cache interface {
	SetOrder(models.Order)
	GetOrder(string) (models.Order, bool)
}

type Service struct {
	r         repository
	c         cache
	log       *slog.Logger
	validator *validator.Validate
}

func NewService(r repository, c cache, log *slog.Logger) *Service {
	return &Service{
		r:         r,
		c:         c,
		log:       log,
		validator: validator.New(),
	}
}

func (s *Service) GetOrderByUID(orderUID string) (*models.Order, error) {
	if order, found := s.c.GetOrder(orderUID); found {
		slog.Info("Got order from cache", slog.String("order_uid", orderUID))
		return &order, nil
	}

	order, err := s.r.GetOrderByUID(orderUID)
	if err != nil {
		if errors.Is(err, errorx.ErrOrderNotFound) {
			s.log.Warn("Order not found", slog.String("order_uid", orderUID))
			return nil, err
		}
		s.log.Error("Failed to get order", slog.String("error", err.Error()))
		return nil, errorx.ErrInternal
	}

	s.log.Info("Got order from db", slog.String("order_uid", orderUID))
	return order, nil
}

func (s *Service) SaveOrder(order *models.Order) error {
	err := s.validator.Struct(order)
	if err != nil {
		s.log.Error("Error of validation order", slog.String("error", err.Error()), slog.String("order_uid", order.OrderUID))
		return errorx.ErrOrderValidation
	}

	err = s.r.SaveOrder(order)
	if err != nil {
		s.log.Error("Failed to save order", slog.String("error", err.Error()))
		return errorx.ErrInternal
	}

	s.c.SetOrder(*order)

	return nil
}

func (s *Service) FillCache(size int) error {
	orders, err := s.r.GetAllOrders(size)
	if err != nil {
		return err
	}

	for _, order := range orders {
		s.c.SetOrder(order)
	}

	return nil
}
