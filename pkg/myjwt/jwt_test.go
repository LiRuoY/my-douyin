package myjwt

import "testing"

func TestCheck(t *testing.T) {
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQ2Nzg4MzEsImlkIjoxLCJvcmlnX2lhdCI6MTY1NDA3NDAzMX0.YMV__O4lwBroDu6Dyq9Pa0t5gFnK3NGyKXVcDjcep-Y"
	claims := Check(tokenString)

	if claims == nil {
		t.Fatal("claims is nil")
	}
	t.Logf("%T %#[1]v", (*claims))
}
