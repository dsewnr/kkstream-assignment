package actions

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/pkg/errors"
)

var UPLOAD_DIR = envy.Get("UPLOAD_DIR", "/tmp/uploads")

func init() {
	if ENV == "test" {
		UPLOAD_DIR = filepath.Join(UPLOAD_DIR, "test")
	}
}

func checkDirExists(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// ApiUpload default implementation.
func ApiUpload(c buffalo.Context) error {
	f, err := c.File("video")
	defer f.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	if err := checkDirExists(UPLOAD_DIR); err != nil {
		return errors.WithStack(err)
	}

	ft, err := os.Create(filepath.Join(UPLOAD_DIR, f.Filename))
	defer ft.Close()
	if err != nil {
		return errors.WithStack(err)
	}
	_, err = io.Copy(ft, f)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.Render(200, r.String("OK"))
}

// ApiUploadMulti default implementation.
func ApiUploadMulti(c buffalo.Context) error {
	if err := checkDirExists(UPLOAD_DIR); err != nil {
		return errors.WithStack(err)
	}

	for _, mf := range c.Request().MultipartForm.File["video"] {

		f, err := mf.Open()
		defer f.Close()
		if err != nil {
			return errors.WithStack(err)
		}

		ft, err := os.Create(filepath.Join(UPLOAD_DIR, mf.Filename))
		defer ft.Close()
		if err != nil {
			return errors.WithStack(err)
		}

		_, err = io.Copy(ft, f)

	}
	return c.Render(200, r.String("OK"))
}

// ApiUploadGcs default implementation.
func ApiUploadGcs(c buffalo.Context) error {
	f, err := c.File("video")
	defer f.Close()
	if err != nil {
		return errors.WithStack(err)
	}

	// Copy file into buffer
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return errors.WithStack(err)
	}

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return errors.WithStack(err)
	}

	// Init client
	wc := client.Bucket(envy.Get("GOOGLE_CLOUD_STORAGE_BUCKET", "")).Object(f.Filename).NewWriter(ctx)
	wc.ContentType = "video/mp4"
	wc.ACL = []storage.ACLRule{{Entity: storage.AllUsers, Role: storage.RoleReader}}

	// Upload to Google Cloud Storage
	if _, err := wc.Write(buf.Bytes()); err != nil {
		return errors.WithStack(err)
	}
	if err := wc.Close(); err != nil {
		return errors.WithStack(err)
	}

	// Object attributes
	fmt.Println("updated object:", wc.Attrs())

	return c.Render(200, r.String("OK"))
}
