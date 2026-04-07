package audit_logs

import (
	"context"
	"encoding/json"

	domain "github.com/OctoetIx/Hotel-Bookings-and-Reservation/pkg/domain/audit_logs"
)

type Publisher interface {
	Publish(topic string, message []byte) error
}

type Service struct {
	repo domain.Repository
	pub  Publisher
}

func NewService(repo domain.Repository, pub Publisher) *Service {
	return &Service{repo: repo, pub: pub}
}

// ComputeDiff - simple JSON diff
func ComputeDiff(before, after json.RawMessage) json.RawMessage {
	var b, a map[string]interface{}
	json.Unmarshal(before, &b)
	json.Unmarshal(after, &a)

	diff := map[string]interface{}{}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			diff[k] = v
		}
	}
	result, _ := json.Marshal(diff)
	return result
}

// EvaluateRisk - naive example
func EvaluateRisk(log *domain.AuditLog) (bool, string, string) {

	switch log.Action {

	case "DELETE_USER":
		return true, string(domain.RiskHigh), "User deletion"

	case "PROMOTE_TO_ADMIN":
		return true, string(domain.RiskCritical), "Privilege escalation"

	case "LOGIN_FAILED":
		return false, string(domain.RiskMedium), "Failed login"

	case "SUSPICIOUS_LOGIN":
		return true, string(domain.RiskHigh), "IP address change"

	default:
		return false, string(domain.RiskLow), ""
	}
}

// LogEvent - main usecase method
func (s *Service) LogEvent(ctx context.Context, event *domain.AuditLog) error {

	// compute field differences
	event.ChangedFields = ComputeDiff(event.BeforeState, event.AfterState)

	// risk evaluation
	suspicious, level, reason := EvaluateRisk(event)
	event.Suspicious = suspicious
	event.RiskLevel = domain.RiskLevel(level)
	event.Reason = reason

	// persist audit log
	if err := s.repo.Save(ctx, event); err != nil {
		return err
	}

	// serialize event for streaming
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// publish asynchronously (do not block request)
	go func() {
		_ = s.pub.Publish("audit_events", data)
	}()

	return nil
}


func (s *Service) FindWithFilter(
	ctx context.Context,
	filter domain.AuditFilter,
	page int,
	limit int,
) ([]domain.AuditLog, int64, error) {

	return s.repo.FindWithFilter(ctx, filter, page, limit)
}

func (s *Service) GetAuditLogs(ctx context.Context, filter domain.AuditFilter, page, limit int) ([]domain.AuditLog, int64, error) {
	return s.repo.FindWithFilter(ctx,filter, page, limit)
}