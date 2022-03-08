package main

import "testing"

func Test_Main(t *testing.T) {
	cmder := buildCmd("echo", "hello", "world")

	cmder.Exec()

	main()
}
