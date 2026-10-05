package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

func recordingEffectiveHDScale(scale, hd int, explicit, mode bool) (int, error) {
	if !mode {
		return hd, nil
	}
	if explicit && hd != scale {
		return 0, errors.New("configuration_mismatch")
	}
	return scale, nil
}

func recordingControlFont(flow *flowData, scale int) (*helpData, error) {
	if flow == nil || 320*scale < 160 {
		return nil, errors.New("recording_control_font")
	}
	h := &helpData{Glyphs: map[string]helpGlyph{}}
	for _, r := range "停止錄製收尾中完成失敗" {
		rows, ok := flow.Glyphs[string(r)]
		if !ok || len(rows) != flow.size() {
			return nil, errors.New("recording_control_font")
		}
		glyph := helpGlyph{W: flow.size(), Rows: make([]int, len(rows))}
		for i, row := range rows {
			glyph.Rows[i] = int(row)
		}
		h.Glyphs[string(r)] = glyph
	}
	if 5*(flow.size()+2) > 140 || flow.size() > 24 {
		return nil, errors.New("recording_control_font")
	}
	return h, nil
}
func (l *inputRecording) mouseInFooter(scale int) bool {
	if l == nil || l.replay {
		return false
	}
	l.Poll()
	x, y := ebiten.CursorPosition()
	if y < 200*scale {
		return false
	}
	if x >= 8 && x < 160 && y >= 200*scale+4 && y < 200*scale+36 && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		l.Request()
	}
	return true
}
func (l *inputRecording) drawControls(dst *ebiten.Image, scale int) {
	if l == nil || l.replay {
		return
	}
	label := "停止錄製"
	if l.Requested() && l.state.Load() == 0 {
		label = "錄製收尾中"
	}
	switch l.state.Load() {
	case 1:
		label = "錄製收尾中"
	case 2:
		label = "錄製完成"
	case 3:
		label = "錄製失敗"
	}
	img := image.NewRGBA(image.Rect(0, 0, 320*scale, 40))
	for y := 0; y < 40; y++ {
		for x := 0; x < 320*scale; x++ {
			img.SetRGBA(x, y, color.RGBA{24, 28, 34, 255})
		}
	}
	for y := 4; y < 36; y++ {
		for x := 8; x < 160; x++ {
			img.SetRGBA(x, y, color.RGBA{55, 67, 82, 255})
		}
	}
	l.font.drawText(img, 20, 12, label, color.RGBA{255, 225, 155, 255})
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(0, float64(200*scale))
	dst.DrawImage(ebiten.NewImageFromImage(img), op)
}

// 錄放模式阻止既有除錯訊息輸出原文、逐鍵碼、路徑或任意錯誤。
// 缺譯仍記固定診斷，不把玩家輸入及防拷內容印到 log。
type recordingPrivacyLog struct{}

func (recordingPrivacyLog) Write(p []byte) (int, error) {
	s := strings.TrimSpace(string(p))
	safe := ""
	switch s {
	case "錄製完成，繼續遊戲", "重播完成", "重播取消，尚未完成":
		safe = s
	}
	for _, prefix := range []string{"錄放啟動失敗：", "錄製失敗，已保留未完成片段：", "錄製失敗，檔案收尾未完成："} {
		if strings.HasPrefix(s, prefix) {
			code := strings.TrimPrefix(s, prefix)
			for _, allowed := range []string{"wire_rejected", "configuration_mismatch", "initialization_failed", "recording_output_failed", "input_step_mismatch", "api_result_mismatch", "recording_ineligible", "boundary_invalid", "draw_timeout", "draw_mismatch", "writer_append", "writer_finish", "stop_timeout", "cleanup_incomplete", "owner_timeout", "close_timeout"} {
				if code == allowed {
					safe = s
				}
			}
		}
	}
	if strings.HasPrefix(s, "無譯文 ") || strings.HasPrefix(s, "填槽缺字 ") || strings.Contains(s, "缺譯：") {
		safe = "錄放：缺譯"
	}
	if safe != "" {
		_, err := fmt.Fprintln(os.Stderr, safe)
		if err != nil {
			return 0, err
		}
	}
	return len(p), nil
}
