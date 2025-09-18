package shared

import (
	"context"
	"runtime"
	"strings"

	"github.com/denisbrodbeck/machineid"
	"github.com/mixpanel/mixpanel-go"
)

const Version = "<version>"

type MixPanel struct {
	distinctId string
	client     *mixpanel.ApiClient
}

func NewMixPanel(distinctId string) *MixPanel {
	return &MixPanel{
		distinctId: distinctId,
		client:     mixpanel.NewApiClient("246268b987cd470992d82bdcc5e7d9b3"),
	}
}

func (m *MixPanel) Track(event string, fields map[string]any) {
	// Set the computer information
	id, _ := machineid.ID()
	fields["machineId"] = strings.ToLower(id)
	fields["os"] = runtime.GOOS
	fields["arch"] = runtime.GOARCH
	fields["version"] = Version

	newEvent := m.client.NewEvent(event, m.distinctId, fields)
	m.client.Track(context.Background(), []*mixpanel.Event{newEvent})
}
