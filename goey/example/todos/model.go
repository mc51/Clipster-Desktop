package main

type TodoItem struct {
	Text      string
	Completed bool
}

var (
	Model []TodoItem
)

func getItemCounts() (int, int) {
	completedCount := 0
	for _, v := range Model {
		if v.Completed {
			completedCount++
		}
	}
	return completedCount, len(Model) - completedCount
}
