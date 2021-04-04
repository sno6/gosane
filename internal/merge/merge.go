package merge

import (
	"encoding/json"
	"errors"
	"reflect"
	"strings"

	"github.com/imdario/mergo"
)

type Map map[string]interface{}

func (m *Map) Init(b []byte) error {
	if m == nil {
		*m = make(Map)
	}

	return json.Unmarshal(b, m)
}

func (m *Map) Get(key string) (interface{}, bool) {
	if m == nil {
		return nil, false
	}

	v, ok := (*m)[key]
	return v, ok
}

// MergeMap will pull values from m where they exist (key) and set in dst.
func MergeMap(dst interface{}, src interface{}, m Map, toCamel bool) error {
	dstV := reflect.ValueOf(dst)
	dstT := reflect.TypeOf(dst)
	if dstV.Kind() != reflect.Ptr || dstV.IsNil() {
		return errors.New("given dst value is nil or non-pointer")
	}

	srcV := reflect.ValueOf(src)
	srcT := reflect.TypeOf(src)
	if srcV.Kind() == reflect.Ptr {
		return errors.New("given src value is nil or a pointer")
	}

	vMap := make(map[string]reflect.Value)
	for i := 0; i < srcT.NumField(); i++ {
		vf := srcV.Field(i)
		tf := srcT.Field(i)

		vMap[tf.Name] = vf
	}

	for i := 0; i < dstT.Elem().NumField(); i++ {
		vf := dstV.Elem().Field(i)
		tf := dstT.Elem().Field(i)

		// Is a key by this name found in our map?
		name := getFieldName(tf, toCamel)
		_, ok := m[name]
		if !ok {
			continue
		}

		// Find the matching key value in src.
		v, ok := vMap[tf.Name]
		if !ok {
			continue
		}

		vf.Set(v)
	}

	return nil
}

// Merge reads from src and merges into dst.
func Merge(dst, src interface{}) error {
	return mergo.MergeWithOverwrite(dst, src)
}

func MergeOverwriteEmptyValue(dst, src interface{}) error {
	return mergo.MergeWithOverwrite(dst, src, mergo.WithOverwriteWithEmptyValue)
}

func getFieldName(tf reflect.StructField, toCamel bool) string {
	name := tf.Name
	if tf.Tag != "" {
		jsonTag, ok := tf.Tag.Lookup("json")
		if ok {
			if cIdx := strings.Index(jsonTag, ","); cIdx > 0 {
				jsonTag = jsonTag[:cIdx]
			}

			if toCamel {
				jsonTag = convertSnakeCase(jsonTag)
			}

			name = jsonTag
		}
	}

	return name
}

func convertSnakeCase(s string) string {
	split := strings.Split(s, "_")

	var sb strings.Builder
	for i, sv := range split {
		if i == 0 {
			sb.WriteString(sv)
		} else {
			sb.WriteString(strings.Title(sv))
		}
	}

	return sb.String()
}
