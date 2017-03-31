package web

import "log"

func partPath() {
	p := "/abc/cde/fgh/ijk/"
	start, l := 0, len(p)
	for {
		s, e := partIn(p, start)
		log.Println(p[s:e])
		if l == e {
			break
		}
		start = e - 1
	}
	log.Println(p)
}
