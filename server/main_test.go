package main

import "testing"

//서버가 켜졌다고 가정하고 실행하는 테스트

func TestCheckDuplication(t *testing.T) {

	result := 1 + 23
	if result != 3 {
		t.Errorf("expected:%d actual:%d", 3, result)
	}
}
