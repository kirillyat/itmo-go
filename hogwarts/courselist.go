//go:build !solution

package hogwarts

func topoSort(sortedCourses *[]string, course string, prerequisites map[string][]string,
	visited map[string]bool, inCurrentPath map[string]bool) {
	if inCurrentPath[course] {
		panic("detected a cycle in prerequisites")
	}
	if visited[course] {
		return
	}

	inCurrentPath[course] = true

	for _, prereq := range prerequisites[course] {
		topoSort(sortedCourses, prereq, prerequisites, visited, inCurrentPath)
	}

	visited[course] = true
	inCurrentPath[course] = false

	*sortedCourses = append(*sortedCourses, course)
}

func GetCourseList(prerequisites map[string][]string) []string {
	visited := make(map[string]bool)
	inCurrentPath := make(map[string]bool)
	var sortedCourses []string

	for course := range prerequisites {
		topoSort(&sortedCourses, course, prerequisites, visited, inCurrentPath)
	}

	return sortedCourses
}
