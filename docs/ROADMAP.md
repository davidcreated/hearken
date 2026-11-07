# Roadmap

## 1. Commands (done)

Parse what was said into an action and a target, match the target to a
home's devices, rooms and modes, and support user shortcuts and nicknames.
Pure Go with no dependencies, tested against the JSON cases in `spec/`.

## 2. Speech recognition in Go

Turn microphone audio into text on the device, with no cloud.

- Run models through [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx)'s
  Go bindings: a streaming recogniser for commands and a keyword spotter
  for the wake phrase.
- Start from an existing open model whose licence allows commercial use.
- Measure it on the phrases in `spec/` before any training, so later
  training has a baseline to beat.

## 3. Training

Fine-tune the recogniser on home-command speech.

- Training runs in Python (the model toolkits are Python), in a
  `training/` folder with its own instructions; Go only runs the result.
- Data follows [DATA.md](DATA.md). Speakers are split between training and
  testing so results reflect voices the model has not heard.
- Each trained model is exported to ONNX and published with its word
  error rate on the test set next to the baseline's.

## 4. Mobile

Builds of the library that iOS and Android apps call directly.

- A small C interface over the Go API, built as a shared library for
  Android (per CPU architecture) and a static library for iOS.
- An example of calling it from a Flutter app through Dart FFI.
