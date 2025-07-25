package deck

import (
	"context"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

const basePresentationID = "1wIik04tlp1U4SBHTLrSu20dPFlAGTbRHxnqdRFF9nPo"

func TestApply(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") == "" {
		t.Skip("skipping integration test, set TEST_INTEGRATION=1 to run")
	}

	ctx := context.Background()
	d, err := CreateFrom(ctx, basePresentationID)
	if err != nil {
		t.Fatal(err)
	}
	presentationID := d.ID()
	t.Cleanup(func() {
		if err := Delete(ctx, presentationID); err != nil {
			t.Fatalf("failed to delete presentation %s: %v", presentationID, err)
		}
	})

	cmpopts := cmp.Options{
		cmpopts.IgnoreFields(Fragment{}, "ClassName", "SoftLineBreak"),
		cmpopts.IgnoreUnexported(Slide{}),
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := New(ctx, WithPresentationID(presentationID))
			if err != nil {
				t.Fatal(err)
			}
			if err := d.Apply(ctx, tt.before); err != nil {
				t.Fatal(err)
			}
			before, err := d.DumpSlides(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(before, tt.before, cmpopts...); diff != "" {
				t.Error(diff)
			}

			if err := d.Apply(ctx, tt.after); err != nil {
				t.Fatal(err)
			}
			after, err := d.DumpSlides(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(after, tt.after, cmpopts...); diff != "" {
				t.Error(diff)
			}
		})
	}
}
