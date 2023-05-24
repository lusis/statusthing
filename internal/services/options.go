package services

// ServiceOption is a functional option for configuring a [StatusThingService]
type ServiceOption func(s *StatusThingService) error

// WithDefaults triggers adding default data at creation time
func WithDefaults() ServiceOption {
	return func(s *StatusThingService) error {
		s.loadDefaults = true
		return nil
	}
}
