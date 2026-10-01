package main

import (
	"encoding/json"
	"reflect"
	"testing"
)

// 字模只用自製識別值，期望序列不由被測函式產生。
func TestExpandSlotsComplete(t *testing.T) {
	ent := entry{Dy: 7, G: [][]uint16{{101}, {0}, {102}, {0}, {103}}, D: []int{1, 3}}
	digits := map[string][]uint16{"1": {11}, "2": {12}}
	ordinals := map[string][]uint16{"SECOND": {22}}
	before, _ := json.Marshal([]any{ent, digits, ordinals})
	got, ok := expandSlots(ent, []string{"12", "SECOND"}, digits, ordinals)
	if !ok || got.Dy != 7 || !reflect.DeepEqual(got.G, [][]uint16{{101}, {11}, {12}, {102}, {22}, {103}}) || !reflect.DeepEqual(got.D, ent.D) {
		t.Fatalf("完整槽位順序不符：%v %v", got, ok)
	}
	after, _ := json.Marshal([]any{ent, digits, ordinals})
	if string(before) != string(after) {
		t.Fatal("展開改動輸入")
	}
}

func TestExpandSlotsMissingGlyphReturnsOriginal(t *testing.T) {
	tests := []struct {
		name     string
		token    string
		digits   map[string][]uint16
		ordinals map[string][]uint16
	}{
		{"首位數字缺項", "12", map[string][]uint16{"2": {12}}, nil},
		{"後位數字缺項", "12", map[string][]uint16{"1": {11}}, nil},
		{"數字nil", "12", map[string][]uint16{"1": {11}, "2": nil}, nil},
		{"數字空字模", "12", map[string][]uint16{"1": {11}, "2": {}}, nil},
		{"序數缺項", "SECOND", nil, nil},
		{"序數nil", "SECOND", nil, map[string][]uint16{"SECOND": nil}},
		{"序數空字模", "SECOND", nil, map[string][]uint16{"SECOND": {}}},
		{"非序數映射nil不得改走數字", "12", map[string][]uint16{"1": {11}, "2": {12}}, map[string][]uint16{"12": nil}},
		{"非序數映射空不得改走數字", "12", map[string][]uint16{"1": {11}, "2": {12}}, map[string][]uint16{"12": {}}},
		{"序數缺項不得改走字母", "SECOND", map[string][]uint16{"S": {1}, "E": {2}, "C": {3}, "O": {4}, "N": {5}, "D": {6}}, nil},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ent := entry{Dy: 3, G: [][]uint16{{101}, {0}, {102}}, D: []int{1}}
			before, _ := json.Marshal([]any{ent, test.digits, test.ordinals})
			got, ok := expandSlots(ent, []string{test.token}, test.digits, test.ordinals)
			if ok || !reflect.DeepEqual(got, ent) {
				t.Fatalf("缺字竟交付部分結果：%v %v", got, ok)
			}
			after, _ := json.Marshal([]any{ent, test.digits, test.ordinals})
			if string(before) != string(after) {
				t.Fatal("缺字改動輸入")
			}
		})
	}
}

func TestExpandSlotsRetainsExistingBoundaries(t *testing.T) {
	tests := []struct {
		name     string
		slots    []int
		tokens   []string
		ordinals map[string][]uint16
		want     [][]uint16
	}{
		{"無槽", nil, []string{"12"}, nil, [][]uint16{{101}, {0}, {102}}},
		{"槽索引集合", []int{1, 1, -1, 99}, []string{"12", "99"}, nil, [][]uint16{{101}, {11}, {12}, {102}}},
		{"缺token沿既有字模", []int{1}, nil, nil, [][]uint16{{101}, {0}, {102}}},
		{"直接序數表映射優先", []int{1}, []string{"12"}, map[string][]uint16{"12": {44}}, [][]uint16{{101}, {44}, {102}}},
		{"非空全零字模保留", []int{1}, []string{"12"}, map[string][]uint16{"12": {0}}, [][]uint16{{101}, {0}, {102}}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, ok := expandSlots(entry{G: [][]uint16{{101}, {0}, {102}}, D: test.slots}, test.tokens,
				map[string][]uint16{"1": {11}, "2": {12}}, test.ordinals)
			if !ok || !reflect.DeepEqual(got.G, test.want) {
				t.Fatalf("既有邊界結果不符：%v %v", got.G, ok)
			}
		})
	}
}
