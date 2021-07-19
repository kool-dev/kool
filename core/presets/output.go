package presets

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type OutputWritter interface {
	WriteYAML(string, *yaml.Node) error
}

type DefaultOutputWritter struct {}

func (o *DefaultOutputWritter) WriteYAML(filePath string, document *yaml.Node) (err error) {
	var (
		buff = new(bytes.Buffer)
		encoder *yaml.Encoder
		file *os.File
	)

	if document.Kind != yaml.DocumentNode {
		err = fmt.Errorf("unexpected yaml.Node; expected document (1), but got %d", document.Kind)
		return
	}

	encoder = yaml.NewEncoder(buff)
	encoder.SetIndent(2)

	if err = encoder.Encode(document); err != nil {
		return
	}

	if err = encoder.Close(); err != nil {
		return
	}

	if file, err = os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, os.ModePerm); err != nil {
		return
	}

	if _, err = io.Copy(file, buff); err != nil {
		return
	}

	buff.Reset()
	buff = nil

	if err = file.Sync(); err != nil {
		return
	}

	err = file.Close()
	return
}
