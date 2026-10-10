"""Send one ordinary X11 command to the bounded normal-world controller."""
import argparse
import json
import time
from pathlib import Path

parser = argparse.ArgumentParser()
parser.add_argument('--number', type=int, required=True)
parser.add_argument('--label', default='player-step')
parser.add_argument('--click', type=int, nargs=2)
parser.add_argument('--move', type=int, nargs=2)
parser.add_argument('--key', choices=['Escape', 'space', 'Return', 'Up', 'Down', 'Left', 'Right', 'Tab'])
parser.add_argument('--text')
parser.add_argument('--wait', type=float, default=1)
parser.add_argument('--finish', action='store_true')
parser.add_argument('--record', choices=['start','stop'])
parser.add_argument('--live', action='store_true')
args = parser.parse_args()
directory = Path('/out')
assert 1 <= args.number <= 30
assert args.label.replace('-', '').isalnum() and len(args.label) <= 32
command = {'finish': True} if args.finish else {'label': args.label, 'wait': args.wait}
if args.move:
    command['move'] = args.move
if args.click:
    command['click'] = args.click
if args.key:
    command['key'] = args.key
if args.text:
    command['text'] = args.text
if args.record:
    command['record'] = args.record
if args.live:
    command['pause_after'] = False
with (directory / f'command-{args.number}.json').open('x') as stream:
    stream.write(json.dumps(command) + '\n')
if args.finish:
    print('normal_finish_requested', flush=True)
else:
    ack = directory / f'ack-{args.number}.json'
    deadline = time.monotonic() + 40
    while not ack.is_file() and time.monotonic() < deadline:
        time.sleep(.2)
    assert ack.is_file(), 'normal_command_ack_pending'
    print(ack.read_text(), flush=True)
