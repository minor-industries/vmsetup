package vmsetup

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/go-resty/resty/v2"
)

func download(ctx context.Context, url, path string) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	c := resty.New()

	resp, err := c.R().
		SetContext(ctx).
		SetDoNotParseResponse(true). // key: stream, don't buffer
		Get(url)
	if err != nil {
		return err
	}
	defer resp.RawBody().Close()

	if resp.IsError() {
		return fmt.Errorf("http %d", resp.StatusCode())
	}

	_, err = io.Copy(out, resp.RawBody())
	return err
}

func HashFileSHA256(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}
