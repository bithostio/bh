package api

// ErrorResponse represents the error response from the API
type ErrorResponse struct {
	Errors []string `json:"errors"`
}

// Pagination represents pagination metadata
type Pagination struct {
	Page            int  `json:"page"`
	HasPreviousPage bool `json:"has_previous_page"`
	HasNextPage     bool `json:"has_next_page"`
}

// Meta represents response metadata
type Meta struct {
	Pagination Pagination `json:"pagination"`
}

// UserResponse represents the user information response
type UserResponse struct {
	FullName    string  `json:"full_name"`
	Email       string  `json:"email"`
	Balance     float64 `json:"balance"`
	ServerLimit int     `json:"server_limit"`
}

// Server represents a server resource
type Server struct {
	ID                 int     `json:"id"`
	Name               string  `json:"name"`
	Pending            bool    `json:"pending"`
	Power              bool    `json:"power"`
	CostSoFar          float64 `json:"cost_so_far"`
	Status             string  `json:"status"`
	BackupsEnabled     bool    `json:"backups_enabled"`
	Message            string  `json:"message"`
	IPAddress          string  `json:"ip_address"`
	PrivateIPAddress   string  `json:"private_ip_address"`
	IPAddressV6        string  `json:"ip_address_v6"`
	PrivateIPAddressV6 string  `json:"private_ip_address_v6"`
	ProviderID         int     `json:"provider_id"`
}

// Provider represents a cloud provider
type Provider struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Region represents a geographic region
type Region struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	Slug               string   `json:"slug"`
	Origin             string   `json:"origin"`
	Number             int      `json:"number"`
	AvailableSizeSlugs []string `json:"available_size_slugs"`
}

// Size represents a server size/plan
type Size struct {
	ID           int     `json:"id"`
	Name         string  `json:"name"`
	Slug         string  `json:"slug"`
	Price        float64 `json:"price"`
	PricePerHour string  `json:"price_per_hour"`
	WindowsFee   float64 `json:"windows_fee"`
	Memory       string  `json:"memory"`
	Processor    string  `json:"processor"`
	Bandwidth    string  `json:"bandwidth"`
	Disk         string  `json:"disk"`
	Kind         string  `json:"kind"`
}

// Image represents an OS image
type Image struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Distribution string `json:"distribution"`
	Windows      bool   `json:"windows"`
	Architecture string `json:"architecture"`
}

// SSHKey represents an SSH public key
type SSHKey struct {
	ID    int    `json:"id"`
	Key   string `json:"key"`
	Label string `json:"label"`
}

// CreateServerRequest represents the request to create a server
type CreateServerRequest struct {
	Name           string `json:"name"`
	SizeID         int    `json:"size_id"`
	RegionID       int    `json:"region_id"`
	ImageID        int    `json:"image_id"`
	ProviderID     int    `json:"provider_id"`
	KeyIDs         []int  `json:"key_ids"`
	BackupsEnabled bool   `json:"backups_enabled"`
	Terms          bool   `json:"terms"`
}

// CreateSSHKeyRequest represents the request to create an SSH key
type CreateSSHKeyRequest struct {
	Label      string `json:"label"`
	Key        string `json:"key"`
	ProviderID int    `json:"provider_id,omitempty"`
}

// Response wrappers for API responses
type ServersResponse struct {
	Servers []Server `json:"servers"`
	Meta    *Meta    `json:"meta,omitempty"`
}

type ServerResponse struct {
	Server Server `json:"server"`
}

type ProvidersResponse struct {
	Providers []Provider `json:"providers"`
	Meta      *Meta      `json:"meta,omitempty"`
}

type RegionsResponse struct {
	Regions []Region `json:"regions"`
	Meta    *Meta    `json:"meta,omitempty"`
}

type SizesResponse struct {
	Sizes []Size `json:"sizes"`
	Meta  *Meta  `json:"meta,omitempty"`
}

type ImagesResponse struct {
	Images []Image `json:"images"`
	Meta   *Meta   `json:"meta,omitempty"`
}

type KeysResponse struct {
	Keys []SSHKey `json:"keys"`
	Meta *Meta    `json:"meta,omitempty"`
}

type KeyResponse struct {
	Key SSHKey `json:"key"`
}
