package models

type DBModel struct {
	ProcessCreate  []ProcessCreate  `json:"process_create"`
	CreateFile     []CreateFile     `json:"create_file"`
	RegCreateKey   []RegCreateKey   `json:"create_reg_key"`
	OpenRegKey     []OpenRegKey     `json:"open_reg_key"`
	GetRegKey      []GetRegKey      `json:"get_reg_key"`
	SetRegValue    []SetRegValue    `json:"set_reg_value"`
	DeleteRegKey   []DeleteRegKey   `json:"deleted_reg_key"`
	DeleteRegValue []DeleteRegValue `json:"deleted_reg_value"`
	RenameFile     []RenameFile     `json:"rename_file"`
	DeleteFile     []DeleteFile     `json:"deleted_file"`
	UDP            []UDP            `json:"udp"`
	TCP            []TCP            `json:"tcp"`
}

type ProcessCreate struct {
	PID         string `json:"pid"`
	ChildPID    string `json:"child_pid"`
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	Operation   string `json:"operation"`
}
type CreateFile struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	CreatePath  string `json:"create_path"`
}
type RenameFile struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	OriginName  string `json:"origin_name"`
	ChangeName  string `json:"change_name"`
}
type DeleteFile struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	ProcessPath string `json:"process_path"`
	DeletePath  string `json:"delete_path"`
}

type RegCreateKey struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Key         string `json:"key"`
}

type OpenRegKey struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Key         string `json:"key"`
}

type GetRegKey struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Key         string `json:"key"`
}

type SetRegValue struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Value       string `json:"value"`
}

type DeleteRegKey struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	RegKey      string `json:"reg_key"`
	ProcessPath string `json:"process_path"`
}
type DeleteRegValue struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	RegValue    string `json:"reg_value"`
	ProcessPath string `json:"process_path"`
}
type UDP struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Action      string `json:"action"`
	Server      string `json:"server"`
}
type TCP struct {
	PID         string `json:"pid"`
	ProcessName string `json:"process_name"`
	Action      string `json:"action"`
	Server      string `json:"server"`
}