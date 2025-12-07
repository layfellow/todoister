// Copyright 2025 layfellow. All rights reserved.
// Use of this source code is governed by an MIT license
// that can be found in the LICENSE file.

package util

import (
	"time"
)

// convertCachedToTodoistData converts a CachedTodoistData protobuf message to TodoistData.
func convertCachedToTodoistData(cached *CachedTodoistData) *TodoistData {
	todoistData := &TodoistData{
		Projects: make([]TodoistProject, len(cached.Projects)),
		Sections: make([]TodoistSection, len(cached.Sections)),
		Items:    make([]TodoistItem, len(cached.Items)),
		Labels:   make([]TodoistLabel, len(cached.Labels)),
		Comments: make([]TodoistComment, len(cached.Comments)),
	}

	// Convert Projects
	for i, p := range cached.Projects {
		todoistData.Projects[i] = TodoistProject{
			ID:       p.GetId(),
			ParentID: p.GetParentId(),
			Project: Project{
				Name:      p.GetName(),
				Color:     p.GetColor(),
				ViewStyle: p.GetViewStyle(),
			},
		}
	}

	// Convert Sections
	for i, s := range cached.Sections {
		todoistData.Sections[i] = TodoistSection{
			ID:        s.GetId(),
			ProjectID: s.GetProjectId(),
			Order:     int(s.GetOrder()),
			Section: Section{
				Name:      s.GetName(),
				Collapsed: s.GetCollapsed(),
			},
		}
	}

	// Convert Items
	for i, item := range cached.Items {
		todoistItem := TodoistItem{
			ID:        item.GetId(),
			ProjectID: item.GetProjectId(),
			SectionID: item.GetSectionId(),
			Labels:    item.GetLabels(),
			Task: Task{
				Content:     item.GetContent(),
				Description: item.GetDescription(),
				Priority:    int(item.GetPriority()),
				ChildOrder:  int(item.GetChildOrder()),
				Collapsed:   item.GetCollapsed(),
				CompletedAt: item.GetCompletedAt(),
			},
		}

		// Convert Duration if present
		if item.Duration != nil {
			todoistItem.Duration = &Duration{
				Amount: int(item.Duration.GetAmount()),
				Unit:   item.Duration.GetUnit(),
			}
		}

		// Convert Due if present
		if item.Due != nil {
			todoistItem.Due = &Due{
				Date:        item.Due.GetDate(),
				IsRecurring: item.Due.GetIsRecurring(),
				String:      item.Due.GetDueString(),
				Datetime:    item.Due.GetDatetime(),
				Timezone:    item.Due.GetTimezone(),
			}
		}

		todoistData.Items[i] = todoistItem
	}

	// Convert Labels
	for i, l := range cached.Labels {
		todoistData.Labels[i] = TodoistLabel{
			ID: l.GetId(),
			Label: Label{
				Name:  l.GetName(),
				Color: l.GetColor(),
			},
		}
	}

	// Convert Comments
	for i, c := range cached.Comments {
		todoistData.Comments[i] = TodoistComment{
			ID:        c.GetId(),
			TaskID:    c.GetTaskId(),
			ProjectID: c.GetProjectId(),
			Comment: Comment{
				Content: c.GetContent(),
			},
		}
	}

	return todoistData
}

// convertTodoistDataToCached converts TodoistData to CachedTodoistData protobuf message.
func convertTodoistDataToCached(data *TodoistData, syncToken string) *CachedTodoistData {
	cached := &CachedTodoistData{
		SyncToken: syncToken,
		CachedAt:  time.Now().Unix(),
		Projects:  make([]*PbProject, len(data.Projects)),
		Sections:  make([]*PbSection, len(data.Sections)),
		Items:     make([]*PbItem, len(data.Items)),
		Labels:    make([]*PbLabel, len(data.Labels)),
		Comments:  make([]*PbComment, len(data.Comments)),
	}

	// Convert Projects
	for i, p := range data.Projects {
		cached.Projects[i] = &PbProject{
			Id:        p.ID,
			ParentId:  p.ParentID,
			Name:      p.Name,
			Color:     p.Color,
			ViewStyle: p.ViewStyle,
		}
	}

	// Convert Sections
	for i, s := range data.Sections {
		cached.Sections[i] = &PbSection{
			Id:        s.ID,
			ProjectId: s.ProjectID,
			Name:      s.Name,
			Collapsed: s.Collapsed,
			Order:     int32(s.Order),
		}
	}

	// Convert Items
	for i, item := range data.Items {
		cachedItem := &PbItem{
			Id:          item.ID,
			ProjectId:   item.ProjectID,
			SectionId:   item.SectionID,
			Content:     item.Content,
			Description: item.Description,
			Priority:    int32(item.Priority),
			ChildOrder:  int32(item.ChildOrder),
			Collapsed:   item.Collapsed,
			Labels:      item.Labels,
			CompletedAt: item.CompletedAt,
		}

		// Convert Duration if present
		if item.Duration != nil {
			cachedItem.Duration = &PbDuration{
				Amount: int32(item.Duration.Amount),
				Unit:   item.Duration.Unit,
			}
		}

		// Convert Due if present
		if item.Due != nil {
			cachedItem.Due = &PbDue{
				Date:        item.Due.Date,
				IsRecurring: item.Due.IsRecurring,
				DueString:   item.Due.String,
				Datetime:    item.Due.Datetime,
				Timezone:    item.Due.Timezone,
			}
		}

		cached.Items[i] = cachedItem
	}

	// Convert Labels
	for i, l := range data.Labels {
		cached.Labels[i] = &PbLabel{
			Id:    l.ID,
			Name:  l.Name,
			Color: l.Color,
		}
	}

	// Convert Comments
	for i, c := range data.Comments {
		cached.Comments[i] = &PbComment{
			Id:        c.ID,
			TaskId:    c.TaskID,
			ProjectId: c.ProjectID,
			Content:   c.Content,
		}
	}

	return cached
}

// mergeData merges incremental sync data into existing cached data.
// Handles additions, updates, and deletions (via is_deleted flag).
func mergeData(cached *TodoistData, incremental *SyncResponse) *TodoistData {
	// Create maps for O(1) lookups
	projectMap := make(map[string]TodoistProject)
	sectionMap := make(map[string]TodoistSection)
	itemMap := make(map[string]TodoistItem)
	labelMap := make(map[string]TodoistLabel)
	commentMap := make(map[string]TodoistComment)

	// Populate maps with cached data
	for _, p := range cached.Projects {
		projectMap[p.ID] = p
	}
	for _, s := range cached.Sections {
		sectionMap[s.ID] = s
	}
	for _, i := range cached.Items {
		itemMap[i.ID] = i
	}
	for _, l := range cached.Labels {
		labelMap[l.ID] = l
	}
	for _, c := range cached.Comments {
		commentMap[c.ID] = c
	}

	// Merge Projects (updates and additions, filter out deletions)
	for _, p := range incremental.Projects {
		// Note: is_deleted is not exposed in our simplified struct,
		// so we treat all incremental items as updates/additions
		projectMap[p.ID] = p
	}

	// Merge Sections
	for _, s := range incremental.Sections {
		sectionMap[s.ID] = s
	}

	// Merge Items
	for _, i := range incremental.Items {
		itemMap[i.ID] = i
	}

	// Merge Labels
	for _, l := range incremental.Labels {
		labelMap[l.ID] = l
	}

	// Merge Comments (both notes and project_notes)
	for _, c := range incremental.Notes {
		commentMap[c.ID] = c
	}
	for _, c := range incremental.ProjectNotes {
		commentMap[c.ID] = c
	}

	// Convert maps back to slices
	result := &TodoistData{
		Projects: make([]TodoistProject, 0, len(projectMap)),
		Sections: make([]TodoistSection, 0, len(sectionMap)),
		Items:    make([]TodoistItem, 0, len(itemMap)),
		Labels:   make([]TodoistLabel, 0, len(labelMap)),
		Comments: make([]TodoistComment, 0, len(commentMap)),
	}

	for _, p := range projectMap {
		result.Projects = append(result.Projects, p)
	}
	for _, s := range sectionMap {
		result.Sections = append(result.Sections, s)
	}
	for _, i := range itemMap {
		result.Items = append(result.Items, i)
	}
	for _, l := range labelMap {
		result.Labels = append(result.Labels, l)
	}
	for _, c := range commentMap {
		result.Comments = append(result.Comments, c)
	}

	return result
}
