"""OPL2 事件層比對：dosgolem -opl-log vs DOSBox-X raw OPL（.dro v2）（docs/spec/005）。

    tools/py.sh tools/opl_events.py <dosgolem opl-log> <DOSBox-X .dro> [--json out.json] [--flip N]

dosgolem 端照 DOSBox-X `Adlib::Capture` 的規則過濾（spec 005 §3）後，與 DRO 的（暫存器, 值）序列逐筆比。
--flip N：反向對照，把 dosgolem 過濾後第 N 筆的值改掉，工具必須回報該位置不一致。
輸出兩種判定：pass（spec 005 §4 原判準）與 revised.pass（§4.1 修訂：不計結束遊戲造成的結尾靜音）；結束碼看修訂判定。
"""
import difflib
import json
import struct
import sys

TIMER_REGS = {0x02, 0x03, 0x04}  # OPL2 模式下被 Chip::Write 攔下，不進快取也不擷取


def capture_table():
    """DOSBox-X MakeTables()：raw 索引 → 暫存器。"""
    regs = [0x01, 0x04, 0x05, 0x08, 0xBD]
    for i in range(24):
        if (i & 7) < 6:
            regs += [0x20 + i, 0x40 + i, 0x60 + i, 0x80 + i, 0xE0 + i]
    for i in range(9):
        regs += [0xA0 + i, 0xB0 + i, 0xC0 + i]
    return regs


TABLE = capture_table()
IN_TABLE = set(TABLE)


def read_golem(path):
    out = []
    for line in open(path, encoding="utf-8"):
        if not line.strip() or line.startswith("#"):
            continue
        step, reg, val = line.split()
        out.append((int(step), int(reg, 16), int(val, 16)))
    return out


def golem_capture(writes):
    """回 [(步數, 暫存器, 值)]：DOSBox-X 在同樣寫入序列下會擷取到的內容（含開頭的快取傾印）。"""
    cache = [0] * 256
    out = []
    started = False
    for step, reg, val in writes:
        if reg in TIMER_REGS:
            continue
        if not started:
            trigger = (0xB0 <= reg <= 0xB8 and val & 0x20) or (reg == 0xBD and (val & 0x3F) > 0x20)
            if trigger:
                started = True
                for r in range(256):  # WriteCache()
                    v = cache[r]
                    if 0xB0 <= r <= 0xB8:
                        v &= ~0x20
                    if r == 0xBD:
                        v &= ~0x1F
                    if v and r in IN_TABLE:
                        out.append((step, r, v))
                if reg in IN_TABLE:
                    out.append((step, reg, val))
        elif reg in IN_TABLE and cache[reg] != val:
            out.append((step, reg, val))
        cache[reg] = val
    return out


def read_dro(path):
    data = open(path, "rb").read()
    if data[:8] != b"DBRAWOPL":
        raise SystemExit("不是 DRO 檔：%s" % path)
    ver_hi, ver_lo, commands, ms, hw, fmt, comp, d256, dshift, tsize = struct.unpack_from("<HHIIBBBBBB", data, 8)
    if (ver_hi, ver_lo) != (2, 0):
        raise SystemExit("只支援 DRO v2.0，這份是 %d.%d" % (ver_hi, ver_lo))
    pos = 26
    to_reg = list(data[pos:pos + tsize])
    pos += tsize
    out = []
    t = 0
    for i in range(commands):
        raw, val = data[pos + 2 * i], data[pos + 2 * i + 1]
        if raw == d256:
            t += val + 1
        elif raw == dshift:
            t += (val + 1) << 8
        else:
            if raw & 0x80:
                raise SystemExit("第 %d 筆寫到第二個埠，OPL2 不該出現" % i)
            out.append((t, to_reg[raw], val))
    meta = {"commands": commands, "milliseconds": ms, "hardware": hw, "table_size": tsize}
    return out, meta


def main(argv):
    args = [a for a in argv]
    out_json = None
    flip = None
    if "--json" in args:
        i = args.index("--json"); out_json = args[i + 1]; del args[i:i + 2]
    if "--flip" in args:
        i = args.index("--flip"); flip = int(args[i + 1]); del args[i:i + 2]
    if len(args) != 2:
        print(__doc__)
        return 2
    golem = golem_capture(read_golem(args[0]))
    dro, meta = read_dro(args[1])
    if flip is not None:
        s, r, v = golem[flip]
        golem[flip] = (s, r, v ^ 0xFF)
    # spec 005 §4.1（修訂）：DOSBox-X 擷取結尾若是結束遊戲造成的全聲道靜音（10 ms 內 A0–A8／B0–B8 全寫 0），另算一份不含它的判定
    quit_tail = 0
    while quit_tail < len(dro):
        t, r, v = dro[len(dro) - 1 - quit_tail]
        if v == 0 and (0xA0 <= r <= 0xA8 or 0xB0 <= r <= 0xB8) and dro[-1][0] - t <= 10:
            quit_tail += 1
        else:
            break
    if quit_tail < 9:
        quit_tail = 0
    g = [(r, v) for _, r, v in golem]
    d = [(r, v) for _, r, v in dro]
    n = min(len(g), len(d))
    prefix = next((i for i in range(n) if g[i] != d[i]), n)
    result = {"layer": "opl-events", "golem_captured": len(g), "dosboxx_captured": len(d),
              "compared": n, "identical_prefix": prefix, "dro": meta}
    if prefix < n:
        ctx = lambda seq: ["%02X=%02X" % x for x in seq[max(0, prefix - 5):prefix + 6]]
        result["first_mismatch"] = {"index": prefix, "golem": ctx(g), "dosboxx": ctx(d)}
        sm = difflib.SequenceMatcher(None, g[:n], d[:n], autojunk=False)
        result["aligned_mismatch_ops"] = sum(max(i2 - i1, j2 - j1) for tag, i1, i2, j1, j2 in sm.get_opcodes() if tag != "equal")
    # 時間診斷：dosgolem 步數 vs DRO 毫秒（取前 n 筆）的最小平方斜率
    if n > 2:
        xs = [golem[i][0] - golem[0][0] for i in range(n)]
        ys = [dro[i][0] for i in range(n)]
        mx, my = sum(xs) / n, sum(ys) / n
        sxx = sum((x - mx) ** 2 for x in xs)
        result["ms_per_million_steps"] = round(sum((x - mx) * (y - my) for x, y in zip(xs, ys)) / sxx * 1e6, 3) if sxx else None
    result["pass"] = prefix == n and n >= 5000
    d2 = d[:len(d) - quit_tail]
    n2 = min(len(g), len(d2))
    prefix2 = next((i for i in range(n2) if g[i] != d2[i]), n2)
    result["quit_silence_tail"] = quit_tail
    result["revised"] = {"compared": n2, "identical_prefix": prefix2, "pass": prefix2 == n2 and n2 >= 5000}
    print(json.dumps(result, ensure_ascii=False, indent=2))
    if out_json:
        with open(out_json, "w", encoding="utf-8") as f:
            json.dump(result, f, ensure_ascii=False, indent=2)
    return 0 if result["revised"]["pass"] else 1


if __name__ == "__main__":
    sys.exit(main(sys.argv[1:]))
