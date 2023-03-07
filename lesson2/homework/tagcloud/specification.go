package tagcloud

import "sort"

// TagCloud aggregates statistics about used tags
type TagCloud struct {
	// TODO: add fields if necessary
	tags        []TagStat
	existedTags map[string]bool
}

// TagStat represents statistics regarding single tag
type TagStat struct {
	Tag             string
	OccurrenceCount int
}

// New should create a valid TagCloud instance
// TODO: You decide whether this function should return a pointer or a value
func New() TagCloud {
	// TODO: Implement this
	return TagCloud{
		tags:        []TagStat{},
		existedTags: make(map[string]bool),
	}
}

// AddTag should add a tag to the cloud if it wasn't present and increase tag occurrence count
// thread-safety is not needed
// TODO: You decide whether receiver should be a pointer or a value
<<<<<<< HEAD
func (tcloud *TagCloud) AddTag(tag string) {
	if tcloud.existedTags[tag] {
		for idx, tg := range tcloud.tags {
			if tg.Tag == tag {
				tcloud.tags[idx].OccurrenceCount++
				sort.Slice(tcloud.tags, func(i, j int) bool {
					return tcloud.tags[i].OccurrenceCount > tcloud.tags[j].OccurrenceCount
=======
func (tc *TagCloud) AddTag(tag string) {
	if tc.existedTags[tag] {
		for _, tg := range tc.tags {
			if tg.Tag == tag {
				tg.OccurrenceCount++
				sort.Slice(tc.tags, func(i, j int) bool {
					return tc.tags[i].OccurrenceCount > tc.tags[j].OccurrenceCount
>>>>>>> cbcb300 ([hw-1, hw-2])
				})
			}
		}
	} else {
		newTag := TagStat{
			Tag:             tag,
			OccurrenceCount: 1,
		}
<<<<<<< HEAD
		tcloud.tags = append(tcloud.tags, newTag)
		tcloud.existedTags[tag] = true
=======
		tc.tags = append(tc.tags, newTag)
		tc.existedTags[tag] = true
>>>>>>> cbcb300 ([hw-1, hw-2])
	}
	// TODO: Implement this
}

// TopN should return top N most frequent tags ordered in descending order by occurrence count
// if there are multiple tags with the same occurrence count then the order is defined by implementation
// if n is greater that TagCloud size then all elements should be returned
// thread-safety is not needed
// there are no restrictions on time complexity
// TODO: You decide whether receiver should be a pointer or a value
<<<<<<< HEAD
func (tcloud *TagCloud) TopN(n int) []TagStat {
	// TODO: Implement this
	if n > len(tcloud.tags) {
		return tcloud.tags
	} else {
		return tcloud.tags[:n]
	}
=======
func (tc *TagCloud) TopN(n int) []TagStat {
	// TODO: Implement this
	return tc.tags[:n]
>>>>>>> cbcb300 ([hw-1, hw-2])
}
