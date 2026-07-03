package config

import sharedconfig "github.com/panduputragit/gym/backend/packages/config"

type Config struct {
	Name            string
	Port            string
	GinMode         string
	AuthURL         string
	EmployeeURL     string
	BranchURL       string
	MemberURL       string
	MembershipURL   string
	AttendanceURL   string
	NotificationURL string
	PaymentURL      string
}

func Load() Config {
	_ = sharedconfig.LoadEnv(".env", "../../.env", "../../../.env")

	return Config{
		Name:            sharedconfig.String("SERVICE_NAME", "api-gateway"),
		Port:            sharedconfig.String("API_GATEWAY_PORT", sharedconfig.String("PORT", "8080")),
		GinMode:         sharedconfig.String("GIN_MODE", "debug"),
		AuthURL:         sharedconfig.String("AUTH_SERVICE_URL", "http://localhost:5001"),
		EmployeeURL:     sharedconfig.String("EMPLOYEE_SERVICE_URL", "http://localhost:5002"),
		BranchURL:       sharedconfig.String("BRANCH_SERVICE_URL", "http://localhost:5003"),
		MemberURL:       sharedconfig.String("MEMBER_SERVICE_URL", "http://localhost:5004"),
		MembershipURL:   sharedconfig.String("MEMBERSHIP_SERVICE_URL", "http://localhost:5005"),
		AttendanceURL:   sharedconfig.String("ATTENDANCE_SERVICE_URL", "http://localhost:5006"),
		NotificationURL: sharedconfig.String("NOTIFICATION_SERVICE_URL", "http://localhost:5007"),
		PaymentURL:      sharedconfig.String("PAYMENT_SERVICE_URL", "http://localhost:5008"),
	}
}

func (c Config) ServiceURLs() map[string]string {
	return map[string]string{
		"auth":          c.AuthURL,
		"employees":     c.EmployeeURL,
		"branches":      c.BranchURL,
		"members":       c.MemberURL,
		"memberships":   c.MembershipURL,
		"attendance":    c.AttendanceURL,
		"notifications": c.NotificationURL,
		"payments":      c.PaymentURL,
	}
}
