sed "s/AAA/$1/g
s/BBB/$2/g
s/CCC/$3/g
s/DDD/$4/g" <<\END
func (AAA *BBB) CCC() DDD {
	return AAA._CCC
}

func (AAA *BBB) setParent(p *controlParent) {
	basesetParent(AAA, p)
}

func (AAA *BBB) allocate(x int, y int, width int, height int, d *sizing) []*allocation {
	return baseallocate(AAA, x, y, width, height, d)
}

func (AAA *BBB) preferredSize(d *sizing) (width, height int) {
	return basepreferredSize(AAA, d)
}

func (AAA *BBB) commitResize(a *allocation, d *sizing) {
	basecommitResize(AAA, a, d)
}

func (AAA *BBB) getAuxResizeInfo(d *sizing) {
	basegetAuxResizeInfo(AAA, d)
}
END
