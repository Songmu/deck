package deck

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/k1LoW/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/slides/v1"
)

// Apply the markdown slides to the presentation.
func (d *Deck) Apply(ctx context.Context, slides Slides) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	pages := make([]int, 0, len(slides))
	for i := range len(slides) {
		pages = append(pages, i+1)
	}
	return d.ApplyPages(ctx, slides, pages)
}

// ApplyPages applies the markdown slides to the presentation with the specified pages.
func (d *Deck) ApplyPages(ctx context.Context, ss Slides, pages []int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	if err := d.refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh presentation: %w", err)
	}
	layoutObjectIdMap := map[string]*slides.Page{}
	for _, l := range d.presentation.Layouts {
		layoutObjectIdMap[l.ObjectId] = l
	}

	before := make(Slides, 0, len(d.presentation.Slides))
	var after Slides
	for _, p := range d.presentation.Slides {
		slide := convertToSlide(p, layoutObjectIdMap)
		before = append(before, slide)
		after = append(after, slide)
	}

	for i, slide := range ss {
		if !slices.Contains(pages, i+1) {
			continue
		}
		if slide.Layout == "" {
			if i == 0 {
				slide.Layout = d.defaultTitleLayout
			} else {
				slide.Layout = d.defaultLayout
			}
		}
		if len(after) < i {
			after[i] = slide
		} else if len(after) == i {
			after = append(after, slide)
		} else {
			after[i] = slide
		}
	}
	if len(after) > len(ss) {
		var deleteIndexes []int
		for i := len(ss); i < len(after); i++ {
			if !slices.Contains(pages, i+1) {
				deleteIndexes = append(deleteIndexes, i)
			}
		}
		slices.Reverse(deleteIndexes)
		for _, i := range deleteIndexes {
			after = slices.Delete(after, i, i+1)
		}
	}

	actions, err := generateActions(before, after)
	if err != nil {
		return fmt.Errorf("failed to generate actions: %w", err)
	}

	var actionDetails []actionDetail
	for _, action := range actions {
		switch action.actionType {
		case actionTypeAppend:
			actionDetails = append(actionDetails, actionDetail{
				ActionType:  actionTypeAppend,
				Titles:      action.slide.Titles,
				Index:       nil,
				MoveToIndex: nil,
			})
		case actionTypeInsert:
			actionDetails = append(actionDetails, actionDetail{
				ActionType:  actionTypeInsert,
				Titles:      action.slide.Titles,
				Index:       &action.index,
				MoveToIndex: nil,
			})
		case actionTypeUpdate:
			actionDetails = append(actionDetails, actionDetail{
				ActionType:  actionTypeUpdate,
				Titles:      action.slide.Titles,
				Index:       &action.index,
				MoveToIndex: nil,
			})
		case actionTypeMove:
			actionDetails = append(actionDetails, actionDetail{
				ActionType:  actionTypeMove,
				Titles:      action.slide.Titles,
				Index:       &action.index,
				MoveToIndex: &action.moveToIndex,
			})
		}
	}
	d.logger.Info("applying actions", slog.Any("actions", actionDetails))

	for _, action := range actions {
		switch action.actionType {
		case actionTypeAppend:
			if err := d.AppendPage(ctx, action.slide); err != nil {
				return fmt.Errorf("failed to append slide: %w", err)
			}
		case actionTypeInsert:
			if err := d.InsertPage(ctx, action.index, action.slide); err != nil {
				return fmt.Errorf("failed to apply page: %w", err)
			}
		case actionTypeUpdate:
			d.logger.Info("appling page", slog.Int("index", action.index))
			if err := d.applyPage(ctx, action.index, action.slide); err != nil {
				return fmt.Errorf("failed to apply page: %w", err)
			}
			d.logger.Info("applied page", slog.Int("index", action.index))
		case actionTypeMove:
			if err := d.MovePage(ctx, action.index, action.moveToIndex); err != nil {
				return fmt.Errorf("failed to move page: %w", err)
			}
		case actionTypeDelete:
			if err := d.DeletePage(ctx, action.index); err != nil {
				return fmt.Errorf("failed to delete page: %w", err)
			}
		}
	}

	return nil
}

func (d *Deck) DeletePage(ctx context.Context, index int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d.logger.Info("deleting page", slog.Int("index", index))
	if err := d.deletePage(ctx, index); err != nil {
		return err
	}
	d.logger.Info("deleted page", slog.Int("index", index))
	return nil
}

func (d *Deck) DeletePageAfter(ctx context.Context, index int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	if len(d.presentation.Slides) <= index {
		return nil
	}
	req := &slides.BatchUpdatePresentationRequest{}
	for i := index + 1; i < len(d.presentation.Slides); i++ {
		req.Requests = append(req.Requests, &slides.Request{
			DeleteObject: &slides.DeleteObjectRequest{
				ObjectId: d.presentation.Slides[i].ObjectId,
			},
		})
	}
	if len(req.Requests) == 0 {
		return nil
	}
	if _, err := d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do(); err != nil {
		return err
	}
	if err := d.refresh(ctx); err != nil {
		return err
	}
	return nil
}

func (d *Deck) AppendPage(ctx context.Context, slide *Slide) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d.logger.Info("appending new page")
	index := len(d.presentation.Slides)
	if err := d.createPage(ctx, index, slide); err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	if err := d.refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh presentation: %w", err)
	}
	if err := d.applyPage(ctx, index, slide); err != nil {
		return fmt.Errorf("failed to apply page: %w", err)
	}
	if err := d.refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh presentation: %w", err)
	}
	d.logger.Info("appended page")
	return nil
}

func (d *Deck) MovePage(ctx context.Context, from_index, to_index int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d.logger.Info("moving page", slog.Int("from_index", from_index), slog.Int("to_index", to_index))
	if err := d.movePage(ctx, from_index, to_index); err != nil {
		return err
	}
	d.logger.Info("moved page", slog.Int("from_index", from_index), slog.Int("to_index", to_index))
	return nil
}

func (d *Deck) InsertPage(ctx context.Context, index int, slide *Slide) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d.logger.Info("inserting page", slog.Int("index", index))
	if len(d.presentation.Slides) <= index {
		return fmt.Errorf("index out of range: %d", index)
	}
	if err := d.createPage(ctx, index, slide); err != nil {
		return fmt.Errorf("failed to create page: %w", err)
	}
	if index == 0 {
		if err := d.movePage(ctx, 1, 0); err != nil {
			return fmt.Errorf("failed to move page: %w", err)
		}
	}
	if err := d.refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh presentation: %w", err)
	}
	if err := d.applyPage(ctx, index, slide); err != nil {
		return fmt.Errorf("failed to apply page: %w", err)
	}
	if err := d.refresh(ctx); err != nil {
		return fmt.Errorf("failed to refresh presentation: %w", err)
	}
	d.logger.Info("inserted page", slog.Int("index", index))
	return nil
}

func (d *Deck) applyPage(ctx context.Context, index int, slide *Slide) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	layoutMap := map[string]*slides.Page{}
	for _, l := range d.presentation.Layouts {
		layoutMap[l.LayoutProperties.DisplayName] = l
	}

	layout, ok := layoutMap[slide.Layout]
	if !ok {
		return fmt.Errorf("layout not found: %q", slide.Layout)
	}

	if len(d.presentation.Slides) <= index {
		return fmt.Errorf("index out of range: %d", index)
	}
	if slide.Freeze {
		d.logger.Info("skip applying page. because freeze:true", slog.Int("index", index))
		return nil
	}
	currentSlide := d.presentation.Slides[index]
	if currentSlide.SlideProperties.LayoutObjectId != layout.ObjectId {
		if err := d.updateLayout(ctx, index, slide); err != nil {
			return err
		}
	}

	var (
		titles                    []placeholder
		subtitles                 []placeholder
		bodies                    []placeholder
		currentImages             []*Image
		currentImageObjectIDMap   = map[*Image]string{} // key: *Image, value: objectID
		currentTextBoxes          []*textBox
		currentTextBoxObjectIDMap = map[*textBox]string{} // key: *textBox, value: objectID
	)
	currentSlide = d.presentation.Slides[index]
	for _, element := range currentSlide.PageElements {
		switch {
		case element.Shape != nil && element.Shape.Placeholder != nil:
			switch element.Shape.Placeholder.Type {
			case "CENTERED_TITLE", "TITLE":
				titles = append(titles, placeholder{
					objectID: element.ObjectId,
					x:        element.Transform.TranslateX,
					y:        element.Transform.TranslateY,
				})
				if err := d.clearPlaceholder(ctx, element.ObjectId); err != nil {
					return err
				}
			case "SUBTITLE":
				subtitles = append(subtitles, placeholder{
					objectID: element.ObjectId,
					x:        element.Transform.TranslateX,
					y:        element.Transform.TranslateY,
				})
				if err := d.clearPlaceholder(ctx, element.ObjectId); err != nil {
					return err
				}
			case "BODY":
				bodies = append(bodies, placeholder{
					objectID: element.ObjectId,
					x:        element.Transform.TranslateX,
					y:        element.Transform.TranslateY,
				})
				if err := d.clearPlaceholder(ctx, element.ObjectId); err != nil {
					return err
				}
			}
		case element.Image != nil:
			var (
				image *Image
				err   error
			)
			if element.Description == descriptionImageFromMarkdown {
				image, err = NewImageFromMarkdown(element.Image.ContentUrl)
				if err != nil {
					return fmt.Errorf("failed to create image from code block %s: %w", element.Image.ContentUrl, err)
				}
			} else {
				image, err = NewImage(element.Image.ContentUrl)
				if err != nil {
					return fmt.Errorf("failed to create image from %s: %w", element.Image.ContentUrl, err)
				}
			}
			currentImages = append(currentImages, image)
			currentImageObjectIDMap[image] = element.ObjectId
		case element.Shape != nil && element.Shape.ShapeType == "TEXT_BOX" && element.Shape.Text != nil:
			tb := &textBox{}
			if element.Description == descriptionTextboxFromMarkdown {
				tb.fromMarkdown = true
			}
			tb.paragraphs = convertToParagraphs(element.Shape.Text)
			currentTextBoxes = append(currentTextBoxes, tb)
			currentTextBoxObjectIDMap[tb] = element.ObjectId
		}
	}
	var speakerNotesID string
	for _, element := range currentSlide.SlideProperties.NotesPage.PageElements {
		if element.Shape != nil && element.Shape.Placeholder != nil {
			if element.Shape.Placeholder.Type == "BODY" {
				speakerNotesID = element.ObjectId
				if err := d.clearPlaceholder(ctx, speakerNotesID); err != nil {
					return err
				}
			}
		}
	}
	if speakerNotesID == "" {
		return fmt.Errorf("speaker notes not found")
	}

	// set titles
	req := &slides.BatchUpdatePresentationRequest{}
	sort.Slice(titles, func(i, j int) bool {
		if titles[i].y == titles[j].y {
			return titles[i].x < titles[j].x
		}
		return titles[i].y < titles[j].y
	})
	for i, title := range slide.Titles {
		if len(titles) <= i {
			continue
		}
		req.Requests = append(req.Requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: titles[i].objectID,
				Text:     title,
			},
		})
	}

	// set subtitles
	sort.Slice(subtitles, func(i, j int) bool {
		if subtitles[i].y == subtitles[j].y {
			return subtitles[i].x < subtitles[j].x
		}
		return subtitles[i].y < subtitles[j].y
	})
	for i, subtitle := range slide.Subtitles {
		if len(subtitles) <= i {
			continue
		}
		req.Requests = append(req.Requests, &slides.Request{
			InsertText: &slides.InsertTextRequest{
				ObjectId: subtitles[i].objectID,
				Text:     subtitle,
			},
		})
	}

	// set speaker notes
	req.Requests = append(req.Requests, &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId: speakerNotesID,
			Text:     slide.SpeakerNote,
		},
	})

	// set bodies
	sort.Slice(bodies, func(i, j int) bool {
		if bodies[i].y == bodies[j].y {
			return bodies[i].x < bodies[j].x
		}
		return bodies[i].y < bodies[j].y
	})
	for i, body := range slide.Bodies {
		if len(bodies) <= i {
			continue
		}
		reqs, styleReqs, err := d.applyParagraphsRequests(bodies[i].objectID, body.Paragraphs)
		if err != nil {
			return fmt.Errorf("failed to apply paragraphs: %w", err)
		}
		req.Requests = append(req.Requests, reqs...)
		req.Requests = append(req.Requests, styleReqs...)
	}

	// set images
	for i, image := range slide.Images {
		found := false
		for _, currentImage := range currentImages {
			if currentImage.Compare(image) {
				found = true
				break
			}
		}
		if found {
			continue
		}

		// upload image to Google Drive
		df := &drive.File{
			Name:     fmt.Sprintf("________tmp-for-deck-%s", time.Now().Format(time.RFC3339)),
			MimeType: string(image.mimeType),
		}
		uploaded, err := d.driveSrv.Files.Create(df).Media(bytes.NewBuffer(image.Bytes())).Do()
		if err != nil {
			return fmt.Errorf("failed to upload image: %w", err)
		}
		if _, err := d.driveSrv.Permissions.Create(uploaded.Id, &drive.Permission{
			Type: "anyone",
			Role: "reader",
		}).Do(); err != nil {
			return fmt.Errorf("failed to set permission for image: %w", err)
		}
		f, err := d.driveSrv.Files.Get(uploaded.Id).Fields("webContentLink").Do()
		if err != nil {
			return fmt.Errorf("failed to get webContentLink for image: %w", err)
		}
		defer func() {
			// Clean up the uploaded image after use
			if err := d.driveSrv.Files.Delete(uploaded.Id).Do(); err != nil {
				d.logger.Error("failed to delete uploaded image", slog.String("id", uploaded.Id), slog.Any("error", err))
			}
		}()
		if f.WebContentLink == "" {
			return fmt.Errorf("webContentLink is empty for image: %s", uploaded.Id)
		}

		imageObjectID := fmt.Sprintf("image-%s", uuid.New().String())
		req.Requests = append(req.Requests, &slides.Request{
			CreateImage: &slides.CreateImageRequest{
				ObjectId: imageObjectID,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: currentSlide.ObjectId,
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: float64(i+1) * 100000,
						TranslateY: float64(i+1) * 100000,
						Unit:       "EMU",
					},
				},
				Url: f.WebContentLink,
			},
		})
		if image.fromMarkdown {
			req.Requests = append(req.Requests, &slides.Request{
				UpdatePageElementAltText: &slides.UpdatePageElementAltTextRequest{
					ObjectId:    imageObjectID,
					Description: descriptionImageFromMarkdown,
				},
			})
		}
	}

	// set text boxes
	for i, bq := range slide.BlockQuotes {
		found := false
		for _, currentTextBox := range currentTextBoxes {
			if paragraphsEqual(currentTextBox.paragraphs, bq.Paragraphs) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		// create new text box
		textBoxObjectID := fmt.Sprintf("textbox-%s", uuid.New().String())
		req.Requests = append(req.Requests, &slides.Request{
			CreateShape: &slides.CreateShapeRequest{
				ObjectId: textBoxObjectID,
				ElementProperties: &slides.PageElementProperties{
					PageObjectId: currentSlide.ObjectId,
					Size: &slides.Size{
						Height: &slides.Dimension{
							Magnitude: float64(500000 * len(bq.Paragraphs)),
							Unit:      "EMU",
						},
						Width: &slides.Dimension{
							Magnitude: 5000000,
							Unit:      "EMU",
						},
					},
					Transform: &slides.AffineTransform{
						ScaleX:     1.0,
						ScaleY:     1.0,
						TranslateX: float64(i+1) * 100000,
						TranslateY: float64(i+1) * 100000,
						Unit:       "EMU",
					},
				},
				ShapeType: "TEXT_BOX",
			},
		})

		sp, ok := d.shapes[styleBlockQuote]
		if ok {
			req.Requests = append(req.Requests, &slides.Request{
				UpdateShapeProperties: &slides.UpdateShapePropertiesRequest{
					ObjectId:        textBoxObjectID,
					ShapeProperties: sp,
					Fields:          "shapeBackgroundFill,outline,shadow",
				},
			})
		}
		reqs, styleReqs, err := d.applyParagraphsRequests(textBoxObjectID, bq.Paragraphs)
		if err != nil {
			return fmt.Errorf("failed to apply paragraphs: %w", err)
		}
		req.Requests = append(req.Requests, reqs...)

		s, ok := d.styles[styleBlockQuote]
		if ok {
			req.Requests = append(req.Requests, &slides.Request{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: textBoxObjectID,
					Style:    s,
					Fields:   "bold,italic,underline,foregroundColor,fontFamily,backgroundColor",
				},
			})
		}

		req.Requests = append(req.Requests, styleReqs...)

		req.Requests = append(req.Requests, &slides.Request{
			UpdatePageElementAltText: &slides.UpdatePageElementAltTextRequest{
				ObjectId:    textBoxObjectID,
				Description: descriptionTextboxFromMarkdown,
			},
		})
	}

	// set skip flag to slide
	req.Requests = append(req.Requests, &slides.Request{
		UpdateSlideProperties: &slides.UpdateSlidePropertiesRequest{
			ObjectId: currentSlide.ObjectId,
			SlideProperties: &slides.SlideProperties{
				IsSkipped: slide.Skip,
			},
			Fields: "isSkipped",
		},
	})

	// prune unmatched images via markdown
	for _, currentImage := range currentImages {
		if !currentImage.fromMarkdown {
			continue
		}
		found := false
		for _, image := range slide.Images {
			if currentImage.Compare(image) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		imageObjectID, ok := currentImageObjectIDMap[currentImage]
		if !ok {
			return fmt.Errorf("image object ID not found for image: %s", currentImage.url)
		}
		req.Requests = append(req.Requests, &slides.Request{
			DeleteObject: &slides.DeleteObjectRequest{
				ObjectId: imageObjectID,
			},
		})
	}

	// prune unmatched text boxes via markdown
	for _, currentTextBox := range currentTextBoxes {
		if !currentTextBox.fromMarkdown {
			continue
		}
		found := false
		for _, bq := range slide.BlockQuotes {
			if paragraphsEqual(currentTextBox.paragraphs, bq.Paragraphs) {
				found = true
				break
			}
		}
		if found {
			continue
		}
		textBoxObjectID, ok := currentTextBoxObjectIDMap[currentTextBox]
		if !ok {
			return fmt.Errorf("text box object ID not found for text box: %v", currentTextBox.paragraphs)
		}
		req.Requests = append(req.Requests, &slides.Request{
			DeleteObject: &slides.DeleteObjectRequest{
				ObjectId: textBoxObjectID,
			},
		})
	}

	if len(req.Requests) > 0 {
		if _, err := d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do(); err != nil {
			return err
		}
	}

	return nil
}

func (d *Deck) deletePage(ctx context.Context, index int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	if len(d.presentation.Slides) <= index {
		return nil
	}
	currentSlide := d.presentation.Slides[index]
	req := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{
			{
				DeleteObject: &slides.DeleteObjectRequest{
					ObjectId: currentSlide.ObjectId,
				},
			},
		},
	}
	if _, err := d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do(); err != nil {
		return err
	}
	if err := d.refresh(ctx); err != nil {
		return err
	}
	return nil
}

func (d *Deck) movePage(ctx context.Context, from_index, to_index int) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	if from_index == to_index || from_index < 0 || to_index < 0 || from_index >= len(d.presentation.Slides) || to_index >= len(d.presentation.Slides) {
		return nil
	}

	currentSlide := d.presentation.Slides[from_index]

	if from_index < to_index {
		to_index++
	}

	req := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{
			{
				UpdateSlidesPosition: &slides.UpdateSlidesPositionRequest{
					SlideObjectIds:  []string{currentSlide.ObjectId},
					InsertionIndex:  int64(to_index),
					ForceSendFields: []string{"InsertionIndex"},
				},
			},
		},
	}

	if _, err := d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do(); err != nil {
		return err
	}
	if err := d.refresh(ctx); err != nil {
		return err
	}
	return nil
}

func (d *Deck) applyParagraphsRequests(objectID string, paragraphs []*Paragraph) (reqs []*slides.Request, styleReqs []*slides.Request, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()

	bulletRanges := map[int]*bulletRange{}
	count := int64(0)
	text := ""
	bulletStartIndex := int64(0) // reset per body
	bulletEndIndex := int64(0)   // reset per body
	currentBullet := BulletNone
	for j, paragraph := range paragraphs {
		plen := 0
		startIndex := count
		if paragraph.Bullet != BulletNone {
			if paragraph.Nesting > 0 {
				text += strings.Repeat("\t", paragraph.Nesting)
				plen += paragraph.Nesting
			}
		}
		for _, fragment := range paragraph.Fragments {
			fValue := fragment.Value
			flen := countString(fragment.Value)
			for _, req := range d.getInlineStyleRequests(fragment) {
				req.ObjectId = objectID
				req.TextRange = &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: ptrInt64(startIndex),
					EndIndex:   ptrInt64(startIndex + int64(flen)),
				}
				styleReqs = append(styleReqs, &slides.Request{
					UpdateTextStyle: req,
				})
			}
			if fragment.SoftLineBreak {
				fValue += "\n"
				flen++
			}
			plen += flen
			text += fValue
			startIndex += int64(flen)
		}

		if len(paragraphs) > j+1 {
			nextParagraph := paragraphs[j+1]
			if paragraph.Bullet != nextParagraph.Bullet || paragraph.Bullet != BulletNone {
				text += "\n"
				plen++
			}
		}

		if paragraph.Bullet != BulletNone {
			if paragraph.Nesting == 0 && currentBullet != paragraph.Bullet {
				bulletStartIndex = count
				bulletEndIndex = count
				bulletRanges[int(bulletStartIndex)] = &bulletRange{
					bullet: paragraph.Bullet,
					start:  bulletStartIndex,
					end:    bulletEndIndex,
				}
			}
			bulletEndIndex += int64(plen)
			bulletRanges[int(bulletStartIndex)].end = bulletEndIndex
		}
		currentBullet = paragraph.Bullet
		count += int64(plen)
	}

	reqs = append(reqs, &slides.Request{
		InsertText: &slides.InsertTextRequest{
			ObjectId: objectID,
			Text:     text,
		},
	})
	var bulletRangeSlice []*bulletRange
	for _, r := range bulletRanges {
		bulletRangeSlice = append(bulletRangeSlice, r)
	}
	// reverse sort
	// Because the Range changes each time it is converted to a list, convert from the end to a list.
	sort.Slice(bulletRangeSlice, func(i, j int) bool {
		return bulletRangeSlice[i].start > bulletRangeSlice[j].start
	})
	for _, r := range bulletRangeSlice {
		startIndex := int64(r.start)
		endIndex := int64(r.end - 1)
		if startIndex <= endIndex {
			endIndex++
		}
		styleReqs = append(styleReqs, &slides.Request{
			CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
				ObjectId:     objectID,
				BulletPreset: convertBullet(r.bullet),
				TextRange: &slides.Range{
					Type:       "FIXED_RANGE",
					StartIndex: ptrInt64(startIndex),
					EndIndex:   ptrInt64(endIndex),
				},
			},
		})
	}

	return reqs, styleReqs, nil
}

func (d *Deck) getInlineStyleRequests(fragment *Fragment) (reqs []*slides.UpdateTextStyleRequest) {
	if fragment.Code {
		s, ok := d.styles[styleCode]
		if ok {
			reqs = append(reqs, buildCustomStyleRequest(s))
		} else {
			reqs = append(reqs, &slides.UpdateTextStyleRequest{
				Style: &slides.TextStyle{
					ForegroundColor: &slides.OptionalColor{
						OpaqueColor: &slides.OpaqueColor{
							RgbColor: &slides.RgbColor{
								Red:   0.0,
								Green: 0.0,
								Blue:  0.0,
							},
						},
					},
					FontFamily: defaultCodeFontFamily,
					BackgroundColor: &slides.OptionalColor{
						OpaqueColor: &slides.OpaqueColor{
							RgbColor: &slides.RgbColor{
								Red:   0.95,
								Green: 0.95,
								Blue:  0.95,
							},
						},
					},
				},
				Fields: "foregroundColor,fontFamily,backgroundColor",
			})
		}
	}

	if fragment.Bold {
		s, ok := d.styles[styleBold]
		if ok {
			reqs = append(reqs, buildCustomStyleRequest(s))
		} else {
			reqs = append(reqs, &slides.UpdateTextStyleRequest{
				Style: &slides.TextStyle{
					Bold: true,
				},
				Fields: "bold",
			})
		}
	}

	if fragment.Italic {
		s, ok := d.styles[styleItalic]
		if ok {
			reqs = append(reqs, buildCustomStyleRequest(s))
		} else {
			reqs = append(reqs, &slides.UpdateTextStyleRequest{
				Style: &slides.TextStyle{
					Italic: true,
				},
				Fields: "italic",
			})
		}
	}

	if fragment.Link != "" {
		s, ok := d.styles[styleLink]
		if ok {
			req := buildCustomStyleRequest(s)
			req.Fields = "link,bold,italic,underline,foregroundColor,fontFamily,backgroundColor"
			req.Style.Link = &slides.Link{
				Url: fragment.Link,
			}
			reqs = append(reqs, req)
		} else {
			reqs = append(reqs, &slides.UpdateTextStyleRequest{
				Style: &slides.TextStyle{
					Link: &slides.Link{
						Url: fragment.Link,
					},
				},
				Fields: "link",
			})
		}
	}

	if fragment.ClassName != "" {
		s, ok := d.styles[fragment.ClassName]
		if ok {
			reqs = append(reqs, buildCustomStyleRequest(s))
		}
	}

	return reqs
}

func buildCustomStyleRequest(s *slides.TextStyle) *slides.UpdateTextStyleRequest {
	return &slides.UpdateTextStyleRequest{
		Style: &slides.TextStyle{
			Bold:            s.Bold,
			Italic:          s.Italic,
			Underline:       s.Underline,
			ForegroundColor: s.ForegroundColor,
			FontFamily:      s.FontFamily,
			BackgroundColor: s.BackgroundColor,
		},
		Fields: "bold,italic,underline,foregroundColor,fontFamily,backgroundColor",
	}
}

func (d *Deck) clearPlaceholder(ctx context.Context, placeholderID string) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	req := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{
			{
				UpdateTextStyle: &slides.UpdateTextStyleRequest{
					ObjectId: placeholderID,
					Style: &slides.TextStyle{
						Bold:   false,
						Italic: false,
					},
					TextRange: &slides.Range{
						Type: "ALL",
					},
					Fields: "*",
				},
			},
			{
				DeleteParagraphBullets: &slides.DeleteParagraphBulletsRequest{
					ObjectId: placeholderID,
					TextRange: &slides.Range{
						Type: "ALL",
					},
				},
			},
			{
				DeleteText: &slides.DeleteTextRequest{
					ObjectId: placeholderID,
					TextRange: &slides.Range{
						Type: "ALL",
					},
				},
			},
		},
	}

	_, _ = d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do()
	return nil
}

func (d *Deck) updateLayout(ctx context.Context, index int, slide *Slide) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	currentSlide := d.presentation.Slides[index]
	// create new page
	if err := d.createPage(ctx, index+1, slide); err != nil {
		return err
	}

	newSlide := d.presentation.Slides[index+1]
	req := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{},
	}
	insertReq := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{},
	}
	var (
		styleReqs  []*slides.Request
		bulletReqs []*slides.Request
	)

	for _, element := range currentSlide.PageElements {
		// copy images from the current slide to the new slide
		if element.Image != nil {
			req.Requests = append(req.Requests, &slides.Request{
				CreateImage: &slides.CreateImageRequest{
					ElementProperties: &slides.PageElementProperties{
						Size:         element.Size,
						Transform:    element.Transform,
						PageObjectId: newSlide.ObjectId,
					},
					Url: element.Image.ContentUrl,
				},
			})
		}
		// copy shapes from the current slide to the new slide
		if element.Shape != nil && element.Shape.Placeholder == nil {
			type paragraphInfo struct {
				startIndex   int64
				endIndex     int64
				bullet       *slides.Bullet
				nestingLevel int64
			}

			var paragraphInfos []paragraphInfo
			currentIndex := int64(0)
			text := ""
			shapeObjectID := fmt.Sprintf("shape-%s", uuid.New().String())

			for _, textElement := range element.Shape.Text.TextElements {
				if textElement.ParagraphMarker != nil {
					pInfo := paragraphInfo{
						startIndex: currentIndex,
					}
					if textElement.ParagraphMarker.Bullet != nil {
						pInfo.bullet = textElement.ParagraphMarker.Bullet
						pInfo.nestingLevel = textElement.ParagraphMarker.Bullet.NestingLevel
					}
					paragraphInfos = append(paragraphInfos, pInfo)
				}

				if textElement.TextRun != nil {
					runText := textElement.TextRun.Content

					// Handle nesting by adding tabs
					if len(paragraphInfos) > 0 && currentIndex == paragraphInfos[len(paragraphInfos)-1].startIndex {
						// This is the start of a bulleted paragraph
						if paragraphInfos[len(paragraphInfos)-1].nestingLevel > 0 {
							// Add tabs for nesting
							tabs := strings.Repeat("\t", int(paragraphInfos[len(paragraphInfos)-1].nestingLevel))
							text += tabs
							currentIndex += int64(countString(tabs))
						}
					}

					text += runText

					// Adjust style indices based on actual position in new text
					if textElement.TextRun.Style != nil {
						startIdx := currentIndex
						endIdx := currentIndex + int64(countString(runText))
						styleReqs = append(styleReqs, &slides.Request{
							UpdateTextStyle: &slides.UpdateTextStyleRequest{
								ObjectId: shapeObjectID,
								Style:    textElement.TextRun.Style,
								TextRange: &slides.Range{
									Type:       "FIXED_RANGE",
									StartIndex: ptrInt64(startIdx),
									EndIndex:   ptrInt64(endIdx),
								},
								Fields: "*",
							},
						})
					}
					currentIndex += int64(countString(runText))
				}
			}

			// Update end indices for paragraphs
			for i := range paragraphInfos {
				if i < len(paragraphInfos)-1 {
					paragraphInfos[i].endIndex = paragraphInfos[i+1].startIndex - 1
				} else {
					paragraphInfos[i].endIndex = currentIndex - 1
				}
			}
			req.Requests = append(req.Requests, &slides.Request{
				CreateShape: &slides.CreateShapeRequest{
					ObjectId: shapeObjectID,
					ElementProperties: &slides.PageElementProperties{
						Size:         element.Size,
						Transform:    element.Transform,
						PageObjectId: newSlide.ObjectId,
					},
					ShapeType: element.Shape.ShapeType,
				},
			})
			styleReqs = append(styleReqs, &slides.Request{
				UpdateShapeProperties: &slides.UpdateShapePropertiesRequest{
					ObjectId:        shapeObjectID,
					ShapeProperties: element.Shape.ShapeProperties,
					Fields:          "contentAlignment,link,outline,shadow,shapeBackgroundFill",
				},
			})

			insertReq.Requests = append(insertReq.Requests, &slides.Request{
				InsertText: &slides.InsertTextRequest{
					ObjectId: shapeObjectID,
					Text:     strings.TrimSuffix(text, "\n"),
				},
			})

			var br *bulletRange
			for _, pInfo := range paragraphInfos {
				if pInfo.bullet == nil {
					if br != nil {
						bulletReqs = append(bulletReqs, &slides.Request{
							CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
								ObjectId:     shapeObjectID,
								BulletPreset: convertBullet(br.bullet),
								TextRange: &slides.Range{
									Type:       "FIXED_RANGE",
									StartIndex: ptrInt64(int64(br.start)),
									EndIndex:   ptrInt64(int64(br.end)),
								},
							},
						})
						br = nil
					}
					continue
				}
				if br == nil {
					br = &bulletRange{
						bullet: getBulletPresetFromSlidesBullet(pInfo.bullet),
						start:  pInfo.startIndex,
						end:    pInfo.endIndex,
					}
				} else {
					br.end = pInfo.endIndex
				}
			}
			if br != nil {
				bulletReqs = append(bulletReqs, &slides.Request{
					CreateParagraphBullets: &slides.CreateParagraphBulletsRequest{
						ObjectId:     shapeObjectID,
						BulletPreset: convertBullet(br.bullet),
						TextRange: &slides.Range{
							Type:       "FIXED_RANGE",
							StartIndex: ptrInt64(int64(br.start)),
							EndIndex:   ptrInt64(int64(br.end)),
						},
					},
				})
			}

			// reverse sort
			// Because the Range changes each time it is converted to a list, convert from the end to a list.
			sort.Slice(bulletReqs, func(i, j int) bool {
				return *bulletReqs[i].CreateParagraphBullets.TextRange.StartIndex > *bulletReqs[j].CreateParagraphBullets.TextRange.StartIndex
			})

			if len(styleReqs) > 0 || len(bulletReqs) > 0 {
				// Apply styles first, then bullets (important for correct rendering)
				insertReq.Requests = append(insertReq.Requests, styleReqs...)
				insertReq.Requests = append(insertReq.Requests, bulletReqs...)
				styleReqs = nil  // reset after adding to requests
				bulletReqs = nil // reset after adding to requests
			}
		}
	}

	if len(req.Requests) > 0 {
		if _, err := d.srv.Presentations.BatchUpdate(d.id, req).Context(ctx).Do(); err != nil {
			return fmt.Errorf("failed to copy images: %w", err)
		}
	}
	if len(insertReq.Requests) > 0 {
		if _, err := d.srv.Presentations.BatchUpdate(d.id, insertReq).Context(ctx).Do(); err != nil {
			return fmt.Errorf("failed to insert text: %w", err)
		}
	}

	if err := d.DeletePage(ctx, index); err != nil {
		return err
	}
	return nil
}

func convertToSlide(p *slides.Page, layoutObjectIdMap map[string]*slides.Page) *Slide {
	slide := &Slide{
		Layout: "",
		Freeze: false,
	}
	if p.SlideProperties != nil {
		page, ok := layoutObjectIdMap[p.SlideProperties.LayoutObjectId]
		if ok {
			slide.Layout = page.LayoutProperties.DisplayName
		}
	}

	var titles []string
	var subtitles []string
	var bodies []*Body
	var images []*Image
	var blockQuotes []*BlockQuote

	// Extract titles, subtitles, and bodies from page elements
	for _, element := range p.PageElements {
		switch {
		case element.Shape != nil && element.Shape.Text != nil && element.Shape.Placeholder != nil:
			switch element.Shape.Placeholder.Type {
			case "CENTERED_TITLE", "TITLE":
				text := extractText(element.Shape.Text)
				if text != "" {
					titles = append(titles, text)
				}
			case "SUBTITLE":
				text := extractText(element.Shape.Text)
				if text != "" {
					subtitles = append(subtitles, text)
				}
			case "BODY":
				paragraphs := convertToParagraphs(element.Shape.Text)
				if len(paragraphs) > 0 {
					bodies = append(bodies, &Body{
						Paragraphs: paragraphs,
					})
				}
			}
		case element.Image != nil:
			var (
				image *Image
				err   error
			)
			if element.Description == descriptionImageFromMarkdown {
				image, err = NewImageFromMarkdown(element.Image.ContentUrl)
				if err != nil {
					continue // Skip if image cannot be created
				}
			} else {
				image, err = NewImage(element.Image.ContentUrl)
				if err != nil {
					continue // Skip if image cannot be created
				}
			}
			images = append(images, image)
		case element.Shape != nil && element.Shape.ShapeType == "TEXT_BOX" && element.Shape.Text != nil:
			if element.Description != descriptionTextboxFromMarkdown {
				continue
			}
			bq := &BlockQuote{
				Paragraphs: convertToParagraphs(element.Shape.Text),
			}
			blockQuotes = append(blockQuotes, bq)
		}
	}

	slide.Titles = titles
	slide.Subtitles = subtitles
	slide.Bodies = bodies
	slide.Images = images
	slide.BlockQuotes = blockQuotes

	// Extract speaker notes
	if p.SlideProperties != nil && p.SlideProperties.NotesPage != nil {
		for _, element := range p.SlideProperties.NotesPage.PageElements {
			if element.Shape != nil && element.Shape.Text != nil && element.Shape.Placeholder != nil {
				if element.Shape.Placeholder.Type == "BODY" {
					slide.SpeakerNote = extractText(element.Shape.Text)
					break
				}
			}
		}
	}

	return slide
}

// extractText extracts plain text from Shape.Text.
func extractText(text *slides.TextContent) string {
	if text == nil || len(text.TextElements) == 0 {
		return ""
	}

	var result strings.Builder
	for _, element := range text.TextElements {
		if element.TextRun != nil {
			result.WriteString(element.TextRun.Content)
		}
	}
	return strings.TrimSpace(result.String())
}

// convertToParagraphs converts TextContent to a slice of Paragraphs.
func convertToParagraphs(text *slides.TextContent) []*Paragraph {
	if text == nil || len(text.TextElements) == 0 {
		return nil
	}

	var paragraphs []*Paragraph
	var currentParagraph *Paragraph
	var currentBullet Bullet

	for _, element := range text.TextElements {
		if element.ParagraphMarker != nil {
			// Start of a new paragraph
			if currentParagraph != nil && len(currentParagraph.Fragments) > 0 {
				// Check if this is a continuation of a non-bullet paragraph
				// If the previous paragraph had no bullet and this one also has no bullet,
				// merge them with a newline fragment
				if currentParagraph.Bullet == BulletNone &&
					(element.ParagraphMarker.Bullet == nil) {
					// Add newline fragment to continue the paragraph
					currentParagraph.Fragments = append(currentParagraph.Fragments, &Fragment{
						Value: "\n",
					})
					// Don't create a new paragraph, continue with the current one
				} else {
					paragraphs = append(paragraphs, currentParagraph)
					currentParagraph = &Paragraph{
						Fragments: []*Fragment{},
						Nesting:   0,
					}
				}
			} else {
				currentParagraph = &Paragraph{
					Fragments: []*Fragment{},
					Nesting:   0,
				}
			}

			// Process bullet points
			if element.ParagraphMarker.Bullet != nil {
				// Determine the type of bullet points based on glyph content
				if element.ParagraphMarker.Bullet.Glyph != "" {
					glyph := element.ParagraphMarker.Bullet.Glyph
					// Check for numbered bullets (1, 2, 3, etc.)
					if strings.Contains(glyph, "1") || strings.Contains(glyph, "2") || strings.Contains(glyph, "3") ||
						strings.Contains(glyph, "4") || strings.Contains(glyph, "5") || strings.Contains(glyph, "6") ||
						strings.Contains(glyph, "7") || strings.Contains(glyph, "8") || strings.Contains(glyph, "9") ||
						strings.Contains(glyph, "0") {
						currentBullet = BulletNumber
					} else {
						currentBullet = BulletDash
					}
				} else {
					// If no glyph, assume it's a dash bullet
					currentBullet = BulletDash
				}
				currentParagraph.Bullet = currentBullet

				// Set nesting level
				currentParagraph.Nesting = int(element.ParagraphMarker.Bullet.NestingLevel)
			} else {
				currentBullet = BulletNone
				currentParagraph.Bullet = currentBullet
			}
		}

		if element.TextRun != nil && currentParagraph != nil {
			// Process text content
			content := element.TextRun.Content

			// Check if this is an empty content that should be treated as SoftLineBreak
			if content == "" {
				fragment := &Fragment{
					Value:         "",
					SoftLineBreak: true,
					ClassName:     "",
				}
				currentParagraph.Fragments = append(currentParagraph.Fragments, fragment)
				continue
			}

			// Handle special case where content is just a newline
			if content == "\n" {
				// Check if the previous fragment exists and can be marked with SoftLineBreak
				if len(currentParagraph.Fragments) > 0 {
					lastFragment := currentParagraph.Fragments[len(currentParagraph.Fragments)-1]
					if lastFragment.Value != "" && !lastFragment.SoftLineBreak {
						lastFragment.SoftLineBreak = true
						continue
					}
				}
				// If no previous fragment or it already has SoftLineBreak, add as newline fragment
				currentParagraph.Fragments = append(currentParagraph.Fragments, &Fragment{
					Value: "\n",
				})
				continue
			}

			// Get styles from TextRun
			var bold, italic, code bool
			var link string
			if element.TextRun.Style != nil {
				bold = element.TextRun.Style.Bold
				italic = element.TextRun.Style.Italic
				if element.TextRun.Style.Link != nil && element.TextRun.Style.Link.Url != "" {
					link = element.TextRun.Style.Link.Url
				}

				// Detect code style (based on font family and background color)
				if element.TextRun.Style.FontFamily == defaultCodeFontFamily ||
					(element.TextRun.Style.BackgroundColor != nil &&
						element.TextRun.Style.BackgroundColor.OpaqueColor != nil &&
						element.TextRun.Style.BackgroundColor.OpaqueColor.RgbColor != nil) {
					code = true
				}
			}

			// Process line breaks
			softLineBreak := false
			if strings.HasSuffix(content, "\n") {
				content = strings.TrimSuffix(content, "\n")
				softLineBreak = true
			}

			fragments := []*Fragment{{
				Value:         content,
				Bold:          bold,
				Italic:        italic,
				Code:          code,
				Link:          link,
				SoftLineBreak: softLineBreak,
			}}

			for _, fragment := range fragments {
				// Only add non-empty fragments
				if fragment.Value != "" {
					currentParagraph.Fragments = append(currentParagraph.Fragments, fragment)
				}
			}
		}
	}

	// Add the last paragraph
	if currentParagraph != nil && len(currentParagraph.Fragments) > 0 {
		paragraphs = append(paragraphs, currentParagraph)
	}

	return paragraphs
}

// countString counts the number of characters in a string, considering UTF-16 surrogate pairs.
// This is because Google Slides' character count is derived from JavaScript.
func countString(s string) int {
	length := 0
	for _, r := range s {
		if r <= 0xFFFF && (r < 0xD800 || r > 0xDFFF) {
			length++
		} else {
			length += 2
		}
	}
	return length
}

func convertBullet(b Bullet) string {
	switch b {
	case BulletDash:
		return "BULLET_DISC_CIRCLE_SQUARE"
	case BulletNumber:
		return "NUMBERED_DIGIT_ALPHA_ROMAN"
	case BulletAlpha:
		return "NUMBERED_DIGIT_ALPHA_ROMAN"
	default:
		return "UNRECOGNIZED"
	}
}

// getBulletPresetFromSlidesBullet converts a slides.Bullet to a BulletPreset string.
func getBulletPresetFromSlidesBullet(bullet *slides.Bullet) Bullet {
	if bullet == nil || bullet.Glyph == "" {
		return BulletNone
	}

	glyph := bullet.Glyph
	// Check for numbered bullets (1, 2, 3, etc.)
	for _, digit := range "0123456789" {
		if strings.Contains(glyph, string(digit)) {
			return BulletNumber
		}
	}

	// Check for alphabetic bullets (a., A., etc.)
	if strings.ContainsAny(glyph, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ") {
		return BulletAlpha
	}

	// Default to disc/circle/square bullets
	return BulletDash
}
