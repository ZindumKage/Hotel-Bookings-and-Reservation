
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
	if log.Action == "DELETE_USER" {
		return true, string(domain.RiskHigh), "Deleting a user"
	}
	return false, string(domain.RiskLow), ""
}

// LogEvent - main usecase method
func (s *Service) LogEvent(event *domain.AuditLog) error {
	event.ChangedFields = ComputeDiff(event.BeforeState, event.AfterState)
	suspicious, level, reason := EvaluateRisk(event)
	event.Suspicious = suspicious
	event.RiskLevel = domain.RiskLevel(level)
	event.Reason = reason

	if err := s.repo.Save(event); err != nil {
		return err
	}

	data, _ := json.Marshal(event)
	s.pub.Publish("audit_events", data)
	return nil
}

func (s *Service) GetAuditLogs(ctx context.Context, filter domain.AuditFilter, page, limit int) ([]domain.AuditLog, int64, error) {
	return s.repo.FindWithFilter(filter, page, limit)
}