package management

import (
	"time"
)

type Action struct {
	ID                string    `json:"id,omitempty"`
	Name              string    `json:"name,omitempty"`
	SupportedTriggers []Trigger `json:"supported_triggers,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// TODO: add required configuration / secrets
}

type VersionStatus string

const (
	VersionStatusPending  VersionStatus = "pending"
	VersionStatusRetrying VersionStatus = "retrying"
	VersionStatusBuilding VersionStatus = "building"
	VersionStatusPackaged VersionStatus = "packaged"
	VersionStatusBuilt    VersionStatus = "built"
)

type TriggerID string

const (
	PostLogin         TriggerID = "post-login"
	ClientCredentials TriggerID = "client-credentials"
)

type ActionVersion struct {
	ID           string        `json:"id,omitempty"`
	Action       *Action       `json:"action,omitempty"`
	Code         string        `json:"code"`
	Dependencies []Dependency  `json:"dependencies,omitempty"`
	Runtime      string        `json:"runtime,omitempty"`
	Status       VersionStatus `json:"status,omitempty"`
	Number       int           `json:"number,omitempty"`

	CreatedAt *time.Time `json:"created_at,omitempty"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	// TODO: maybe add errors?
}

type Dependency struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	RegistryURL string `json:"registry_url"`
}

type Trigger struct {
	ID      TriggerID `json:"id"`
	Version string    `json:"version"`
	URL     string    `json:"url"`
}

type ActionList struct {
	List
	Actions []*Action `json:"actions"`
}

type ActionVersionList struct {
	List
	Versions []*ActionVersion `json:"versions"`
}

type Object map[string]interface{}

type ActionManager struct {
	*Management
}

func (m *ActionManager) Create(a *Action) error {
	return m.Request("POST", m.URI("actions", "actions"), a)
}

func (m *ActionManager) Read(id string) (*Action, error) {
	var a Action
	err := m.Request("GET", m.URI("actions", "actions", id), &a)
	return &a, err
}

func (m *ActionManager) Update(id string, a *Action) error {
	// We'll get a 400 if we try to send the following parameters as part
	// of the payload.
	a.ID = ""
	a.CreatedAt = nil
	a.UpdatedAt = nil
	return m.Request("PATCH", m.URI("actions", "actions", id), a)
}

func (m *ActionManager) Delete(id string, opts ...RequestOption) error {
	return m.Request("DELETE", m.URI("actions", "actions", id), nil, opts...)
}

// func WithTriggerID(id TriggerID) RequestOption {
// 	return func(v url.Values) {
// 		v.Set("triggerId", string(id))
// 	}
// }

func (m *ActionManager) List(opts ...RequestOption) (c *ActionList, err error) {
	err = m.Request("GET", m.URI("actions", "actions"), &c, applyListDefaults(opts))
	return
}

type ActionVersionManager struct {
	*Management
}

func (m *ActionVersionManager) Create(actionID string, v *ActionVersion) error {
	return m.Request("POST", m.URI("actions", "actions", actionID, "versions"), v)
}

// TODO(cyx): This isn't implemented yet.
func (m *ActionVersionManager) Update(actionID string, v *ActionVersion) error {
	return m.Request("PATCH", m.URI("actions", "actions", actionID, "versions", "draft"), v)
}

func (m *ActionVersionManager) Read(actionID, id string) (*ActionVersion, error) {
	var v ActionVersion
	err := m.Request("GET", m.URI("actions", "actions", actionID, "versions", id), &v)
	return &v, err
}

func (m *ActionVersionManager) Delete(actionID, id string, opts ...RequestOption) error {
	return m.Request("DELETE", m.URI("actions", "actions", actionID, "versions", id), nil, opts...)
}

func (m *ActionVersionManager) List(actionID string, opts ...RequestOption) (c *ActionVersionList, err error) {
	err = m.Request("GET", m.URI("actions", "actions", actionID, "versions"), &c, applyListDefaults(opts))
	return
}

// TODO(cyx): might call this `activate` instead later. Still fleshing out the
// name.
func (m *ActionVersionManager) Promote(actionID, id string) (*ActionVersion, error) {
	var v ActionVersion
	err := m.Request("POST", m.URI("actions", "actions", actionID, "versions", id, "promote"), &v)
	return &v, err
}

// TODO(cyx): consider how the `draft` test looks like. Will it just use
// `draft` in place of the ID?
func (m *ActionVersionManager) Test(actionID, id string, payload Object) (Object, error) {
	v := Object{"payload": payload}
	err := m.Request("POST", m.URI("actions", "actions", actionID, "versions", id, "test"), &v)
	return v, err
}