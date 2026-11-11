# tideglass

Next high and low tide from a local predictions file. No network.

Author: synthalorian

## Build

```sh
go build -o tideglass .
```

## Test

```sh
go test ./...
```

## Usage

Predictions are a JSON array. `kind` is `H` or `L`. Times are RFC3339 UTC.

```sh
tideglass --file predictions.json --now 2026-11-10T12:00:00Z
```

```
next_high 2026-11-10T15:00:00Z 4.2
next_low 2026-11-10T21:30:00Z 0.4
```

The command prints the earliest high and the earliest low strictly after `--now`, height to one decimal. It omits a line when that extreme does not exist. If neither exists, it prints `no upcoming tide` and exits 0.

Bad JSON, a kind other than `H` or `L`, or missing `--file` / `--now` exits 2.
