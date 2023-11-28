package v2

import (
	"github.com/aishuchen/goctl-swagger/render/types"
	"github.com/zeromicro/go-zero/core/stringx"
	"github.com/zeromicro/go-zero/tools/goctl/api/spec"
	"strings"
)

func lookupGozeroTag(tags []*spec.Tag) *spec.Tag {
	for _, tag := range tags {
		switch tag.Key {
		case types.PathTagKey, types.FormTagKey, types.HeaderTagKey, types.JsonTagKey:
			return tag
		}
	}
	return nil
}

func isOptionalTag(tag *spec.Tag) bool {
	if len(tag.Options) == 0 {
		return false
	}
	is := stringx.Contains(tag.Options, "optional") || stringx.Contains(tag.Options, "omitempty")
	if is {
		return is
	}
	_, is = defaultTag(tag)
	return is
}

func defaultTag(tag *spec.Tag) (string, bool) {
	if len(tag.Options) == 0 {
		return "", false
	}
	for _, each := range tag.Options {
		s := strings.Split(each, "default=")
		if len(s) != 2 {
			return "", false
		}
		return s[1], true
	}
	return "", false
}
