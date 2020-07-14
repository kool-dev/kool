package tgz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// TarGz holds configuration for generating a new
type TarGz struct {
	sourceDir string
	file      *os.File

	ignorePatterns []*regexp.Regexp
	ignoreFiles    []string
	ignorePrefixes []string

	g *gzip.Writer
	t *tar.Writer
}

// NewTemp allocates and opens files for generating a new tarball
// with Gzip compression in a temporary file.
func NewTemp() (tgz *TarGz, err error) {
	tgz = new(TarGz)
	tgz.file, err = ioutil.TempFile(os.TempDir(), "*.tgz")
	if err != nil {
		tgz = nil
		return
	}
	tgz.g = gzip.NewWriter(tgz.file)
	tgz.t = tar.NewWriter(tgz.g)
	return
}

// Compress adds the given folder to the tarball archive
func (tgz *TarGz) Compress(dir string) (tmpfile string, err error) {
	tgz.sourceDir = dir

	if err = filepath.Walk(tgz.sourceDir, tgz.add); err != nil {
		return
	}

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
func (tgz *TarGz) SetIgnoreList(ignore [][]byte) {
	var (
		i, l int
		p    string
	)
	l = len(ignore)
	for i = 0; i < l; i++ {
		if len(ignore[i]) == 0 {
			continue
		}

		if bytes.Contains(ignore[i], []byte("*")) {
			p = strings.ReplaceAll(regexp.QuoteMeta(strings.ReplaceAll(string(ignore[i]), "*", "__ANYTHING__")), "__ANYTHING__", ".*")
			rx, err := regexp.Compile(p)
			if err != nil {
				fmt.Println("ERROR: could not parse ignore pattern: ", string(ignore[i]), err)
			} else {
				tgz.ignorePatterns = append(tgz.ignorePatterns, rx)
			}
		} else if bytes.HasPrefix(ignore[i], []byte("/")) {
			tgz.ignorePrefixes = append(tgz.ignorePrefixes, string(ignore[i]))
		} else {
			tgz.ignoreFiles = append(tgz.ignoreFiles, "/"+string(ignore[i]))
		}
	}
}

func (tgz *TarGz) shouldIgnore(relPath string) bool {
	var (
		suffix, prefix string
		pattern        *regexp.Regexp
	)

	for _, pattern = range tgz.ignorePatterns {
		if pattern.MatchString(relPath) {
			// skip adding into the tarball
			return true
		}
	}
	for _, suffix = range tgz.ignoreFiles {
		if strings.HasSuffix(relPath, suffix) {
			// skip adding into the tarball
			return true
		}
	}
	for _, prefix = range tgz.ignorePrefixes {
		if strings.HasPrefix(relPath, prefix) {
			// skip adding into the tarball
			return true
		}
	}

	return false
}
