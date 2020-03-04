package main

func main() {
	test("wow","awesome")
}

func test(str ...string) {
	st := ""
	for _,v := range str{
		st+=v
	}
	println(st)
}




