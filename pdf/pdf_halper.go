package pdf

import (
	"strings"

	"github.com/jung-kurt/gofpdf"
)

func fitText(pdf *gofpdf.Fpdf, txt string, maxWidth float64) string {
	if txt == "" {
		return ""
	}

	if pdf.GetStringWidth(txt) <= maxWidth {
		return txt
	}

	runes := []rune(txt)
	for len(runes) > 0 {
		runes = runes[:len(runes)-1]
		candidate := string(runes) + "…"
		if pdf.GetStringWidth(candidate) <= maxWidth {
			return candidate
		}
	}
	return ""
}

func formatTrajanje(traj string, tipd string) string {
	parts := strings.Split(traj, ":")

	switch tipd {

	case "1":
		// dd:hh:mm → hh:mm
		if len(parts) == 3 {
			return parts[1] + ":" + parts[2]
		}
		return traj

	case "3":
		// dd:hh:mm → dddd:mm
		if len(parts) == 3 {
			return parts[0] + ":" + parts[2]
		}
		return traj

	case "2":
		// dd:hh:mm → hh:mm
		if len(parts) == 3 {
			return parts[1] + ":" + parts[2]
		}
		return traj

	default:
		return traj
	}
}

func resetX(pdf *gofpdf.Fpdf) {
	left, _, _, _ := pdf.GetMargins()
	pdf.SetX(left)
}

func withLayout(pdf *gofpdf.Fpdf, fn func()) {
	x, y := pdf.GetXY()
	fn()
	pdf.SetXY(x, y)
}
