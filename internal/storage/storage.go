package storage

import (
	"errors"
	"sync"
)

var (
	ErrNotinitializedGlobalStorage = errors.New("the global storage instance is not initialized")
	ErrNeedUpdateOptions           = errors.New("must need update options")

	once   sync.Once
	global Storage
)

func Initialization(storage Storage) {
	once.Do(func() {
		if global == nil {
			global = storage
		}
	})
}

type Storage interface {
	Create(value interface{}) error
	Update(values interface{}, options *UpdateOptions) error
	FindByID(dest interface{}, conditions ...interface{}) error
	FindByQuery(dest interface{}, options *QueryOptions) error
	Delete(value interface{}, conditions ...interface{}) error
}

func Create(value interface{}) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.Create(value)
}

func Update(values interface{}, options *UpdateOptions) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.Update(values, options)
}

func FindByID(dest interface{}, conditions ...interface{}) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.FindByID(dest, conditions)
}

func FindByQuery(dest interface{}, options *QueryOptions) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.FindByQuery(dest, options)
}

func Delete(value interface{}, conditions ...interface{}) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.Delete(value, conditions)
}

func checkGlobal() error {
	if global == nil {
		return ErrNotinitializedGlobalStorage
	}

	return nil
}

type UpdateOptions struct {
	*condition
}

func NewUpdateOptions(query interface{}, args ...interface{}) *UpdateOptions {
	return &UpdateOptions{
		condition: &condition{Query: query, Args: args},
	}
}

type QueryOptions struct {
	where  *condition
	not    *condition
	or     *condition
	order  interface{}
	limit  int
	offset int
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		limit:  -1,
		offset: -1,
	}
}

func (q *QueryOptions) Where() *condition {
	return q.where
}

func (q *QueryOptions) SetWhere(query interface{}, args ...interface{}) *QueryOptions {
	q.where = &condition{Query: query, Args: args}
	return q
}

func (q *QueryOptions) Or() *condition {
	return q.or
}

func (q *QueryOptions) SetOr(query interface{}, args ...interface{}) *QueryOptions {
	q.or = &condition{Query: query, Args: args}
	return q
}

func (q *QueryOptions) Not() *condition {
	return q.not
}

func (q *QueryOptions) SetNot(query interface{}, args ...interface{}) *QueryOptions {
	q.not = &condition{Query: query, Args: args}
	return q
}

func (q *QueryOptions) Order() interface{} {
	return q.order
}

func (q *QueryOptions) SetOrder(order interface{}) *QueryOptions {
	q.order = order
	return q
}

func (q *QueryOptions) Limit() int {
	return q.limit
}

func (q *QueryOptions) SetLimit(limit int) *QueryOptions {
	q.limit = limit
	return q
}

func (q *QueryOptions) Offset() int {
	return q.offset
}

func (q *QueryOptions) SetOffset(offset int) *QueryOptions {
	q.offset = offset
	return q
}

type condition struct {
	Query interface{}
	Args  []interface{}
}
