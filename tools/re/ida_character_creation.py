"""只解P6的128-byte runtime片段；線性項目是候選，末端call另有執行錨點。"""
from pathlib import Path
import hashlib
import json
import ida_auto
import ida_bytes
import ida_ida
import ida_nalt
import ida_pro
import ida_segment
import ida_ua
import idaapi
import idc

ida_auto.auto_wait()
source = Path(ida_nalt.get_input_file_path())
raw = source.read_bytes()
metadata = json.loads(source.with_suffix('.json').read_text())
assert len(raw) == 128 and raw[-3:] == bytes.fromhex('e8e2a2')
assert metadata['runtime_segment'] == 0x0110 and metadata['runtime_offset'] == 0x9515
assert metadata['length'] == 128 and metadata['sha256'] == hashlib.sha256(raw).hexdigest()
start = ida_ida.inf_get_min_ea()
segment = ida_segment.getseg(start)
assert segment is not None
ida_segment.set_segm_addressing(segment, 0)

def decode(delta):
    address = start + delta
    instruction = ida_ua.insn_t()
    length = ida_ua.decode_insn(instruction, address)
    if length <= 0 or delta + length > len(raw):
        return None
    ida_bytes.del_items(address, ida_bytes.DELIT_SIMPLE, length)
    ida_ua.create_insn(address)
    return {'ida_blob_ea': address, 'runtime_segment': 0x0110, 'runtime_offset': 0x9515 + delta,
            'bytes': ida_bytes.get_bytes(address, length).hex(), 'length': length,
            'mnemonic': instruction.get_canon_mnem(),
            'operands': [idc.print_operand(address, i) for i in range(3) if idc.print_operand(address, i)]}

rows = []
delta = 0
while delta < len(raw):
    row = decode(delta)
    if row:
        rows.append(row)
        delta += row['length']
    else:
        delta += 1
anchor = decode(125)
assert anchor and anchor['bytes'] == 'e8e2a2' and anchor['mnemonic'] == 'call' and anchor['length'] == 3
anchor['evidence_level'] = 'confirmed call instruction; executed source paired by P6 runtime receipt'
anchor['target_runtime_segment'] = 0x0110
anchor['target_runtime_offset'] = (0x9595 + int.from_bytes(raw[-2:], 'little', signed=True)) & 0xffff
assert anchor['target_runtime_offset'] == 0x3877
report = {'schema': 'wiz7-p6-ida-snippet/1', 'kernel_version': idaapi.get_kernel_version(),
          'input_name': source.name, 'input_sha256': metadata['sha256'], 'input_length': len(raw),
          'ida_address_space': 'raw binary blob; branch operand text retains IDA blob address space',
          'runtime_address_space': 'dosgolem runtime CS:IP, 0110:9515; not an IDA linear address',
          'linear_decode_level': 'candidate boundaries only; unobserved control flow not confirmed',
          'linear_candidates': rows, 'executed_call_anchor': anchor,
          'limit': 'Only this128-byte slice; no full-function, F-path or formal-route claim.'}
with Path('/out/ida-creation-source.json').open('x') as handle:
    handle.write(json.dumps(report, ensure_ascii=False, indent=2) + '\n')
ida_pro.qexit(0)
