# tideglass

Go standard library only. No network. Module path `github.com/synthalorian/tideglass`.

## Test command

`go test ./...`

## Input

JSON array of objects: `{"time":"2026-11-10T15:00:00Z","height_ft":4.2,"kind":"H"}`. `kind` is `H` or `L`. Times are RFC3339 UTC.

## Command

`tideglass --file predictions.json --now 2026-11-10T12:00:00Z`

Print exactly two lines, or fewer if that extreme does not exist after `--now`:

```
next_high 2026-11-10T15:00:00Z 4.2
next_low 2026-11-10T21:30:00Z 0.4
```

Height is printed with one decimal. Choose the earliest event of each kind strictly after `--now`. If neither exists, print `no upcoming tide` and exit 0.

Bad JSON, bad kind, or missing flags: exit 2.

## Tests

Use a fixture in the repo. Cover the example, a file whose next event is only a low, and a file with nothing after `--now`.
