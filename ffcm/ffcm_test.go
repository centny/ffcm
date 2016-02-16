package main

import (
	"os"
	"testing"
)

func TestDim(t *testing.T) {
	ef = func(int) {
	}
	//
	os.Args = []string{"ffcm", "-d", "100", "100", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "300", "100", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "100", "400", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "210", "500", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "21x0", "500", "200", "300"}
	main()
	//
	os.Args = []string{"ffcm", "-d", "21x0", "500"}
	main()
}
