"""030：斜線原字模、ASCII拒絕與整則更新邊界。"""
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

import bake_flow as flow


class FlowSlashTests(unittest.TestCase):
    def setUp(self):
        flow.slash_glyph.cache_clear()

    def tearDown(self):
        flow.slash_glyph.cache_clear()

    def test_unifont_native_rows_keep_all_eight_bits(self):
        rows = [128, 64, 32, 16, 8, 4, 2, 1] * 2
        with patch.multiple(flow, BACKEND="unifont", SIZE=16), patch.object(
                flow, "unicode_glyph", return_value=(8, rows)):
            self.assertEqual(flow.slash_glyph(), tuple(row << 4 for row in rows))

    def test_eten_uses_only_ascii_slash_slot(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "ascii.15"
            raw = bytearray([255] * (256 * 15))
            rows = bytes([128, 64, 32, 16, 8, 4, 2, 1, 3, 6, 12, 24, 48, 96, 192])
            raw[47 * 15:48 * 15] = rows
            path.write_bytes(raw)
            with patch.multiple(flow, BACKEND="eten", SIZE=15, ASC15=path):
                self.assertEqual(flow.slash_glyph(), tuple(row << 4 for row in rows))

    def test_bad_unifont_shapes_fail(self):
        for width, rows in [(16, [1] * 16), (8, [1] * 15), (8, [0] * 16),
                            (8, [256] * 16), (8, [-1] * 16)]:
            with self.subTest(width=width, rows=rows):
                flow.slash_glyph.cache_clear()
                with patch.multiple(flow, BACKEND="unifont", SIZE=16), patch.object(
                        flow, "unicode_glyph", return_value=(width, rows)):
                    with self.assertRaises(ValueError):
                        flow.slash_glyph()

    def test_missing_or_bad_eten_font_fails(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / "ascii.15"
            with patch.multiple(flow, BACKEND="eten", SIZE=15, ASC15=path):
                with self.assertRaises(FileNotFoundError):
                    flow.slash_glyph()
                for raw in [bytes(3839), bytes(3841), bytes(3840)]:
                    flow.slash_glyph.cache_clear()
                    path.write_bytes(raw)
                    with self.assertRaises(ValueError):
                        flow.slash_glyph()

    def test_only_slash_is_new_ascii_exception(self):
        with patch.object(flow, "glyph_char", return_value=[1] * 15), patch.object(
                flow, "slash_glyph", return_value=tuple([16] * 15)):
            glyphs = {"舊": [2] * 15}
            self.assertEqual(flow.usable_zh("\x04否/\x05好/", glyphs), "否/好/")
            self.assertEqual(set(glyphs), {"舊", "否", "/", "好"})
            for cp in range(32, 127):
                if cp == 47:
                    continue
                with self.subTest(cp=cp):
                    glyphs = {"舊": [2] * 15}
                    self.assertEqual(flow.usable_zh("否/" + chr(cp), glyphs), "")
                    self.assertEqual(glyphs, {"舊": [2] * 15})
            glyphs = {}
            self.assertEqual(flow.usable_zh("否/😀", glyphs), "")
            self.assertEqual(glyphs, {})

    def test_slash_failure_is_not_swallowed_by_eten(self):
        glyphs = {"舊": [2] * 15}
        with patch.object(flow, "BACKEND", "eten"), patch.object(
                flow, "glyph_char", return_value=[1] * 15), patch.object(
                flow, "slash_glyph", side_effect=ValueError("bad slash")):
            with self.assertRaises(ValueError):
                flow.usable_zh("否/", glyphs)
        self.assertEqual(glyphs, {"舊": [2] * 15})

    def test_original_cjk_failure_policy_is_preserved(self):
        for backend in ["eten", "unifont"]:
            with self.subTest(backend=backend), patch.object(flow, "BACKEND", backend), patch.object(
                    flow, "glyph_char", side_effect=ValueError("missing")):
                glyphs = {}
                if backend == "eten":
                    self.assertEqual(flow.usable_zh("否", glyphs), "")
                else:
                    with self.assertRaises(ValueError):
                        flow.usable_zh("否", glyphs)
                self.assertEqual(glyphs, {})

    def test_font_failure_precedes_catalog_read_and_output(self):
        with tempfile.TemporaryDirectory() as tmp:
            target = Path(tmp) / "flow.json"
            target.write_bytes(b"existing output")
            with patch.object(flow, "BACKEND", "eten"), patch.object(
                    flow, "slash_glyph", side_effect=ValueError("bad slash")), patch.object(
                    flow, "load_player_catalog") as catalog, patch.object(
                    flow.sys, "argv", ["bake_flow.py", str(target)]):
                with self.assertRaises(ValueError):
                    flow.main()
                catalog.assert_not_called()
            self.assertEqual(target.read_bytes(), b"existing output")


if __name__ == "__main__":
    unittest.main()
