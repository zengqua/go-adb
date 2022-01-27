package adb

import "fmt"

type Error struct {
	Code int
	Err  error
}

func (r *RequestError) Error() string {
	return fmt.Sprintf("status %d: err %v", r.StatusCode, r.Err)
}
