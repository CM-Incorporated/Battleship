package main

import (
	//"net/http"
	"fmt"

	"cmincorporated.com/gamecode"
)

func main() {
	gamecode.MakeWork()
	isWorking := gamecode.IsWorking()
	fmt.Println(isWorking)
}

/*
	fs:=http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)
	http.ListenAndServe(":8080", nil)
*/
