package tgz

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"kool-dev/kool/core/shell"
	"os"
	"path/filepath"
	"strings"
)

// TarGz holds configuration for generating a new
type TarGz struct {
	sourceDir string
	file      *os.File

	ignoreFilesMap map[string]bool

	g *gzip.Writer
	t *tar.Writer
}

// NewTemp allocates and opens files for generating a new tarball
// with Gzip compression in a temporary file.
func NewTemp() (tgz *TarGz, err error) {
	tgz = new(TarGz)
	tgz.file, err = os.CreateTemp(os.TempDir(), "*.tgz")
	if err != nil {
		tgz = nil
		return
	}
	tgz.g = gzip.NewWriter(tgz.file)
	tgz.t = tar.NewWriter(tgz.g)
	return
}

// CompressFiles creates the tarball with the given files list
func (tgz *TarGz) CompressFiles(files []string) (tmpfile string, err error) {
	var (
		file string
		fi   os.FileInfo
	)

	for _, file = range files {
		if file == "" {
			continue
		}

		if fi, err = os.Stat(file); os.IsNotExist(err) {
			// if we listed a file that do not exist, that is probably
			// the case where a versioned file was removed but not yet
			// commited. Hence let's only warn about it on Stderr.
			shell.Warning(fmt.Errorf("file not found, not including on tarball: %s", file))
		} else if err = tgz.add(file, fi, err); err != nil {
			shell.Error(fmt.Errorf("failed to add file into archive: %v", err))
		}
	}

	tmpfile, err = tgz.finishCompress()
	return
}

// SourceDir defines the source directory that is going to be compressed
func (tgz *TarGz) SourceDir(dir string) {
	tgz.sourceDir = dir
}

// CompressFolder adds the given folder to the tarball archive
func (tgz *TarGz) CompressFolder(dir string) (tmpfile string, err error) {
	tgz.sourceDir = dir

	if err = filepath.Walk(tgz.sourceDir, tgz.add); err != nil {
		return
	}

	tmpfile, err = tgz.finishCompress()
	return
}

func (tgz *TarGz) finishCompress() (tmpfile string, err error) {
	if err = tgz.t.Close(); err != nil {
		return
	}
	if err = tgz.g.Close(); err != nil {
		return
	}
	if err = tgz.file.Sync(); err != nil {
		return
	}

	tmpfile = tgz.file.Name()
	err = tgz.file.Close()
	return
}

func (tgz *TarGz) add(file string, fi os.FileInfo, err error) error {
	var (
		relPath string
		header  *tar.Header
		fh      *os.File
	)

	if err != nil {
		return err
	}

	if fi.Mode()&os.ModeSymlink == os.ModeSymlink {
		// ignore symlink
		return nil
	}

	relPath = strings.TrimPrefix(file, tgz.sourceDir)

	if relPath == "" || relPath == "/" {
		return nil
	}

	if tgz.shouldIgnore(relPath) {
		return nil
	}

	header, err = tar.FileInfoHeader(fi, file)
	if err != nil {
		return err
	}
	header.Name = strings.TrimPrefix(relPath, "/")
	if err = tgz.t.WriteHeader(header); err != nil {
		return err
	}

	if !fi.IsDir() {
		fh, err = os.Open(file)
		if err != nil {
			return err
		}
		if _, err = io.Copy(tgz.t, fh); err != nil {
			return err
		}
	}

	return nil
}

// SetIgnoreList defines the list of file patterns
// that must be ignored from the tarball created.
func (tgz *TarGz) SetIgnoreList(ignoreList []string) {
	var file string
	tgz.ignoreFilesMap = make(map[string]bool, len(ignoreList))
	for _, file = range ignoreList {
		tgz.ignoreFilesMap[strings.Trim(file, string(os.PathSeparator))] = true
	}
}

func (tgz *TarGz) shouldIgnore(relPath string) (ignores bool) {
	_, ignores = tgz.ignoreFilesMap[strings.Trim(relPath, string(os.PathSeparator))]
	return
}
