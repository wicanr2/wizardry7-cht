package main

// expandSlots 依spec015展開一般文字的數字／序數槽。
// 缺少所需字模時回原entry及false，呼叫端整則保留英文。
// 不修改輸入；字模row沿用既有唯讀資料，不新增深拷契約。
func expandSlots(ent entry, nums []string, digits, ordinals map[string][]uint16) (entry, bool) {
	if len(ent.D) == 0 {
		return ent, true
	}
	slots := make(map[int]bool, len(ent.D))
	for _, index := range ent.D {
		slots[index] = true
	}
	glyphs := make([][]uint16, 0, len(ent.G)+4)
	next := 0
	for index, glyph := range ent.G {
		if !slots[index] || next >= len(nums) {
			glyphs = append(glyphs, glyph)
			continue
		}
		token := nums[next]
		next++
		if ordinal, exists := ordinals[token]; exists {
			if len(ordinal) == 0 {
				return ent, false
			}
			glyphs = append(glyphs, ordinal)
			continue
		}
		if ordinalsPub[token] {
			return ent, false
		}
		for _, digit := range token {
			bitmap, exists := digits[string(digit)]
			if !exists || len(bitmap) == 0 {
				return ent, false
			}
			glyphs = append(glyphs, bitmap)
		}
	}
	ent.G = glyphs
	return ent, true
}
