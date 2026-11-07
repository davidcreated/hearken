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

## Use it

```sh
go get github.com/davidcreated/hearken
```

```go
parser, err := hearken.NewParser(`^hey (home )?flex `) // your wake phrase, or ""
if err != nil {
	return err
}
command, ok := parser.Parse("Hey Flex, set the kitchen to Away mode")
if !ok {
	return errors.New("not a command")
}
// command: {Action: "on", Kind: "mode", Target: "away", Space: "kitchen"}

switch t := hearken.Resolve(command, home).(type) {
case hearken.ModeTarget:
	// turn t.Mode on
case hearken.Ambiguous:
	// ask which of t.Options was meant
case hearken.NoMatch:
	// say nothing has that name
}
```

`home` lists your devices, rooms and modes (`hearken.Home`). Users can add
their own phrases with `Parser.FindShortcut` and other names for things with
`FindNickname`.

### From a terminal

```sh
go install github.com/davidcreated/hearken/cmd/hearken@latest
hearken parse "turn off the living room heater"
hearken match --home examples/home.json "set the kitchen to away mode"
```

## The spec

[`spec/`](spec) holds every case the library is tested against, as plain
JSON: what was said and what it should mean. It is language-neutral, so a
port to another language can run the same cases and stay in step.

## Licence

Copyright 2026 David Roberts. Licensed under the
[Apache License 2.0](LICENSE).
