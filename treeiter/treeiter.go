//go:build !solution

package treeiter

func DoInOrder[E interface {
	Left() *E
	Right() *E
}](root *E, action func(t *E)) {
	if root == nil {
		return
	}
	DoInOrder((*root).Left(), action)
	action(root)
	DoInOrder((*root).Right(), action)
}
