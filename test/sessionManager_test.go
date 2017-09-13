package test

import "testing"
import "session"
import _ "session/memory"

func TestNew(t *testing.T)  {
	_, err := session.NewSessionMgr("memory", "test-cookie", 3600)
	if err != nil {
		t.Error(err)
	}
}
