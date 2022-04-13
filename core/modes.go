package core

var (
	ExternallySet     = &Mode{0, "External", false}
	DmxIn             = &Mode{1, "DMX in", false}
	DmxOut            = &Mode{2, "DMX out", false}
	DmxToPixelsWithIr = &Mode{3, "DMX > Pixels, IR enabled", true}
	DmxToPixels       = &Mode{4, "DMX > Pixels", false}
	RgbToPixelsWithIr = &Mode{5, "RGB > Pixels, IR enabled", true}
	RgbToPixels       = &Mode{6, "RGB > Pixels", false}
	FxToPixelsWithIr  = &Mode{7, "Fx > Pixels, IR enabled", true}
	FxToPixels        = &Mode{8, "Fx > Pixels", false}
	AutoFxWithIr      = &Mode{9, "Auto Fx, IR enabled", true}
	AutoFx            = &Mode{10, "Auto Fx", false}
	Emergency         = &Mode{11, "Emergency", false}
)

type Mode struct {
	Code            int    `json:"code"`
	Name            string `json:"name"`
	InfraredEnabled bool   `json:"infraredEnabled"`
}

var Modes = []*Mode{
	ExternallySet, DmxIn, DmxOut, DmxToPixelsWithIr, DmxToPixels, RgbToPixelsWithIr, RgbToPixels, FxToPixelsWithIr, FxToPixels, AutoFxWithIr, AutoFx, Emergency,
}
