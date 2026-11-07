# Recordings and data

How speech recordings for training are prepared. Audio is never committed to
this repository: `/data/`, `*.wav`, `*.flac` and `*.onnx` are ignored.

## Consent

Every recording needs the speaker's consent to be used for training, and
separate consent before it is ever published. Recordings without consent
are not used.

## Audio

- WAV, 16 kHz, mono, 16-bit.
- One command per file, 1 to 8 seconds, with a little silence at each end.
- Recorded on the kind of device the model will run on (phones, a speaker
  across the room), in real rooms: quiet, with a TV on, in a kitchen.

## Manifest

Each set of recordings comes with a `manifest.jsonl`, one line per file:

```json
{"audio": "speaker-07/0003.wav", "text": "turn off the kitchen lights", "speaker": "speaker-07", "consent": true, "device": "Pixel 8", "room": "kitchen", "noise": "tv"}
```

- `text` is what was actually said, written normally. Tooling normalises it
  the same way the library does (`hearken.Normalize`).
- `speaker` is an anonymous id, never a name. The same person keeps the same
  id across sessions so their recordings stay on one side of the
  training/test split.
- `device`, `room` and `noise` are optional but help explain where the model
  does well or badly.

## What to record

Commands people really say to a home: switching devices and rooms on and
off, asking whether something is on, turning modes on, "turn off
everything", and the wake phrase on its own. Use the phrases in `spec/` as a
starting list, then the names of real devices and rooms in each home.
