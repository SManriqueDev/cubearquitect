package cubepath

type ProjectResponse []ProjectItem

type ProjectItem struct {
	Project    ProjectInfo   `json:"project"`
	Networks   []interface{} `json:"networks"`
	Baremetals []interface{} `json:"baremetals"`
	VPS        []interface{} `json:"vps"`
}

type ProjectInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type VPSCreateRequest struct {
	Name             string   `json:"name"`
	PlanName         string   `json:"plan_name"`
	TemplateName     string   `json:"template_name"`
	LocationName     string   `json:"location_name"`
	Label            string   `json:"label,omitempty"`
	SSHKeyNames      []string `json:"ssh_key_names,omitempty"`
	NetworkID        int      `json:"network_id,omitempty"`
	Password         string   `json:"password,omitempty"`
	IPv4             bool     `json:"ipv4"`
	EnableBackups    bool     `json:"enable_backups"`
	FirewallGroupIDs []int    `json:"firewall_group_ids,omitempty"`
	CustomCloudinit  string   `json:"custom_cloudinit,omitempty"`
}

// FloatingIP represents a floating IP in CubePath response
type FloatingIP struct {
	Address        string `json:"address"`
	Netmask        string `json:"netmask"`
	Type           string `json:"type"`
	ProtectionType string `json:"protection_type"`
	IsPrimary      bool   `json:"is_primary"`
}

// FloatingIPList represents the floating_ips structure in CubePath response
type FloatingIPList struct {
	List         []FloatingIP `json:"list"`
	PricePerHour float64      `json:"price_per_hour"`
}

// VPS matches the actual CubePath /vps/ response structure
type VPS struct {
	ID          int            `json:"id"`
	Label       string         `json:"label"`
	Name        string         `json:"name"`
	Hostname    string         `json:"hostname"`
	Status      string         `json:"status"`
	User        string         `json:"user"`
	Plan        interface{}    `json:"plan"`
	Template    interface{}    `json:"template"`
	FloatingIPs FloatingIPList `json:"floating_ips"`
	SSHKeys     []interface{}  `json:"ssh_keys"`
	Location    interface{}    `json:"location"`
	Network     interface{}    `json:"network"`
	Firewall    []interface{}  `json:"firewall_groups"`

	// Computed fields for convenience (not in JSON)
	IPv4 string `json:"-"`
	IPv6 string `json:"-"`
}

// ExtractIPs extracts IPv4 and IPv6 from FloatingIPs
func (v *VPS) ExtractIPs() {
	for _, ip := range v.FloatingIPs.List {
		if ip.Type == "IPv4" && ip.IsPrimary {
			v.IPv4 = ip.Address
		} else if ip.Type == "IPv6" && ip.IsPrimary {
			v.IPv6 = ip.Address
		}
	}
}

// VPSCreateResponse is the response from CubePath /vps/create endpoint
type VPSCreateResponse struct {
	Detail      string `json:"detail"`
	VPSID       int    `json:"vps_id"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Plan        string `json:"plan"`
	Location    string `json:"location"`
	IPv4Address string `json:"ipv4_address"`
	IPv6Address string `json:"ipv6_address"`
}

// ============================================
// Pricing - New Endpoints /vps/plans & /vps/templates
// ============================================

// PlansResponse from /vps/plans endpoint
type PlansResponse struct {
	Locations []VPSLocationNew `json:"locations"`
}

// VPSLocationNew represents location structure in /vps/plans
type VPSLocationNew struct {
	LocationName string    `json:"location_name"`
	Description  string    `json:"description"`
	Clusters     []Cluster `json:"clusters"`
}

// Cluster represents cluster in /vps/plans
type Cluster struct {
	ClusterName string `json:"cluster_name"`
	Type        string `json:"type"`
	Plans       []Plan `json:"plans"`
}

// Plan represents a VPS plan in /vps/plans
type Plan struct {
	PlanName     string `json:"plan_name"`
	RAM          int    `json:"ram"`
	CPU          int    `json:"cpu"`
	Storage      int    `json:"storage"`
	Bandwidth    int    `json:"bandwidth"`
	PricePerHour string `json:"price_per_hour"`
	Status       int    `json:"status"`
}

// TemplatesResponse from /vps/templates endpoint
type TemplatesResponse struct {
	OperatingSystems []TemplateOS  `json:"operating_systems"`
	Applications     []TemplateApp `json:"applications"`
}

// TemplateOS represents an operating system template
type TemplateOS struct {
	TemplateName string `json:"template_name"`
	OSName       string `json:"os_name"`
	Version      string `json:"version"`
}

// TemplateApp represents an application template
type TemplateApp struct {
	AppName         string `json:"app_name"`
	Version         string `json:"version"`
	RecommendedPlan string `json:"recommended_plan"`
}
