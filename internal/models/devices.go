package models

// DuplexMode represents the duplex mode of a port
type DuplexMode string

const (
	DuplexModeAuto DuplexMode = "auto"
	DuplexModeFull DuplexMode = "full"
	DuplexModeHalf DuplexMode = "half"
)

// Speed represents the speed of a port in Mbps
type Speed string

const (
	SpeedAuto  Speed = "auto"
	Speed10    Speed = "10"
	Speed100   Speed = "100"
	Speed1000  Speed = "1000"
	Speed10000 Speed = "10000"
)

type Router struct {
	ID        uint   `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	Status    string `json:"status"`
	Ports     []Port `json:"ports" gorm:"foreignKey:RouterID"`
	Connected bool   `json:"connected" gorm:"default:false"`
}

// Port represents a network port configuration
type Port struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	RouterID    uint       `json:"router_id"`
	Number      int        `json:"number" gorm:"check:number >= 1 AND number <= 65535"`
	Protocol    string     `json:"protocol" gorm:"check:protocol IN ('tcp', 'udp')"`
	Status      string     `json:"status" gorm:"default:'closed'"` // open, closed, filtered
	PortNumber  int        `json:"portNumber"`
	Speed       Speed      `json:"speed"`
	DuplexMode  DuplexMode `json:"duplexMode"`
	Description string     `json:"description"`
}

type CreateRouterRequest struct {
	Name  string    `json:"name" binding:"required"`
	Ports []PortReq `json:"ports"`
}

type ConnectRouterRequest struct {
	IPAddress string `json:"ip_address" binding:"required"`
}

type ConnectRouterResponse struct {
	RouterID  uint   `json:"router_id"`
	Name      string `json:"name"`
	IPAddress string `json:"ip_address"`
	LocalIP   string `json:"local_ip"`
	Status    string `json:"status"`
	Connected bool   `json:"connected"`
}

type PortReq struct {
	Number   int    `json:"number" binding:"required,min=1,max=65535"`
	Protocol string `json:"protocol" binding:"required,oneof=tcp udp"`
}

type PingRequest struct {
	IPAddress string `json:"ip_address"`
}

type PingResult struct {
	IPAddress string  `json:"ip_address"`
	Latency   float64 `json:"latency"`
	Status    string  `json:"status"`
}

type PacketRequest struct {
	SourceIP      string `json:"source_ip"`
	DestinationIP string `json:"destination_ip"`
	Protocol      string `json:"protocol" binding:"required,oneof=tcp udp"`
	Port          int    `json:"port" binding:"required,min=1,max=65535"`
	Data          string `json:"data"`
}

type PacketResponse struct {
	SourceIP      string  `json:"source_ip"`
	DestinationIP string  `json:"destination_ip"`
	Protocol      string  `json:"protocol"`
	Port          int     `json:"port"`
	Status        string  `json:"status"`
	Latency       float64 `json:"latency"`
	Error         string  `json:"error,omitempty"`
}

type ConfigureRouterRequest struct {
	RouterID uint   `json:"router_id" binding:"required"`
	Name     string `json:"name,omitempty"`
	Status   string `json:"status,omitempty" binding:"oneof=active inactive maintenance"`
}

// ConfigurePortRequest represents the request to configure a port
type ConfigurePortRequest struct {
	RouterID    string     `json:"routerId"`
	PortNumber  int        `json:"portNumber"`
	Status      string     `json:"status"`
	Speed       Speed      `json:"speed"`
	DuplexMode  DuplexMode `json:"duplexMode"`
	Description string     `json:"description"`
}

type ConfigureResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
