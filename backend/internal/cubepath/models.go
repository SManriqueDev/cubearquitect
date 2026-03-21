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

type VPS struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Status   string `json:"status"`
	Username string `json:"username"`
	Label    string `json:"label"`
	IPv4     string `json:"ipv4"`
	IPv6     string `json:"ipv6"`
}
