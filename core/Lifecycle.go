package core

type Lifecycle interface {
	Init()
	BeforeStart()
	Starting()
	AfterStart()
	BeforeDestroy()
	Destroy()
	AfterDestroy()
}
