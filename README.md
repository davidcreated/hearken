# hearken

Offline voice control for the home, in Go.

hearken turns what someone says into something a home can do: *"turn off the
living room heater"*, *"set the kitchen to Away mode"*, *"turn off
everything"*. It matches the spoken name to your devices, rooms and scenes,
and asks which one when a name is ambiguous.

It is being built in the open, in stages:

1. **Commands (in progress):** a parser for spoken commands, a name matcher,
   and user shortcuts and nicknames. Pure Go, no dependencies.
2. **Speech recognition:** an on-device recogniser run from Go through
   [sherpa-onnx](https://github.com/k2-fsa/sherpa-onnx).
3. **Training:** fine-tuning that recogniser on home-command speech.
4. **Mobile:** builds of the library that iOS and Android apps can call.

See [docs/ROADMAP.md](docs/ROADMAP.md) for detail.

## Licence

Copyright 2026 David Roberts. Licensed under the
[Apache License 2.0](LICENSE).
