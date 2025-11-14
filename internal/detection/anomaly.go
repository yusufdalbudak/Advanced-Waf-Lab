package detection

// AnomalyScore tracks the anomaly score for a request
type AnomalyScore struct {
	Total int
	Tags  map[string]int
}

// NewAnomalyScore creates a new AnomalyScore
func NewAnomalyScore() *AnomalyScore {
	return &AnomalyScore{
		Total: 0,
		Tags:  make(map[string]int),
	}
}

// Add adds a score value to the total and specific tags
func (a *AnomalyScore) Add(score int, tags []string) {
	a.Total += score
	for _, tag := range tags {
		a.Tags[tag] += score
	}
}

// GetTagScore returns the score for a specific tag
func (a *AnomalyScore) GetTagScore(tag string) int {
	return a.Tags[tag]
}

