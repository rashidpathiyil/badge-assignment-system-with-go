package main

// predefinedEventTypes contains predefined event type definitions required for the badges
var predefinedEventTypes = map[string]NewEventTypeRequest{
	"check-in": {
		Name:        "check-in",
		Description: "User check-in event, used for attendance tracking",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"format":      "time",
					"description": "The time of check-in in HH:MM:SS format",
				},
				"date": map[string]interface{}{
					"type":   "string",
					"format": "date",
				},
				"location": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"user_id", "time", "date"},
		},
	},
	"check-out": {
		Name:        "check-out",
		Description: "User check-out event, used for attendance tracking",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"time": map[string]interface{}{
					"type":        "string",
					"format":      "time",
					"description": "The time of check-out in HH:MM:SS format",
				},
				"date": map[string]interface{}{
					"type":   "string",
					"format": "date",
				},
				"location": map[string]interface{}{
					"type": "string",
				},
			},
			"required": []string{"user_id", "time", "date"},
		},
	},
	"work-log": {
		Name:        "work-log",
		Description: "User work log entry for tracking hours worked",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"hours": map[string]interface{}{
					"type":        "number",
					"minimum":     0.25,
					"description": "Number of hours worked",
				},
				"date": map[string]interface{}{
					"type":   "string",
					"format": "date",
				},
				"task_id": map[string]interface{}{
					"type": "string",
				},
				"description": map[string]interface{}{
					"type": "string",
				},
				"is_overtime": map[string]interface{}{
					"type":    "boolean",
					"default": false,
				},
			},
			"required": []string{"user_id", "hours", "date"},
		},
	},
	"meeting-attendance": {
		Name:        "meeting-attendance",
		Description: "User attendance at a meeting",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"meeting_id": map[string]interface{}{
					"type": "string",
				},
				"meeting_title": map[string]interface{}{
					"type": "string",
				},
				"start_time": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
				"end_time": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
				"duration_minutes": map[string]interface{}{
					"type":    "integer",
					"minimum": 1,
				},
			},
			"required": []string{"user_id", "meeting_id", "start_time"},
		},
	},
	"bug-report": {
		Name:        "bug-report",
		Description: "User submitted a bug report",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"bug_id": map[string]interface{}{
					"type": "string",
				},
				"title": map[string]interface{}{
					"type": "string",
				},
				"description": map[string]interface{}{
					"type": "string",
				},
				"severity": map[string]interface{}{
					"type": "string",
					"enum": []string{"low", "medium", "high", "critical"},
				},
				"status": map[string]interface{}{
					"type": "string",
					"enum": []string{"reported", "in_progress", "fixed", "closed", "wont_fix"},
				},
				"reported_at": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
			},
			"required": []string{"user_id", "bug_id", "title", "severity", "status", "reported_at"},
		},
	},
	"task-completion": {
		Name:        "task-completion",
		Description: "User completed a task",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"task_id": map[string]interface{}{
					"type": "string",
				},
				"title": map[string]interface{}{
					"type": "string",
				},
				"completed_at": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
				"priority": map[string]interface{}{
					"type": "string",
					"enum": []string{"low", "medium", "high"},
				},
			},
			"required": []string{"user_id", "task_id", "completed_at"},
		},
	},
	"project-collaboration": {
		Name:        "project-collaboration",
		Description: "User collaborated on a team project",
		Schema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"user_id": map[string]interface{}{
					"type": "string",
				},
				"project_id": map[string]interface{}{
					"type": "string",
				},
				"project_name": map[string]interface{}{
					"type": "string",
				},
				"contribution_type": map[string]interface{}{
					"type": "string",
					"enum": []string{"code", "review", "documentation", "design", "testing", "meeting"},
				},
				"timestamp": map[string]interface{}{
					"type":   "string",
					"format": "date-time",
				},
			},
			"required": []string{"user_id", "project_id", "contribution_type", "timestamp"},
		},
	},
}
