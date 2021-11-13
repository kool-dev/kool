package yamler

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Merger interface {
	Merge(*yaml.Node, *yaml.Node) error
}

type DefaultMerger struct{}

func (m *DefaultMerger) Merge(src *yaml.Node, dst *yaml.Node) (err error) {
	// allow only for merging similar structures for now
	if src.Kind != dst.Kind {
		if dst.IsZero() || (dst.Kind == yaml.ScalarNode && dst.Value == "" && len(dst.Content) == 0) {
			// it's an empty
			*dst = *src
			return
		}

		err = fmt.Errorf("trying to merge different kinds: %d into %d", src.Kind, dst.Kind)
		return
	}

	m.mergeComments(src, dst)

	var (
		l = len(src.Content)
		i int
	)

	if dst.Kind == yaml.DocumentNode {
		// we got two whole files, nice!
		// each content node should be a whole document by itselt
		// (assuming, in case of multiple documents per file, order
		// must match)
		for i = 0; i < l; i++ {
			if err = m.Merge(src.Content[i], dst.Content[i]); err != nil {
				return
			}
		}
		return
	}

	if dst.Kind == yaml.MappingNode {
		// a mapping, aka an object!
		// we wanna do a first simple merge, matching
		// existing keys (and merging them) them.

		var mapInfo = make(map[string]int)

		for i = 0; i < l; i += 2 {
			mapInfo[src.Content[i].Value] = i
		}

		for tag, index := range mapInfo {
			var exists = false
			for i, nd := range dst.Content {
				if tag == nd.Value {
					m.mergeComments(src.Content[index], nd)

					if err = m.Merge(src.Content[index+1], dst.Content[i+1]); err != nil {
						return
					}

					exists = true
					break
				}
			}

			if !exists {
				dst.Content = append(dst.Content, src.Content[index], src.Content[index+1])
			}
		}

		return
	}

	if dst.Kind == yaml.ScalarNode {
		dst.Value = src.Value
		return
	}

	if dst.Kind == yaml.SequenceNode {
		dst.Content = append(dst.Content, src.Content...)
		return
	}

	return
}

func (m *DefaultMerger) mergeComments(src *yaml.Node, dst *yaml.Node) {
	if dst.HeadComment == "" && src.HeadComment != "" {
		dst.HeadComment = src.HeadComment
	}

	if dst.LineComment == "" && src.LineComment != "" {
		dst.LineComment = src.LineComment
	}
}
