package runtime

type LifeCycleHandler interface {
	Exit() // Expected to be call when system_exit object is return when eval
}
