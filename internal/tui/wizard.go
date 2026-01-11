package tui

import (
	"fmt"
	"slices"
	"strings"

	"github.com/bithostio/bh/internal/api"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

const (
	stepProvider = iota
	stepRegion
	stepSize
	stepImage
	stepSSHKeys
	stepBackups
	stepName
	stepConfirm
	stepCount
)

type wizardModel struct {
	step   int
	cursor int
	offset int // scroll offset for long lists
	width  int
	height int

	// Data from API
	providers []api.Provider
	regions   []api.Region
	sizes     []api.Size
	images    []api.Image
	sshkeys   []api.SSHKey

	// Selections
	selectedProvider *api.Provider
	selectedRegion   *api.Region
	selectedSize     *api.Size
	selectedImage    *api.Image
	selectedKeys     []int // IDs of selected SSH keys
	backupsEnabled   bool
	serverName       string

	// Text input for name
	nameInput textinput.Model

	// State
	loading   bool
	completed bool
}

func newWizardModel() wizardModel {
	ti := textinput.New()
	ti.Placeholder = "my-server"
	ti.CharLimit = 64
	ti.Width = 30

	return wizardModel{
		step:         stepProvider,
		selectedKeys: []int{},
		nameInput:    ti,
		loading:      true,
	}
}

func (m wizardModel) SetProviders(providers []api.Provider) wizardModel {
	m.providers = providers
	m.loading = false
	m.cursor = 0
	m.offset = 0
	return m
}

func (m wizardModel) SetRegions(regions []api.Region) wizardModel {
	m.regions = regions
	m.loading = false
	m.cursor = 0
	m.offset = 0
	return m
}

func (m wizardModel) SetSizes(sizes []api.Size) wizardModel {
	m.sizes = sizes
	m.loading = false
	m.cursor = 0
	m.offset = 0
	return m
}

func (m wizardModel) SetImages(images []api.Image) wizardModel {
	m.images = images
	m.loading = false
	m.cursor = 0
	m.offset = 0
	return m
}

func (m wizardModel) SetSSHKeys(keys []api.SSHKey) wizardModel {
	m.sshkeys = keys
	m.loading = false
	m.cursor = 0
	m.offset = 0
	// Auto-select all keys by default
	if len(m.selectedKeys) == 0 {
		for _, k := range keys {
			m.selectedKeys = append(m.selectedKeys, k.ID)
		}
	}
	return m
}

func (m wizardModel) SetSize(width, height int) wizardModel {
	m.width = width
	m.height = height
	m.nameInput.Width = min(40, width-20)
	return m
}

func (m wizardModel) prevStep() wizardModel {
	if m.step > 0 {
		m.step--
		m.cursor = 0
		m.offset = 0
		m.loading = false
	}
	return m
}

// visibleLines returns how many list items can be displayed
func (m wizardModel) visibleLines() int {
	// Calculate used space:
	// - Title line: "NEW SERVER  Step X of Y" = 1
	// - Divider line = 1
	// - Blank line = 1
	// - Step title (e.g., "Select Size/Plan:") = 1
	// - Blank line after title = 1
	// - Scroll indicator (up) = 1 (if needed)
	// - Scroll indicator (down) = 1 (if needed)
	// - Blank lines before help = 2
	// - Help line = 1
	// Total fixed overhead: ~8-10 lines
	const headerLines = 5  // title, divider, blank, step title, blank
	const footerLines = 3  // blank, blank, help
	const scrollIndicators = 2 // potential up/down indicators

	reserved := headerLines + footerLines + scrollIndicators
	available := m.height - reserved

	// Ensure at least 5 items are visible for usability
	if available < 5 {
		available = 5
	}

	return available
}

// adjustOffset ensures the cursor is visible within the viewport
func (m wizardModel) adjustOffset() wizardModel {
	visible := m.visibleLines()
	listLen := m.getListLen()

	// If list fits on screen, no offset needed
	if listLen <= visible {
		m.offset = 0
		return m
	}

	// Cursor above viewport
	if m.cursor < m.offset {
		m.offset = m.cursor
	}

	// Cursor below viewport
	if m.cursor >= m.offset+visible {
		m.offset = m.cursor - visible + 1
	}

	// Clamp offset
	maxOffset := listLen - visible
	if m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}

	return m
}

func (m wizardModel) Update(msg tea.KeyMsg, root Model) (wizardModel, tea.Cmd) {
	var cmds []tea.Cmd

	// Handle name input step specially
	if m.step == stepName {
		switch msg.String() {
		case "enter":
			m.serverName = m.nameInput.Value()
			if m.serverName == "" {
				m.serverName = "server"
			}
			m.step++
			return m, nil
		default:
			var cmd tea.Cmd
			m.nameInput, cmd = m.nameInput.Update(msg)
			return m, cmd
		}
	}

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
			m = m.adjustOffset()
		}
	case "down", "j":
		if m.cursor < m.getListLen()-1 {
			m.cursor++
			m = m.adjustOffset()
		}
	case " ", "tab":
		// Toggle for SSH keys multi-select
		if m.step == stepSSHKeys && m.cursor < len(m.sshkeys) {
			keyID := m.sshkeys[m.cursor].ID
			m.selectedKeys = toggleSelection(m.selectedKeys, keyID)
		}
		// Toggle for backups
		if m.step == stepBackups {
			m.backupsEnabled = !m.backupsEnabled
		}
	case "enter":
		switch m.step {
		case stepProvider:
			if m.cursor < len(m.providers) {
				m.selectedProvider = &m.providers[m.cursor]
				m.step++
				m.loading = true
				m.cursor = 0
				// Clear dependent selections
				m.selectedRegion = nil
				m.selectedSize = nil
				m.selectedImage = nil
				m.regions = nil
				m.sizes = nil
				m.images = nil
				cmds = append(cmds, root.fetchRegions(m.selectedProvider.ID))
			}
		case stepRegion:
			if m.cursor < len(m.regions) {
				m.selectedRegion = &m.regions[m.cursor]
				m.step++
				m.loading = true
				m.cursor = 0
				// Clear dependent selections
				m.selectedSize = nil
				m.sizes = nil
				cmds = append(cmds, root.fetchSizes(m.selectedRegion.ID, m.selectedProvider.ID))
			}
		case stepSize:
			if m.cursor < len(m.sizes) {
				m.selectedSize = &m.sizes[m.cursor]
				m.step++
				m.loading = true
				m.cursor = 0
				cmds = append(cmds, root.fetchImages(m.selectedProvider.ID))
			}
		case stepImage:
			if m.cursor < len(m.images) {
				m.selectedImage = &m.images[m.cursor]
				m.step++
				m.loading = true
				m.cursor = 0
				cmds = append(cmds, root.fetchSSHKeys())
			}
		case stepSSHKeys:
			if len(m.selectedKeys) > 0 {
				m.step++
				m.cursor = 0
			}
		case stepBackups:
			m.step++
			m.nameInput.Focus()
		case stepConfirm:
			// Create server
			req := &api.CreateServerRequest{
				Name:           m.serverName,
				ProviderID:     m.selectedProvider.ID,
				RegionID:       m.selectedRegion.ID,
				SizeID:         m.selectedSize.ID,
				ImageID:        m.selectedImage.ID,
				KeyIDs:         m.selectedKeys,
				BackupsEnabled: m.backupsEnabled,
				Terms:          true,
			}
			m.completed = true
			cmds = append(cmds, root.createServer(req))
		}
	}

	return m, tea.Batch(cmds...)
}

func (m wizardModel) getListLen() int {
	switch m.step {
	case stepProvider:
		return len(m.providers)
	case stepRegion:
		return len(m.regions)
	case stepSize:
		return len(m.sizes)
	case stepImage:
		return len(m.images)
	case stepSSHKeys:
		return len(m.sshkeys)
	case stepBackups:
		return 2 // Yes/No
	default:
		return 0
	}
}

func (m wizardModel) View() string {
	var b strings.Builder

	// Title and step indicator
	title := titleStyle.Render("NEW SERVER")
	stepInfo := stepStyle.Render(fmt.Sprintf("Step %d of %d", m.step+1, stepCount))
	b.WriteString(title + "  " + stepInfo)
	b.WriteString("\n")
	b.WriteString(subtleStyle.Render(strings.Repeat("─", 60)))
	b.WriteString("\n\n")

	if m.loading {
		b.WriteString(subtleStyle.Render("Loading..."))
		b.WriteString("\n\n")
		b.WriteString(WizardHelp())
		return b.String()
	}

	// Render current step
	switch m.step {
	case stepProvider:
		b.WriteString(m.renderListStep("Select Provider:", m.renderProviders()))
	case stepRegion:
		b.WriteString(m.renderListStep("Select Region:", m.renderRegions()))
	case stepSize:
		b.WriteString(m.renderListStep("Select Size/Plan:", m.renderSizes()))
	case stepImage:
		b.WriteString(m.renderListStep("Select Operating System:", m.renderImages()))
	case stepSSHKeys:
		b.WriteString(m.renderSSHKeysStep())
	case stepBackups:
		b.WriteString(m.renderBackupsStep())
	case stepName:
		b.WriteString(m.renderNameStep())
	case stepConfirm:
		b.WriteString(m.renderConfirmStep())
	}

	b.WriteString("\n\n")
	switch m.step {
	case stepName:
		b.WriteString(helpStyle.Render("enter: confirm • esc: back"))
	case stepSSHKeys:
		b.WriteString(helpStyle.Render("↑/↓: navigate • space: toggle • enter: continue • esc: back"))
	default:
		b.WriteString(WizardHelp())
	}

	return b.String()
}

func (m wizardModel) renderListStep(title string, items string) string {
	return fmt.Sprintf("%s\n\n%s", subtleStyle.Render(title), items)
}

func (m wizardModel) renderProviders() string {
	var b strings.Builder
	visible := m.visibleLines()
	end := min(m.offset+visible, len(m.providers))
	needsScroll := len(m.providers) > visible

	// Always reserve space for scroll indicator (keeps list stable)
	if needsScroll {
		if m.offset > 0 {
			b.WriteString(subtleStyle.Render("  ↑ more above"))
		}
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		p := m.providers[i]
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + selectedStyle.Render(p.Name) + "\n")
		} else {
			b.WriteString("  " + p.Name + "\n")
		}
	}

	// Always reserve space for scroll indicator
	if needsScroll {
		if end < len(m.providers) {
			b.WriteString(subtleStyle.Render("  ↓ more below"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m wizardModel) renderRegions() string {
	var b strings.Builder
	visible := m.visibleLines()
	end := min(m.offset+visible, len(m.regions))
	needsScroll := len(m.regions) > visible

	if needsScroll {
		if m.offset > 0 {
			b.WriteString(subtleStyle.Render("  ↑ more above"))
		}
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		r := m.regions[i]
		line := fmt.Sprintf("%s (%s)", r.Name, r.Origin)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + selectedStyle.Render(line) + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	if needsScroll {
		if end < len(m.regions) {
			b.WriteString(subtleStyle.Render("  ↓ more below"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m wizardModel) renderSizes() string {
	var b strings.Builder
	visible := m.visibleLines()
	end := min(m.offset+visible, len(m.sizes))
	needsScroll := len(m.sizes) > visible

	if needsScroll {
		if m.offset > 0 {
			b.WriteString(subtleStyle.Render("  ↑ more above"))
		}
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		s := m.sizes[i]
		line := fmt.Sprintf("%-20s %s/mo  %s  %sMB RAM  %sGB SSD",
			s.Name, FormatMoney(s.Price), s.Processor, s.Memory, s.Disk)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + selectedStyle.Render(line) + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	if needsScroll {
		if end < len(m.sizes) {
			b.WriteString(subtleStyle.Render("  ↓ more below"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m wizardModel) renderImages() string {
	var b strings.Builder
	visible := m.visibleLines()
	end := min(m.offset+visible, len(m.images))
	needsScroll := len(m.images) > visible

	if needsScroll {
		if m.offset > 0 {
			b.WriteString(subtleStyle.Render("  ↑ more above"))
		}
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		img := m.images[i]
		line := fmt.Sprintf("%s (%s)", img.Name, img.Architecture)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + selectedStyle.Render(line) + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	if needsScroll {
		if end < len(m.images) {
			b.WriteString(subtleStyle.Render("  ↓ more below"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m wizardModel) renderSSHKeysStep() string {
	var b strings.Builder
	b.WriteString(subtleStyle.Render("Select SSH Keys (space to toggle):"))
	b.WriteString("\n\n")

	if len(m.sshkeys) == 0 {
		b.WriteString(failedStyle.Render("No SSH keys found. Add keys at https://dashboard.bithost.io/keys"))
		return b.String()
	}

	visible := m.visibleLines()
	end := min(m.offset+visible, len(m.sshkeys))
	needsScroll := len(m.sshkeys) > visible

	if needsScroll {
		if m.offset > 0 {
			b.WriteString(subtleStyle.Render("  ↑ more above"))
		}
		b.WriteString("\n")
	}

	for i := m.offset; i < end; i++ {
		k := m.sshkeys[i]
		checked := "[ ]"
		if contains(m.selectedKeys, k.ID) {
			checked = activeStyle.Render("[x]")
		}
		line := fmt.Sprintf("%s %s", checked, k.Label)
		if i == m.cursor {
			b.WriteString(selectedStyle.Render("> ") + line + "\n")
		} else {
			b.WriteString("  " + line + "\n")
		}
	}

	if needsScroll {
		if end < len(m.sshkeys) {
			b.WriteString(subtleStyle.Render("  ↓ more below"))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func (m wizardModel) renderBackupsStep() string {
	var b strings.Builder
	b.WriteString(subtleStyle.Render("Enable automatic backups? (+20% cost)"))
	b.WriteString("\n\n")

	yesCheck := "( )"
	noCheck := "( )"
	if m.backupsEnabled {
		yesCheck = activeStyle.Render("(•)")
	} else {
		noCheck = activeStyle.Render("(•)")
	}

	b.WriteString(fmt.Sprintf("  %s Yes\n", yesCheck))
	b.WriteString(fmt.Sprintf("  %s No\n", noCheck))
	b.WriteString("\n")
	b.WriteString(helpStyle.Render("space: toggle • enter: continue"))
	return b.String()
}

func (m wizardModel) renderNameStep() string {
	var b strings.Builder
	b.WriteString(subtleStyle.Render("Enter server name:"))
	b.WriteString("\n\n")
	b.WriteString(m.nameInput.View())
	return b.String()
}

func (m wizardModel) renderConfirmStep() string {
	var b strings.Builder
	b.WriteString(titleStyle.Render("Confirm Server Configuration"))
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Name:") + " " + valueStyle.Render(m.serverName) + "\n")
	b.WriteString(labelStyle.Render("Provider:") + " " + valueStyle.Render(m.selectedProvider.Name) + "\n")
	b.WriteString(labelStyle.Render("Region:") + " " + valueStyle.Render(m.selectedRegion.Name) + "\n")
	b.WriteString(labelStyle.Render("Size:") + " " + valueStyle.Render(m.selectedSize.Name) + "\n")
	b.WriteString(labelStyle.Render("OS:") + " " + valueStyle.Render(m.selectedImage.Name) + "\n")
	b.WriteString(labelStyle.Render("SSH Keys:") + " " + valueStyle.Render(fmt.Sprintf("%d selected", len(m.selectedKeys))) + "\n")
	b.WriteString(labelStyle.Render("Backups:") + " " + FormatEnabled(m.backupsEnabled) + "\n")

	cost := FormatMoney(m.selectedSize.Price)
	if m.backupsEnabled {
		cost += fmt.Sprintf(" + %s (backups)", FormatMoney(m.selectedSize.Price*0.2))
	}
	b.WriteString(labelStyle.Render("Monthly Cost:") + " " + valueStyle.Render(cost) + "\n")

	b.WriteString("\n")
	b.WriteString(helpStyle.Render("enter: create server • esc: go back"))

	return b.String()
}

// Helper functions

func toggleSelection(slice []int, item int) []int {
	for i, v := range slice {
		if v == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return append(slice, item)
}

func contains(slice []int, item int) bool {
	return slices.Contains(slice, item)
}
