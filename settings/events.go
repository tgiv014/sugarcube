package settings

type SettingsUpdatedEvent struct {
}

func (s SettingsUpdatedEvent) Topic() string {
	return "settingsUpdated"
}

func (s SettingsUpdatedEvent) Decode(v any) error {
	return nil
}
