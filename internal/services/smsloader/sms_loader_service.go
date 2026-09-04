package smsloader

import (
	"context"
	"log"
	"sync"
	"ticket-api/internal/db/sms_type_messages_relation"
	"time"
)

type SMSLoaderService struct {
	queries  *sms_type_messages_relation.Queries
	messages map[string]string // sms_type.title -> sms_message.body
	mu       sync.RWMutex
}

func NewSMSLoaderService(queries *sms_type_messages_relation.Queries) *SMSLoaderService {
	s := &SMSLoaderService{
		queries:  queries,
		messages: make(map[string]string),
	}

	// Initial synchronous load from DB
	s.loadFromDB()

	// Periodic refresh every 5 minutes in background
	if queries != nil {
		go s.autoRefresh(5 * time.Minute)
	}

	return s
}

// loadFromDB fetches all SMS type-message relations and populates memory cache atomically
func (s *SMSLoaderService) loadFromDB() {
	if s.queries == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	records, err := s.queries.GetAllSMSTypeMessageRelations(ctx)
	if err != nil {
		log.Printf("[SMSLoader] Failed to load SMS type-message relations: %v", err)
		return
	}

	temp := make(map[string]string, len(records))
	for _, r := range records {
		temp[r.SmsTypeTitle] = r.SmsMessageBody
	}

	s.mu.Lock()
	s.messages = temp
	s.mu.Unlock()

	log.Printf("[SMSLoader] Loaded %d SMS message templates from DB", len(temp))
}

// autoRefresh periodically updates the in-memory SMS templates
func (s *SMSLoaderService) autoRefresh(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		s.loadFromDB()
	}
}

// GetMessageByType returns the SMS message template body for a given SMS type title
func (s *SMSLoaderService) GetMessageByType(typeTitle string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	msg, ok := s.messages[typeTitle]
	return msg, ok
}
