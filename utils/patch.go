package utils

import (
	"github.com/xiiot/xlink-plugin-sdk-go/comctx/models"
	"reflect"
)

type Patch[T any] struct {
	Added   []T
	Removed []T
	Updated []T
}
type FullPatch struct {
	NodeChanges    Patch[*models.Node]    `json:"node_changes"`
	GroupChanges   Patch[*models.Group]   `json:"group_changes"`
	SettingChanges Patch[*models.Setting] `json:"setting_changes"`
	TagChanges     Patch[*models.Tag]     `json:"tag_changes"`
}

// DiffAll 比较同一类资源，处理 []*T
func DiffAll[T any](oldList, newList []T) Patch[T] {
	patch := Patch[T]{}

	oldMap := make(map[string]T)
	newMap := make(map[string]T)

	for _, v := range oldList {
		oldMap[getID(v)] = v
	}
	for _, v := range newList {
		newMap[getID(v)] = v
	}

	// 找新增和更新
	for id, newItem := range newMap {
		if oldItem, exists := oldMap[id]; !exists {
			patch.Added = append(patch.Added, newItem)
		} else if !deepEqualPtr(oldItem, newItem) {
			patch.Updated = append(patch.Updated, newItem)
		}
	}

	// 找删除
	for id, oldItem := range oldMap {
		if _, exists := newMap[id]; !exists {
			patch.Removed = append(patch.Removed, oldItem)
		}
	}

	return patch
}

// ==== 提取唯一标识 ====
func getID[T any](v T) string {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typeName := val.Type().Name()

	switch typeName {
	case "Node":
		return val.FieldByName("Name").String()
	case "Group":
		d := val.FieldByName("DriverName").String()
		n := val.FieldByName("Name").String()
		return d + "::" + n
	case "Tag":
		d := val.FieldByName("DriverName").String()
		g := val.FieldByName("GroupName").String()
		n := val.FieldByName("Name").String()
		return d + "::" + g + "::" + n
	case "Setting":
		return val.FieldByName("NodeName").String()
	case "Subscription":
		a := val.FieldByName("AppName").String()
		d := val.FieldByName("DriverName").String()
		g := val.FieldByName("GroupName").String()
		return a + "::" + d + "::" + g
	default:
		if f := val.FieldByName("Name"); f.IsValid() {
			return f.String()
		}
	}
	return ""
}

// ==== 指针对比真实内容 ====
func deepEqualPtr[T any](a, b T) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	if va.Kind() == reflect.Ptr {
		va = va.Elem()
	}
	if vb.Kind() == reflect.Ptr {
		vb = vb.Elem()
	}
	return reflect.DeepEqual(va.Interface(), vb.Interface())
}
func HasChanges[T any](p Patch[T]) bool {
	return len(p.Added) > 0 || len(p.Updated) > 0 || len(p.Removed) > 0
}
