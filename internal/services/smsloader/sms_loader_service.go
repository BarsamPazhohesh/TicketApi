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
	typeIDs  map[string]int64  // sms_type.title -> sms_type.id
	mu       sync.RWMutex
}

func NewSMSLoaderService(queries *sms_type_messages_relation.Queries) *SMSLoaderService {
	s := &SMSLoaderService{
		queries:  queries,
		messages: make(map[string]string),
		typeIDs:  make(map[string]int64),
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

	tempMessages := make(map[string]string, len(records))
	tempTypeIDs := make(map[string]int64, len(records))
	for _, r := range records {
		tempMessages[r.SmsTypeTitle] = r.SmsMessageBody
		tempTypeIDs[r.SmsTypeTitle] = r.SmsTypeID
	}

	s.mu.Lock()
	s.messages = tempMessages
	s.typeIDs = tempTypeIDs
	s.mu.Unlock()

	log.Printf("[SMSLoader] Loaded %d SMS message templates from DB", len(tempMessages))
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

// GetTypeIDByTitle returns the SMS type ID for a given SMS type title
func (s *SMSLoaderService) GetTypeIDByTitle(typeTitle string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	id, ok := s.typeIDs[typeTitle]
	return id, ok
}

