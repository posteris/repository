package repository

import "gorm.io/gorm"

type Model[E any] interface {
	ToEntity() E
	FromEntity(entity E) interface{}
}

type Repository[M Model[E], E any] struct {
	db *gorm.DB
}

func New[M Model[E], E any](db *gorm.DB) *Repository[M, E] {
	return &Repository[M, E]{
		db: db,
	}
}

func (r *Repository[M, E]) Create(entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.Create(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()

	return nil
}

func (r *Repository[M, E]) Updates(entity *E) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.Updates(&model).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()

	return nil
}

func (r *Repository[M, E]) Delete(entity *E, conds ...interface{}) error {
	var start M
	model := start.FromEntity(*entity).(M)

	err := r.db.Delete(&model, conds...).Error
	if err != nil {
		return err
	}

	*entity = model.ToEntity()

	return nil
}
