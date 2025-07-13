package deck

import (
	"context"
	"fmt"
	"io"

	"github.com/k1LoW/errors"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/slides/v1"
)

// Create Google Slides presentation.
func Create(ctx context.Context) (_ *Deck, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d := &Deck{
		styles: map[string]*slides.TextStyle{},
		shapes: map[string]*slides.ShapeProperties{},
	}
	if err := d.initialize(ctx); err != nil {
		return nil, err
	}
	title := "Untitled"
	file := &drive.File{
		Name:     title,
		MimeType: "application/vnd.google-apps.presentation",
	}
	f, err := d.driveSrv.Files.Create(file).Do()
	if err != nil {
		return nil, err
	}
	d.id = f.Id
	if err := d.refresh(ctx); err != nil {
		return nil, err
	}
	return d, nil
}

// CreateFrom creates a new Deck from the presentation ID.
func CreateFrom(ctx context.Context, id string) (_ *Deck, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d := &Deck{
		styles: map[string]*slides.TextStyle{},
		shapes: map[string]*slides.ShapeProperties{},
	}
	if err := d.initialize(ctx); err != nil {
		return nil, err
	}
	// copy presentation
	file := &drive.File{
		Name:     "Untitled",
		MimeType: "application/vnd.google-apps.presentation",
	}
	f, err := d.driveSrv.Files.Copy(id, file).Do()
	if err != nil {
		return nil, err
	}
	d.id = f.Id
	if err := d.refresh(ctx); err != nil {
		return nil, err
	}
	// delete all slides
	if err := d.DeletePageAfter(ctx, -1); err != nil {
		return nil, err
	}
	// create first slide
	if err := d.createPage(ctx, 0, &Slide{
		Layout: d.defaultTitleLayout,
	}); err != nil {
		return nil, err
	}
	return d, nil
}

// Export the presentation as PDF.
func (d *Deck) Export(ctx context.Context, w io.Writer) (err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	req, err := d.driveSrv.Files.Export(d.id, "application/pdf").Context(ctx).Download()
	if err != nil {
		return err
	}
	if err := req.Write(w); err != nil {
		return fmt.Errorf("unable to create PDF file: %w", err)
	}
	return nil
}

// List Google Slides presentations.
func List(ctx context.Context) (_ []*Presentation, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	d := &Deck{
		styles: map[string]*slides.TextStyle{},
		shapes: map[string]*slides.ShapeProperties{},
	}
	if err := d.initialize(ctx); err != nil {
		return nil, err
	}
	var presentations []*Presentation

	r, err := d.driveSrv.Files.List().Q("mimeType='application/vnd.google-apps.presentation'").Fields("files(id, name)").Do()
	if err != nil {
		return nil, err
	}

	for _, f := range r.Files {
		presentations = append(presentations, &Presentation{
			ID:    f.Id,
			Title: f.Name,
		})
	}

	return presentations, nil
}

// List Google Slides presentations.
func (d *Deck) List() (_ []*Presentation, err error) {
	defer func() {
		err = errors.WithStack(err)
	}()
	var presentations []*Presentation

	r, err := d.driveSrv.Files.List().Q("mimeType='application/vnd.google-apps.presentation'").Fields("files(id, name)").Do()
	if err != nil {
		return nil, err
	}

	for _, f := range r.Files {
		presentations = append(presentations, &Presentation{
			ID:    f.Id,
			Title: f.Name,
		})
	}

	return presentations, nil
}

// ListLayouts lists layouts of the presentation.
func (d *Deck) ListLayouts() []string {
	var layouts []string
	for _, l := range d.presentation.Layouts {
		layouts = append(layouts, l.LayoutProperties.DisplayName)
	}
	return layouts
}
