/*
 * Copyright (c) 2025 yakumioto <yaku.mioto@gmail.com>
 * All rights reserved.
 */

package alkaid

type Middleware interface {
	Name() string
	HandlerFunc() HandlerFunc
}
