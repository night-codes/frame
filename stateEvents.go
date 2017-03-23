package frame

const (
	// Withdrawn - window is not shown
	cWithdrawn = 1 << iota
	cIconified
	cMaximized
	cSticky
	cFullscreen
	cAbove
	cBelow
	cFocused
	cTiled
)

// State struct
type State struct {
	Hidden     bool
	Iconified  bool
	Maximized  bool
	Sticky     bool
	Fullscreen bool
	Above      bool
	Below      bool
	Focused    bool
	Tiled      bool
}
