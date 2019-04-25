package rules

import (
	"k8s.io/gengo/types"
)

type ListTypeMissing struct{}

func (l *ListTypeMissing) Validate(t *types.Type) ([]string, error) {
	fields := make([]string, 0)

	switch t.Kind {
	case types.Struct:
		for _, m := range t.Members {
			goName := m.Name
			if m.Type.Kind != types.Slice {
				continue
			}
			if types.ExtractCommentTags("+", m.CommentLines)["listType"] == nil {
				fields = append(fields, goName)
				continue
			}
		}
	}

	return fields, nil

}

func (l *ListTypeMissing) Name() string {
	return "list_type_missing"
}
