package main

// Zingage Attendance Tracking System
// Recording different types of shifts and bonuses:
// 1. Short Shift Bonus (20 points)
// 2. Weekend Warriors Consistency Reward (500 points)
// 3. Last Minute Open Shift (75 points)

import (
    "time"
    "net/http"
    "github.com/gorilla/mux"
    "log"
    "encoding/json"
    "os"
    "sync"
)

type Shift struct {
    WorkerID string `json:"worker_id"`
    ShiftType string `json:"shift_type"` // "short", "weekend", "lastMinute"
    Points int `json:"points"`
    Timestamp time.Time `json:"timestamp"`
}

type ShiftManager struct {
    mu sync.Mutex
    ShortShifts []Shift     `json:"short_shifts"`
    WeekendShifts []Shift   `json:"weekend_shifts"`
    LastMinuteShifts []Shift `json:"last_minute_shifts"`
}

func NewShiftManager() *ShiftManager {
    return &ShiftManager{
        ShortShifts: make([]Shift, 0),
        WeekendShifts: make([]Shift, 0),
        LastMinuteShifts: make([]Shift, 0),
    }
}

func (sm *ShiftManager) loadShifts(filename string) error {
    data, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    sm.mu.Lock()
    defer sm.mu.Unlock()
    return json.Unmarshal(data, sm)
}

func (sm *ShiftManager) saveShifts(filename string) error {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    data, err := json.Marshal(sm)
    if err != nil {
        return err
    }
    
    return os.WriteFile(filename, data, 0644)
}

func (sm *ShiftManager) logShift(w http.ResponseWriter, r *http.Request) {
    var shift Shift
    err := json.NewDecoder(r.Body).Decode(&shift)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    shift.Timestamp = time.Now()
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Assign points and store in appropriate slice based on shift type
    switch shift.ShiftType {
    case "short":
        shift.Points = 20
        sm.ShortShifts = append(sm.ShortShifts, shift)
    case "weekend":
        shift.Points = 500
        sm.WeekendShifts = append(sm.WeekendShifts, shift)
    case "lastMinute":
        shift.Points = 75
        sm.LastMinuteShifts = append(sm.LastMinuteShifts, shift)
    default:
        http.Error(w, "Invalid shift type", http.StatusBadRequest)
        return
    }
    
    // Save shifts after each update
    go sm.saveShifts("shifts.json")
    
    w.WriteHeader(http.StatusCreated)
}

func (sm *ShiftManager) getWeekendWarriors(w http.ResponseWriter, r *http.Request) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Group shifts by worker
    workerShifts := make(map[string]int)
    for _, shift := range sm.WeekendShifts {
        workerShifts[shift.WorkerID]++
    }
    
    // Find workers with consistent weekend shifts
    warriors := make([]string, 0)
    for workerID, count := range workerShifts {
        if count >= 2 { // Define your consistency criteria here
            warriors = append(warriors, workerID)
        }
    }
    
    json.NewEncoder(w).Encode(warriors)
}

func (sm *ShiftManager) getLastMinuteHeroes(w http.ResponseWriter, r *http.Request) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Get unique workers who took last-minute shifts
    heroes := make(map[string]bool)
    for _, shift := range sm.LastMinuteShifts {
        heroes[shift.WorkerID] = true
    }
    
    // Convert to slice
    heroList := make([]string, 0, len(heroes))
    for workerID := range heroes {
        heroList = append(heroList, workerID)
    }
    
    json.NewEncoder(w).Encode(heroList)
}

func (sm *ShiftManager) getShortShiftWorkers(w http.ResponseWriter, r *http.Request) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    // Get unique workers who took short shifts
    workers := make(map[string]bool)
    for _, shift := range sm.ShortShifts {
        workers[shift.WorkerID] = true
    }
    
    // Convert to slice
    workerList := make([]string, 0, len(workers))
    for workerID := range workers {
        workerList = append(workerList, workerID)
    }
    
    json.NewEncoder(w).Encode(workerList)
}

func (sm *ShiftManager) getWorkerStats(w http.ResponseWriter, r *http.Request) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    
    vars := mux.Vars(r)
    workerID := vars["id"]
    
    stats := struct {
        TotalPoints int `json:"total_points"`
        ShortShifts int `json:"short_shifts"`
        WeekendShifts int `json:"weekend_shifts"`
        LastMinuteShifts int `json:"last_minute_shifts"`
    }{}
    
    // Calculate stats
    for _, shift := range sm.ShortShifts {
        if shift.WorkerID == workerID {
            stats.ShortShifts++
            stats.TotalPoints += shift.Points
        }
    }
    
    for _, shift := range sm.WeekendShifts {
        if shift.WorkerID == workerID {
            stats.WeekendShifts++
            stats.TotalPoints += shift.Points
        }
    }
    
    for _, shift := range sm.LastMinuteShifts {
        if shift.WorkerID == workerID {
            stats.LastMinuteShifts++
            stats.TotalPoints += shift.Points
        }
    }
    
    json.NewEncoder(w).Encode(stats)
}

func main() {
    router := mux.NewRouter()
    shiftManager := NewShiftManager()
    
    // Load existing shifts if any
    err := shiftManager.loadShifts("shifts.json")
    if err != nil {
        log.Printf("No existing shifts found or error loading: %v", err)
    }
    
    // Routes
    router.HandleFunc("/shift", shiftManager.logShift).Methods("POST")
    router.HandleFunc("/weekend-warriors", shiftManager.getWeekendWarriors).Methods("GET")
    router.HandleFunc("/last-minute-heroes", shiftManager.getLastMinuteHeroes).Methods("GET")
    router.HandleFunc("/short-shift-workers", shiftManager.getShortShiftWorkers).Methods("GET")
    router.HandleFunc("/worker/{id}/stats", shiftManager.getWorkerStats).Methods("GET")
    
    log.Println("Server starting on port 8080...")
    log.Fatal(http.ListenAndServe(":8080", router))
}
