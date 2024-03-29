/*
 * Copyright (c) 2021. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package storage

import (
	"errors"
	"sync"

	"github.com/gin-gonic/gin"
)

var (
	once   sync.Once
	global Storage

	ErrNotinitializedGlobalStorage = errors.New("the global storage instance is not initialized")
	ErrNeedUpdateOptions           = errors.New("must need update options")

	ErrNotFound = errors.New("not found")
)

func Initialize(storage Storage) {
	once.Do(func() {
		if global == nil {
			global = storage
		}
	})
}

type Storage interface {
	AutoMigrate(dst ...interface{}) error
	Create(value interface{}) error
	Update(values interface{}, options *UpdateOptions) error
	FindByID(dest interface{}, conditions ...interface{}) error
	FindByQuery(dest interface{}, options *QueryOptions) error
	Delete(value interface{}, conditions ...interface{}) error
	Begin() Storage
	Commit() error
}

func AutoMigrate(dst ...interface{}) error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.AutoMigrate(dst...)
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

func Begin() Storage {
	if err := checkGlobal(); err != nil {
		return nil
	}

	return global.Begin()
}

func Commit() error {
	if err := checkGlobal(); err != nil {
		return err
	}

	return global.Commit()
}

func checkGlobal() error {
	if global == nil {
		return ErrNotinitializedGlobalStorage
	}

	return nil
}

type condition struct {
	Query interface{}
	Args  []interface{}
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
	ors    []*condition
	order  interface{}
	limit  int
	offset int
}

// NewQueryOptionsWithCtx 通过 gin context 获取 options
// 示例：example.com/users?type=name&q=mioto&createdAt=1649902088&page=20&limit=20
func NewQueryOptionsWithCtx(ctx *gin.Context) *QueryOptions {
	// todo: 目前只支持一个字段的模糊查询以及翻页查询
	// typ := ctx.Query("type")
	// q := ctx.Query("query")
	// ctx.BindQuery()

	return NewQueryOptions()
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		limit:  -1,
		offset: -1,
	}
}

func (q *QueryOptions) GetWhere() *condition {
	return q.where
}

func (q *QueryOptions) Where(query interface{}, args ...interface{}) *QueryOptions {
	q.where = &condition{Query: query, Args: args}
	return q
}

func (q *QueryOptions) GetOrs() []*condition {
	return q.ors
}

func (q *QueryOptions) Or(query interface{}, args ...interface{}) *QueryOptions {
	if q.ors == nil {
		q.ors = make([]*condition, 0)
	}

	q.ors = append(q.ors, &condition{Query: query, Args: args})

	return q
}

func (q *QueryOptions) GetNot() *condition {
	return q.not
}

func (q *QueryOptions) Not(query interface{}, args ...interface{}) *QueryOptions {
	q.not = &condition{Query: query, Args: args}
	return q
}

func (q *QueryOptions) GetOrder() interface{} {
	return q.order
}

func (q *QueryOptions) Order(order interface{}) *QueryOptions {
	q.order = order
	return q
}

func (q *QueryOptions) GetLimit() int {
	return q.limit
}

func (q *QueryOptions) Limit(limit int) *QueryOptions {
	q.limit = limit
	return q
}

func (q *QueryOptions) GetOffset() int {
	return q.offset
}

func (q *QueryOptions) Offset(offset int) *QueryOptions {
	q.offset = offset
	return q
}
