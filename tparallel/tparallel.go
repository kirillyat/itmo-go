//go:build !solution

package tparallel

type T struct {
	finished chan struct{}
	wg       chan struct{}
	parallel bool
	parent   *T
	subTests []*T
}

func newT(parent *T) *T {
	return &T{
		finished: make(chan struct{}),
		wg:       make(chan struct{}),
		parent:   parent,
		subTests: make([]*T, 0),
	}
}

func (t *T) Parallel() {
	if t.parallel {
		panic("This test is already parallel")
	}
	t.parallel = true
	t.parent.subTests = append(t.parent.subTests, t)

	t.finished <- struct{}{}
	<-t.parent.wg
}

func (t *T) tRunner(subtest func(t *T)) {
	subtest(t)
	if len(t.subTests) > 0 {
		close(t.wg)

		for _, sub := range t.subTests {
			<-sub.finished
		}
	}
	if t.parallel {
		t.parent.finished <- struct{}{}
	}
	t.finished <- struct{}{}

}

func (t *T) Run(subtest func(t *T)) {
	subT := newT(t)
	go subT.tRunner(subtest)
	<-subT.finished
}

func Run(topTests []func(t *T)) {
	root := newT(nil)
	for _, fn := range topTests {
		root.Run(fn)
	}
	close(root.wg)

	if len(root.subTests) > 0 {
		<-root.finished
	}
}
