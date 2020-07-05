/*
 * Copyright (c) 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 */

package storage

//go:generate mockgen -destination mock/provider.go -source storage.go -package mock

// DBProvider database provider
type DBProvider interface {
	Create(interface{}) error
	Exist(interface{}) (bool, error)
	Update(interface{}, interface{}) error
	Get(string, interface{}) error
	Query(interface{}, ...interface{}) error
	Delete(interface{}) error
}
