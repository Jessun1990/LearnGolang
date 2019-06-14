package chapter4

import (
	"fmt"
	"net/http"
)

func errHandleExample() {
	type cumsomRes struct {
		Error    error
		Response *http.Response
	}

	checkStatus := func(done <-chan interface{},
		urls ...string) <-chan cumsomRes {
		results := make(chan cumsomRes)
		go func() {
			defer close(results)

			for _, url := range urls {
				rsp, err := http.Get(url)
				res := cumsomRes{Error: err, Response: rsp}
				select {
				case <-done:
					return
				case results <- res:
				}
			}
		}()
		return results
	}
	done := make(chan interface{})
	defer close(done)
	urls := []string{"https://www.google.com", "https://badhost"}
	for res := range checkStatus(done, urls...) {
		if res.Error != nil {
			fmt.Printf("err: %+v", res.Error)
		}
		fmt.Printf("Response: %+v\n", res.Response.Status)
	}
}
