# Asciiverse changes

Upstream: github.com/fezcode/atlas.games (MIT, fezcode), remote `upstream`,
branch `main` at `59f2de8` (2026-09-04). Our work is on the `asciiverse`
branch. Atlas Games is a collection behind one menu; per the catalog's
collection rule each title gets its own binary, and upstream's own menu
(`main.go`) is left as it is.

## Titles built from this repo

- **Wilson's Revenge** (`cmd/wilsons-revenge`): the runner on its own, no menu.
  The game's rules (`internal/wilson/engine.go`) are unchanged. Added:
  - `--demo`, an attract mode for the catalog preview: Wilson plays himself
    (`internal/wilson/autopilot.go`, looking ahead in each lane to shoot,
    change lane or jump, and going for pick-ups when the way is clear),
    starts after two seconds on the title, runs again three seconds after a
    crash, and any key quits.
  - Session logging, the Asciiverse convention, with Go's `log/slog`
    (`internal/wilson/log.go`): off unless `LOG_LEVEL` is set, then to
    stderr. Debug logs every key as it arrives, start, pause, crash, the
    gun, ammo, power-ups and jumping lost or regained, and a snapshot of
    the run twice a second. `LOG_LEVEL=debug wilsons-revenge 2>wilson.log`.
  - In `internal/wilson/ui.go` only the update loop changed, to drive the
    two above.
- **WFC Worlds** (`cmd/wfc-worlds`, `internal/worlds`), for Art & Demos: the land
  and city generators (`internal/wfc`, `internal/city`, unchanged) played as a
  demo. Each world fills the window (upstream's are a fixed 200x60), grows in
  about four seconds whatever the window size, holds for four, then the other
  kind follows; a resize starts a fresh world at the new size. Space, Enter or
  Right skips to the next world, any other key quits. Rows are styled in runs
  rather than cell by cell, so a full window keeps up at 15 ms a frame. A world
  that hits a contradiction stops after twice its cell count in steps instead
  of spinning.

Build: `go build -o build/ ./cmd/...`. Previews: `asciiverse/record-preview.sh`.
Releases: `.github/workflows/release.yml` on a `v<upstream>-av.<n>` tag builds
both titles for Linux (x86_64, aarch64, static), macOS and Windows, with .deb,
.rpm, .apk and Arch packages (`asciiverse/nfpm-<id>.yaml`).
