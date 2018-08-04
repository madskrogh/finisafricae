//Package util defines simple utility functions, used frequently throughout the application
package util

//HandleError panics if the passed error isn't nil
func HandleError(err error) {
	if err != nil {
		panic(err)
	}
}
