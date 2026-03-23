package tui

import (
	"fmt"
	"sort"
	"time"

	"github.com/System9-Software/sysix/internal/collector"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type focusedModel struct {
	procs     []collector.Process
	ports     []collector.Port
	showProcs bool
	showPorts bool
	sortBy    string
	err       error
	width     int
	height    int
}

func tickFocused(rate int) tea.Cmd {
	return tea.Tick(time.Duration(rate)*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func initialFocusedModel(showProcs, showPorts bool) focusedModel {
	snap, _ := collector.GetSnapshot()
	ports, _ := collector.GetPorts()
	m := focusedModel{
		showProcs: showProcs,
		showPorts: showPorts,
		sortBy:    "mem",
	}
	if snap != nil {
		m.procs = snap.Processes
	}
	m.ports = ports
	return m
}

func (m focusedModel) Init() tea.Cmd {
	return tickFocused(2)
}

func (m focusedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "c":
			m.sortBy = "cpu"
		case "m":
			m.sortBy = "mem"
		case "n":
			m.sortBy = "name"
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tickMsg:
		snap, _ := collector.GetSnapshot()
		ports, _ := collector.GetPorts()
		if snap != nil {
			m.procs = snap.Processes
		}
		if ports != nil {
			m.ports = ports
		}
		return m, tickFocused(2)
	}
	return m, nil
}

func (m focusedModel) sortedProcs() []collector.Process {
	procs := make([]collector.Process, len(m.procs))
	copy(procs, m.procs)
	switch m.sortBy {
	case "cpu":
		sort.Slice(procs, func(i, j int) bool { return procs[i].CPUPercent > procs[j].CPUPercent })
	case "name":
		sort.Slice(procs, func(i, j int) bool { return procs[i].Name < procs[j].Name })
	default:
		sort.Slice(procs, func(i, j int) bool { return procs[i].MemMB > procs[j].MemMB })
	}
	filtered := []collector.Process{}
	for _, p := range procs {
		if p.MemMB >= 1 {
			filtered = append(filtered, p)
		}
	}
	return filtered
}

func (m focusedModel) renderProcs() string {
	// logo is 6 lines, title+sort+header = 4 lines, footer = 1, padding = 3
	maxPerCol := m.height - 14
	if maxPerCol < 1 {
		maxPerCol = 10
	}

	filtered := m.sortedProcs()
	// cap total to 2 columns worth
	if len(filtered) > maxPerCol*2 {
		filtered = filtered[:maxPerCol*2]
	}

	colW := m.width/2 - 2
	rowFmt := "%-6d %-18s %6.1f%% %6.0fMB"
	hdr := lipgloss.NewStyle().Foreground(colorMuted).Render(
		fmt.Sprintf("%-6s %-18s %7s %7s", "PID", "Name", "CPU%", "Mem"),
	)
	sortHint := lipgloss.NewStyle().Foreground(colorMuted).Render(
		"sort: [m]em  [c]pu  [n]ame",
	)

	// split into two columns
	col1 := filtered
	col2 := []collector.Process{}
	if len(filtered) > maxPerCol {
		col1 = filtered[:maxPerCol]
		col2 = filtered[maxPerCol:]
	}

	leftRows := hdr + "\n"
	for _, p := range col1 {
		leftRows += valueStyle.Render(fmt.Sprintf(rowFmt, p.PID, truncate(p.Name, 18), p.CPUPercent, p.MemMB)) + "\n"
	}

	rightRows := hdr + "\n"
	for _, p := range col2 {
		rightRows += valueStyle.Render(fmt.Sprintf(rowFmt, p.PID, truncate(p.Name, 18), p.CPUPercent, p.MemMB)) + "\n"
	}

	leftTitle := panelTitleStyle.Render("[ Processes ]") + "\n" + sortHint
	rightTitle := panelTitleStyle.Render("[ cont. ]")

	left := lipgloss.NewStyle().Width(colW).Render(leftTitle + "\n\n" + leftRows)
	right := lipgloss.NewStyle().Width(colW).Render(rightTitle + "\n\n" + rightRows)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m focusedModel) renderPorts() string {
	maxRows := m.height - 14
	if maxRows < 1 {
		maxRows = 10
	}

	ports := m.ports
	if len(ports) > maxRows*2 {
		ports = ports[:maxRows*2]
	}

	colW := m.width/2 - 2
	rowFmt := "%-8d %-6s %-12s %d"
	hdr := lipgloss.NewStyle().Foreground(colorMuted).Render(
		fmt.Sprintf("%-8s %-6s %-12s %s", "Port", "Type", "Status", "PID"),
	)

	col1 := ports
	col2 := []collector.Port{}
	if len(ports) > maxRows {
		col1 = ports[:maxRows]
		col2 = ports[maxRows:]
	}

	leftRows := hdr + "\n"
	for _, p := range col1 {
		leftRows += valueStyle.Render(fmt.Sprintf(rowFmt, p.Port, p.Type, p.Status, p.PID)) + "\n"
	}

	rightRows := hdr + "\n"
	for _, p := range col2 {
		rightRows += valueStyle.Render(fmt.Sprintf(rowFmt, p.Port, p.Type, p.Status, p.PID)) + "\n"
	}

	left := lipgloss.NewStyle().Width(colW).Render(panelTitleStyle.Render("[ Ports ]") + "\n\n" + leftRows)
	right := lipgloss.NewStyle().Width(colW).Render(panelTitleStyle.Render("[ cont. ]") + "\n\n" + rightRows)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m focusedModel) renderBoth() string {
	maxRows := m.height - 14
	if maxRows < 1 {
		maxRows = 10
	}

	colW := m.width/2 - 2
	rowFmt := "%-6d %-16s %5.1f%% %5.0fMB"
	hdr := lipgloss.NewStyle().Foreground(colorMuted).Render(
		fmt.Sprintf("%-6s %-16s %6s %6s", "PID", "Name", "CPU%", "Mem"),
	)

	filtered := m.sortedProcs()
	if len(filtered) > maxRows {
		filtered = filtered[:maxRows]
	}

	procRows := hdr + "\n"
	for _, p := range filtered {
		procRows += valueStyle.Render(fmt.Sprintf(rowFmt, p.PID, truncate(p.Name, 16), p.CPUPercent, p.MemMB)) + "\n"
	}

	portHdr := lipgloss.NewStyle().Foreground(colorMuted).Render(
		fmt.Sprintf("%-8s %-6s %-10s %s", "Port", "Type", "Status", "PID"),
	)
	ports := m.ports
	if len(ports) > maxRows {
		ports = ports[:maxRows]
	}
	portRows := portHdr + "\n"
	for _, p := range ports {
		portRows += valueStyle.Render(fmt.Sprintf("%-8d %-6s %-10s %d", p.Port, p.Type, p.Status, p.PID)) + "\n"
	}

	left := lipgloss.NewStyle().Width(colW).Render(panelTitleStyle.Render("[ Processes ]") + "\n\n" + procRows)
	right := lipgloss.NewStyle().Width(colW).Render(panelTitleStyle.Render("[ Ports ]") + "\n\n" + portRows)

	return lipgloss.JoinHorizontal(lipgloss.Top, left, right)
}

func (m focusedModel) View() string {
	if m.width == 0 {
		return "loading..."
	}

	header := logoStyle.Render(logo) + "\n"
	footer := footerStyle.Render("q: quit")

	var content string
	if m.showProcs && m.showPorts {
		content = m.renderBoth()
	} else if m.showProcs {
		content = m.renderProcs()
	} else if m.showPorts {
		content = m.renderPorts()
	} else {
		content = "nothing to show"
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, content, footer)
}

func StartFocused(showProcs, showPorts bool) error {
	p := tea.NewProgram(initialFocusedModel(showProcs, showPorts), tea.WithAltScreen())
	_, err := p.Run()
	return err
}
