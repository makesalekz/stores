package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	v1 "gitlab.calendaria.team/services/stores/api/stores/v1"
	"gitlab.calendaria.team/services/stores/ent"
	"gitlab.calendaria.team/services/stores/internal/biz"
	"gitlab.calendaria.team/services/stores/internal/data/mock"
	"gitlab.calendaria.team/services/stores/internal/service"
	"gitlab.calendaria.team/services/utils/v2/auth"
	"gitlab.calendaria.team/services/utils/v2/zap"
)

func beforeStoresTest(t *testing.T) (
	context.Context,
	*service.StoresService,
	*gomock.Controller,
	*mock.MockStoresRepo,
	int64,
	int64,
) {
	logger := zap.NewZapLogger(true)
	ctrl := gomock.NewController(t)
	storesRepo := mock.NewMockStoresRepo(ctrl)

	storesUsecase, err := biz.NewStoresUsecase(logger, storesRepo)
	require.NoError(t, err)

	storesService := service.NewStoresService(storesUsecase)

	var tenantID int64 = 12
	var actorID int64 = 332
	ctx := auth.NewTenantContext(auth.NewActorContext(context.Background(), actorID), tenantID)

	return ctx, storesService, ctrl, storesRepo, tenantID, actorID
}

func TestCreateStore(t *testing.T) {
	ctx, svc, ctrl, repo, tenantID, _ := beforeStoresTest(t)
	defer ctrl.Finish()

	lat := 43.238949
	lon := 76.945465

	repo.EXPECT().CreateStore(gomock.Any(), gomock.Any()).DoAndReturn(
		func(_ context.Context, dto interface{}) (*ent.Store, error) {
			return &ent.Store{
				ID:        1,
				TenantID:  tenantID,
				Name:      "Store 1",
				Address:   "Main St 1",
				Lat:       &lat,
				Lon:       &lon,
				Phone:     "+7777",
				WorkHours: "9-18",
				IsActive:  true,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}, nil
		},
	)

	reply, err := svc.CreateStore(ctx, &v1.CreateStoreRequest{
		Name:      "Store 1",
		Address:   "Main St 1",
		Lat:       &lat,
		Lon:       &lon,
		Phone:     "+7777",
		WorkHours: "9-18",
	})

	require.NoError(t, err)
	require.Equal(t, int64(1), reply.Store.Id)
	require.Equal(t, "Store 1", reply.Store.Name)
	require.Equal(t, tenantID, reply.Store.TenantId)
	require.NotNil(t, reply.Store.Lat)
	require.Equal(t, lat, *reply.Store.Lat)
}

func TestCreateStore_EmptyTenantID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := zap.NewZapLogger(true)
	storesRepo := mock.NewMockStoresRepo(ctrl)
	storesUsecase, _ := biz.NewStoresUsecase(logger, storesRepo)
	storesService := service.NewStoresService(storesUsecase)

	// Context without tenant
	ctx := context.Background()

	_, err := storesService.CreateStore(ctx, &v1.CreateStoreRequest{
		Name: "Store",
	})
	require.Error(t, err)
}

func TestSetStoreResponsible(t *testing.T) {
	ctx, svc, ctrl, repo, tenantID, _ := beforeStoresTest(t)
	defer ctrl.Finish()

	var memberID int64 = 5
	repo.EXPECT().SetResponsible(gomock.Any(), int64(1), tenantID, memberID).Return(
		&ent.Store{
			ID:            1,
			TenantID:      tenantID,
			Name:          "Store 1",
			ResponsibleID: &memberID,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}, nil,
	)

	reply, err := svc.SetStoreResponsible(ctx, &v1.SetStoreResponsibleRequest{
		StoreId:  1,
		MemberId: memberID,
	})

	require.NoError(t, err)
	require.NotNil(t, reply.Store.ResponsibleId)
	require.Equal(t, memberID, *reply.Store.ResponsibleId)
}

func TestSetStoreResponsible_StoreNotFound(t *testing.T) {
	ctx, svc, ctrl, repo, tenantID, _ := beforeStoresTest(t)
	defer ctrl.Finish()

	var memberID int64 = 999
	repo.EXPECT().SetResponsible(gomock.Any(), int64(1), tenantID, memberID).Return(
		nil, &ent.NotFoundError{},
	)

	_, err := svc.SetStoreResponsible(ctx, &v1.SetStoreResponsibleRequest{
		StoreId:  1,
		MemberId: memberID,
	})

	require.Error(t, err)
}

func TestGetStoresByCoordinates(t *testing.T) {
	ctx, svc, ctrl, repo, _, _ := beforeStoresTest(t)
	defer ctrl.Finish()

	lat1 := 43.238949
	lon1 := 76.945465

	repo.EXPECT().GetStoresByIDs(gomock.Any(), []int64{1, 2}).Return(
		[]*ent.Store{
			{ID: 1, Name: "Store 1", Lat: &lat1, Lon: &lon1, TenantID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()},
			{ID: 2, Name: "Store 2", TenantID: 12, CreatedAt: time.Now(), UpdatedAt: time.Now()},
		}, nil,
	)

	reply, err := svc.GetStoresByCoordinates(ctx, &v1.GetStoresByCoordinatesRequest{
		StoreIds: []int64{1, 2},
	})

	require.NoError(t, err)
	require.Len(t, reply.Stores, 2)
	require.Equal(t, int64(1), reply.Stores[0].Id)
	require.NotNil(t, reply.Stores[0].Lat)
	require.Equal(t, lat1, *reply.Stores[0].Lat)
	require.Nil(t, reply.Stores[1].Lat)
}
