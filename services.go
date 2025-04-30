/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

type Service struct {
	Name        string
	Controllers []Controller
}

func NewService(name string, cs ...Controller) *Service {
	return &Service{
		Name:        name,
		Controllers: cs,
	}
}
