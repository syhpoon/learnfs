// Copyright (c) 2022 Max Kuznetsov <syhpoon@syhpoon.ca>

package fs

import "errors"

var (
	ErrorNotFound      = errors.New("not found")
	ErrorAlreadyExists = errors.New("already exists")
)
