package main

// predefinedBadges contains the badge definitions and criteria for the standard badges
var predefinedBadges = map[string]NewBadgeRequest{
	"early-bird": {
		Name:        "Early Bird",
		Description: "Checked in before 9 AM for 5 consecutive days.",
		ImageURL:    "https://example.com/badges/early-bird.png",
		FlowDefinition: map[string]interface{}{
			"event": "check-in",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType":      "day",
					"count":           map[string]interface{}{"$gte": 5},
					"excludeWeekends": false,
					"excludeHolidays": false,
				},
				"payload": map[string]interface{}{
					"time": map[string]interface{}{"$lt": "09:00:00"},
				},
			},
		},
	},
	"workaholic": {
		Name:        "Workaholic",
		Description: "Logged 40+ hours in a week.",
		ImageURL:    "https://example.com/badges/workaholic.png",
		FlowDefinition: map[string]interface{}{
			"event": "work-log",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "week",
					"count":      map[string]interface{}{"$gte": 1},
				},
				"$aggregate": map[string]interface{}{
					"function": "sum",
					"property": "hours",
					"result":   map[string]interface{}{"$gte": 40},
				},
			},
		},
	},
	"consistency-king": {
		Name:        "Consistency King",
		Description: "Checked in & out without missing a day for a month.",
		ImageURL:    "https://example.com/badges/consistency-king.png",
		FlowDefinition: map[string]interface{}{
			"$and": []map[string]interface{}{
				{
					"event": "check-in",
					"criteria": map[string]interface{}{
						"$timePeriod": map[string]interface{}{
							"periodType":      "day",
							"count":           map[string]interface{}{"$gte": 20},
							"excludeWeekends": true,
							"excludeHolidays": true,
						},
					},
				},
				{
					"event": "check-out",
					"criteria": map[string]interface{}{
						"$timePeriod": map[string]interface{}{
							"periodType":      "day",
							"count":           map[string]interface{}{"$gte": 20},
							"excludeWeekends": true,
							"excludeHolidays": true,
						},
					},
				},
			},
		},
	},
	"team-player": {
		Name:        "Team Player",
		Description: "Collaborated on 5+ team projects in a month.",
		ImageURL:    "https://example.com/badges/team-player.png",
		FlowDefinition: map[string]interface{}{
			"event": "project-collaboration",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "month",
					"count":      map[string]interface{}{"$gte": 1},
				},
				"count": map[string]interface{}{"$gte": 5},
			},
		},
	},
	"overtime-warrior": {
		Name:        "Overtime Warrior",
		Description: "Logged extra 10 hours in a week.",
		ImageURL:    "https://example.com/badges/overtime-warrior.png",
		FlowDefinition: map[string]interface{}{
			"event": "work-log",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "week",
					"count":      map[string]interface{}{"$gte": 1},
				},
				"payload": map[string]interface{}{
					"is_overtime": true,
				},
				"$aggregate": map[string]interface{}{
					"function": "sum",
					"property": "hours",
					"result":   map[string]interface{}{"$gte": 10},
				},
			},
		},
	},
	"meeting-maestro": {
		Name:        "Meeting Maestro",
		Description: "Attended 10+ meetings in a month.",
		ImageURL:    "https://example.com/badges/meeting-maestro.png",
		FlowDefinition: map[string]interface{}{
			"event": "meeting-attendance",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "month",
					"count":      map[string]interface{}{"$gte": 1},
				},
				"count": map[string]interface{}{"$gte": 10},
			},
		},
	},
	"bug-hunter": {
		Name:        "Bug Hunter",
		Description: "Reported 5+ issues that got fixed.",
		ImageURL:    "https://example.com/badges/bug-hunter.png",
		FlowDefinition: map[string]interface{}{
			"event": "bug-report",
			"criteria": map[string]interface{}{
				"payload": map[string]interface{}{
					"status": "fixed",
				},
				"count": map[string]interface{}{"$gte": 5},
			},
		},
	},
	"task-master": {
		Name:        "Task Master",
		Description: "Completed 50+ tasks in a month.",
		ImageURL:    "https://example.com/badges/task-master.png",
		FlowDefinition: map[string]interface{}{
			"event": "task-completion",
			"criteria": map[string]interface{}{
				"$timePeriod": map[string]interface{}{
					"periodType": "month",
					"count":      map[string]interface{}{"$gte": 1},
				},
				"count": map[string]interface{}{"$gte": 50},
			},
		},
	},
}
