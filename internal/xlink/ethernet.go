package xlink

type Ethernet struct {
	Id         string `json:"id"`
	Ip         string `json:"ip"`
	Gate       string `json:"gate"`
	Mask       string `json:"mask"`
	Link       *bool  `json:"link"`
	Admin      *bool  `json:"admin"`
	Enabled    *bool  `json:"enabled"`
	DefaultLan *bool  `json:"defaultLan"`
	AdminOnly  *bool  `json:"adminOnly"`
	Igmp       *bool  `json:"igmp"`
	Ndi        *bool  `json:"ndi"`
	Default    *bool  `json:"default"`
	Backup     *bool  `json:"backup"`
	Active     *bool  `json:"active"`
}
