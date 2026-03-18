package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// --- Data Models ---

type Status string

const (
	StatusActive      Status = "active"
	StatusIdle        Status = "idle"
	StatusCharging    Status = "charging"
	StatusError       Status = "error"
	StatusMaintenance Status = "maintenance"
)

type Robot struct {
	ID                 string  `json:"id"`
	Category           string  `json:"category"`
	Status             Status  `json:"status"`
	BatteryLevel       float64 `json:"battery_level"`
	CPUTemp            float64 `json:"cpu_temp"`
	TaskCompletionRate float64 `json:"task_completion_rate"`
	DistanceTraveled   float64 `json:"distance_traveled_today"`
	ErrorCount         int     `json:"error_count"`
	UptimeHours        float64 `json:"uptime_hours"`
	Lat                float64 `json:"lat"`
	Lng                float64 `json:"lng"`
	LastEvent          string  `json:"last_event,omitempty"`
	LastEventTime      int64   `json:"last_event_time,omitempty"`
}

type CategoryConfig struct {
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	Count  int    `json:"count"`
	Icon   string `json:"icon"`
	LatMin float64
	LatMax float64
	LngMin float64
	LngMax float64
}

type CategorySummary struct {
	Name            string  `json:"name"`
	Icon            string  `json:"icon"`
	Count           int     `json:"count"`
	ActiveCount     int     `json:"active_count"`
	IdleCount       int     `json:"idle_count"`
	ChargingCount   int     `json:"charging_count"`
	ErrorCount      int     `json:"error_count"`
	MaintenanceCount int    `json:"maintenance_count"`
	AvgBattery      float64 `json:"avg_battery"`
	AvgCPUTemp      float64 `json:"avg_cpu_temp"`
	AvgTaskRate     float64 `json:"avg_task_rate"`
	TotalErrors     int     `json:"total_errors"`
}

type FleetSnapshot struct {
	Timestamp        int64   `json:"timestamp"`
	TotalRobots      int     `json:"total_robots"`
	ActiveCount      int     `json:"active_count"`
	IdleCount        int     `json:"idle_count"`
	ChargingCount    int     `json:"charging_count"`
	ErrorCount       int     `json:"error_count"`
	MaintenanceCount int     `json:"maintenance_count"`
	AvgBattery       float64 `json:"avg_battery"`
	AvgTaskRate      float64 `json:"avg_task_rate"`
	TotalErrorCount  int     `json:"total_error_count"`
	EventsPerSec     float64 `json:"events_per_sec"`
	TasksCompleted   int     `json:"tasks_completed"`
	AlertsActive     int     `json:"alerts_active"`
	AvgCPUTemp       float64 `json:"avg_cpu_temp"`
	NetworkLatency   float64 `json:"network_latency"`
	Events           []Event `json:"events"`
}

type Event struct {
	Timestamp int64  `json:"timestamp"`
	RobotID   string `json:"robot_id"`
	Category  string `json:"category"`
	Message   string `json:"message"`
	Type      string `json:"type"` // info, warning, error
}

// --- Global State ---

var (
	robots     []*Robot
	robotIndex map[string]*Robot
	categories []CategoryConfig
	mu         sync.RWMutex
	events     []Event
	eventMu    sync.RWMutex
	startTime  time.Time

	// Per-tick volatile counters (written by simulationTick, read by computeFleetSnapshot)
	tickStatusChanges int
	tickTasksDone     int
)

func initCategories() []CategoryConfig {
	return []CategoryConfig{
		{Name: "Happy", Prefix: "HAP", Count: 300, Icon: "😄", LatMin: 40.70, LatMax: 40.80, LngMin: -74.02, LngMax: -73.95},
		{Name: "Curious", Prefix: "CUR", Count: 200, Icon: "🧐", LatMin: 40.72, LatMax: 40.78, LngMin: -74.00, LngMax: -73.96},
		{Name: "Depressed", Prefix: "DEP", Count: 250, Icon: "😞", LatMin: 40.74, LatMax: 40.82, LngMin: -73.99, LngMax: -73.93},
		{Name: "Smug", Prefix: "SMG", Count: 150, Icon: "😏", LatMin: 40.63, LatMax: 40.66, LngMin: -73.80, LngMax: -73.76},
		{Name: "Existential Dread", Prefix: "DRD", Count: 400, Icon: "😱", LatMin: 40.65, LatMax: 40.72, LngMin: -74.05, LngMax: -73.98},
		{Name: "Angry", Prefix: "ANG", Count: 250, Icon: "🤬", LatMin: 40.75, LatMax: 40.80, LngMin: -73.98, LngMax: -73.93},
		{Name: "Confused", Prefix: "CNF", Count: 200, Icon: "😵‍💫", LatMin: 40.76, LatMax: 40.79, LngMin: -73.96, LngMax: -73.94},
		{Name: "Delirious", Prefix: "DLR", Count: 350, Icon: "🤪", LatMin: 40.68, LatMax: 40.82, LngMin: -74.03, LngMax: -73.90},
	}
}

func generateFleet() {
	categories = initCategories()
	robots = make([]*Robot, 0, 2100)
	robotIndex = make(map[string]*Robot)

	for _, cat := range categories {
		for i := 0; i < cat.Count; i++ {
			id := fmt.Sprintf("%s-%04d", cat.Prefix, i+1)
			status := randomStatus()
			r := &Robot{
				ID:                 id,
				Category:           cat.Name,
				Status:             status,
				BatteryLevel:       30 + rand.Float64()*70,
				CPUTemp:            35 + rand.Float64()*40,
				TaskCompletionRate: 50 + rand.Float64()*50,
				DistanceTraveled:   rand.Float64() * 30,
				ErrorCount:         rand.Intn(3),
				UptimeHours:        rand.Float64() * 720,
				Lat:                cat.LatMin + rand.Float64()*(cat.LatMax-cat.LatMin),
				Lng:                cat.LngMin + rand.Float64()*(cat.LngMax-cat.LngMin),
			}
			if status == StatusError {
				r.ErrorCount = 1 + rand.Intn(5)
			}
			robots = append(robots, r)
			robotIndex[id] = r
		}
	}
}

func randomStatus() Status {
	r := rand.Float64()
	switch {
	case r < 0.70:
		return StatusActive
	case r < 0.82:
		return StatusIdle
	case r < 0.92:
		return StatusCharging
	case r < 0.97:
		return StatusError
	default:
		return StatusMaintenance
	}
}

// --- Simulation ---

func simulationTick() {
	mu.Lock()
	defer mu.Unlock()

	elapsed := time.Since(startTime).Seconds()
	// Daily curve factor (simulates activity patterns)
	dailyCurve := 0.7 + 0.3*math.Sin(elapsed/30*math.Pi)

	var newEvents []Event
	now := time.Now().UnixMilli()
	statusChanges := 0
	tasksDone := 0

	for _, r := range robots {
		// Battery drift
		switch r.Status {
		case StatusActive:
			r.BatteryLevel -= 0.02 + rand.Float64()*0.08
			r.DistanceTraveled += rand.Float64() * 0.05
			r.CPUTemp += (rand.Float64() - 0.45) * 0.5
		case StatusCharging:
			r.BatteryLevel += 0.3 + rand.Float64()*0.5
		case StatusIdle:
			r.BatteryLevel -= 0.005
			r.CPUTemp -= 0.1
		case StatusError:
			r.CPUTemp += rand.Float64() * 0.3
		}

		// Clamp values
		r.BatteryLevel = clamp(r.BatteryLevel, 0, 100)
		r.CPUTemp = clamp(r.CPUTemp, 30, 95)
		r.TaskCompletionRate = clamp(r.TaskCompletionRate+(rand.Float64()-0.48)*0.5, 0, 100)

		// Status transitions
		oldStatus := r.Status
		if r.BatteryLevel < 10 && r.Status == StatusActive {
			r.Status = StatusCharging
			newEvents = append(newEvents, Event{
				Timestamp: now, RobotID: r.ID, Category: r.Category,
				Message: "battery critical, taking a power nap", Type: "warning",
			})
		} else if r.BatteryLevel > 95 && r.Status == StatusCharging {
			r.Status = StatusActive
			newEvents = append(newEvents, Event{
				Timestamp: now, RobotID: r.ID, Category: r.Category,
				Message: "fully charged, feeling excited!", Type: "info",
			})
		}

		// Random status changes (4x more frequent)
		if rand.Float64() < 0.008*dailyCurve {
			if r.Status == StatusActive {
				if rand.Float64() < 0.3 {
					r.Status = StatusError
					r.ErrorCount++
					newEvents = append(newEvents, Event{
						Timestamp: now, RobotID: r.ID, Category: r.Category,
						Message: "entered Marvin mode, feeling very depressed", Type: "error",
					})
				} else {
					r.Status = StatusIdle
				}
			} else if r.Status == StatusIdle && rand.Float64() < 0.5 {
				r.Status = StatusActive
				newEvents = append(newEvents, Event{
					Timestamp: now, RobotID: r.ID, Category: r.Category,
					Message: "found new data, feeling excited again!", Type: "info",
				})
			} else if r.Status == StatusError && rand.Float64() < 0.3 {
				r.Status = StatusMaintenance
				newEvents = append(newEvents, Event{
					Timestamp: now, RobotID: r.ID, Category: r.Category,
					Message: "heading to the spa for some self-care", Type: "warning",
				})
			} else if r.Status == StatusMaintenance && rand.Float64() < 0.4 {
				r.Status = StatusActive
				r.ErrorCount = 0
				newEvents = append(newEvents, Event{
					Timestamp: now, RobotID: r.ID, Category: r.Category,
					Message: "spa day complete, feeling refreshed and excited!", Type: "info",
				})
			}
		}

		// Task completion events (3x more frequent)
		if r.Status == StatusActive && rand.Float64() < 0.015 {
			tasksDone++
			newEvents = append(newEvents, Event{
				Timestamp: now, RobotID: r.ID, Category: r.Category,
				Message: "task crushed it, living the dream", Type: "info",
			})
		}

		if oldStatus != r.Status {
			statusChanges++
			r.LastEvent = fmt.Sprintf("status changed: %s → %s", oldStatus, r.Status)
			r.LastEventTime = now
		}
	}

	// Occasional cascading failure pattern (4x more frequent, up to 12 affected)
	if rand.Float64() < 0.004 {
		cat := categories[rand.Intn(len(categories))]
		affected := 0
		for _, r := range robots {
			if r.Category == cat.Name && r.Status == StatusActive && affected < 12 {
				r.Status = StatusError
				r.ErrorCount++
				affected++
				statusChanges++
				newEvents = append(newEvents, Event{
					Timestamp: now, RobotID: r.ID, Category: r.Category,
					Message: fmt.Sprintf("existential crisis spreading through %s fleet", cat.Name), Type: "error",
				})
			}
		}
	}

	// Update per-tick counters for volatile metrics
	tickStatusChanges = statusChanges
	tickTasksDone = tasksDone

	// Store events
	if len(newEvents) > 0 {
		eventMu.Lock()
		events = append(newEvents, events...)
		if len(events) > 200 {
			events = events[:200]
		}
		eventMu.Unlock()
	}
}

func clamp(val, min, max float64) float64 {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}

// --- HTTP Handlers ---

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func robotsHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	categoryFilter := r.URL.Query().Get("category")
	statusFilter := r.URL.Query().Get("status")

	if perPage <= 0 || perPage > 500 {
		perPage = 50
	}
	if page <= 0 {
		page = 1
	}

	filtered := make([]*Robot, 0)
	for _, robot := range robots {
		if categoryFilter != "" && robot.Category != categoryFilter {
			continue
		}
		if statusFilter != "" && robot.Status != Status(statusFilter) {
			continue
		}
		filtered = append(filtered, robot)
	}

	total := len(filtered)
	start := (page - 1) * perPage
	if start > total {
		start = total
	}
	end := start + perPage
	if end > total {
		end = total
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"robots":   filtered[start:end],
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func robotByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	mu.RLock()
	defer mu.RUnlock()

	robot, ok := robotIndex[id]
	if !ok {
		http.Error(w, "robot not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(robot)
}

func categoriesHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	defer mu.RUnlock()

	summaries := computeCategorySummaries()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(summaries)
}

func computeCategorySummaries() []CategorySummary {
	catMap := make(map[string]*CategorySummary)
	catOrder := make([]string, 0)

	for _, cat := range categories {
		catMap[cat.Name] = &CategorySummary{Name: cat.Name, Icon: cat.Icon}
		catOrder = append(catOrder, cat.Name)
	}

	for _, r := range robots {
		s := catMap[r.Category]
		s.Count++
		s.AvgBattery += r.BatteryLevel
		s.AvgCPUTemp += r.CPUTemp
		s.AvgTaskRate += r.TaskCompletionRate
		s.TotalErrors += r.ErrorCount

		switch r.Status {
		case StatusActive:
			s.ActiveCount++
		case StatusIdle:
			s.IdleCount++
		case StatusCharging:
			s.ChargingCount++
		case StatusError:
			s.ErrorCount++
		case StatusMaintenance:
			s.MaintenanceCount++
		}
	}

	result := make([]CategorySummary, 0, len(catOrder))
	for _, name := range catOrder {
		s := catMap[name]
		if s.Count > 0 {
			s.AvgBattery /= float64(s.Count)
			s.AvgCPUTemp /= float64(s.Count)
			s.AvgTaskRate /= float64(s.Count)
		}
		s.AvgBattery = math.Round(s.AvgBattery*10) / 10
		s.AvgCPUTemp = math.Round(s.AvgCPUTemp*10) / 10
		s.AvgTaskRate = math.Round(s.AvgTaskRate*10) / 10
		result = append(result, *s)
	}
	return result
}

func computeFleetSnapshot() FleetSnapshot {
	now := time.Now()
	snap := FleetSnapshot{
		Timestamp:   now.UnixMilli(),
		TotalRobots: len(robots),
	}

	var totalBattery, totalTaskRate, totalCPUTemp float64
	alertsActive := 0
	for _, r := range robots {
		totalBattery += r.BatteryLevel
		totalTaskRate += r.TaskCompletionRate
		totalCPUTemp += r.CPUTemp
		snap.TotalErrorCount += r.ErrorCount

		// Alerts: robots in error/maintenance with high CPU
		if (r.Status == StatusError || r.Status == StatusMaintenance) && r.CPUTemp > 70 {
			alertsActive++
		}

		switch r.Status {
		case StatusActive:
			snap.ActiveCount++
		case StatusIdle:
			snap.IdleCount++
		case StatusCharging:
			snap.ChargingCount++
		case StatusError:
			snap.ErrorCount++
		case StatusMaintenance:
			snap.MaintenanceCount++
		}
	}

	if snap.TotalRobots > 0 {
		snap.AvgBattery = math.Round(totalBattery/float64(snap.TotalRobots)*10) / 10
		snap.AvgTaskRate = math.Round(totalTaskRate/float64(snap.TotalRobots)*10) / 10
		snap.AvgCPUTemp = math.Round(totalCPUTemp/float64(snap.TotalRobots)*10) / 10
	}

	// Volatile metrics with wave patterns for organic oscillation
	elapsed := now.Sub(startTime).Seconds()

	// Events per sec: base from tick counters + sinusoidal waves + noise
	baseEvents := float64(tickStatusChanges+tickTasksDone) * 5.0 // scale up since tick is 200ms
	wave1 := 8.0 * math.Sin(elapsed/6.0*2*math.Pi)
	wave2 := 5.0 * math.Sin(elapsed/10.0*2*math.Pi)
	wave3 := 3.0 * math.Sin(elapsed/26.0*2*math.Pi)
	noise := (rand.Float64() - 0.5) * 10.0
	snap.EventsPerSec = math.Round((math.Max(0, baseEvents+wave1+wave2+wave3+noise+20.0))*10) / 10

	// Tasks completed: tick counter + wave + noise
	taskWave := 4.0*math.Sin(elapsed/8.0*2*math.Pi) + 3.0*math.Sin(elapsed/15.0*2*math.Pi)
	taskNoise := (rand.Float64() - 0.5) * 6.0
	snap.TasksCompleted = int(math.Max(0, float64(tickTasksDone)*5.0+taskWave+taskNoise+12.0))

	// Alerts active: real count + wave for drama
	alertWave := 3.0*math.Sin(elapsed/10.0*2*math.Pi) + 2.0*math.Sin(elapsed/22.0*2*math.Pi)
	alertNoise := (rand.Float64() - 0.5) * 4.0
	snap.AlertsActive = int(math.Max(0, float64(alertsActive)+alertWave+alertNoise))

	// Network latency: synthetic volatile metric with built-in oscillation
	latBase := 45.0
	latWave := 15.0*math.Sin(elapsed/6.0*2*math.Pi) + 10.0*math.Sin(elapsed/14.0*2*math.Pi) + 8.0*math.Sin(elapsed/26.0*2*math.Pi)
	latNoise := (rand.Float64() - 0.5) * 20.0
	snap.NetworkLatency = math.Round(math.Max(5, latBase+latWave+latNoise)*10) / 10

	eventMu.RLock()
	if len(events) > 20 {
		snap.Events = events[:20]
	} else {
		snap.Events = events
	}
	eventMu.RUnlock()

	return snap
}

func fleetStreamHandler(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			mu.RLock()
			snap := computeFleetSnapshot()
			mu.RUnlock()

			data, _ := json.Marshal(snap)
			fmt.Fprintf(w, "data: %s\n\n", data)
			flusher.Flush()
		}
	}
}

func categoryStreamHandler(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	// URL decode: replace hyphens/underscores with spaces for matching
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")

	// Find matching category (case-insensitive)
	var found bool
	for _, cat := range categories {
		if strings.EqualFold(cat.Name, name) {
			name = cat.Name
			found = true
			break
		}
	}
	if !found {
		http.Error(w, "category not found", http.StatusNotFound)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ticker.C:
			mu.RLock()
			summaries := computeCategorySummaries()
			mu.RUnlock()

			for _, s := range summaries {
				if s.Name == name {
					data, _ := json.Marshal(s)
					fmt.Fprintf(w, "data: %s\n\n", data)
					flusher.Flush()
					break
				}
			}
		}
	}
}

func main() {
	startTime = time.Now()
	generateFleet()

	// Start simulation (200ms ticks for high-frequency data)
	go func() {
		ticker := time.NewTicker(200 * time.Millisecond)
		defer ticker.Stop()
		for range ticker.C {
			simulationTick()
		}
	}()

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173", "http://localhost:*", "http://web:5173"},
		AllowedMethods:   []string{"GET", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/api/health", healthHandler)
	r.Get("/api/robots", robotsHandler)
	r.Get("/api/robots/{id}", robotByIDHandler)
	r.Get("/api/categories", categoriesHandler)
	r.Get("/api/stream/fleet", fleetStreamHandler)
	r.Get("/api/stream/category/{name}", categoryStreamHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Robot HQ API starting on :%s with %d robots", port, len(robots))
	log.Fatal(http.ListenAndServe(":"+port, r))
}
