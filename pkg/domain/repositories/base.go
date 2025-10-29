package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/victor-lima-142/oak-bank/pkg/domain/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Repository é o repositório genérico base
type Repository[T models.Entity] struct {
	db *gorm.DB
}

// NewRepository cria uma nova instância do repositório genérico
func NewRepository[T models.Entity](db *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: db,
	}
}

// Create cria uma nova entidade no banco de dados
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Create(entity).Error
}

// CreateBatch cria múltiplas entidades em batch
func (r *Repository[T]) CreateBatch(ctx context.Context, entities []T, batchSize int) error {
	return r.db.WithContext(ctx).CreateInBatches(entities, batchSize).Error
}

// FindByID busca uma entidade por ID
func (r *Repository[T]) FindByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).First(&entity, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// FindByIDWithPreload busca uma entidade por ID com preload de relações
func (r *Repository[T]) FindByIDWithPreload(ctx context.Context, id string, preloads ...string) (*T, error) {
	var entity T
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(&entity, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// FindAll busca todas as entidades
func (r *Repository[T]) FindAll(ctx context.Context) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Find(&entities).Error
	return entities, err
}

// FindAllWithPreload busca todas as entidades com preload
func (r *Repository[T]) FindAllWithPreload(ctx context.Context, preloads ...string) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx)

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.Find(&entities).Error
	return entities, err
}

// FindByCondition busca entidades com condições customizadas
func (r *Repository[T]) FindByCondition(ctx context.Context, condition string, args ...interface{}) ([]T, error) {
	var entities []T
	err := r.db.WithContext(ctx).Where(condition, args...).Find(&entities).Error
	return entities, err
}

// FindOneByCondition busca uma única entidade com condição
func (r *Repository[T]) FindOneByCondition(ctx context.Context, condition string, args ...interface{}) (*T, error) {
	var entity T
	err := r.db.WithContext(ctx).Where(condition, args...).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Update atualiza uma entidade existente
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	return r.db.WithContext(ctx).Save(entity).Error
}

// UpdateFields atualiza campos específicos de uma entidade
func (r *Repository[T]) UpdateFields(ctx context.Context, id string, fields map[string]interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(fields).Error
}

// Delete deleta uma entidade por ID
func (r *Repository[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

// DeleteByCondition deleta entidades que atendem uma condição
func (r *Repository[T]) DeleteByCondition(ctx context.Context, condition string, args ...interface{}) error {
	var entity T
	return r.db.WithContext(ctx).Where(condition, args...).Delete(&entity).Error
}

// SoftDelete realiza soft delete (se a entidade tiver gorm.DeletedAt)
func (r *Repository[T]) SoftDelete(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

// Restore restaura uma entidade com soft delete
func (r *Repository[T]) Restore(ctx context.Context, id string) error {
	var entity T
	return r.db.WithContext(ctx).Model(&entity).Unscoped().Where("id = ?", id).Update("deleted_at", nil).Error
}

// Count conta o número de registros
func (r *Repository[T]) Count(ctx context.Context) (int64, error) {
	var entity T
	var count int64
	err := r.db.WithContext(ctx).Model(&entity).Count(&count).Error
	return count, err
}

// CountByCondition conta registros com condição
func (r *Repository[T]) CountByCondition(ctx context.Context, condition string, args ...interface{}) (int64, error) {
	var entity T
	var count int64
	err := r.db.WithContext(ctx).Model(&entity).Where(condition, args...).Count(&count).Error
	return count, err
}

// Exists verifica se uma entidade existe
func (r *Repository[T]) Exists(ctx context.Context, id string) (bool, error) {
	var entity T
	err := r.db.WithContext(ctx).Select("id").First(&entity, "id = ?", id).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// ExistsByCondition verifica se existe alguma entidade com a condição
func (r *Repository[T]) ExistsByCondition(ctx context.Context, condition string, args ...interface{}) (bool, error) {
	var entity T
	err := r.db.WithContext(ctx).Select("id").Where(condition, args...).First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	return true, nil
}

// PaginationOptions define opções de paginação
type PaginationOptions struct {
	Page     int
	PageSize int
	OrderBy  string
	Order    string // "asc" ou "desc"
}

// PaginatedResult é o resultado paginado
type PaginatedResult[T models.Entity] struct {
	Data       []T   `json:"data"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalPages int   `json:"total_pages"`
}

// FindWithPagination busca entidades com paginação
func (r *Repository[T]) FindWithPagination(ctx context.Context, opts PaginationOptions) (*PaginatedResult[T], error) {
	var entities []T
	var total int64

	// Validar opções
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}
	if opts.Order == "" {
		opts.Order = "desc"
	}

	offset := (opts.Page - 1) * opts.PageSize

	// Contar total
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Count(&total).Error; err != nil {
		return nil, err
	}

	// Buscar dados paginados
	query := r.db.WithContext(ctx).Offset(offset).Limit(opts.PageSize)

	if opts.OrderBy != "" {
		orderClause := fmt.Sprintf("%s %s", opts.OrderBy, opts.Order)
		query = query.Order(orderClause)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / opts.PageSize
	if int(total)%opts.PageSize > 0 {
		totalPages++
	}

	return &PaginatedResult[T]{
		Data:       entities,
		Total:      total,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
}

// FindWithPaginationAndCondition busca entidades com paginação e condição
func (r *Repository[T]) FindWithPaginationAndCondition(
	ctx context.Context,
	opts PaginationOptions,
	condition string,
	args ...interface{},
) (*PaginatedResult[T], error) {
	var entities []T
	var total int64

	// Validar opções
	if opts.Page <= 0 {
		opts.Page = 1
	}
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}
	if opts.Order == "" {
		opts.Order = "desc"
	}

	offset := (opts.Page - 1) * opts.PageSize

	// Contar total com condição
	var entity T
	if err := r.db.WithContext(ctx).Model(&entity).Where(condition, args...).Count(&total).Error; err != nil {
		return nil, err
	}

	// Buscar dados paginados com condição
	query := r.db.WithContext(ctx).Where(condition, args...).Offset(offset).Limit(opts.PageSize)

	if opts.OrderBy != "" {
		orderClause := fmt.Sprintf("%s %s", opts.OrderBy, opts.Order)
		query = query.Order(orderClause)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	totalPages := int(total) / opts.PageSize
	if int(total)%opts.PageSize > 0 {
		totalPages++
	}

	return &PaginatedResult[T]{
		Data:       entities,
		Total:      total,
		Page:       opts.Page,
		PageSize:   opts.PageSize,
		TotalPages: totalPages,
	}, nil
}

// Transaction executa operações em uma transação
func (r *Repository[T]) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return r.db.WithContext(ctx).Transaction(fn)
}

// GetDB retorna a instância do GORM DB
func (r *Repository[T]) GetDB() *gorm.DB {
	return r.db
}

// WithTx retorna um novo repository com uma transação
func (r *Repository[T]) WithTx(tx *gorm.DB) *Repository[T] {
	return &Repository[T]{
		db: tx,
	}
}

// Upsert cria ou atualiza uma entidade (insert ou update baseado em conflito)
func (r *Repository[T]) Upsert(ctx context.Context, entity *T, conflictColumns []string, updateColumns []string) error {
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   convertToColumns(conflictColumns),
			DoUpdates: clause.AssignmentColumns(updateColumns),
		}).
		Create(entity).Error
}

// BulkUpdate atualiza múltiplas entidades
func (r *Repository[T]) BulkUpdate(ctx context.Context, entities []T) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, entity := range entities {
			if err := tx.Save(&entity).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// FirstOrCreate busca ou cria uma entidade
func (r *Repository[T]) FirstOrCreate(ctx context.Context, entity *T, condition string, args ...interface{}) error {
	return r.db.WithContext(ctx).Where(condition, args...).FirstOrCreate(entity).Error
}

// Helper function
func convertToColumns(columnNames []string) []clause.Column {
	columns := make([]clause.Column, len(columnNames))
	for i, name := range columnNames {
		columns[i] = clause.Column{Name: name}
	}
	return columns
}

// ===========================
// Query Builder (Opcional - para queries mais complexas)
// ===========================

type QueryBuilder[T models.Entity] struct {
	db         *gorm.DB
	conditions []string
	args       []interface{}
	preloads   []string
	orderBy    string
	limit      int
	offset     int
}

// NewQueryBuilder cria um novo query builder
func (r *Repository[T]) NewQueryBuilder() *QueryBuilder[T] {
	return &QueryBuilder[T]{
		db:         r.db,
		conditions: []string{},
		args:       []interface{}{},
		preloads:   []string{},
	}
}

// Where adiciona uma condição WHERE
func (qb *QueryBuilder[T]) Where(condition string, args ...interface{}) *QueryBuilder[T] {
	qb.conditions = append(qb.conditions, condition)
	qb.args = append(qb.args, args...)
	return qb
}

// Preload adiciona um preload
func (qb *QueryBuilder[T]) Preload(relation string) *QueryBuilder[T] {
	qb.preloads = append(qb.preloads, relation)
	return qb
}

// OrderBy define a ordenação
func (qb *QueryBuilder[T]) OrderBy(order string) *QueryBuilder[T] {
	qb.orderBy = order
	return qb
}

// Limit define o limite
func (qb *QueryBuilder[T]) Limit(limit int) *QueryBuilder[T] {
	qb.limit = limit
	return qb
}

// Offset define o offset
func (qb *QueryBuilder[T]) Offset(offset int) *QueryBuilder[T] {
	qb.offset = offset
	return qb
}

// Find executa a query e retorna os resultados
func (qb *QueryBuilder[T]) Find(ctx context.Context) ([]T, error) {
	var entities []T
	query := qb.db.WithContext(ctx)

	// Aplicar condições
	for i, condition := range qb.conditions {
		query = query.Where(condition, qb.args[i])
	}

	// Aplicar preloads
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// Aplicar ordenação
	if qb.orderBy != "" {
		query = query.Order(qb.orderBy)
	}

	// Aplicar limit e offset
	if qb.limit > 0 {
		query = query.Limit(qb.limit)
	}
	if qb.offset > 0 {
		query = query.Offset(qb.offset)
	}

	err := query.Find(&entities).Error
	return entities, err
}

// First executa a query e retorna o primeiro resultado
func (qb *QueryBuilder[T]) First(ctx context.Context) (*T, error) {
	var entity T
	query := qb.db.WithContext(ctx)

	// Aplicar condições
	for i, condition := range qb.conditions {
		query = query.Where(condition, qb.args[i])
	}

	// Aplicar preloads
	for _, preload := range qb.preloads {
		query = query.Preload(preload)
	}

	// Aplicar ordenação
	if qb.orderBy != "" {
		query = query.Order(qb.orderBy)
	}

	err := query.First(&entity).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &entity, nil
}

// Count conta os resultados
func (qb *QueryBuilder[T]) Count(ctx context.Context) (int64, error) {
	var entity T
	var count int64
	query := qb.db.WithContext(ctx).Model(&entity)

	// Aplicar condições
	for i, condition := range qb.conditions {
		query = query.Where(condition, qb.args[i])
	}

	err := query.Count(&count).Error
	return count, err
}
