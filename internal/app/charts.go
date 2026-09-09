// This file holds the generic chart primitives used by every dashboard widget.
//
// Nothing in here knows where its numbers came from: each renderer takes plain
// `model` values plus layout numbers, so the current seed data and a future API
// response render through exactly the same code path. Binding data to a chart
// happens only in widgets.go.

package app

import (
	"fmt"
	"strings"

	"image/color"

	"charm.land/lipgloss/v2"
	"github.com/three-thirds/hackatime-tui/internal/model"
)

const (
	// minChartWidth clamps card width so Lipgloss never panics on very small
	// terminals.
	minChartWidth = 12
	// chartChrome is the horizontal cells a card spends on its border and
	// padding. Lipgloss counts both inside Width, so every renderer subtracts
	// it before laying out its own columns.
	chartChrome = 4
)

// ValueFormat selects what a chart prints in its trailing value column.
type ValueFormat int

const (
	ValueDuration ValueFormat = iota // "3h 04m"
	ValuePercent                     // "38%"
)

// BarChartConfig describes a horizontal bar chart. Only Title, Width and Items
// are required; the zero value of every other field is a usable default.
type BarChartConfig struct {
	Title     string
	Width     int                   // content width handed to the slot border
	MinHeight int                   // box height floor, keeps a row aligned
	Items     []model.BreakdownItem // rows, already ordered by the caller
	Palette   []color.Color         // cycled per bar, falls back to ChartPalette
	Format    ValueFormat           // what the value column shows
	LabelCap  int                   // longest label column, 0 means 18
	Empty     string                // message shown when Items is empty
}

// RenderBarChart draws one labelled bar per item, scaled against the largest
// value in the set, with a right-hand value column that never wraps.
func RenderBarChart(cfg BarChartConfig) string {
	width := max(cfg.Width, minChartWidth)

	lines := []string{SlotTitleStyle.Render(cfg.Title)}
	if len(cfg.Items) == 0 {
		lines = append(lines, DimText.Render(emptyMessage(cfg.Empty)))
		return chartBox(width, cfg.MinHeight, lines)
	}

	labelCap := cfg.LabelCap
	if labelCap <= 0 {
		labelCap = 18
	}
	palette := cfg.Palette
	if len(palette) == 0 {
		palette = ChartPalette
	}

	labelWidth := labelColumnWidth(cfg.Items, labelCap)
	valueWidth := valueColumnWidth(cfg.Items, cfg.Format)
	// label + gap + bar + gap + value has to fit inside the card content.
	barWidth := max(width-chartChrome-labelWidth-valueWidth-2, 1)
	largest := largestValue(cfg.Items, cfg.Format)

	for i, item := range cfg.Items {
		label := truncateLabel(item.Name, labelWidth)
		filled := scaleToWidth(barValue(item, cfg.Format), largest, barWidth)
		bar := lipgloss.NewStyle().Foreground(palette[i%len(palette)]).Render(strings.Repeat("█", filled))
		track := DimText.Render(strings.Repeat("·", barWidth-filled))
		// Pad the value before styling it: fmt widths count ANSI bytes, so a
		// styled string would never line up.
		value := fmt.Sprintf("%*s", valueWidth, formatValue(item, cfg.Format))

		lines = append(lines, fmt.Sprintf("%-*s %s%s %s",
			labelWidth, label, bar, track, DimText.Render(value)))
	}

	return chartBox(width, cfg.MinHeight, lines)
}

// ColumnChartConfig describes a vertical bar chart plotted over time.
type ColumnChartConfig struct {
	Title      string
	Width      int
	PlotHeight int // rows of plot area, excluding the title and axis labels
	Points     []model.TimelinePoint
	Color      color.Color
	Empty      string
}

// RenderColumnChart plots each point as a vertical bar. Bars use eighth-block
// characters so a column can express a fractional row instead of snapping to
// the nearest whole one.
func RenderColumnChart(cfg ColumnChartConfig) string {
	width := max(cfg.Width, minChartWidth)
	plotHeight := max(cfg.PlotHeight, 3)

	lines := []string{SlotTitleStyle.Render(cfg.Title)}
	if len(cfg.Points) == 0 {
		lines = append(lines, DimText.Render(emptyMessage(cfg.Empty)))
		return chartBox(width, plotHeight+4, lines)
	}

	barColor := cfg.Color
	if barColor == nil {
		barColor = DraculaPurple
	}

	available := width - chartChrome
	columnWidth := max(available/len(cfg.Points), 3)
	barWidth := columnWidth - 1 // one trailing cell keeps the columns apart

	var largest float64
	for _, point := range cfg.Points {
		largest = max(largest, float64(point.Duration))
	}

	// Height of every column in fractional plot rows, computed once and reused
	// by each row of the plot below.
	heights := make([]float64, len(cfg.Points))
	for i, point := range cfg.Points {
		if largest > 0 {
			heights[i] = float64(point.Duration) / largest * float64(plotHeight)
		}
	}

	barStyle := lipgloss.NewStyle().Foreground(barColor)
	for row := plotHeight - 1; row >= 0; row-- {
		var line strings.Builder
		for _, height := range heights {
			line.WriteString(barStyle.Render(strings.Repeat(verticalBlock(height-float64(row)), barWidth)))
			line.WriteString(" ")
		}
		lines = append(lines, line.String())
	}

	var axis strings.Builder
	for _, point := range cfg.Points {
		axis.WriteString(centerLabel(point.Label, columnWidth))
	}
	lines = append(lines, DimText.Render(axis.String()))

	return chartBox(width, plotHeight+4, lines)
}

// HeatmapConfig describes a GitHub-style contribution grid.
type HeatmapConfig struct {
	Title string
	Width int
	Days  []model.HeatmapDay // consecutive and chronological, oldest first
	Empty string
}

// RenderHeatmap lays days out as weekday rows by week columns, keeping only the
// most recent weeks that fit the available width.
func RenderHeatmap(cfg HeatmapConfig) string {
	width := max(cfg.Width, minChartWidth)

	lines := []string{SlotTitleStyle.Render(cfg.Title)}
	if len(cfg.Days) == 0 {
		lines = append(lines, DimText.Render(emptyMessage(cfg.Empty)))
		return chartBox(width, 12, lines)
	}

	const weekdayLabelWidth = 4
	grid := buildHeatmapGrid(cfg.Days)

	// Prefer wide cells, and fall back to single-cell columns once the history
	// is long enough that two-wide columns would overflow.
	available := max(width-chartChrome-weekdayLabelWidth, 1)
	cellWidth := 2
	if len(grid)*cellWidth > available {
		cellWidth = 1
	}
	if visible := available / cellWidth; len(grid) > visible {
		grid = grid[len(grid)-visible:]
	}

	weekdayLabels := [7]string{"", "Mon", "", "Wed", "", "Fri", ""}
	for row := range 7 {
		var line strings.Builder
		line.WriteString(DimText.Render(fmt.Sprintf("%-*s", weekdayLabelWidth, weekdayLabels[row])))
		for _, week := range grid {
			line.WriteString(heatmapSwatch(week[row], cellWidth))
		}
		lines = append(lines, line.String())
	}

	var legend strings.Builder
	legend.WriteString(DimText.Render("Less "))
	for _, level := range HeatmapLevels {
		legend.WriteString(lipgloss.NewStyle().Foreground(level).Render("█ "))
	}
	legend.WriteString(DimText.Render("More"))
	lines = append(lines, "", legend.String())

	return chartBox(width, 12, lines)
}

// MeterConfig describes a single-value progress meter.
type MeterConfig struct {
	Title     string
	Width     int
	Percent   int
	Status    string // short headline shown next to the percentage
	Detail    string // supporting line under the bar
	Color     color.Color
	MinHeight int // box height floor, keeps a row aligned
}

// RenderMeter draws a percentage headline over a filled progress track.
func RenderMeter(cfg MeterConfig) string {
	width := max(cfg.Width, minChartWidth)
	percent := min(max(cfg.Percent, 0), 100)

	fillColor := cfg.Color
	if fillColor == nil {
		fillColor = DraculaGreen
	}

	barWidth := max(width-chartChrome, 1)
	filled := percent * barWidth / 100
	bar := lipgloss.NewStyle().Foreground(fillColor).Render(strings.Repeat("█", filled)) +
		DimText.Render(strings.Repeat("░", barWidth-filled))

	lines := []string{
		SlotTitleStyle.Render(cfg.Title),
		fmt.Sprintf("%s %s",
			MetricValueStyle.Render(fmt.Sprintf("%d%%", percent)),
			DimText.Render(cfg.Status)),
		bar,
		DimText.Render(cfg.Detail),
	}

	return chartBox(width, cfg.MinHeight, lines)
}

// SplitSegment is one slice of a proportional split bar.
type SplitSegment struct {
	Label   string
	Percent int
	Detail  string // free-form trailing text, e.g. a formatted duration
	Color   color.Color
}

// SplitBarConfig describes a single bar divided between competing segments.
type SplitBarConfig struct {
	Title     string
	Width     int
	Segments  []SplitSegment
	Empty     string
	MinHeight int // box height floor, keeps a row aligned
}

// RenderSplitBar draws one bar shared by all segments plus a legend line each.
// Rounding leftovers go to the largest segment so the bar always fills exactly.
func RenderSplitBar(cfg SplitBarConfig) string {
	width := max(cfg.Width, minChartWidth)

	lines := []string{SlotTitleStyle.Render(cfg.Title)}
	if len(cfg.Segments) == 0 {
		lines = append(lines, DimText.Render(emptyMessage(cfg.Empty)))
		return chartBox(width, cfg.MinHeight, lines)
	}

	barWidth := max(width-chartChrome, 1)
	widths := make([]int, len(cfg.Segments))
	used, largest := 0, 0
	for i, segment := range cfg.Segments {
		widths[i] = max(segment.Percent, 0) * barWidth / 100
		used += widths[i]
		if segment.Percent > cfg.Segments[largest].Percent {
			largest = i
		}
	}
	widths[largest] = max(widths[largest]+barWidth-used, 0)

	var bar strings.Builder
	for i, segment := range cfg.Segments {
		bar.WriteString(lipgloss.NewStyle().
			Foreground(segmentColor(segment, i)).
			Render(strings.Repeat("█", widths[i])))
	}
	lines = append(lines, bar.String())

	for i, segment := range cfg.Segments {
		lines = append(lines, fmt.Sprintf("%s %-8s %s",
			lipgloss.NewStyle().Foreground(segmentColor(segment, i)).Render("█"),
			truncateLabel(segment.Label, 8),
			DimText.Render(fmt.Sprintf("%3d%%  %s", segment.Percent, segment.Detail))))
	}

	return chartBox(width, cfg.MinHeight, lines)
}

// chartBox wraps rendered chart lines in the shared slot border, growing the
// box to fit its content but never shrinking below minHeight. Lipgloss counts
// the border inside both dimensions, so minHeight is a whole-box height and the
// content floor is two rows below it.
func chartBox(width, minHeight int, lines []string) string {
	return SlotBorder.
		Width(max(width, minChartWidth)).
		Height(max(minHeight, len(lines)+2)).
		Render(strings.Join(lines, "\n"))
}

func emptyMessage(message string) string {
	if message == "" {
		return "No data yet"
	}
	return message
}

func segmentColor(segment SplitSegment, index int) color.Color {
	if segment.Color != nil {
		return segment.Color
	}
	return ChartPalette[index%len(ChartPalette)]
}

func barValue(item model.BreakdownItem, format ValueFormat) float64 {
	if format == ValuePercent {
		return item.Percentage
	}
	return float64(item.Duration)
}

func formatValue(item model.BreakdownItem, format ValueFormat) string {
	if format == ValuePercent {
		return fmt.Sprintf("%.0f%%", item.Percentage)
	}
	return formatDuration(item.Duration)
}

func largestValue(items []model.BreakdownItem, format ValueFormat) float64 {
	var largest float64
	for _, item := range items {
		largest = max(largest, barValue(item, format))
	}
	return largest
}

func labelColumnWidth(items []model.BreakdownItem, limit int) int {
	width := 0
	for _, item := range items {
		width = max(width, lipgloss.Width(item.Name))
	}
	return min(max(width, 6), limit)
}

func valueColumnWidth(items []model.BreakdownItem, format ValueFormat) int {
	width := 0
	for _, item := range items {
		width = max(width, lipgloss.Width(formatValue(item, format)))
	}
	return width
}

// scaleToWidth maps a value onto the bar track, guaranteeing that any non-zero
// value stays visible as at least one cell.
func scaleToWidth(value, largest float64, available int) int {
	if value <= 0 || largest <= 0 {
		return 0
	}
	return min(max(1, int(value/largest*float64(available))), available)
}

func truncateLabel(label string, width int) string {
	if lipgloss.Width(label) <= width {
		return label
	}
	if width <= 3 {
		return strings.Repeat(".", width)
	}
	return label[:width-3] + "..."
}

// verticalBlock picks the block character for how much of one plot row a bar
// fills: the whole row, a fraction of it, or nothing.
func verticalBlock(fill float64) string {
	if fill >= 1 {
		return "█"
	}
	if fill <= 0 {
		return " "
	}
	blocks := []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇"}
	return blocks[min(int(fill*float64(len(blocks))), len(blocks)-1)]
}

func centerLabel(label string, width int) string {
	if lipgloss.Width(label) >= width {
		return truncateLabel(label, width)
	}
	left := (width - lipgloss.Width(label)) / 2
	return strings.Repeat(" ", left) + label + strings.Repeat(" ", width-left-lipgloss.Width(label))
}

// heatmapCell is one day of the grid; absent days pad the first and last week.
type heatmapCell struct {
	present bool
	level   int
}

// buildHeatmapGrid bins consecutive days into week columns of seven weekday
// rows, padding the leading days of the first week. Days are assumed to be
// chronological and gapless, which is how both the seed data and the API
// summaries deliver them.
func buildHeatmapGrid(days []model.HeatmapDay) [][7]heatmapCell {
	lead := int(days[0].Date.Weekday())
	columns := (lead + len(days) + 6) / 7

	grid := make([][7]heatmapCell, columns)
	for i, day := range days {
		slot := lead + i
		grid[slot/7][slot%7] = heatmapCell{present: true, level: day.Level}
	}
	return grid
}

func heatmapSwatch(cell heatmapCell, cellWidth int) string {
	if !cell.present {
		return strings.Repeat(" ", cellWidth)
	}
	level := min(max(cell.level, 0), len(HeatmapLevels)-1)
	return lipgloss.NewStyle().Foreground(HeatmapLevels[level]).Render(strings.Repeat("█", cellWidth))
}
