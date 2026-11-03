package wilson

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Session logging, the Asciiverse convention: off unless LOG_LEVEL is set,
// then written to stderr, e.g. `LOG_LEVEL=debug wilsons-revenge 2>wilson.log`.
// Debug logs every key as it arrives, each change of state and a snapshot
// twice a second; info logs only the changes.
var logger = newLogger()

func newLogger() *slog.Logger {
	var level slog.Level
	switch strings.ToLower(os.Getenv("LOG_LEVEL")) {
	case "debug", "trace":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		return slog.New(slog.NewTextHandler(io.Discard, nil))
	}
	return slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
}

// logChanges compares the state before and after a frame and logs what
// changed: the run ending, the gun, ammo, power-ups, jumping lost or back.
func logChanges(before GameState, s *GameState) {
	if !before.GameOver && s.GameOver {
		logger.Info("game over", "score", s.Score, "speed", s.Speed, "lane", s.ChickenLane, "frame", s.FrameCount)
	}
	if !before.HasGun && s.HasGun {
		logger.Info("gun picked up", "ammo", s.Ammo)
	} else if s.Ammo > before.Ammo {
		logger.Info("ammo picked up", "ammo", s.Ammo)
	}
	if before.InvincibleTimer <= 0 && s.InvincibleTimer > 0 {
		logger.Info("power-up: invincible and double points")
	}
	if s.Score < before.Score {
		logger.Info("power-down", "score", s.Score)
	}
	if before.JumpDisabled != s.JumpDisabled {
		logger.Info("jumping", "disabled", s.JumpDisabled)
	}
}

func logSnapshot(s *GameState) {
	logger.Debug("state", "frame", s.FrameCount, "score", s.Score, "speed", s.Speed,
		"lane", s.ChickenLane, "jumping", s.Jumping, "gun", s.HasGun, "ammo", s.Ammo,
		"invincible", s.InvincibleTimer, "objects", len(s.Objects))
}
