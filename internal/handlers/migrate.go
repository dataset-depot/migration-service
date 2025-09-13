package handlers

import (
	"archive/zip"
	"context"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func (h *Handlers) migrateUpload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 32<<20)

	if err := r.ParseMultipartForm(64 << 20); err != nil {
		http.Error(w, "bad form: "+err.Error(), 400); return
	}
	file, hdr, err := r.FormFile("bundle")
	if err != nil { http.Error(w, "missing file field 'bundle'", 400); return }
	defer file.Close()

	tmpZip, err := os.CreateTemp("", "migs-*.zip")
	if err != nil { http.Error(w, err.Error(), 500); return }
	defer os.Remove(tmpZip.Name())

	if _, err := io.Copy(tmpZip, file); err != nil { http.Error(w, err.Error(), 500); return }
	info, _ := tmpZip.Stat()
	rdr, err := zip.NewReader(tmpZip, info.Size())
	if err != nil { http.Error(w, "not a zip: "+err.Error(), 400); return }

	dst, err := os.MkdirTemp("", "migs-unpack")
	if err != nil { http.Error(w, err.Error(), 500); return }
	defer os.RemoveAll(dst)

	if err := unzipTo(dst, rdr); err != nil { http.Error(w, err.Error(), 500); return }

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Minute); defer cancel()

	if err := h.migrator.UpFromDir(ctx, filepath.Join(dst, "db", "migrations")); err != nil {
		http.Error(w, err.Error(), 500); return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok: " + hdr.Filename + "\n"))
}

func unzipTo(dst string, zr *zip.Reader) error {
	for _, f := range zr.File {
		target := filepath.Join(dst, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(target, 0o755); err != nil { return err }
			continue
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil { return err }
		rc, err := f.Open(); if err != nil { return err }
		out, err := os.Create(target); if err != nil { rc.Close(); return err }
		if _, err := io.Copy(out, rc); err != nil { rc.Close(); out.Close(); return err }
		rc.Close(); out.Close()
	}
	return nil
}
