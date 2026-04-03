package repositories

import (
	"context"

	"gorm.io/gorm"
)

type repository[T any] interface {
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, id uint) error
	DeleteByUUID(ctx context.Context, uuid string) error
	FindByID(ctx context.Context, id uint) (*T, error)
	FindByUUID(ctx context.Context, uuid string) (*T, error)
}

type baseRepository[T any] struct {
	db *gorm.DB
}

func newBaseRepository[T any](db *gorm.DB) baseRepository[T] {
	return baseRepository[T]{db: db}
}

func (r *baseRepository[T]) Create(ctx context.Context, entity *T) error {
	return gorm.G[T](r.db).Create(ctx, entity)
}

func (r *baseRepository[T]) Update(ctx context.Context, entity *T) error {
	_, err := gorm.G[T](r.db).Updates(ctx, *entity)
	return err
}

func (r *baseRepository[T]) Delete(ctx context.Context, id uint) error {
	_, err := gorm.G[T](r.db).Where("id = ?", id).Delete(ctx)
	return err
}

func (r *baseRepository[T]) DeleteByUUID(ctx context.Context, uuid string) error {
	_, err := gorm.G[T](r.db).Where("uuid = ?", uuid).Delete(ctx)
	return err
}

func (r *baseRepository[T]) FindByID(ctx context.Context, id uint) (*T, error) {
	entity, err := gorm.G[T](r.db).Where("id = ?", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

func (r *baseRepository[T]) FindByUUID(ctx context.Context, uuid string) (*T, error) {
	entity, err := gorm.G[T](r.db).Where("uuid = ?", uuid).First(ctx)
	if err != nil {
		return nil, err
	}
	return &entity, nil
}
