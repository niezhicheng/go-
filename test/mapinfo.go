package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main()  {
	rand.Seed(time.Now().UnixNano())
	var scoreMap = make(map[string]int,200)
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("stud%02d",i)
		value := rand.Intn(100)
		scoreMap[key] = value
	}

	var key = make([]string,0,200)
	for s, _ := range scoreMap {
		key = append(key,s)
	}

	sort.Strings(key)
	for i, s := range key {
		fmt.Println(i,s)
	}

}
