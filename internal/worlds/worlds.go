// Package worlds plays the land and city generators as a looping demo:
// each world is sized to the window, collapses in front of you, holds for a
// moment once complete, and gives way to the next.
package worlds

import (
	"strings"
	"time"

	"atlas.games/internal/city"
	"atlas.games/internal/wfc"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	frame = 15 * time.Millisecond
	hold  = 4 * time.Second
)

type tickMsg time.Time

type cell struct {
	ch    string
	color string
	bold  bool
}

var unknown = cell{"?", "235", false}

// world is one generator, land or city, behind a common face.
type world interface {
	step() bool
	cell(x, y int) cell
}

type land struct{ g *wfc.WFC }

func (l land) step() bool { return l.g.Step() }

func (l land) cell(x, y int) cell {
	t := l.g.Grid[y][x]
	if !t.Collapsed {
		return unknown
	}
	switch t.Type {
	case wfc.Water:
		return cell{"~", "33", false}
	case wfc.Land:
		return cell{"█", "185", false}
	case wfc.Forest:
		return cell{"↑", "34", false}
	case wfc.Mountain:
		return cell{"▲", "255", true}
	case wfc.Lava:
		return cell{"░", "196", true}
	}
	return cell{" ", "", false}
}

type town struct{ g *city.WFC }

func (c town) step() bool { return c.g.Step() }

var roads = map[city.TileType]string{
	city.RoadV: "║", city.RoadH: "═", city.RoadTL: "╔", city.RoadTR: "╗",
	city.RoadBL: "╚", city.RoadBR: "╝", city.RoadTU: "╩", city.RoadTD: "╦",
	city.RoadTLT: "╠", city.RoadTRT: "╣", city.RoadCross: "╬",
}

func (c town) cell(x, y int) cell {
	t := c.g.Grid[y][x]
	if !t.Collapsed {
		return unknown
	}
	if r, ok := roads[t.Type]; ok {
		return cell{r, "244", false}
	}
	switch t.Type {
	case city.Building:
		return cell{"█", "255", true}
	case city.Commercial:
		return cell{"S", "220", true}
	case city.Park:
		return cell{"♣", "34", false}
	case city.Water:
		return cell{"~", "33", false}
	}
	return cell{" ", "", false}
}

type Model struct {
	w, h     int
	scene    int // even: land, odd: city
	cur      world
	steps    int
	done     bool
	doneAt   time.Time
	perFrame int
}

func NewModel() Model { return Model{} }

func (m Model) Init() tea.Cmd { return tick() }

func tick() tea.Cmd {
	return tea.Tick(frame, func(t time.Time) tea.Msg { return tickMsg(t) })
}

// next starts a fresh world for the current window size.
func (m *Model) next() {
	if m.w < 4 || m.h < 4 {
		m.cur = nil
		return
	}
	if m.scene%2 == 0 {
		m.cur = land{wfc.NewWFC(m.w, m.h)}
	} else {
		m.cur = town{city.NewWFC(m.w, m.h)}
	}
	m.scene++
	m.steps = 0
	m.done = false
	// About four seconds to collapse a world, whatever the window size.
	m.perFrame = m.w * m.h / 260
	if m.perFrame < 10 {
		m.perFrame = 10
	}
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height
		m.scene-- // keep the same kind of world after a resize
		if m.scene < 0 {
			m.scene = 0
		}
		m.next()
	case tea.KeyMsg:
		switch msg.String() {
		case " ", "enter", "right":
			m.next()
		default:
			return m, tea.Quit
		}
	case tickMsg:
		if m.cur == nil {
			return m, tick()
		}
		if m.done {
			if time.Since(m.doneAt) >= hold {
				m.next()
			}
			return m, tick()
		}
		for i := 0; i < m.perFrame; i++ {
			m.steps++
			// A contradiction can leave a cell that never collapses; stop
			// after a generous budget rather than spin on it.
			if m.cur.step() || m.steps > m.w*m.h*2 {
				m.done = true
				m.doneAt = time.Now()
				break
			}
		}
		return m, tick()
	}
	return m, nil
}

// View draws runs of same-styled cells in one go; styling every cell on
// its own is too slow for a full window at this frame rate.
func (m Model) View() string {
	if m.cur == nil {
		return ""
	}
	var sb strings.Builder
	for y := 0; y < m.h; y++ {
		var run strings.Builder
		var style cell
		flush := func() {
			if run.Len() == 0 {
				return
			}
			s := lipgloss.NewStyle().Bold(style.bold)
			if style.color != "" {
				s = s.Foreground(lipgloss.Color(style.color))
			}
			sb.WriteString(s.Render(run.String()))
			run.Reset()
		}
		for x := 0; x < m.w; x++ {
			c := m.cur.cell(x, y)
			if c.color != style.color || c.bold != style.bold {
				flush()
				style = c
			}
			run.WriteString(c.ch)
		}
		flush()
		if y < m.h-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
