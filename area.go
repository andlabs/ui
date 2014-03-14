// 14 march 2014

package ui

import (
	"sync"
	"image"
)

// Area represents a blank canvas upon which programs may draw anything and receive arbitrary events from the user.
// An Area has an explicit size, represented in pixels, that may be different from the size shown in its Window; scrollbars are placed automatically should they be needed.
// To facilitate development and debugging, for the time being, Areas have a fixed size of 320x240 and only work on GTK+.
type Area struct {
	// Paint is signaled when the Area needs to be redrawn.
	// You MUST handle Paint signals; failure to do so will result in the UI task hanging.
	// See the documentation of PaintRequest for details.
	Paint		chan PaintRequest

	lock		sync.Mutex
	created	bool
	sysData	*sysData
}

// PaintRequest contains the information needed to redraw an Area.
// On each Paint event, an Area will receive a full request on its Paint channel.
// It must send something back on Out in order to complete the painting.
// Example:
// 	imgFromFile, _, err := image.Decode(file)
// 	if err != nil { panic(err) }
// 	img := image.NewNRGBA(imgFromFile.Rect)
// 	draw.Draw(img, img.Rect, imgFromFile, image.ZP, draw.Over)
// 	for req := range area.Paint {
// 		req.Out <- img.SubImage(req.Rect).(*image.NRGBA)
// 	}
type PaintRequest struct {
	// Rect is the clip rectangle of the whole Area that needs to be redrawn.
	// The image sent on Out must have the same size as Rect (but does not have to have the same Rect.Min/Rect.Max points).
	Rect		image.Rectangle

	// Out is where you send the image to draw.
	// Only one image per PaintRequest may be sent; you must send an image.
	// Do not close Out; the package will do this itself.
	Out		chan<- *image.NRGBA
}

// NewArea creates a new Area.
func NewArea() *Area {
	return &Area{
		sysData:	mksysdata(c_area),
		Paint:	make(chan PaintRequest),
	}
}

func (a *Area) make(window *sysData) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	a.sysData.paint = a.Paint
	err := a.sysData.make("", window)
	if err != nil {
		return err
	}
	a.created = true
	return nil
}

func (a *Area) setRect(x int, y int, width int, height int, winheight int) error {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.sysData.setRect(x, y, width, height, winheight)
}

func (a *Area) preferredSize() (width int, height int) {
	a.lock.Lock()
	defer a.lock.Unlock()

	return a.sysData.preferredSize()
}
