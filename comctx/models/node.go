package models

type Node struct {
	Name          string         `json:"name"`
	Type          int64          `json:"type"`
	State         int64          `json:"state"`
	PluginName    string         `json:"plugin_name"`
	Groups        []Group        `gorm:"foreignKey:DriverName;references:Name"`
	Setting       Setting        `gorm:"foreignKey:NodeName;references:Name"`
	Subscriptions []Subscription `gorm:"foreignKey:AppName;references:Name"`
}

func FindGroupByTag(groups []Group, tag *Tag) *Group {
	for _, group := range groups {
		if group.DriverName == tag.DriverName && group.Name == tag.GroupName {
			return &group
		}
	}
	return nil
}
func FindGroupByName(groups []Group, g *Group) *Group {
	for _, group := range groups {
		if group.DriverName == g.DriverName && group.Name == g.OldName {
			return &group
		}
	}
	return nil
}

func (n *Node) AddTags(tags []*Tag) map[string]*Group {
	addedGroups := make(map[string]*Group) // 初始化 map

	for _, newTag := range tags {
		group := FindGroupByTag(n.Groups, newTag)
		if group == nil {
			continue // 没找到 group，跳过，或者可以加日志
		}

		group.Tags = append(group.Tags, *newTag) //  安全地新增
		group.OldName = group.Name
		addedGroups[group.Name] = group
	}

	return addedGroups
}

func (n *Node) UpdateTags(tags []*Tag) map[string]*Group {
	updatedGroups := make(map[string]*Group)

	for _, newTag := range tags {
		group := FindGroupByTag(n.Groups, newTag)
		if group == nil {
			continue // 找不到 Group 跳过
		}

		// 在 group.Tags 里找到这个 tag，并更新
		for i, oldTag := range group.Tags {
			if oldTag.DriverName == newTag.DriverName &&
				oldTag.GroupName == newTag.GroupName &&
				oldTag.Name == newTag.Name {

				// 找到了，更新内容
				group.Tags[i] = *newTag
				break
			}
		}
		group.OldName = group.Name
		// 收集这次有更新的 group
		updatedGroups[group.Name] = group
	}

	return updatedGroups
}

func (n *Node) RemoveTags(removed []*Tag) map[string]*Group {
	updatedGroups := make(map[string]*Group) // 初始化

	for _, tagToRemove := range removed {
		group := FindGroupByTag(n.Groups, tagToRemove)
		if group == nil {
			continue // 找不到 Group，跳过
		}

		newTags := make([]Tag, 0, len(group.Tags))
		for _, existingTag := range group.Tags {
			if existingTag.DriverName == tagToRemove.DriverName &&
				existingTag.GroupName == tagToRemove.GroupName &&
				existingTag.Name == tagToRemove.Name {
				// 匹配到了这个 tag，跳过（也就是删除）
				continue
			}
			newTags = append(newTags, existingTag)
		}

		group.Tags = newTags // 更新 group 的 tags
		group.OldName = group.Name
		updatedGroups[group.Name] = group
	}

	return updatedGroups
}
func (n *Node) UpdatedGroups(groups []*Group) map[string]*Group {
	updatedGroups := make(map[string]*Group) //  初始化 map

	for _, g := range groups {
		group := FindGroupByName(n.Groups, g)
		if group != nil {
			group.Name = g.Name
			group.OldName = g.OldName
			group.Interval = g.Interval
			updatedGroups[g.Name] = group
		} else {
			return nil
		}
	}

	return updatedGroups
}

func (n *Node) AddedGroups(groups []*Group) map[string]*Group {
	addedGroups := make(map[string]*Group)
	for _, g := range groups {
		n.Groups = append(n.Groups, *g)
		addedGroups[g.Name] = g
	}
	return addedGroups
}
