import json
import unittest

from message_catalog import decode_player_catalog


class MessageCatalogTests(unittest.TestCase):
    def test_reserved_payload_is_never_utf8_or_json_decoded(self):
        raw = (b'{"records":[{"id":2499,"blk":1,"off":7,"en":"ordinary"},'
               b'{"id":2500,"blk":2,"off":9,"en":"\xff\\q"},'
               b'{"id":2574,"blk":3,"off":11,"en":"\xff"},'
               b'{"id":2575,"blk":4,"off":13,"en":"ordinary\\nend"}]}')
        data, skipped = decode_player_catalog(raw)
        self.assertEqual(skipped, 2)
        self.assertEqual([r['en'] for r in data['records']],
                         ['ordinary', '', '', 'ordinary\nend'])
        self.assertEqual(data['records'][1], {'id':2500,'blk':2,'off':9,'en':''})

    def test_nonreserved_data_and_metadata_are_preserved(self):
        original = {'schema': 'synthetic', 'records': [
            {'id':100,'blk':0,'off':0,'en':'braces { } and "quotes"'},
            {'id':2575,'blk':1,'off':2,'en':'\\ and 控制\x1f'}]}
        data, skipped = decode_player_catalog(json.dumps(original, ensure_ascii=False).encode())
        self.assertEqual(data, original)
        self.assertEqual(skipped, 0)

    def test_corrupt_nonreserved_or_unknown_shape_fails(self):
        for raw in [b'{"records":[{"id":100,"blk":0,"off":0,"en":"\xff"}]}',
                    b'{"records":[{"id":02500,"blk":0,"off":0,"en":"ordinary"}]}',
                    b'{"records":[{"en":"ordinary","id":100}]}',
                    b'{"records":[{"id":100,"blk":0,"off":0,"en":"unfinished}',
                    b'{"records":[{"id":100,"blk":0,"off":0,"en":"ordinary"},'
                    b'{"id":100,"blk":0,"off":0,"en":"ordinary"}]}']:
            with self.assertRaises(ValueError):
                decode_player_catalog(raw)

    def test_nested_duplicate_and_nonstring_records_fail_before_decode(self):
        record=b'{"id":2500,"blk":0,"off":0,"en":"\xff"}'
        for raw in [b'{"outer":{"records":[]},"records":['+record+b']}',
                    b'{"records":[],"records":['+record+b']}',
                    b'{"tree":{"records":['+record+b']},"records":[]}',
                    b'{"records":[{"id":2500,"blk":0,"off":0,"en":{}}]}',
                    b'{"records":[{"id":2500,"blk":0,"off":0,"en":"\xff","extra":0}]}']:
            with self.assertRaises(ValueError):
                decode_player_catalog(raw)


if __name__ == '__main__':
    unittest.main()
