package awtree

// DetectConfig controls configurable detection thresholds and diagnostics.
type DetectConfig struct {
	MajorityThresholdPct int  `json:"majority_threshold_pct"`
	MaxButtonWidth       int  `json:"max_button_width"`
	MaxButtonLabelLen    int  `json:"max_button_label_len"`
	Debug                bool `json:"debug"`
}

// Option configures Detect.
type Option func(*DetectConfig)

// DefaultDetectConfig returns the default detection settings.
func DefaultDetectConfig() DetectConfig {
	return DetectConfig{
		MajorityThresholdPct: 60,
		MaxButtonWidth:       30,
		MaxButtonLabelLen:    20,
	}
}

// WithMajorityThresholdPct sets the percentage threshold used for bar detection.
func WithMajorityThresholdPct(pct int) Option {
	return func(cfg *DetectConfig) {
		if pct >= 1 && pct <= 100 {
			cfg.MajorityThresholdPct = pct
		}
	}
}

// WithMaxButtonWidth sets the maximum scan width for bracketed buttons.
func WithMaxButtonWidth(width int) Option {
	return func(cfg *DetectConfig) {
		if width >= 2 {
			cfg.MaxButtonWidth = width
		}
	}
}

// WithMaxButtonLabelLen sets the maximum allowed button label length.
func WithMaxButtonLabelLen(length int) Option {
	return func(cfg *DetectConfig) {
		if length >= 1 {
			cfg.MaxButtonLabelLen = length
		}
	}
}

// WithDebug enables detection diagnostics on the returned ElementMap.
func WithDebug(enabled bool) Option {
	return func(cfg *DetectConfig) {
		cfg.Debug = enabled
	}
}

func applyDetectOptions(opts []Option) DetectConfig {
	cfg := DefaultDetectConfig()
	for _, opt := range opts {
		if opt != nil {
			opt(&cfg)
		}
	}
	return cfg
}
