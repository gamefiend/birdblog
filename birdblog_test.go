package birdblog_test

import (
	"birdblog"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGivenJSONDataShouldDecodeIntoStruct(t *testing.T) {
	file, err := os.Open("testdata/response.json")
	if err != nil {
		t.Fatal(err)
	}
	response, err := birdblog.DecodeJSONIntoStruct(file)
	if err != nil {
		t.Fatal(err)
	}
	if len(response.Data) != 10 {
		t.Fatalf("response is the wrong length. Wanted 10, got %d", len(response.Data))
	}
	want := "@qh_murphy 100%. The Attack of Opportunity (and ESPECIALLY the Full Attack) led to just standing still in melee and wailing on each other."
	got := response.Data[0].Text

	if !(cmp.Equal(want, got)) {
		t.Error(cmp.Diff(want, got))
	}
}
