package data

import (
	"context"

	"github.com/makesalekz/stores/ent"
	"github.com/makesalekz/stores/ent/store"
	utils_v1 "github.com/makesalekz/utils/api/utils/v1"
)

type StoreDto struct {
	ID            int64
	TenantID      int64
	Name          string
	Address       string
	Lat           *float64
	Lon           *float64
	Phone         string
	WorkHours     string
	IsActive      bool
	ResponsibleID *int64
}

type StoresListFilter struct {
	TenantID   int64
	OnlyActive bool
}

//go:generate mockgen -source=stores.go -destination=mock/stores.go -package=mock

// StoresRepo interface for stores repository.
type StoresRepo interface {
	CreateStore(ctx context.Context, dto StoreDto) (*ent.Store, error)
	UpdateStore(ctx context.Context, dto StoreDto) (*ent.Store, error)
	DeleteStore(ctx context.Context, id int64, tenantID int64) error
	GetStore(ctx context.Context, id int64, tenantID int64) (*ent.Store, error)
	ListStores(ctx context.Context, filter StoresListFilter, paginate *utils_v1.PaginateRequest) ([]*ent.Store, error)
	CountListStores(ctx context.Context, filter StoresListFilter) (int32, error)
	SetResponsible(ctx context.Context, storeID int64, tenantID int64, memberID int64) (*ent.Store, error)
	GetStoresByIDs(ctx context.Context, ids []int64) ([]*ent.Store, error)
}

type storesRepo struct {
	db *ent.Client
}

// NewStoresRepo creates a new stores repository.
func NewStoresRepo(d *Data) StoresRepo {
	return &storesRepo{
		db: d.db,
	}
}

func (r *storesRepo) CreateStore(ctx context.Context, dto StoreDto) (*ent.Store, error) {
	create := r.db.Store.Create().
		SetTenantID(dto.TenantID).
		SetName(dto.Name).
		SetAddress(dto.Address).
		SetPhone(dto.Phone).
		SetWorkHours(dto.WorkHours).
		SetIsActive(true)

	if dto.Lat != nil {
		create.SetLat(*dto.Lat)
	}
	if dto.Lon != nil {
		create.SetLon(*dto.Lon)
	}

	return create.Save(ctx)
}

func (r *storesRepo) UpdateStore(ctx context.Context, dto StoreDto) (*ent.Store, error) {
	update := r.db.Store.UpdateOneID(dto.ID).
		Where(store.TenantID(dto.TenantID)).
		SetName(dto.Name).
		SetAddress(dto.Address).
		SetPhone(dto.Phone).
		SetWorkHours(dto.WorkHours).
		SetIsActive(dto.IsActive)

	if dto.Lat != nil {
		update.SetLat(*dto.Lat)
	} else {
		update.ClearLat()
	}
	if dto.Lon != nil {
		update.SetLon(*dto.Lon)
	} else {
		update.ClearLon()
	}

	return update.Save(ctx)
}

func (r *storesRepo) DeleteStore(ctx context.Context, id int64, tenantID int64) error {
	_, err := r.db.Store.Delete().
		Where(store.ID(id), store.TenantID(tenantID)).
		Exec(ctx)
	return err
}

func (r *storesRepo) GetStore(ctx context.Context, id int64, tenantID int64) (*ent.Store, error) {
	return r.db.Store.Query().
		Where(store.ID(id), store.TenantID(tenantID)).
		Only(ctx)
}

func (r *storesRepo) ListStores(
	ctx context.Context, filter StoresListFilter, paginate *utils_v1.PaginateRequest,
) ([]*ent.Store, error) {
	query := r.db.Store.Query().Where(store.TenantID(filter.TenantID))

	if filter.OnlyActive {
		query.Where(store.IsActive(true))
	}

	if paginate.GetFromId() != 0 {
		query.Where(store.IDGT(paginate.GetFromId()))
	}

	if paginate.GetLimit() == 0 {
		paginate.Limit = 100
	}

	return query.Limit(int(paginate.GetLimit())).Order(ent.Asc(store.FieldID)).All(ctx)
}

func (r *storesRepo) CountListStores(ctx context.Context, filter StoresListFilter) (int32, error) {
	query := r.db.Store.Query().Where(store.TenantID(filter.TenantID))

	if filter.OnlyActive {
		query.Where(store.IsActive(true))
	}

	count, err := query.Count(ctx)
	return int32(count), err
}

func (r *storesRepo) SetResponsible(ctx context.Context, storeID int64, tenantID int64, memberID int64) (*ent.Store, error) {
	// Note: member validation is trusted from the caller (tenants service owns member data).
	return r.db.Store.UpdateOneID(storeID).
		Where(store.TenantID(tenantID)).
		SetResponsibleID(memberID).
		Save(ctx)
}

func (r *storesRepo) GetStoresByIDs(ctx context.Context, ids []int64) ([]*ent.Store, error) {
	return r.db.Store.Query().
		Where(store.IDIn(ids...)).
		All(ctx)
}
