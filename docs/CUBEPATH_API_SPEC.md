# CubePath API - Especificación Técnica para CubeArchitect

> **Documento extraído del código fuente de cubecli**  
> Última actualización: Hackathon CubeCLI

---

## Índice

1. [Configuración Base](#configuración-base)
2. [API Spec - Endpoints Completos](#sección-1-api-spec)
3. [Modelos Go (Structs)](#sección-2-modelos-go-structs)
4. [Ejemplos de Petición](#sección-3-ejemplos-de-petición)
5. [Estados de Recursos](#estados-de-recursos)

---

## Configuración Base

### Base URL
```
https://api.cubepath.com
```

### Headers Obligatorios

| Header | Valor | Descripción |
|--------|-------|-------------|
| `Authorization` | `Bearer <CUBE_API_TOKEN>` | Token de autenticación |
| `Content-Type` | `application/json` | Formato de cuerpo de petición |
| `User-Agent` | `CubeCLI/<version>` | Identificación del cliente |

### Variables de Entorno

| Variable | Prioridad | Descripción |
|----------|-----------|-------------|
| `CUBE_API_TOKEN` | 1 (más alta) | Token de API directamente |
| `CUBE_API_URL` | 1 (más alta) | URL base personalizada |
| `~/.cubecli/config.json` | 2 | Archivo de configuración local |

### Timeout por Defecto
- **30 segundos** para todas las peticiones HTTP

---

## Sección 1: API Spec

### VPS (Virtual Private Servers)

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `vps create` | POST | `/vps/create/{project_id}` | Crear VPS |
| `vps list` | GET | `/projects/` | Listar todos los VPS |
| `vps show` | GET | `/projects/` | Ver detalles de VPS |
| `vps destroy` | POST | `/vps/destroy/{vps_id}` | Destruir VPS |
| `vps update` | PATCH | `/vps/update/{vps_id}` | Actualizar VPS |
| `vps resize` | POST | `/vps/resize/vps_id/{vps_id}/resize_plan/{plan_name}` | Redimensionar VPS |
| `vps change-password` | POST | `/vps/{vps_id}/change-password` | Cambiar contraseña |
| `vps reinstall` | POST | `/vps/reinstall/{vps_id}` | Reinstalar SO |
| `vps power start` | POST | `/vps/{vps_id}/power/start_vps` | Iniciar VPS |
| `vps power stop` | POST | `/vps/{vps_id}/power/stop_vps` | Detener VPS |
| `vps power restart` | POST | `/vps/{vps_id}/power/restart_vps` | Reiniciar VPS |
| `vps power reset` | POST | `/vps/{vps_id}/power/reset_vps` | Resetear VPS |
| `vps template list` | GET | `/pricing` | Listar templates |
| `vps plan list` | GET | `/pricing` | Listar planes |
| `vps backup list` | GET | `/vps/{vps_id}/backups` | Listar backups |
| `vps backup create` | POST | `/vps/{vps_id}/backups` | Crear backup |
| `vps backup restore` | POST | `/vps/{vps_id}/backups/{backup_id}/restore` | Restaurar backup |
| `vps backup delete` | DELETE | `/vps/{vps_id}/backups/{backup_id}` | Eliminar backup |
| `vps backup settings` | GET | `/vps/{vps_id}/backup/settings` | Ver settings de backup |
| `vps backup configure` | PUT | `/vps/{vps_id}/backup/settings` | Configurar backups |

### Load Balancer

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `lb list` | GET | `/loadbalancer/` | Listar LBs |
| `lb show` | GET | `/loadbalancer/` | Ver detalle LB |
| `lb create` | POST | `/loadbalancer/` | Crear LB |
| `lb update` | PATCH | `/loadbalancer/{lb_uuid}` | Actualizar LB |
| `lb delete` | DELETE | `/loadbalancer/{lb_uuid}` | Eliminar LB |
| `lb resize` | POST | `/loadbalancer/{lb_uuid}/resize` | Redimensionar LB |
| `lb listener create` | POST | `/loadbalancer/{lb_uuid}/listeners` | Crear listener |
| `lb listener update` | PATCH | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}` | Actualizar listener |
| `lb listener delete` | DELETE | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}` | Eliminar listener |
| `lb target add` | POST | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/targets` | Añadir target |
| `lb target update` | PATCH | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/targets/{target_uuid}` | Actualizar target |
| `lb target remove` | DELETE | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/targets/{target_uuid}` | Remover target |
| `lb target drain` | POST | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/targets/{target_uuid}/drain` | Drenar target |
| `lb health-check configure` | PUT | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/health-check` | Configurar health check |
| `lb health-check delete` | DELETE | `/loadbalancer/{lb_uuid}/listeners/{listener_uuid}/health-check` | Eliminar health check |
| `lb plan list` | GET | `/loadbalancer/plans` | Listar planes LB |

### CDN (Content Delivery Network)

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `cdn zone list` | GET | `/cdn/zones` | Listar zonas |
| `cdn zone show` | GET | `/cdn/zones/{zone_uuid}` | Ver zona |
| `cdn zone create` | POST | `/cdn/zones` | Crear zona |
| `cdn zone update` | PATCH | `/cdn/zones/{zone_uuid}` | Actualizar zona |
| `cdn zone delete` | DELETE | `/cdn/zones/{zone_uuid}` | Eliminar zona |
| `cdn zone pricing` | GET | `/cdn/zones/{zone_uuid}/pricing` | Ver pricing |
| `cdn origin list` | GET | `/cdn/zones/{zone_uuid}/origins` | Listar origins |
| `cdn origin create` | POST | `/cdn/zones/{zone_uuid}/origins` | Crear origin |
| `cdn origin update` | PATCH | `/cdn/zones/{zone_uuid}/origins/{origin_uuid}` | Actualizar origin |
| `cdn origin delete` | DELETE | `/cdn/zones/{zone_uuid}/origins/{origin_uuid}` | Eliminar origin |
| `cdn rule list` | GET | `/cdn/zones/{zone_uuid}/rules` | Listar reglas |
| `cdn rule show` | GET | `/cdn/zones/{zone_uuid}/rules/{rule_uuid}` | Ver regla |
| `cdn rule create` | POST | `/cdn/zones/{zone_uuid}/rules` | Crear regla |
| `cdn rule update` | PATCH | `/cdn/zones/{zone_uuid}/rules/{rule_uuid}` | Actualizar regla |
| `cdn rule delete` | DELETE | `/cdn/zones/{zone_uuid}/rules/{rule_uuid}` | Eliminar regla |
| `cdn waf list` | GET | `/cdn/zones/{zone_uuid}/waf-rules` | Listar reglas WAF |
| `cdn waf create` | POST | `/cdn/zones/{zone_uuid}/waf-rules` | Crear regla WAF |
| `cdn waf update` | PATCH | `/cdn/zones/{zone_uuid}/waf-rules/{rule_uuid}` | Actualizar WAF |
| `cdn waf delete` | DELETE | `/cdn/zones/{zone_uuid}/waf-rules/{rule_uuid}` | Eliminar WAF |
| `cdn plan list` | GET | `/cdn/plans` | Listar planes CDN |
| `cdn metrics summary` | GET | `/cdn/zones/{zone_uuid}/metrics/summary?minutes={n}` | Métricas resumidas |
| `cdn metrics requests` | GET | `/cdn/zones/{zone_uuid}/metrics/requests` | Métricas de requests |
| `cdn metrics bandwidth` | GET | `/cdn/zones/{zone_uuid}/metrics/bandwidth` | Métricas de bandwidth |
| `cdn metrics cache` | GET | `/cdn/zones/{zone_uuid}/metrics/cache` | Métricas de cache |
| `cdn metrics top-urls` | GET | `/cdn/zones/{zone_uuid}/metrics/top-urls` | Top URLs |
| `cdn metrics top-countries` | GET | `/cdn/zones/{zone_uuid}/metrics/top-countries` | Top países |
| `cdn metrics top-asn` | GET | `/cdn/zones/{zone_uuid}/metrics/top-asn` | Top ASNs |

### DNS

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `dns zone list` | GET | `/dns/zones` | Listar zonas |
| `dns zone list` | GET | `/dns/zones?project_id={id}` | Listar por proyecto |
| `dns zone show` | GET | `/dns/zones/{zone_uuid}` | Ver zona |
| `dns zone create` | POST | `/dns/zones` | Crear zona |
| `dns zone delete` | DELETE | `/dns/zones/{zone_uuid}` | Eliminar zona |
| `dns zone verify` | POST | `/dns/zones/{zone_uuid}/verify` | Verificar zona |
| `dns zone scan` | POST | `/dns/zones/{zone_uuid}/scan` | Escanear zona |
| `dns record list` | GET | `/dns/zones/{zone_uuid}/records` | Listar registros |
| `dns record list` | GET | `/dns/zones/{zone_uuid}/records?record_type={type}` | Filtrar por tipo |
| `dns record create` | POST | `/dns/zones/{zone_uuid}/records` | Crear registro |
| `dns record update` | PUT | `/dns/zones/{zone_uuid}/records/{record_uuid}` | Actualizar registro |
| `dns record delete` | DELETE | `/dns/zones/{zone_uuid}/records/{record_uuid}` | Eliminar registro |

### Baremetal

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `baremetal deploy` | POST | `/baremetal/deploy/{project_id}` | Desplegar servidor |
| `baremetal list` | GET | `/projects/` | Listar servidores |
| `baremetal show` | GET | `/projects/` | Ver detalles |
| `baremetal sensors` | GET | `/baremetal/{id}/bmc-sensors` | Sensores BMC |
| `baremetal rescue` | POST | `/baremetal/{id}/rescue` | Modo rescue |
| `baremetal reset-bmc` | POST | `/baremetal/{id}/reset-bmc` | Resetear BMC |
| `baremetal update` | PATCH | `/baremetal/update/{id}` | Actualizar servidor |
| `baremetal ipmi` | POST | `/ipmi-proxy/create-session/{id}` | Crear sesión IPMI |
| `baremetal power start` | POST | `/baremetal/{id}/power/start_metal` | Encender |
| `baremetal power stop` | POST | `/baremetal/{id}/power/stop_metal` | Apagar |
| `baremetal power restart` | POST | `/baremetal/{id}/power/restart_metal` | Reiniciar |
| `baremetal reinstall start` | POST | `/baremetal/{id}/reinstall` | Reinstalar |
| `baremetal reinstall status` | GET | `/baremetal/{id}/reinstall/status` | Estado reinstalación |
| `baremetal monitoring enable` | PUT | `/baremetal/{id}/monitoring?enable=true` | Habilitar monitoreo |
| `baremetal monitoring disable` | PUT | `/baremetal/{id}/monitoring?enable=false` | Deshabilitar monitoreo |
| `baremetal model list` | GET | `/pricing` | Listar modelos |

### Network

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `network create` | POST | `/networks/create_network` | Crear red |
| `network list` | GET | `/projects/` | Listar redes |
| `network update` | PUT | `/networks/{network_id}` | Actualizar red |
| `network delete` | DELETE | `/networks/{network_id}` | Eliminar red |

### Projects

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `project create` | POST | `/projects/` | Crear proyecto |
| `project list` | GET | `/projects/` | Listar proyectos |
| `project show` | GET | `/projects/` | Ver proyecto |
| `project delete` | DELETE | `/projects/{project_id}` | Eliminar proyecto |

### SSH Keys

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `ssh-key create` | POST | `/sshkey/create` | Crear clave SSH |
| `ssh-key list` | GET | `/sshkey/user/sshkeys` | Listar claves |
| `ssh-key delete` | DELETE | `/sshkey/{key_id}` | Eliminar clave |

### Floating IPs

| Comando CLI | Método HTTP | Endpoint | Descripción |
|------------|-------------|----------|-------------|
| `floating-ip list` | GET | `/floating_ips/organization` | Listar IPs |
| `floating-ip acquire` | POST | `/floating_ips/acquire?ip_type={type}&location_name={location}` | Adquirir IP |
| `floating-ip release` | POST | `/floating_ips/release/{address}` | Liberar IP |
| `floating-ip assign` | POST | `/floating_ips/assign/vps/{id}?address={addr}` | Asignar a VPS |
| `floating-ip assign` | POST | `/floating_ips/assign/baremetal/{id}?address={addr}` | Asignar a baremetal |
| `floating-ip unassign` | POST | `/floating_ips/unassign/{address}` | Desasignar IP |
| `floating-ip reverse-dns` | POST | `/floating_ips/reverse_dns/configure?ip={ip}&reverse_dns={hostname}` | Configurar rDNS |

### Otros Endpoints

| Recurso | Método HTTP | Endpoint | Descripción |
|---------|-------------|----------|-------------|
| Locations | GET | `/pricing` | Listar ubicaciones |
| DDoS Attacks | GET | `/ddos-attacks/attacks` | Listar ataques DDoS |

---

## Sección 2: Modelos Go (Structs)

```go
// ============================================
// CONFIGURACIÓN
// ============================================

// Config representa la configuración del cliente
type Config struct {
    APIToken string `json:"api_token"`
}

// APIError representa un error de la API
type APIError struct {
    StatusCode int
    Detail     string
}

// PricingResponse representa la respuesta de pricing
type PricingResponse struct {
    VPS       VPSPricingSection       `json:"vps"`
    Baremetal BaremetalPricingSection `json:"baremetal"`
}

type VPSPricingSection struct {
    Locations []VPSLocation `json:"locations"`
}

type VPSLocation struct {
    LocationName string        `json:"location_name"`
    Description  string        `json:"description"`
    Clusters     []VPSCluster  `json:"clusters"`
    Templates    []VPSTemplate `json:"templates"`
}

type VPSCluster struct {
    Plans []VPSPlan `json:"plans"`
}

type VPSPlan struct {
    Name         string `json:"plan_name"`
    CPU          int    `json:"cpu"`
    RAM          int    `json:"ram"`
    Storage      int    `json:"storage"`
    Bandwidth    int    `json:"bandwidth"`
    PricePerHour string `json:"price_per_hour"`
}

type VPSTemplate struct {
    Name    string `json:"template_name"`
    OSName  string `json:"os_name"`
    Version string `json:"version"`
}

type BaremetalPricingSection struct {
    Locations []BaremetalLocation `json:"locations"`
}

type BaremetalLocation struct {
    LocationName      string              `json:"location_name"`
    Description      string              `json:"description"`
    BaremetalModels  []BaremetalModel   `json:"baremetal_models"`
}

type BaremetalModel struct {
    ModelName       string  `json:"model_name"`
    CPU             string  `json:"cpu"`
    CPUSpecs        string  `json:"cpu_specs"`
    CPUBench         float64 `json:"cpu_bench"`
    RAMSize          int     `json:"ram_size"`
    RAMType          string  `json:"ram_type"`
    DiskSize         string  `json:"disk_size"`
    DiskType         string  `json:"disk_type"`
    Port             int     `json:"port"`
    KVM              string  `json:"kvm"`
    Price            float64 `json:"price"`
    Setup            float64 `json:"setup"`
    StockAvailable   int     `json:"stock_available"`
}

// ============================================
// VPS
// ============================================

// VPSCreateRequest para crear un VPS
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

// VPSUpdateRequest para actualizar un VPS
type VPSUpdateRequest struct {
    Name  string `json:"name,omitempty"`
    Label string `json:"label,omitempty"`
}

// VPSResizeRequest para redimensionar un VPS
type VPSResizeRequest struct {
    PlanName string `json:"plan_name"`
}

// VPSBackupCreateRequest para crear un backup
type VPSBackupCreateRequest struct {
    Notes string `json:"notes,omitempty"`
}

// VPSBackupSettings para configurar backups automáticos
type VPSBackupSettings struct {
    Enabled       bool `json:"enabled"`
    ScheduleHour  int  `json:"schedule_hour"`
    RetentionDays int  `json:"retention_days"`
    MaxBackups    int  `json:"max_backups"`
}

// VPSBackup representa un backup
type VPSBackup struct {
    ID         int     `json:"id"`
    Type       string  `json:"backup_type"`
    Status     string  `json:"status"`
    Progress   int     `json:"progress"`
    SizeGB     float64 `json:"size_gb"`
    Notes      string  `json:"notes"`
    CreatedAt  string  `json:"created_at"`
}

// VPS representa la respuesta de un VPS
type VPS struct {
    ID         int              `json:"id"`
    Name       string           `json:"name"`
    Status     string           `json:"status"`
    Username   string           `json:"username"`
    Label      string           `json:"label"`
    SSHKeys    []SSHKeyInfo     `json:"ssh_keys"`
    FloatingIPs FloatingIPList   `json:"floating_ips"`
    Plan       VPSPlanInfo      `json:"plan"`
    Template   VPSTemplateInfo  `json:"template"`
    Location   VPSLocationInfo   `json:"location"`
    IPv4       string           `json:"ipv4"`
    IPv6       string           `json:"ipv6"`
    Network    VPSNetworkInfo   `json:"network"`
}

type SSHKeyInfo struct {
    Name string `json:"name"`
}

type FloatingIPList struct {
    List []FloatingIPInfo `json:"list"`
}

type FloatingIPInfo struct {
    Address string `json:"address"`
}

type VPSPlanInfo struct {
    Name      string `json:"plan_name"`
    VCPUs     int    `json:"cpu"`
    RAM       int    `json:"ram"`
    Storage   int    `json:"storage"`
    Bandwidth int    `json:"bandwidth"`
}

type VPSTemplateInfo struct {
    Name   string `json:"template_name"`
    OSName string `json:"os_name"`
}

type VPSLocationInfo struct {
    Name string `json:"location_name"`
}

type VPSNetworkInfo struct {
    Name string `json:"name"`
}

// ============================================
// LOAD BALANCER
// ============================================

// LoadBalancerCreateRequest para crear un LB
type LoadBalancerCreateRequest struct {
    Name         string `json:"name"`
    PlanName     string `json:"plan_name"`
    LocationName string `json:"location_name"`
    ProjectID    int    `json:"project_id,omitempty"`
    Label        string `json:"label,omitempty"`
}

// LoadBalancerUpdateRequest para actualizar un LB
type LoadBalancerUpdateRequest struct {
    Name  string `json:"name,omitempty"`
    Label string `json:"label,omitempty"`
}

// ListenerCreateRequest para crear un listener
type ListenerCreateRequest struct {
    Name           string `json:"name"`
    Protocol       string `json:"protocol"`
    SourcePort     int    `json:"source_port"`
    TargetPort     int    `json:"target_port"`
    Algorithm      string `json:"algorithm"`
    StickySessions bool   `json:"sticky_sessions"`
}

// ListenerUpdateRequest para actualizar un listener
type ListenerUpdateRequest struct {
    Name       string `json:"name,omitempty"`
    TargetPort int    `json:"target_port,omitempty"`
    Algorithm  string `json:"algorithm,omitempty"`
    Enabled    *bool  `json:"enabled,omitempty"`
}

// TargetAddRequest para añadir un target
type TargetAddRequest struct {
    TargetType string `json:"target_type"`
    TargetUUID string `json:"target_uuid"`
    Port       int    `json:"port,omitempty"`
    Weight     int    `json:"weight"`
}

// TargetUpdateRequest para actualizar un target
type TargetUpdateRequest struct {
    Port    int    `json:"port,omitempty"`
    Weight  int    `json:"weight,omitempty"`
    Enabled *bool  `json:"enabled,omitempty"`
}

// HealthCheckConfigureRequest para configurar health check
type HealthCheckConfigureRequest struct {
    Protocol           string `json:"protocol"`
    Path               string `json:"path"`
    IntervalSeconds    int    `json:"interval_seconds"`
    TimeoutSeconds     int    `json:"timeout_seconds"`
    HealthyThreshold   int    `json:"healthy_threshold"`
    UnhealthyThreshold int    `json:"unhealthy_threshold"`
    ExpectedCodes      string `json:"expected_codes"`
}

// LoadBalancer representa la respuesta de un LB
type LoadBalancer struct {
    UUID               string     `json:"uuid"`
    Name               string     `json:"name"`
    Status             string     `json:"status"`
    PlanName           string     `json:"plan_name"`
    FloatingIPAddress  string     `json:"floating_ip_address"`
    LocationName       string     `json:"location_name"`
    Label              string     `json:"label,omitempty"`
    Listeners          []Listener `json:"listeners"`
}

// Listener representa un listener
type Listener struct {
    UUID            string   `json:"uuid"`
    Name            string   `json:"name"`
    Protocol        string   `json:"protocol"`
    SourcePort      int      `json:"source_port"`
    TargetPort      int      `json:"target_port"`
    Algorithm       string   `json:"algorithm"`
    Enabled         bool     `json:"enabled"`
    StickySessions  bool     `json:"sticky_sessions"`
    Targets         []Target `json:"targets"`
}

// Target representa un target del LB
type Target struct {
    UUID        string `json:"uuid"`
    TargetType  string `json:"target_type"`
    TargetUUID  string `json:"target_uuid"`
    Port        int    `json:"port"`
    Weight      int    `json:"weight"`
    Status      string `json:"status"`
    Enabled     bool   `json:"enabled"`
}

// LBPlan representa un plan de LB
type LBPlan struct {
    Name                  string  `json:"name"`
    PricePerHour         float64 `json:"price_per_hour"`
    PricePerMonth        float64 `json:"price_per_month"`
    MaxListeners         int     `json:"max_listeners"`
    MaxTargets           int     `json:"max_targets"`
    ConnectionsPerSecond  int     `json:"connections_per_second"`
}

// ============================================
// CDN
// ============================================

// CDNZoneCreateRequest para crear una zona CDN
type CDNZoneCreateRequest struct {
    Name         string `json:"name"`
    PlanName     string `json:"plan_name"`
    CustomDomain string `json:"custom_domain,omitempty"`
    ProjectID    int    `json:"project_id,omitempty"`
}

// CDNZoneUpdateRequest para actualizar una zona CDN
type CDNZoneUpdateRequest struct {
    Name         string `json:"name,omitempty"`
    CustomDomain string `json:"custom_domain,omitempty"`
    SSLType      string `json:"ssl_type,omitempty"`
    Certificate  string `json:"certificate,omitempty"`
}

// CDNOriginCreateRequest para crear un origin
type CDNOriginCreateRequest struct {
    Name               string `json:"name"`
    OriginURL          string `json:"origin_url,omitempty"`
    Address            string `json:"address,omitempty"`
    Port               int    `json:"port,omitempty"`
    Protocol           string `json:"protocol,omitempty"`
    Weight             int    `json:"weight"`
    Priority           int    `json:"priority"`
    IsBackup           bool   `json:"is_backup"`
    HealthCheckEnabled bool   `json:"health_check_enabled"`
    HealthCheckPath    string `json:"health_check_path"`
    VerifySSL          bool   `json:"verify_ssl"`
    HostHeader         string `json:"host_header,omitempty"`
    BasePath           string `json:"base_path,omitempty"`
    Enabled            bool   `json:"enabled"`
}

// CDNOrigin representa un origin
type CDNOrigin struct {
    UUID     string `json:"uuid"`
    Name     string `json:"name"`
    Address  string `json:"address"`
    Port     int    `json:"port"`
    Protocol string `json:"protocol"`
    Weight   int    `json:"weight"`
    Priority int    `json:"priority"`
    IsBackup bool   `json:"is_backup"`
    Health   string `json:"health"`
    Enabled  bool   `json:"enabled"`
}

// CDNRuleCreateRequest para crear una regla
type CDNRuleCreateRequest struct {
    Name             string                 `json:"name"`
    RuleType         string                 `json:"rule_type"`
    Priority         int                    `json:"priority"`
    ActionConfig     map[string]interface{} `json:"action_config"`
    MatchConditions  map[string]interface{} `json:"match_conditions,omitempty"`
    Enabled          bool                   `json:"enabled"`
}

// CDNRule representa una regla
type CDNRule struct {
    UUID             string                 `json:"uuid"`
    Name             string                 `json:"name"`
    Type             string                 `json:"rule_type"`
    Priority         int                    `json:"priority"`
    Enabled          bool                   `json:"enabled"`
    ActionConfig     map[string]interface{} `json:"action_config"`
    MatchConditions  map[string]interface{} `json:"match_conditions"`
}

// CDNZone representa una zona CDN
type CDNZone struct {
    UUID         string `json:"uuid"`
    Name         string `json:"name"`
    Domain       string `json:"domain"`
    CustomDomain string `json:"custom_domain"`
    Status       string `json:"status"`
    SSLType      string `json:"ssl_type"`
    Certificate  string `json:"certificate"`
    CreatedAt    string `json:"created_at"`
    UpdatedAt    string `json:"updated_at"`
}

// CDNPlan representa un plan CDN
type CDNPlan struct {
    Name        string  `json:"name"`
    BasePriceHr float64 `json:"base_price_hr"`
    MaxZones    int     `json:"max_zones"`
    MaxOrigins  int     `json:"max_origins"`
    CustomSSL   bool    `json:"custom_ssl"`
}

// ============================================
// DNS
// ============================================

// DNSZoneCreateRequest para crear una zona DNS
type DNSZoneCreateRequest struct {
    Domain    string `json:"domain"`
    ProjectID int    `json:"project_id"`
}

// DNSZone representa una zona DNS
type DNSZone struct {
    UUID        string   `json:"uuid"`
    Domain      string   `json:"domain"`
    Status      string   `json:"status"`
    Records     int      `json:"records"`
    Nameservers []string `json:"nameservers"`
}

// DNSRecordCreateRequest para crear un registro
type DNSRecordCreateRequest struct {
    Name       string `json:"name"`
    RecordType string `json:"record_type"`
    Content    string `json:"content"`
    TTL        int    `json:"ttl"`
    Priority   int    `json:"priority,omitempty"`
    Weight     int    `json:"weight,omitempty"`
    Port       int    `json:"port,omitempty"`
    Comment    string `json:"comment,omitempty"`
}

// DNSRecordUpdateRequest para actualizar un registro
type DNSRecordUpdateRequest struct {
    Content  string `json:"content,omitempty"`
    TTL      int    `json:"ttl,omitempty"`
    Priority int    `json:"priority,omitempty"`
    Weight   int    `json:"weight,omitempty"`
    Port     int    `json:"port,omitempty"`
    Comment  string `json:"comment,omitempty"`
}

// DNSRecord representa un registro DNS
type DNSRecord struct {
    UUID      string `json:"uuid"`
    Name      string `json:"name"`
    Type      string `json:"record_type"`
    Content   string `json:"content"`
    TTL       int    `json:"ttl"`
    Priority  *int   `json:"priority,omitempty"`
}

// ============================================
// BARE METAL
// ============================================

// BaremetalDeployRequest para desplegar un servidor
type BaremetalDeployRequest struct {
    LocationName   string   `json:"location_name"`
    ModelName      string   `json:"model_name"`
    Hostname       string   `json:"hostname"`
    User           string   `json:"user"`
    Password       string   `json:"password"`
    Label          string   `json:"label,omitempty"`
    OSName         string   `json:"os_name,omitempty"`
    DiskLayoutName string   `json:"disk_layout_name,omitempty"`
    SSHKeyNames    []string `json:"ssh_key_names,omitempty"`
}

// BaremetalUpdateRequest para actualizar un servidor
type BaremetalUpdateRequest struct {
    Hostname string `json:"hostname,omitempty"`
    Tags     string `json:"tags,omitempty"`
}

// BaremetalReinstallRequest para reinstallar
type BaremetalReinstallRequest struct {
    OSName         string `json:"os_name"`
    Hostname       string `json:"hostname"`
    User           string `json:"user"`
    Password       string `json:"password"`
    DiskLayoutName string `json:"disk_layout_name,omitempty"`
}

// BaremetalServer representa un servidor baremetal
type BaremetalServer struct {
    ID                int               `json:"id"`
    Hostname          string            `json:"hostname"`
    Status            string            `json:"status"`
    Label             string            `json:"label"`
    SSHUsername       string            `json:"ssh_username"`
    MonitoringEnable  bool              `json:"monitoring_enable"`
    Model            BaremetalModelInfo `json:"baremetal_model"`
    OS               BaremetalOSInfo    `json:"os"`
    Location         BaremetalLocInfo   `json:"location"`
    FloatingIPs       []FloatingIPFull   `json:"floating_ips"`
}

type BaremetalModelInfo struct {
    ModelName      string  `json:"model_name"`
    CPU            string  `json:"cpu"`
    CPUSpecs       string  `json:"cpu_specs"`
    CPUBench       float64 `json:"cpu_bench"`
    RAMSize        int     `json:"ram_size"`
    RAMType        string  `json:"ram_type"`
    DiskSize       string  `json:"disk_size"`
    DiskType       string  `json:"disk_type"`
    Port           int     `json:"port"`
    KVM            string  `json:"kvm"`
    Price          float64 `json:"price"`
}

type BaremetalOSInfo struct {
    Name string `json:"name"`
}

type BaremetalLocInfo struct {
    LocationName string `json:"location_name"`
}

type FloatingIPFull struct {
    Type           string `json:"type"`
    Address        string `json:"address"`
    ProtectionType string `json:"protection_type"`
}

// ============================================
// NETWORK
// ============================================

// NetworkCreateRequest para crear una red
type NetworkCreateRequest struct {
    Name         string `json:"name"`
    LocationName string `json:"location_name"`
    IPRange      string `json:"ip_range"`
    Prefix       int    `json:"prefix"`
    ProjectID    int    `json:"project_id"`
    Label        string `json:"label,omitempty"`
}

// NetworkUpdateRequest para actualizar una red
type NetworkUpdateRequest struct {
    Name  string `json:"name,omitempty"`
    Label string `json:"label,omitempty"`
}

// Network representa una red
type Network struct {
    ID           int    `json:"id"`
    Name         string `json:"name"`
    IPRange      string `json:"ip_range"`
    LocationName string `json:"location_name"`
}

// ============================================
// PROJECT
// ============================================

// ProjectCreateRequest para crear un proyecto
type ProjectCreateRequest struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

// Project representa un proyecto
type Project struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description"`
}

// ============================================
// SSH KEY
// ============================================

// SSHKeyCreateRequest para crear una clave SSH
type SSHKeyCreateRequest struct {
    Name   string `json:"name"`
    SSHKey string `json:"ssh_key"`
}

// SSHKey representa una clave SSH
type SSHKey struct {
    ID          int    `json:"id"`
    Name        string `json:"name"`
    KeyType     string `json:"key_type"`
    Fingerprint string `json:"fingerprint"`
}

// ============================================
// FLOATING IP
// ============================================

// FloatingIP representa una IP flotante
type FloatingIP struct {
    Address        string  `json:"address"`
    Type           string  `json:"type"`
    Status         string  `json:"status"`
    IsPrimary      bool    `json:"is_primary"`
    VPSName        *string `json:"vps_name,omitempty"`
    BaremetalName  *string `json:"baremetal_name,omitempty"`
    LocationName   string  `json:"location_name"`
    ProtectionType string  `json:"protection_type"`
}

// FloatingIPsResponse representa la respuesta de IPs
type FloatingIPsResponse struct {
    SingleIPs []FloatingIP        `json:"single_ips"`
    Subnets   []FloatingIPSubnet  `json:"subnets"`
}

type FloatingIPSubnet struct {
    Prefix         int          `json:"prefix"`
    ProtectionType string       `json:"protection_type"`
    IPAddresses    []FloatingIP `json:"ip_addresses"`
}
```

---

## Sección 3: Ejemplos de Petición

### Go - Cliente HTTP Básico

```go
package cubepath

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "time"
)

type Client struct {
    BaseURL    string
    Token      string
    HTTPClient *http.Client
}

func NewClient(baseURL, token string) *Client {
    return &Client{
        BaseURL: baseURL,
        Token:   token,
        HTTPClient: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (c *Client) doRequest(method, path string, body interface{}) (json.RawMessage, error) {
    var reqBody io.Reader
    if body != nil {
        data, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        reqBody = bytes.NewReader(data)
    }

    url := c.BaseURL + path
    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Authorization", "Bearer "+c.Token)
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("User-Agent", "CubeArchitect/1.0.0")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("connection error: %w", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("API error: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
    }

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }

    if len(respBody) == 0 {
        return json.RawMessage("{}"), nil
    }

    return json.RawMessage(respBody), nil
}

func (c *Client) Get(path string) (json.RawMessage, error) {
    return c.doRequest(http.MethodGet, path, nil)
}

func (c *Client) Post(path string, body interface{}) (json.RawMessage, error) {
    return c.doRequest(http.MethodPost, path, body)
}

func (c *Client) Put(path string, body interface{}) (json.RawMessage, error) {
    return c.doRequest(http.MethodPut, path, body)
}

func (c *Client) Patch(path string, body interface{}) (json.RawMessage, error) {
    return c.doRequest(http.MethodPatch, path, body)
}

func (c *Client) Delete(path string) (json.RawMessage, error) {
    return c.doRequest(http.MethodDelete, path, nil)
}
```

### Go - Crear VPS

```go
package main

import (
    "fmt"
    "log"

    "your-package/cubepath"
)

type VPSCreateRequest struct {
    Name           string   `json:"name"`
    PlanName       string   `json:"plan_name"`
    TemplateName   string   `json:"template_name"`
    LocationName   string   `json:"location_name"`
    SSHKeyNames    []string `json:"ssh_key_names,omitempty"`
    Password       string   `json:"password,omitempty"`
    IPv4           bool     `json:"ipv4"`
    EnableBackups  bool     `json:"enable_backups"`
}

func main() {
    client := cubepath.NewClient(
        "https://api.cubepath.com",
        "YOUR_API_TOKEN",
    )

    req := VPSCreateRequest{
        Name:         "my-web-server",
        PlanName:     "vc-2",
        TemplateName: "ubuntu-22.04",
        LocationName: "ams1",
        SSHKeyNames:  []string{"my-key"},
        IPv4:         true,
        EnableBackups: false,
    }

    resp, err := client.Post(fmt.Sprintf("/vps/create/%d", 123), req)
    if err != nil {
        log.Fatalf("Error creating VPS: %v", err)
    }

    fmt.Printf("VPS created: %s\n", string(resp))
}
```

### cURL - Crear VPS

```bash
curl -X POST "https://api.cubepath.com/vps/create/123" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -H "User-Agent: CubeArchitect/1.0.0" \
  -d '{
    "name": "my-web-server",
    "plan_name": "vc-2",
    "template_name": "ubuntu-22.04",
    "location_name": "ams1",
    "ssh_key_names": ["my-key"],
    "ipv4": true,
    "enable_backups": false
  }'
```

### cURL - Crear Load Balancer

```bash
curl -X POST "https://api.cubepath.com/loadbalancer/" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -H "User-Agent: CubeArchitect/1.0.0" \
  -d '{
    "name": "my-load-balancer",
    "plan_name": "lb-small",
    "location_name": "ams1",
    "project_id": 123
  }'
```

### cURL - Crear CDN Zone

```bash
curl -X POST "https://api.cubepath.com/cdn/zones" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -H "User-Agent: CubeArchitect/1.0.0" \
  -d '{
    "name": "my-cdn-zone",
    "plan_name": "cdn-basic",
    "custom_domain": "cdn.example.com"
  }'
```

### cURL - Listar DNS Records

```bash
curl -X GET "https://api.cubepath.com/dns/zones/ZONE_UUID/records" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "User-Agent: CubeArchitect/1.0.0"
```

### cURL - Crear DNS Record

```bash
curl -X POST "https://api.cubepath.com/dns/zones/ZONE_UUID/records" \
  -H "Authorization: Bearer YOUR_API_TOKEN" \
  -H "Content-Type: application/json" \
  -H "User-Agent: CubeArchitect/1.0.0" \
  -d '{
    "name": "api",
    "record_type": "A",
    "content": "192.0.2.1",
    "ttl": 3600
  }'
```

### Go Fiber - Middleware de Proxy

```go
package main

import (
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    "github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
    app := fiber.New(fiber.Config{
        ErrorHandler: func(c *fiber.Ctx, e error) error {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "error": e.Error(),
            })
        },
    })

    app.Use(logger.New())
    app.Use(cors.New())

    apiToken := os.Getenv("CUBE_API_TOKEN")
    if apiToken == "" {
        log.Fatal("CUBE_API_TOKEN not set")
    }

    // Middleware para añadir headers de CubePath
    app.Use(func(c *fiber.Ctx) error {
        c.Request().Header.Set("Authorization", "Bearer "+apiToken)
        c.Request().Header.Set("Content-Type", "application/json")
        c.Request().Header.Set("User-Agent", "CubeArchitect/1.0.0")
        return c.Next()
    })

    // Health check
    app.Get("/health", func(c *fiber.Ctx) error {
        return c.JSON(fiber.Map{"status": "ok"})
    })

    // Proxy a CubePath API
    app.All("/api/*", func(c *fiber.Ctx) error {
        targetURL := "https://api.cubepath.com" + c.OriginalURL()[4:] // Remove "/api"

        req, err := fiber.NewRequest(c.Method(), targetURL, c.Body())
        if err != nil {
            return err
        }

        // Copiar headers
        c.Request().Header.VisitAll(func(key, value []byte) {
            req.Header.Set(string(key), string(value))
        })

        resp, err := req.Response()
        if err != nil {
            return err
        }

        c.Type("json")
        return c.Send(resp.Body())
    })

    log.Fatal(app.Listen(":8080"))
}
```

---

## Estados de Recursos

El CLI formatea los estados de la siguiente manera:

| Estado | Color | Significado |
|--------|-------|-------------|
| `active`, `running`, `healthy`, `enabled`, `verified` | Verde | Operativo |
| `stopped`, `paused`, `pending`, `provisioning`, `inactive` | Gris | Detenido/Pendiente |
| `error`, `failed`, `unhealthy` | Rojo | Error |
| Otros | Normal | Estado genérico |

### Estados Comunes de VPS
- `provisioning` - Creando
- `running` - En ejecución
- `stopped` - Detenido
- `reinstalling` - Reinstalando
- `resizing` - Redimensionando

### Estados de Load Balancer
- `creating` - Creando
- `active` - Activo
- `error` - Error

### Estados de CDN Zone
- `creating` - Creando
- `active` - Activo
- `pending` - Pendiente

### Estados de Baremetal
- `provisioning` - Provisionando
- `running` - En ejecución
- `stopped` - Detenido

---

## Notas Importantes

### Polling y Monitoreo
**El CLI NO implementa polling activo ni WebSockets.** Solo hace peticiones HTTP individuales y muestra el estado en la respuesta.

Para implementar monitoreo de despliegue en tu backend:
1. Hacer petición POST inicial
2. Implementar polling periódico (ej: cada 5 segundos)
3. Verificar estado hasta que sea `running` o `error`
4. Timeout máximo recomendado: 10 minutos

### Internal URLs / Conexiones entre Servicios
El CLI no maneja conexiones internas explícitamente. La información de IPs y URLs se obtiene de los campos `floating_ips` y `ipv4` de cada recurso.

Para conectar servicios internamente:
1. Obtener la IP flotante o interna del servicio destino
2. Usar esa IP como `address` en origins o targets
3. Para bases de datos, usar la IP del servidor de base de datos

### Rate Limiting
No se encontró evidencia de rate limiting en el CLI. Implementar backoff exponencial en caso de errores 429.

---

## Changelog

| Fecha | Descripción |
|-------|-------------|
| 2026-03-20 | Extracción inicial del código fuente de cubecli |

---

*Generado automáticamente para CubeArchitect - Hackathon CubeCLI*
