package deck_test

import (
	"bytes"
	"context"
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/corona10/goimagehash"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/k1LoW/deck"
	"github.com/k1LoW/deck/md"
)

const basePresentationID = "1wIik04tlp1U4SBHTLrSu20dPFlAGTbRHxnqdRFF9nPo"

func TestApplyMarkdown(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") == "" {
		t.Skip("skipping integration test, set TEST_INTEGRATION=1 to run")
	}

	ctx := context.Background()
	d, err := deck.CreateFrom(ctx, basePresentationID)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateTitle(ctx, "Presentation for deck integration test (Unless you are testing the deck, you can delete this file without any problems)"); err != nil {
		t.Fatalf("failed to update title: %v", err)
	}
	presentationID := d.ID()
	if err := d.AllowReadingByAnyone(ctx); err != nil {
		t.Fatal(err)
	}
	t.Logf("Presentation URL for test: https://docs.google.com/presentation/d/%s", presentationID)

	t.Cleanup(func() {
		if err := deck.Delete(ctx, presentationID); err != nil {
			t.Fatalf("failed to delete presentation %s: %v", presentationID, err)
		}
	})

	tests := []struct {
		in string
	}{
		{"testdata/slide.md"},
		{"testdata/cap.md"},
		{"testdata/br.md"},
		{"testdata/list_and_paragraph.md"},
		{"testdata/paragraph_and_list.md"},
		{"testdata/bold_and_italic.md"},
		{"testdata/emoji.md"},
		{"testdata/code.md"},
		{"testdata/style.md"},
		{"testdata/empty_list.md"},
		{"testdata/empty_link.md"},
		{"testdata/single_list.md"},
		{"testdata/nested_list.md"},
		{"testdata/images.md"},
		{"testdata/blockquote.md"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := os.ReadFile(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			markdownData, err := md.Parse("testdata", b)
			if err != nil {
				t.Fatal(err)
			}
			fromMd, err := markdownData.Contents.ToSlides(ctx, "")
			if err != nil {
				t.Fatal(err)
			}
			d, err := deck.New(ctx, deck.WithPresentationID(presentationID))
			if err != nil {
				t.Fatal(err)
			}
			if err := d.Apply(ctx, fromMd); err != nil {
				t.Fatal(err)
			}
			urls := d.ListSlideURLs()
			for i, url := range urls {
				page := i + 1
				got := deck.Screenshot(t, url)
				p := fmt.Sprintf("%s-%d.golden.png", tt.in, page)
				if os.Getenv("UPDATE_GOLDEN") != "" {
					if err := os.WriteFile(p, got, 0600); err != nil {
						t.Fatalf("failed to update golden file %s: %v", p, err)
					}
					continue
				}
				want, err := os.ReadFile(p)
				if err != nil {
					t.Fatalf("failed to read golden file %s: %v", p, err)
				}
				gotImage, err := png.Decode(bytes.NewReader(got))
				if err != nil {
					t.Fatalf("failed to decode screenshot: %v", err)
				}
				gotPHash, err := goimagehash.PerceptionHash(gotImage)
				if err != nil {
					t.Fatalf("failed to compute perception hash for screenshot: %v", err)
				}
				wantImage, err := png.Decode(bytes.NewReader(want))
				if err != nil {
					t.Fatalf("failed to decode golden screenshot: %v", err)
				}
				wantPHash, err := goimagehash.PerceptionHash(wantImage)
				if err != nil {
					t.Fatalf("failed to compute perception hash for golden screenshot: %v", err)
				}
				distance, err := gotPHash.Distance(wantPHash)
				if err != nil {
					t.Fatalf("failed to compute distance between screenshots: %v", err)
				}
				if distance > 3 { // threshold for similarity
					diffpath := fmt.Sprintf("%s-%d.diff.png", tt.in, page)
					_ = os.WriteFile(diffpath, got, 0600)
					t.Errorf("screenshot %s does not match golden file %s: distance %d, see %s for diff", p, tt.in, distance, diffpath)
				}
			}
		})
	}
}

func TestRoundTripSlidesToGoogleSlidesPresentationAndBack(t *testing.T) {
	if os.Getenv("TEST_INTEGRATION") == "" {
		t.Skip("skipping integration test, set TEST_INTEGRATION=1 to run")
	}

	ctx := context.Background()
	d, err := deck.CreateFrom(ctx, basePresentationID)
	if err != nil {
		t.Fatal(err)
	}
	presentationID := d.ID()
	t.Logf("Presentation URL for test: https://docs.google.com/presentation/d/%s", presentationID)

	t.Cleanup(func() {
		if err := deck.Delete(ctx, presentationID); err != nil {
			t.Fatalf("failed to delete presentation %s: %v", presentationID, err)
		}
	})

	tests := []struct {
		in string
	}{
		{"testdata/slide.md"},
		{"testdata/cap.md"},
		{"testdata/br.md"},
		{"testdata/bold_and_italic.md"},
		//{"testdata/list_and_paragraph.md"}, // FIXME: paragraph is not supported yet
		//{"testdata/paragraph_and_list.md"},  // FIXME: paragraph is not supported yet
		{"testdata/emoji.md"},
		{"testdata/code.md"},
		//{"testdata/style.md"},  // FIXME: class is not supported yet
		{"testdata/empty_list.md"},
		{"testdata/empty_link.md"},
		{"testdata/nested_list.md"},
		// {"testdata/images.md"},
		// {"testdata/blockquote.md"},
	}

	cmpopts := cmp.Options{
		cmpopts.IgnoreFields(deck.Fragment{}, "ClassName", "SoftLineBreak"),
		cmpopts.IgnoreUnexported(deck.Slide{}, deck.Image{}),
		cmpopts.EquateEmpty(),
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			b, err := os.ReadFile(tt.in)
			if err != nil {
				t.Fatal(err)
			}
			markdownData, err := md.Parse("testdata", b)
			if err != nil {
				t.Fatal(err)
			}
			base, err := markdownData.Contents.ToSlides(ctx, "")
			if err != nil {
				t.Fatal(err)
			}
			d, err := deck.New(ctx, deck.WithPresentationID(presentationID))
			if err != nil {
				t.Fatal(err)
			}
			// Clear existing slides before applying new ones
			if err := d.DeletePageAfter(ctx, 0); err != nil {
				t.Fatal(err)
			}
			if err := d.Apply(ctx, base); err != nil {
				t.Fatal(err)
			}
			applied, err := d.DumpSlides(ctx)
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(base, applied, cmpopts...); diff != "" {
				t.Errorf("diff after apply: %s", diff)
			}
			for i, slide := range applied {
				for _, image := range slide.Images {
					found := false
					for _, baseImage := range base[i].Images {
						if image.Compare(baseImage) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("image not found in slide %d", i+1)
					}
				}
			}
			if err := d.Apply(ctx, applied); err != nil {
				t.Fatal(err)
			}
			applied2, err := d.DumpSlides(ctx)
			if err != nil {
				t.Fatal(err)
			}
			for i, slide := range applied2 {
				for _, image := range slide.Images {
					found := false
					for _, mdImage := range applied[i].Images {
						if image.Compare(mdImage) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("image not found in slide %d", i+1)
					}
				}
			}
			if diff := cmp.Diff(applied, applied2, cmpopts...); diff != "" {
				t.Errorf("diff after re-apply: %s", diff)
			}
		})
	}
}
