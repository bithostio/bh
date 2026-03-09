package tui

import (
	"fmt"
	"time"

	"github.com/bithostio/bh/internal/api"
	"github.com/bithostio/bh/internal/config"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var Version string

// View represents the current screen
type view int

const (
	viewAuth view = iota
	viewDashboard
	viewServerDetail
	viewWizard
	viewSSHKeys
)

// Model is the root TUI model
type Model struct {
	cfg    *config.Config
	client *api.Client
	user   *api.UserResponse

	// Current view
	view view

	// Sub-models
	auth      authModel
	dashboard dashboardModel
	detail    serverDetailModel
	wizard    wizardModel
	sshkeys   sshkeysModel
	modal     modalModel

	// Terminal dimensions
	width  int
	height int

	// State
	err       error
	status    string
	statusAt  time.Time
	deleting  bool
	quitting  bool
	needsAuth bool
}

// Message types

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type userMsg *api.UserResponse
type serversMsg []api.Server
type providersMsg []api.Provider
type regionsMsg []api.Region
type sizesMsg []api.Size
type imagesMsg []api.Image
type sshkeysMsg []api.SSHKey
type serverCreatedMsg api.Server
type serverDeletedMsg int
type keyCreatedMsg api.SSHKey
type tickMsg time.Time
type authSuccessMsg struct{}
type clearStatusMsg struct{}

// Run starts the TUI application
func Run() error {
	m := Model{
		view:      viewAuth,
		auth:      newAuthModel(),
		dashboard: newDashboardModel(),
		wizard:    newWizardModel(),
		sshkeys:   newSSHKeysModel(),
		needsAuth: true,
	}

	// Try to load existing config
	cfg, err := config.Load()
	if err == nil && cfg.APIKey != "" {
		m.cfg = cfg
		m.client = api.NewClient(cfg.BaseURL, cfg.APIKey)
		m.view = viewDashboard
		m.needsAuth = false
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	_, err = p.Run()
	return err
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	if m.needsAuth {
		return textinput.Blink
	}
	return tea.Batch(
		m.fetchUser(),
		m.fetchServers(),
		tickCmd(),
	)
}

// Update handles messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global quit
		if msg.String() == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}

		// Handle modal first if active
		if m.modal.active {
			var cmd tea.Cmd
			m.modal, cmd = m.modal.Update(msg)
			if m.modal.confirmed {
				switch m.modal.modalType {
				case modalDelete:
					if m.modal.serverID > 0 {
						m.deleting = true
						m.status = "Deleting server..."
						cmds = append(cmds, m.deleteServer(m.modal.serverID))
					}
				case modalTopup:
					_ = openBrowser("https://dashboard.bithost.io")
				}
				m.modal = modalModel{}
			} else if !m.modal.active {
				// Modal was cancelled
				m.modal = modalModel{}
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}
			return m, tea.Batch(cmds...)
		}

		// View-specific key handling
		switch m.view {
		case viewAuth:
			var cmd tea.Cmd
			m.auth, cmd = m.auth.Update(msg)
			if m.auth.loading {
				// User pressed enter with valid input
				cmds = append(cmds, m.authenticate(m.auth.Value()))
			}
			if cmd != nil {
				cmds = append(cmds, cmd)
			}

		case viewDashboard:
			// Clear any existing error when user interacts
			if m.err != nil {
				m.err = nil
			}

			switch msg.String() {
			case "q":
				m.quitting = true
				return m, tea.Quit
			case "n":
				m.view = viewWizard
				m.wizard = newWizardModel()
				return m, m.fetchProviders()
			case "s":
				m.view = viewSSHKeys
				return m, m.fetchSSHKeys()
			case "r":
				m.status = "Refreshing..."
				cmds = append(cmds, m.fetchServers(), m.fetchUser(), clearStatusCmd())
			case "enter":
				if server := m.dashboard.selectedServer(); server != nil {
					m.view = viewServerDetail
					m.detail = newServerDetailModel(*server)
				}
			case "d":
				if server := m.dashboard.selectedServer(); server != nil {
					m.modal = newModalModel(
						fmt.Sprintf("Delete server '%s'?", server.Name),
						server.ID,
					)
				}
			case "t":
				m.modal = newTopupModal()
			default:
				var cmd tea.Cmd
				m.dashboard, cmd = m.dashboard.Update(msg)
				cmds = append(cmds, cmd)
			}

		case viewServerDetail:
			switch msg.String() {
			case "esc", "q":
				m.view = viewDashboard
			case "d":
				m.modal = newModalModel(
					fmt.Sprintf("Delete server '%s'?", m.detail.server.Name),
					m.detail.server.ID,
				)
			case "c":
				if m.detail.server.IPAddress != "" {
					if err := m.detail.CopySSHCommand(); err == nil {
						m.status = "SSH command copied to clipboard"
						m.statusAt = time.Now()
					}
				}
			}

		case viewWizard:
			switch msg.String() {
			case "esc", "q":
				if m.wizard.step == 0 {
					m.view = viewDashboard
				} else {
					m.wizard = m.wizard.prevStep()
				}
			default:
				var cmd tea.Cmd
				m.wizard, cmd = m.wizard.Update(msg, m)
				cmds = append(cmds, cmd)

				// Check if wizard completed
				if m.wizard.completed {
					m.status = "Creating server..."
				}
			}

		case viewSSHKeys:
			switch msg.String() {
			case "esc", "q":
				if m.sshkeys.adding {
					m.sshkeys = m.sshkeys.CancelAdding()
				} else {
					m.view = viewDashboard
				}
			case "a":
				if !m.sshkeys.adding {
					m.sshkeys = m.sshkeys.StartAdding()
					return m, textinput.Blink
				}
			case "enter":
				if m.sshkeys.adding {
					label, key, ok := m.sshkeys.GetKeyToAdd()
					if ok {
						m.status = "Creating SSH key..."
						cmds = append(cmds, m.createSSHKey(label, key))
					}
				}
				// Also pass to sshkeys for step navigation
				var cmd tea.Cmd
				m.sshkeys, cmd = m.sshkeys.Update(msg)
				cmds = append(cmds, cmd)
			default:
				var cmd tea.Cmd
				m.sshkeys, cmd = m.sshkeys.Update(msg)
				cmds = append(cmds, cmd)
			}
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.auth = m.auth.SetSize(msg.Width, msg.Height)
		m.dashboard = m.dashboard.SetSize(msg.Width, msg.Height)
		m.wizard = m.wizard.SetSize(msg.Width, msg.Height)
		m.sshkeys = m.sshkeys.SetSize(msg.Width, msg.Height)

	case tickMsg:
		// Clear old status messages
		if m.status != "" && time.Since(m.statusAt) > 3*time.Second {
			m.status = ""
		}

		// Auto-refresh dashboard data when on dashboard view
		if m.view == viewDashboard && !m.modal.active && !m.needsAuth {
			cmds = append(cmds, m.fetchServers(), m.fetchUser())
		}
		cmds = append(cmds, tickCmd())

	case clearStatusMsg:
		// Will be cleared on next tick if older than 3s
		// This just triggers a refresh

	case authSuccessMsg:
		m.needsAuth = false
		m.view = viewDashboard
		m.auth = authModel{} // Clear auth model
		cmds = append(cmds, m.fetchUser(), m.fetchServers(), tickCmd())

	case userMsg:
		m.user = msg
		m.dashboard = m.dashboard.SetUser(msg)
		m.err = nil

	case serversMsg:
		m.dashboard = m.dashboard.SetServers(msg)
		if m.status == "Refreshing..." {
			m.status = ""
		}

	case providersMsg:
		m.wizard = m.wizard.SetProviders(msg)

	case regionsMsg:
		m.wizard = m.wizard.SetRegions(msg)

	case sizesMsg:
		m.wizard = m.wizard.SetSizes(msg)

	case imagesMsg:
		m.wizard = m.wizard.SetImages(msg)

	case sshkeysMsg:
		m.wizard = m.wizard.SetSSHKeys(msg)
		m.sshkeys = m.sshkeys.SetKeys(msg)

	case serverCreatedMsg:
		m.view = viewDashboard
		m.wizard = newWizardModel()
		m.status = fmt.Sprintf("Server '%s' created!", msg.Name)
		m.statusAt = time.Now()
		cmds = append(cmds, m.fetchServers())

	case serverDeletedMsg:
		m.modal = modalModel{}
		m.deleting = false
		m.view = viewDashboard
		m.status = "Server deleted successfully"
		m.statusAt = time.Now()
		cmds = append(cmds, m.fetchServers())

	case keyCreatedMsg:
		m.sshkeys = m.sshkeys.CancelAdding()
		m.status = fmt.Sprintf("SSH key '%s' created!", msg.Label)
		m.statusAt = time.Now()
		cmds = append(cmds, m.fetchSSHKeys())

	case errMsg:
		m.err = msg.err
		m.deleting = false
		m.status = ""

		// Handle auth errors
		if m.view == viewAuth {
			m.auth = m.auth.SetError(msg.err.Error())
			m.auth = m.auth.SetLoading(false)
		}

		// Handle SSH key errors
		if m.view == viewSSHKeys && m.sshkeys.adding {
			m.sshkeys = m.sshkeys.SetAddError(msg.err.Error())
		}

		// Handle wizard errors
		if m.view == viewWizard {
			m.wizard = m.wizard.SetLoading(false)
		}
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m Model) View() string {
	if m.quitting {
		return ""
	}

	var content string
	switch m.view {
	case viewAuth:
		content = m.auth.View()
	case viewDashboard:
		content = m.dashboard.View()
		// Add status/error bar
		content += m.renderStatusBar()
	case viewServerDetail:
		content = m.detail.View()
		content += m.renderStatusBar()
	case viewWizard:
		content = m.wizard.View()
		content += m.renderStatusBar()
	case viewSSHKeys:
		content = m.sshkeys.View()
		content += m.renderStatusBar()
	}

	// Overlay modal if active
	if m.modal.active {
		content = m.modal.View(content, m.width, m.height)
	}

	return content
}

func (m Model) renderStatusBar() string {
	if m.err != nil {
		return "\n" + failedStyle.Render("Error: "+m.err.Error())
	}
	if m.status != "" {
		if m.deleting {
			return "\n" + pendingStyle.Render(m.status)
		}
		return "\n" + activeStyle.Render("✓ "+m.status)
	}
	return ""
}

// Commands

func tickCmd() tea.Cmd {
	return tea.Tick(10*time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func clearStatusCmd() tea.Cmd {
	return tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
		return clearStatusMsg{}
	})
}

func (m *Model) authenticate(apiKey string) tea.Cmd {
	return func() tea.Msg {
		// Create config
		cfg := config.New()
		cfg.APIKey = apiKey

		// Test the API key by fetching user
		client := api.NewClient(cfg.BaseURL, cfg.APIKey)
		_, err := client.GetUser()
		if err != nil {
			return errMsg{err}
		}

		// Save config
		if err := config.Save(cfg); err != nil {
			return errMsg{err}
		}

		// Update model's client and config
		m.cfg = cfg
		m.client = client

		return authSuccessMsg{}
	}
}

func (m Model) fetchUser() tea.Cmd {
	return func() tea.Msg {
		user, err := m.client.GetUser()
		if err != nil {
			return errMsg{err}
		}
		return userMsg(user)
	}
}

func (m Model) fetchServers() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListServers(0)
		if err != nil {
			return errMsg{err}
		}
		return serversMsg(resp.Servers)
	}
}

func (m Model) fetchProviders() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListProviders(0)
		if err != nil {
			return errMsg{err}
		}
		return providersMsg(resp.Providers)
	}
}

func (m Model) fetchRegions(provider string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListRegions(provider, 0)
		if err != nil {
			return errMsg{err}
		}
		return regionsMsg(resp.Regions)
	}
}

func (m Model) fetchSizes(regionID int, provider string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListSizes(regionID, provider, 0)
		if err != nil {
			return errMsg{err}
		}
		return sizesMsg(resp.Sizes)
	}
}

func (m Model) fetchImages(provider string) tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListImages(provider, "x86", 0)
		if err != nil {
			return errMsg{err}
		}
		return imagesMsg(resp.Images)
	}
}

func (m Model) fetchSSHKeys() tea.Cmd {
	return func() tea.Msg {
		resp, err := m.client.ListSSHKeys(0)
		if err != nil {
			return errMsg{err}
		}
		return sshkeysMsg(resp.Keys)
	}
}

func (m Model) createSSHKey(label, key string) tea.Cmd {
	return func() tea.Msg {
		req := &api.CreateSSHKeyRequest{
			Label: label,
			Key:   key,
		}
		created, err := m.client.CreateSSHKey(req)
		if err != nil {
			return errMsg{err}
		}
		return keyCreatedMsg(*created)
	}
}

func (m Model) deleteServer(id int) tea.Cmd {
	return func() tea.Msg {
		err := m.client.DeleteServer(id)
		if err != nil {
			return errMsg{err}
		}
		return serverDeletedMsg(id)
	}
}

func (m Model) createServer(req *api.CreateServerRequest) tea.Cmd {
	return func() tea.Msg {
		server, err := m.client.CreateServer(req)
		if err != nil {
			return errMsg{err}
		}
		return serverCreatedMsg(*server)
	}
}
