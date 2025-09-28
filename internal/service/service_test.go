package service_test

import (
	"fmt"
	"log/slog"
	"math/rand/v2"
	"order-manager/internal/models"
	"order-manager/internal/service"
	mocks "order-manager/mock"
	"order-manager/pkg/errorx"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

func TestGetOrderByUID_FoundInCache(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()
	in := orderIn.OrderUID

	repo := mocks.NewMockrepository(ctl)
	cache := mocks.NewMockcache(ctl)

	cache.EXPECT().GetOrder(in).Return(*orderIn, true)

	service := service.NewService(repo, cache, logger)

	_, err := service.GetOrderByUID(in)
	require.NoError(t, err)
}

func TestGetOrderByUID_FoundInDB(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()
	in := orderIn.OrderUID

	repo := mocks.NewMockrepository(ctl)
	cache := mocks.NewMockcache(ctl)

	cache.EXPECT().GetOrder(in).Return(models.Order{}, false)
	repo.EXPECT().GetOrderByUID(in).Return(orderIn, nil)

	service := service.NewService(repo, cache, logger)

	_, err := service.GetOrderByUID(in)
	require.NoError(t, err)
}

func TestGetOrderByUID_OrderNotFoundInDB(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()
	in := orderIn.OrderUID

	repo := mocks.NewMockrepository(ctl)
	repo.EXPECT().GetOrderByUID(in).Return(nil, assert.AnError)

	cache := mocks.NewMockcache(ctl)
	cache.EXPECT().GetOrder(in).Return(models.Order{}, false)

	service := service.NewService(repo, cache, logger)
	_, err := service.GetOrderByUID(in)

	require.ErrorIs(t, err, errorx.ErrInternal)
}

func TestSaveOrder_Success(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()

	repo := mocks.NewMockrepository(ctl)
	repo.EXPECT().SaveOrder(orderIn).Return(nil)
	cache := mocks.NewMockcache(ctl)
	cache.EXPECT().SetOrder(*orderIn)

	service := service.NewService(repo, cache, logger)
	err := service.SaveOrder(orderIn)

	require.NoError(t, err)
}

func TestSaveOrder_OrderValidationError(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()
	orderIn.Payment.DeliveryCost = -100

	repo := mocks.NewMockrepository(ctl)
	cache := mocks.NewMockcache(ctl)

	service := service.NewService(repo, cache, logger)
	err := service.SaveOrder(orderIn)

	require.ErrorIs(t, err, errorx.ErrOrderValidation)
}

func TestSaveOrder_DBError(t *testing.T) {
	t.Parallel()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	orderIn := MakeRandomOrder()

	repo := mocks.NewMockrepository(ctl)
	repo.EXPECT().SaveOrder(orderIn).Return(fmt.Errorf("Error from repository"))
	cache := mocks.NewMockcache(ctl)

	service := service.NewService(repo, cache, logger)

	err := service.SaveOrder(orderIn)

	require.ErrorIs(t, err, errorx.ErrInternal)
}

func MakeRandomOrder() *models.Order {
	item := models.Item{
		ChrtID:      1000000 + rand.IntN(100000),
		TrackNumber: "WBTESTTRACK",
		Price:       1 + rand.IntN(1000),
		Rid:         uuid.New().String(),
		NameItem:    "Test Name Item",
		Sale:        rand.IntN(100),
		Size:        0,
		TotalPrice:  rand.IntN(10000),
		NmID:        0,
		Brand:       "Test Brand",
		Status:      202,
	}
	payment := models.Payment{
		Transaction:  uuid.New().String(),
		RequestID:    "",
		Currency:     "USD",
		Provider:     "wbpay",
		Amount:       rand.IntN(100),
		PaymentDt:    1000000 + rand.IntN(100000),
		Bank:         "test bank",
		DeliveryCost: rand.IntN(1000),
		GoodsTotal:   rand.IntN(10000),
		CustomFee:    0,
	}
	delivery := models.Delivery{
		Name:    "Test Testov",
		Phone:   "+79000000000",
		Zip:     "123456",
		City:    "Moscow",
		Address: "Lenina 10",
		Region:  "Moscow",
		Email:   "test-testov@gmail.com",
	}
	order := models.Order{
		OrderUID:          uuid.New().String(),
		TrackNumber:       "new",
		Entry:             "WBIL",
		Locate:            "en",
		InternalSignature: " ",
		CustomerID:        uuid.New().String(),
		DeliveryService:   "meest",
		Shardkey:          "0",
		SmID:              rand.IntN(100),
		DateCreated:       time.Now().UTC(),
		OffShard:          "1",
		Delivery:          delivery,
		Payment:           payment,
		Item:              []models.Item{item},
	}
	return &order
}
