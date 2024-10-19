package main

import (
	"fmt"
	"math/rand"
	"sync"
)

type PolicyPaymentEvent struct {
	EventID   string
	PolicyID  string
	ProductID string
	AgentID   string
	Premium   int32
}

type AgentPolicyStatistics struct {
	AgentID       string
	TotalPremiums int32
	AvgPremium    int32
}

type PolicyStatistics struct {
	Agents []AgentPolicyStatistics
}

type PolicyService interface {
	ConsumeEvent(event PolicyPaymentEvent)
	Stats() PolicyStatistics
}

type ConcurrentPolicyService struct {
	mu   sync.RWMutex
	data map[string][]PolicyPaymentEvent
}

func NewConcurrentPolicyService() *ConcurrentPolicyService {
	return &ConcurrentPolicyService{
		data: make(map[string][]PolicyPaymentEvent),
	}
}

func (s *ConcurrentPolicyService) ConsumeEvent(event PolicyPaymentEvent) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.data[event.AgentID] = append(s.data[event.AgentID], event)
}

func (s *ConcurrentPolicyService) Stats() PolicyStatistics {
	s.mu.RLock()
	defer s.mu.RUnlock()

	stats := PolicyStatistics{}
	for agentID, events := range s.data {
		totalPremiums := int32(0)
		for _, event := range events {
			totalPremiums += event.Premium
		}
		avgPremium := int32(0)
		if len(events) > 0 {
			avgPremium = totalPremiums / int32(len(events))
		}
		stats.Agents = append(stats.Agents, AgentPolicyStatistics{
			AgentID:       agentID,
			TotalPremiums: totalPremiums,
			AvgPremium:    avgPremium,
		})
	}
	return stats
}

func main() {
	service := NewConcurrentPolicyService()

	// Simulate concurrent event consumption
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			event := PolicyPaymentEvent{
				EventID:   fmt.Sprintf("evt_%d", rand.Intn(1000)),
				PolicyID:  fmt.Sprintf("pol_%d", rand.Intn(100)),
				ProductID: fmt.Sprintf("prod_%d", rand.Intn(10)),
				AgentID:   fmt.Sprintf("agent_%d", rand.Intn(5)),
				Premium:   int32(rand.Intn(1000) + 100),
			}
			service.ConsumeEvent(event)
		}()
	}
	wg.Wait()

	// Print statistics
	stats := service.Stats()
	for _, agentStats := range stats.Agents {
		fmt.Printf("Agent: %s, Total Premiums: %d, Avg Premium: %d\n",
			agentStats.AgentID, agentStats.TotalPremiums, agentStats.AvgPremium)
	}
}