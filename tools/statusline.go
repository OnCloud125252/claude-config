// Package main provides a status line tool for Claude Code CLI.
// It displays model info, git branch, context usage, and session time.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Constants
const (
	// ANSI color definitions - P10k Rainbow theme inspired
	ColorReset  = "\033[0m"
	ColorGold   = "\033[38;2;214;196;161m" // Warm gold
	ColorCyan   = "\033[38;2;122;162;247m" // Bright cyan-blue
	ColorPink   = "\033[38;2;247;118;142m" // Soft pink
	ColorGreen  = "\033[38;2;158;206;106m" // Bright green
	ColorGray   = "\033[38;2;86;95;137m"   // Muted gray-blue
	ColorSilver = "\033[38;2;192;202;245m" // Light blue-gray
	ColorOrange = "\033[38;2;255;158;100m" // Warm orange
	ColorPurple = "\033[38;2;187;154;247m" // Soft purple
	ColorYellow = "\033[38;2;224;175;104m" // Warm yellow

	// Frame colors matching P10k style
	ColorFrame   = "\033[38;2;86;95;137m"   // Frame connectors
	ColorBracket = "\033[38;2;122;162;247m" // Brackets and separators

	// Context colors with rainbow theme
	ColorCtxGreen = "\033[38;2;158;206;106m"
	ColorCtxGold  = "\033[38;2;224;175;104m"
	ColorCtxRed   = "\033[38;2;247;118;142m"

	// Configuration
	SessionTimeoutSeconds = 600 // 10 minutes
	MaxContextTokens      = 200000
	ProgressBarWidth      = 10
	GitBranchCacheSeconds = 5
	MaxUserMessageLines   = 3
	UserMessageLineWidth  = 60
	MaxTranscriptLines    = 100
	MaxUserSearchLines    = 200
	MaxScanTokenSize      = 1024 * 1024 // 1MB for JSON lines
)

// Model configurations with rainbow colors
var modelConfig = map[string][1]string{
	"Opus":   {ColorGold},
	"Sonnet": {ColorCyan},
	"Haiku":  {ColorPink},
	"4":      {ColorPurple}, // For Sonnet 4
}

// Input data structure
type Input struct {
	Model struct {
		DisplayName string `json:"display_name"`
	} `json:"model"`
	SessionID string `json:"session_id"`
	Workspace struct {
		CurrentDir string `json:"current_dir"`
	} `json:"workspace"`
	TranscriptPath string `json:"transcript_path,omitempty"`
}

// Session data structure
type Session struct {
	ID            string     `json:"id"`
	Date          string     `json:"date"`
	Start         int64      `json:"start"`
	LastHeartbeat int64      `json:"last_heartbeat"`
	TotalSeconds  int64      `json:"total_seconds"`
	Intervals     []Interval `json:"intervals"`
}

type Interval struct {
	Start int64  `json:"start"`
	End   *int64 `json:"end"`
}

// Result channel data
type Result struct {
	Type string
	Data any
}

// Simple cache for git branch
var (
	gitBranchCache   string
	gitBranchExpires time.Time
	cacheMutex       sync.RWMutex
)

// Session cache for write debouncing
var (
	sessionCache      = make(map[string]*sessionCacheEntry)
	sessionCacheMutex sync.RWMutex
)

type sessionCacheEntry struct {
	session   Session
	lastWrite time.Time
	dirty     bool
}

const sessionWriteDebounce = 2 * time.Second

func main() {
	var input Input
	if err := json.NewDecoder(os.Stdin).Decode(&input); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to decode input: %v\n", err)
		os.Exit(1)
	}

	// Create result channel
	results := make(chan Result, 4)
	var wg sync.WaitGroup

	// Fetch information in parallel
	wg.Add(4)

	go func() {
		defer wg.Done()
		branch := getGitBranch()
		results <- Result{"git", branch}
	}()

	go func() {
		defer wg.Done()
		totalHours := calculateTotalHours(input.SessionID)
		results <- Result{"hours", totalHours}
	}()

	go func() {
		defer wg.Done()
		contextInfo := analyzeContext(input.TranscriptPath)
		results <- Result{"context", contextInfo}
	}()

	go func() {
		defer wg.Done()
		userMsg := extractUserMessage(input.TranscriptPath, input.SessionID)
		results <- Result{"message", userMsg}
	}()

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var gitBranch, totalHours, contextUsage, userMessage string

	for result := range results {
		switch result.Type {
		case "git":
			gitBranch = result.Data.(string)
		case "hours":
			totalHours = result.Data.(string)
		case "context":
			contextUsage = result.Data.(string)
		case "message":
			userMessage = result.Data.(string)
		}
	}

	// Update session (synchronous to avoid race conditions)
	updateSession(input.SessionID)

	// Get display values (without colors)
	modelName := input.Model.DisplayName
	modelColor := getModelColor(modelName)
	projectName := filepath.Base(input.Workspace.CurrentDir)

	// Output user message with frame continuation
	if userMessage != "" {
		fmt.Printf("%s%s", ColorReset, userMessage)
	}

	// Output status line with all colors applied here
	// First line: model, project, git branch
	fmt.Printf("╭─%s[%s%s%s]  %s %s%s  %s %s%s\n",
		ColorReset, modelColor, modelName, ColorReset,
		ColorSilver, projectName, ColorReset,
		ColorYellow, gitBranch, ColorReset)

	// Second line: context usage and total hours
	fmt.Printf("╰─%s │ %s%s%s\n",
		contextUsage,
		ColorGreen, totalHours, ColorReset)
}

// Get model color based on model name
func getModelColor(model string) string {
	for key, config := range modelConfig {
		if strings.Contains(model, key) {
			return config[0]
		}
	}
	return ColorReset
}

// Get git branch with caching (optimized single command)
func getGitBranch() string {
	cacheMutex.RLock()
	if time.Now().Before(gitBranchExpires) && gitBranchCache != "" {
		result := gitBranchCache
		cacheMutex.RUnlock()
		return result
	}
	cacheMutex.RUnlock()

	// Single command - if not a git repo or other error, this will fail gracefully
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	branch := strings.TrimSpace(string(output))
	if branch == "" {
		return ""
	}

	// Update cache
	cacheMutex.Lock()
	gitBranchCache = branch
	gitBranchExpires = time.Now().Add(GitBranchCacheSeconds * time.Second)
	cacheMutex.Unlock()

	return branch
}

// Update session with heartbeat (optimized with write debouncing)
func updateSession(sessionID string) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	if err := os.MkdirAll(sessionsDir, 0755); err != nil {
		return
	}

	sessionFile := filepath.Join(sessionsDir, sessionID+".json")
	currentTime := time.Now().Unix()
	today := time.Now().Format("2006-01-02")

	sessionCacheMutex.Lock()
	defer sessionCacheMutex.Unlock()

	entry, exists := sessionCache[sessionID]

	// Load from disk if not in cache
	if !exists {
		var session Session
		if data, err := os.ReadFile(sessionFile); err == nil {
			if err := json.Unmarshal(data, &session); err != nil {
				session = createNewSession(sessionID, today, currentTime)
			}
		} else {
			session = createNewSession(sessionID, today, currentTime)
		}

		entry = &sessionCacheEntry{
			session:   session,
			lastWrite: time.Now(),
			dirty:     false,
		}
		sessionCache[sessionID] = entry
	}

	session := &entry.session

	// Update heartbeat
	gap := currentTime - session.LastHeartbeat
	session.LastHeartbeat = currentTime

	if gap < SessionTimeoutSeconds {
		// Extend current interval
		if len(session.Intervals) > 0 {
			session.Intervals[len(session.Intervals)-1].End = &currentTime
		}
	} else {
		// Add new interval
		session.Intervals = append(session.Intervals, Interval{
			Start: currentTime,
			End:   &currentTime,
		})
	}

	// Calculate total time
	var total int64
	for _, interval := range session.Intervals {
		if interval.End != nil {
			total += *interval.End - interval.Start
		}
	}
	session.TotalSeconds = total

	entry.dirty = true

	// Only write if debounce time has passed
	if time.Since(entry.lastWrite) >= sessionWriteDebounce {
		writeSessionToDisk(sessionFile, session)
		entry.lastWrite = time.Now()
		entry.dirty = false
	}
}

// Helper function to write session to disk
func writeSessionToDisk(path string, session *Session) {
	if data, err := json.Marshal(session); err == nil {
		if err := os.WriteFile(path, data, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to save session: %v\n", err)
		}
	}
}

// Calculate total hours for today
func calculateTotalHours(_ string) string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "0m"
	}

	sessionsDir := filepath.Join(homeDir, ".claude", "session-tracker", "sessions")
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		return "0m"
	}

	var totalSeconds int64
	activeSessions := 0
	today := time.Now().Format("2006-01-02")
	currentTime := time.Now().Unix()

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		sessionFile := filepath.Join(sessionsDir, entry.Name())
		data, err := os.ReadFile(sessionFile)
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		// Only count today's sessions
		if session.Date == today {
			totalSeconds += session.TotalSeconds

			// Check if active
			if currentTime-session.LastHeartbeat < SessionTimeoutSeconds {
				activeSessions++
			}
		}
	}

	// Format output
	hours := totalSeconds / 3600
	minutes := (totalSeconds % 3600) / 60

	var timeStr string
	if hours > 0 {
		timeStr = fmt.Sprintf("%dh", hours)
		if minutes > 0 {
			timeStr += fmt.Sprintf("%dm", minutes)
		}
	} else {
		timeStr = fmt.Sprintf("%dm", minutes)
	}

	if activeSessions > 1 {
		return fmt.Sprintf("%s [%d sessions]", timeStr, activeSessions)
	}
	return timeStr
}

// Analyze context usage
func analyzeContext(transcriptPath string) string {
	var contextLength int

	if transcriptPath == "" {
		// When transcriptPath is empty (conversation just started), show initial state
		contextLength = 0
	} else {
		contextLength = calculateContextUsage(transcriptPath)
	}

	// Always show progress bar even when contextLength is 0

	// Calculate percentage
	percentage := int(float64(contextLength) * 100.0 / float64(MaxContextTokens))
	if percentage > 100 {
		percentage = 100
	}

	// Generate progress bar
	progressBar := generateProgressBar(percentage)
	formattedNum := formatNumber(contextLength)
	color := getContextColor(percentage)

	return fmt.Sprintf("%s%s%s %s%d%% (%s)%s",
		color, progressBar, ColorReset, color, percentage, formattedNum, ColorReset)
}

// Calculate context usage from transcript
func calculateContextUsage(transcriptPath string) int {
	lines, err := readLastLines(transcriptPath, MaxTranscriptLines)
	if err != nil {
		return 0
	}

	// Analyze from last to first
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]

		// Skip empty lines
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Try to parse JSON
		var data map[string]any
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		// Check if this is a side-chain message (agent/tool output)
		if sidechain, ok := data["isSidechain"]; ok {
			// Skip side-chain messages
			if isSide, ok := sidechain.(bool); ok && isSide {
				continue
			}
		}

		// Check and extract usage data
		if message, ok := data["message"].(map[string]any); ok {
			if usage, ok := message["usage"].(map[string]any); ok {
				var total float64

				// Calculate all token types
				if input, ok := usage["input_tokens"].(float64); ok {
					total += input
				}
				if cacheRead, ok := usage["cache_read_input_tokens"].(float64); ok {
					total += cacheRead
				}
				if cacheCreation, ok := usage["cache_creation_input_tokens"].(float64); ok {
					total += cacheCreation
				}

				// Return immediately if valid token count found
				if total > 0 {
					return int(total)
				}
			}
		}
	}

	return 0
}

// Generate progress bar visualization
func generateProgressBar(percentage int) string {
	width := ProgressBarWidth
	filled := percentage * width / 100
	filled = min(filled, width)

	empty := width - filled
	color := getContextColor(percentage)

	var bar strings.Builder

	// Left bracket - colored as filled if any progress, gray if empty
	if percentage > 0 {
		bar.WriteString(color)
	} else {
		bar.WriteString(ColorGray)
	}
	bar.WriteString("")
	bar.WriteString(ColorReset)

	// Filled portion
	if filled > 0 {
		bar.WriteString(color)
		bar.WriteString(strings.Repeat("█", filled))
		bar.WriteString(ColorReset)
	}

	// Empty portion
	if empty > 0 {
		bar.WriteString(ColorGray)
		bar.WriteString(strings.Repeat("█", empty))
		bar.WriteString(ColorReset)
	}

	// Right bracket - colored as filled if 100%, gray otherwise
	if percentage >= 100 {
		bar.WriteString(color)
	} else {
		bar.WriteString(ColorGray)
	}
	bar.WriteString("")
	bar.WriteString(ColorReset)

	return bar.String()
}

// Get context color based on percentage
func getContextColor(percentage int) string {
	if percentage < 60 {
		return ColorCtxGreen
	} else if percentage < 80 {
		return ColorCtxGold
	}
	return ColorCtxRed
}

// Format number with units (k, M)
func formatNumber(num int) string {
	if num == 0 {
		return "--"
	}

	if num >= 1000000 {
		return fmt.Sprintf("%dM", num/1000000)
	} else if num >= 1000 {
		return fmt.Sprintf("%dk", num/1000)
	}
	return strconv.Itoa(num)
}

// Extract user message from transcript
func extractUserMessage(transcriptPath, sessionID string) string {
	if transcriptPath == "" {
		return ""
	}

	lines, err := readLastLines(transcriptPath, MaxUserSearchLines)
	if err != nil {
		return ""
	}

	// Search for user message from last to first
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]

		if strings.TrimSpace(line) == "" {
			continue
		}

		var data map[string]any
		if err := json.Unmarshal([]byte(line), &data); err != nil {
			continue
		}

		// Check if it's a user message for current session
		isSidechain, _ := data["isSidechain"].(bool) // side-chain messages are from agents/tools
		sessionMatch := false
		if sid, ok := data["sessionId"].(string); ok && sid == sessionID {
			sessionMatch = true
		}

		if !isSidechain && sessionMatch {
			if message, ok := data["message"].(map[string]any); ok {
				role, _ := message["role"].(string)
				msgType, _ := data["type"].(string)

				if role == "user" && msgType == "user" {
					if content, ok := message["content"].(string); ok {
						// Filter system messages
						if isSystemMessage(content) {
							continue
						}

						// Format and return
						return formatUserMessage(content)
					}
				}
			}
		}
	}

	return ""
}

// Check if message is a system message
func isSystemMessage(content string) bool {
	// Filter JSON format
	if strings.HasPrefix(content, "[") && strings.HasSuffix(content, "]") {
		return true
	}
	if strings.HasPrefix(content, "{") && strings.HasSuffix(content, "}") {
		return true
	}

	// Filter XML tags (but NOT command tags - those are valid user inputs)
	xmlTags := []string{
		"<local-command-stdout>",
		"<bash-stdout>", "<bash-stderr>", "<bash-input>",
	}
	for _, tag := range xmlTags {
		if strings.Contains(content, tag) {
			return true
		}
	}

	// Filter Caveat messages
	if strings.HasPrefix(content, "Caveat:") {
		return true
	}

	return false
}

// Format user message for display
// Helper function to create new session
func createNewSession(sessionID, date string, currentTime int64) Session {
	return Session{
		ID:            sessionID,
		Date:          date,
		Start:         currentTime,
		LastHeartbeat: currentTime,
		TotalSeconds:  0,
		Intervals:     []Interval{{Start: currentTime, End: nil}},
	}
}

// Helper function to read lines from file (optimized reverse reading)
func readLastLines(filePath string, maxLines int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get file size
	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileSize := stat.Size()

	// Handle empty file
	if fileSize == 0 {
		return []string{}, nil
	}

	// Read in chunks from end backwards
	const chunkSize = 8192
	var lines []string
	var partial []byte

	for offset := fileSize; offset > 0 && len(lines) < maxLines; {
		readSize := min(int64(chunkSize), offset)
		offset -= readSize

		buf := make([]byte, readSize)
		if _, err := file.ReadAt(buf, offset); err != nil {
			return nil, err
		}

		// Prepend to partial buffer from previous iteration
		if len(partial) > 0 {
			buf = append(buf, partial...)
		}
		partial = nil

		// Find newlines from end to start
		lineEnd := len(buf)
		for i := len(buf) - 1; i >= 0; i-- {
			if buf[i] == '\n' {
				// Extract line (excluding newline)
				if i+1 < lineEnd {
					line := string(buf[i+1 : lineEnd])
					// Prepend to lines (we're reading backwards)
					lines = append([]string{line}, lines...)
					if len(lines) >= maxLines {
						return lines, nil
					}
				}
				lineEnd = i
			}
		}

		// Save remaining bytes for next iteration
		if lineEnd > 0 {
			partial = make([]byte, lineEnd)
			copy(partial, buf[:lineEnd])
		}
	}

	// Handle remaining partial line (first line of file)
	if len(partial) > 0 && len(lines) < maxLines {
		lines = append([]string{string(partial)}, lines...)
	}

	return lines, nil
}

func formatUserMessage(message string) string {
	if message == "" {
		return ""
	}

	// Extract command name if this is a command message
	commandName := extractCommandName(message)
	if commandName != "" {
		// Display the command in purple
		return fmt.Sprintf("%s%s❯ %s%s\n",
			ColorReset, ColorPurple, commandName, ColorReset)
	}

	maxLines := MaxUserMessageLines
	lineWidth := UserMessageLineWidth

	lines := strings.Split(message, "\n")
	var result []string

	// Detect if this is a command (starts with /)
	isCommand := false
	if len(lines) > 0 {
		trimmed := strings.TrimSpace(lines[0])
		isCommand = strings.HasPrefix(trimmed, "/")
	}

	// Choose color based on message type
	promptColor := ColorGreen
	if isCommand {
		promptColor = ColorPurple
	}

	for i, line := range lines {
		if i >= maxLines {
			break
		}

		line = strings.TrimSpace(line)
		if len(line) > lineWidth {
			line = line[:lineWidth-3] + "..."
		}

		result = append(result, fmt.Sprintf("%s%s❯ %s%s",
			ColorReset, promptColor, line, ColorReset))
	}

	if len(lines) > maxLines {
		result = append(result, fmt.Sprintf("%s❯ %s... (%d more lines)%s",
			promptColor, ColorGray, len(lines)-maxLines, ColorReset))
	}

	if len(result) > 0 {
		return strings.Join(result, "\n") + "\n"
	}

	return ""
}

// Extract command name from XML tags
func extractCommandName(content string) string {
	// Look for <command-name>/command-name</command-name>
	start := strings.Index(content, "<command-name>")
	if start == -1 {
		return ""
	}
	start += len("<command-name>")

	end := strings.Index(content[start:], "</command-name>")
	if end == -1 {
		return ""
	}

	return strings.TrimSpace(content[start : start+end])
}
