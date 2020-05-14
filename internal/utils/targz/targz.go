/*
 * Copyright 2020. The Alkaid Authors. All rights reserved.
 * Use of this source code is governed by a MIT-style
 * license that can be found in the LICENSE file.
 *
 * Alkaid is a BaaS service based on Hyperledger Fabric.
 *
 */

package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type file struct {
	Name string
	Mode int64
	Body []byte
}

type Targz struct {
	Files []*file
}

func New() *Targz {
	return &Targz{
		Files: make([]*file, 0),
	}
}

func (t *Targz) AddFile(name string, mode int64, body []byte) *Targz {
	t.Files = append(t.Files, &file{
		Name: name,
		Mode: mode,
		Body: body,
	})

	return t
}

func (t *Targz) Generate() ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	gw := gzip.NewWriter(buf)
	tw := tar.NewWriter(gw)

	for _, file := range t.Files {
		if err := fileCopy(tw, gw, file); err != nil {
			if cerr := closeStream(tw, gw); cerr != nil {
				return nil, errors.Wrap(err, fmt.Sprintf("file copy failed and close stream error: %s", cerr))
			}
			return nil, errors.Wrap(err, "file copy failed")
		}
	}

	err := closeStream(tw, gw)
	if err != nil {
		return nil, errors.Wrap(err, "close stream failed")
	}

	return buf.Bytes(), nil
}

func fileCopy(tw *tar.Writer, gw *gzip.Writer, file *file) error {
	header := new(tar.Header)

	header.Name = file.Name
	header.Mode = file.Mode
	header.Size = int64(len(file.Body))

	if err := tw.WriteHeader(header); err != nil {
		return err
	}

	if _, err := tw.Write(file.Body); err != nil {
		return err
	}

	if err := tw.Flush(); err != nil {
		return err
	}

	if err := gw.Flush(); err != nil {
		return err
	}

	return nil
}

func closeStream(tw, gw io.Closer) error {
	err := tw.Close()
	if err != nil {
		return err
	}
	err = gw.Close()
	return err
}
