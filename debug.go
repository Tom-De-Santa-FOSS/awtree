package awtree

type DebugEvent struct {
	Detector string `json:"detector"`
	Accepted bool   `json:"accepted"`
	Reason   string `json:"reason"`
	Label    string `json:"label,omitempty"`
	Bounds   *Rect  `json:"bounds,omitempty"`
}

type DebugInfo struct {
	Config DetectConfig `json:"config"`
	Events []DebugEvent `json:"events"`
}

type debugCollector struct {
	info *DebugInfo
}

func newDebugCollector(cfg DetectConfig) *debugCollector {
	if !cfg.Debug {
		return nil
	}
	return &debugCollector{info: &DebugInfo{Config: cfg, Events: make([]DebugEvent, 0, 32)}}
}

func (d *debugCollector) accept(detector string, el Element, reason string) {
	if d == nil {
		return
	}
	bounds := el.Bounds
	d.info.Events = append(d.info.Events, DebugEvent{
		Detector: detector,
		Accepted: true,
		Reason:   reason,
		Label:    el.Label,
		Bounds:   &bounds,
	})
}

func (d *debugCollector) reject(detector, reason string, bounds Rect, label string) {
	if d == nil {
		return
	}
	b := bounds
	d.info.Events = append(d.info.Events, DebugEvent{
		Detector: detector,
		Accepted: false,
		Reason:   reason,
		Label:    label,
		Bounds:   &b,
	})
}
