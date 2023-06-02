package utils

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Placeholder struct {
	idx  int
	Args []interface{}
}

func (pl *Placeholder) Build(length int, cap int) {
	*pl = Placeholder{
		idx:  0,
		Args: make([]interface{}, length, cap),
	}
}

func (pl *Placeholder) Get(arg interface{}) string {
	pl.Args = append(pl.Args, arg)
	pl.idx += 1
	return fmt.Sprintf("$%d", pl.idx)
}

type SqlBuilder interface {
	Build(pl *Placeholder) string
}
type builderualBuilder struct {
	Key   string
	Value any
}

func (builder builderualBuilder) Build(pl *Placeholder) string {
	return fmt.Sprintf("%s = %s", builder.Key, pl.Get(builder.Value))
}

type LikeBuilder struct {
	Key   string
	Value any
}

func (builder LikeBuilder) Build(pl *Placeholder) string {
	return fmt.Sprintf("%s LIKE %s || '%%'", builder.Key, pl.Get(builder.Value))
}

type ILikeBuilder struct {
	Key   string
	Value any
}

func (builder ILikeBuilder) Build(pl *Placeholder) string {
	return fmt.Sprintf("%s ILIKE %s || '%%'", builder.Key, pl.Get(builder.Value))
}

type InBuilder struct {
	Key   string
	Value []any
}

func (builder InBuilder) Build(pl *Placeholder) string {
	array := make([]string, len(builder.Value))
	for i, value := range builder.Value {
		array[i] = pl.Get(value)
	}
	return fmt.Sprintf("%s IN (%s)", builder.Key, strings.Join(array, ","))
}

type AndBuilder struct {
	Value []SqlBuilder
}

func (builder AndBuilder) Build(pl *Placeholder) string {
	array := make([]string, len(builder.Value))
	for i, value := range builder.Value {
		array[i] = value.Build(pl)
	}
	return fmt.Sprintf("(%s)", strings.Join(array, " AND "))
}

type OrBuilder struct {
	Value []SqlBuilder
}

func (builder OrBuilder) Build(pl *Placeholder) string {
	array := make([]string, len(builder.Value))
	for i, value := range builder.Value {
		array[i] = value.Build(pl)
	}
	return fmt.Sprintf("(%s)", strings.Join(array, " OR "))
}

func BuildSQLResponse(row pgx.CollectableRow, response_struct any) error {

	elements := reflect.ValueOf(response_struct).Elem()
	order := make(map[string]int)
	for i, field_descriptor := range row.FieldDescriptions() {
		order[field_descriptor.Name] = i
	}
	pointers := make([]interface{}, len(row.FieldDescriptions()))
	for i := 0; i < elements.NumField(); i++ {
		tag := elements.Type().Field(i).Tag
		key, ok := tag.Lookup("sql")
		if !ok {
			continue
		}
		if index, ok := order[key]; ok {
			pointers[index] = elements.Field(i).Addr().Interface()
		}
	}

	err := row.Scan(pointers...)
	return err
}
