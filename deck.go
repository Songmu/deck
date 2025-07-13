package deck

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/k1LoW/deck/config"
	"github.com/k1LoW/errors"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/slides/v1"
)

const (
	layoutNameForStyle             = "style"
	styleCode                      = "code"
	styleBold                      = "bold"
	styleItalic                    = "italic"
	styleLink                      = "link"
	styleBlockQuote                = "blockquote"
	defaultCodeFontFamily          = "Noto Sans Mono"
	descriptionImageFromMarkdown   = "Image generated from markdown"
	descriptionTextboxFromMarkdown = "Textbox generated from markdown"
)

type Deck struct {
	id                 string
	srv                *slides.Service
	driveSrv           *drive.Service
	presentation       *slides.Presentation
	defaultTitleLayout string
	defaultLayout      string
	styles             map[string]*slides.TextStyle
	shapes             map[string]*slides.ShapeProperties
	logger             *slog.Logger
}

type Option func(*Deck) error

func WithPresentationID(id string) Option {
	return func(d *Deck) error {
		d.id = id
		return nil
	}
}

func WithLogger(logger *slog.Logger) Option {
	return func(d *Deck) error {
		d.logger = logger
		return nil
	}
}

type placeholder struct {
	objectID string
	x        float64
	y        float64
}

type bulletRange struct {
	bullet Bullet
	start  int64
	end    int64
}

type textBox struct {
	paragraphs   []*Paragraph
	fromMarkdown bool
}

type actionDetail struct {
	ActionType  actionType `json:"action_type"`
	Titles      []string   `json:"titles,omitempty"`
	Index       *int       `json:"index,omitempty"`
	MoveToIndex *int       `json:"move_to_index,omitempty"`
}

// Presentation represents a Google Slides presentation.
type Presentation struct {
	ID    string
	Title string
}

// New creates a new Deck.
func New(ctx context.Context, opts ...Option) (_ *Deck, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d := &Deck{
		styles: map[string]*slides.TextStyle{},
		shapes: map[string]*slides.ShapeProperties{},
	}
	for _, opt := range opts {
		if err := opt(d); err != nil {
			return nil, err
		}
	}
	if err := d.initialize(ctx); err != nil {
		return nil, err
	}
	if err := d.refresh(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

// Delete deletes a Google Slides presentation by ID.
func Delete(ctx context.Context, id string) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d := &Deck{
		styles: map[string]*slides.TextStyle{},
	}
	if err := d.initialize(ctx); err != nil {
		return err
	}
	if err := d.driveSrv.Files.Delete(id).Context(ctx).Do(); err != nil {
		return fmt.Errorf("failed to delete presentation: %w", err)
	}
	return nil
}

// ID returns the ID of the presentation.
func (d *Deck) ID() string {
	return d.id
}

// UpdateTitle updates the title of the presentation.
func (d *Deck) UpdateTitle(ctx context.Context, title string) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	file := &drive.File{
		Name: title,
	}
	if _, err := d.driveSrv.Files.Update(d.id, file).Context(ctx).Do(); err != nil {
		return err
	}
	return nil
}

func (d *Deck) initialize(ctx context.Context) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	if d.logger == nil {
		d.logger = slog.New(slog.NewJSONHandler(io.Discard, nil))
	}
	if err := os.MkdirAll(config.DataHomePath(), 0700); err != nil {
		return err
	}
	if err := os.MkdirAll(config.StateHomePath(), 0700); err != nil {
		return err
	}

	creds := filepath.Join(config.DataHomePath(), "credentials.json")
	b, err := os.ReadFile(creds)
	if err != nil {
		return err
	}

	config, err := google.ConfigFromJSON(b, slides.PresentationsScope, slides.DriveScope)
	if err != nil {
		return err
	}

	client, err := d.getHTTPClient(ctx, config)
	if err != nil {
		return err
	}
	srv, err := slides.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	srv.UserAgent = userAgent
	d.srv = srv
	driveSrv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return err
	}
	driveSrv.UserAgent = userAgent
	d.driveSrv = driveSrv
	return nil
}

func (d *Deck) createPage(ctx context.Context, index int, slide *Slide) (err error) {
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

	// create new page
	req := &slides.BatchUpdatePresentationRequest{
		Requests: []*slides.Request{
			{
				CreateSlide: &slides.CreateSlideRequest{
					InsertionIndex: int64(index),
					SlideLayoutReference: &slides.LayoutReference{
						LayoutId: layout.ObjectId,
					},
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

func (d *Deck) refresh(ctx context.Context) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	presentation, err := d.srv.Presentations.Get(d.id).Context(ctx).Do()
	if err != nil {
		return err
	}
	d.presentation = presentation

	// set default layouts and detect style
	for _, l := range d.presentation.Layouts {
		layout := l.LayoutProperties.Name
		switch {
		case strings.HasPrefix(layout, "TITLE_AND_BODY"):
			if d.defaultLayout == "" {
				d.defaultLayout = l.LayoutProperties.DisplayName
			}
		case strings.HasPrefix(layout, "TITLE"):
			if d.defaultTitleLayout == "" {
				d.defaultTitleLayout = l.LayoutProperties.DisplayName
			}
		}

		if l.LayoutProperties.DisplayName == layoutNameForStyle {
			for _, e := range l.PageElements {
				if e.Shape == nil || e.Shape.Text == nil {
					continue
				}
				for _, t := range e.Shape.Text.TextElements {
					if t.TextRun == nil {
						continue
					}
					className := strings.Trim(t.TextRun.Content, " \n")
					if className == "" {
						continue
					}
					d.styles[className] = t.TextRun.Style
					d.shapes[className] = e.Shape.ShapeProperties
				}
			}
		}
	}
	if d.defaultTitleLayout == "" {
		d.defaultTitleLayout = d.presentation.Layouts[0].LayoutProperties.DisplayName
	}
	if d.defaultLayout == "" {
		if len(d.presentation.Layouts) > 1 {
			d.defaultLayout = d.presentation.Layouts[1].LayoutProperties.DisplayName
		} else {
			d.defaultLayout = d.presentation.Layouts[0].LayoutProperties.DisplayName
		}
	}

	return nil
}

func ptrInt64(i int64) *int64 {
	return &i
}
