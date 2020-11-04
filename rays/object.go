package rays

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"

	"github.com/valentinstoecker/GoLib/math/vector"
)

// Dister defines objects that have a distance function
type Dister interface {
	Dist(vector.Vec3) float64
}

// Disters slice of Dister
type Disters []Dister

// Colorer defines objects that have a color function
type Colorer interface {
	Color(Ray) color.Color
}

// ColorDister Colorer + Dister
type ColorDister interface {
	Colorer
	Dister
}

// ColorDisters slice of ColorDister
type ColorDisters []ColorDister

type Cap ColorDisters

func (c Cap) Dist(v vector.Vec3) float64 {
	max := float64(-1000)
	for _, cd := range c {
		max = math.Max(max, cd.Dist(v))
	}
	return max
}

func (c Cap) Color(r Ray) color.Color {
	return ColorDisters(c).Color(r)
}

//Dist -> Dist of Disters
func (ds Disters) Dist(v vector.Vec3) float64 {
	d := math.MaxFloat64
	for _, cd := range ds {
		d = math.Min(cd.Dist(v), d)
	}
	return d
}

// Dist of ColorDisters
func (cds ColorDisters) Dist(v vector.Vec3) float64 {
	d := math.MaxFloat64
	for _, cd := range cds {
		d = math.Min(cd.Dist(v), d)
	}
	return d
}

// Color of ColorDisters
func (cds ColorDisters) Color(r Ray) color.Color {
	minD := math.MaxFloat64
	var minCD ColorDister
	for _, cd := range cds {
		d := cd.Dist(r.Pos)
		if d < minD {
			minCD = cd
			minD = d
		}
	}
	if minCD == nil {
		return color.RGBA{
			R: 127,
			G: 127,
			B: 255,
			A: 255,
		}
	}
	return minCD.Color(r)
}

// Sphere defines Sphere
type Sphere struct {
	Pos vector.Vec3
	Rad float64
	Col color.Color
}

// Dist gets distance to the sphere
func (s Sphere) Dist(v vector.Vec3) float64 {
	return s.Pos.Dist(v) - s.Rad
}

// Color gets Color of sphere
func (s Sphere) Color(_ Ray) color.Color {
	return s.Col
}

// Ray position with direction
type Ray struct {
	Pos vector.Vec3
	Dir vector.Vec3
	X   int
	Y   int
}

// March marches a Ray
func (r Ray) March(cd ColorDister, maxIt int, epsilon float64, maxDist float64) color.Color {
	for i := 0; i < maxIt; i++ {
		d := cd.Dist(r.Pos)
		if d < epsilon {
			return cd.Color(r)
		}
		if d > maxDist {
			break
		}
		r.Pos = r.Pos.Add(r.Dir.Mult(d))
	}
	return color.RGBA{
		R: 127,
		G: 127,
		B: 255,
		A: 255,
	}
}

// MarchVec ray march and return collision vector
func (r Ray) MarchVec(cd Dister, maxIt int, epsilon float64, maxDist float64) vector.Vec3 {
	for i := 0; i < maxIt; i++ {
		d := cd.Dist(r.Pos)
		if d < epsilon {
			return r.Pos
		}
		if d > maxDist {
			break
		}
		r.Pos = r.Pos.Add(r.Dir.Mult(d))
	}
	return vector.Vec3{}
}

// RayMarcher describes a ray marcher
type RayMarcher struct {
	Epsilon  float64
	MaxIt    int
	MaxDist  float64
	Width    int
	Height   int
	FOVRatio float64
	Pos      vector.Vec3
}

// Plane a 3d plane
type Plane struct {
	Pos  vector.Vec3
	Norm vector.Vec3
	Col  color.Color
}

// Dist dist of plane
func (p Plane) Dist(v vector.Vec3) float64 {
	return v.Sub(p.Pos).Dot(p.Norm)
}

// Color of Plane
func (p Plane) Color(r Ray) color.Color {
	return p.Col
}

// Light a point light src
type Light struct {
	Pos vector.Vec3
	Lum float64
}

func (l Light) Dist(v vector.Vec3) float64 {
	return l.Pos.Dist(v)
}

type Lighter struct {
	ColorDister
	Lights []Light
}

func clamp(v float64) float64 {
	return math.Max(0, math.Min(v, 0.99))
}

func (l Lighter) Color(r Ray) color.Color {
	var lum float64
	for _, light := range l.Lights {
		sbPos := r.Pos.Sub(r.Dir.Mult(0.001))
		lr := Ray{
			Pos: sbPos,
			Dir: light.Pos.Sub(sbPos).Normalize(),
		}
		cds := Disters{light, l.ColorDister}
		p := lr.MarchVec(cds, 1000, 0.001, 100)
		if p.Dist(light.Pos) < 0.002 {
			v := light.Pos.Sub(r.Pos)
			lum += (light.Lum / v.Dot(v)) * math.Abs(Normal(cds, r.Pos, 0.001).Dot(lr.Dir))
		}
	}
	cr, cg, cb, ca := l.ColorDister.Color(r).RGBA()
	fr := clamp(float64(cr) / 65535 * lum)
	fg := clamp(float64(cg) / 65535 * lum)
	fb := clamp(float64(cb) / 65535 * lum)
	return color.RGBA{
		R: uint8(fr * 255),
		G: uint8(fg * 255),
		B: uint8(fb * 255),
		A: uint8(ca / 256),
	}
}

// XZPlane <-
type XZPlane struct {
	Y   float64
	Col color.Color
}

// Color Black
func (p XZPlane) Color(r Ray) color.Color {
	x, z := int(math.Floor(r.Pos[0]))%2, int(math.Floor(r.Pos[2]))%2
	if (x == 0) != (z == 0) {
		return color.Gray{Y: 127}
	}
	return p.Col
}

// Dist (v.y - p)
func (p XZPlane) Dist(v vector.Vec3) float64 {
	return v[1] - p.Y
}

// Normal calcs normal of cd at v
func Normal(cd Dister, v vector.Vec3, epsilon float64) vector.Vec3 {
	ex, ey, ez := vector.VecX(epsilon), vector.VecY(epsilon), vector.VecZ(epsilon)
	return vector.Vec3{
		cd.Dist(v.Add(ex)) - cd.Dist(v.Sub(ex)),
		cd.Dist(v.Add(ey)) - cd.Dist(v.Sub(ey)),
		cd.Dist(v.Add(ez)) - cd.Dist(v.Sub(ez)),
	}.Normalize()
}

// March stuffs
func (rm RayMarcher) March(cd ColorDister) image.Image {
	img := image.NewRGBA(image.Rect(0, 0, rm.Width, rm.Height))
	var il sync.Mutex
	var wg sync.WaitGroup
	rs := make(chan Ray, 256)
	wg.Add(runtime.NumCPU())
	for i := 0; i < runtime.NumCPU(); i++ {
		go func() {
			for r := range rs {
				c := r.March(cd, rm.MaxIt, rm.Epsilon, rm.MaxDist)
				il.Lock()
				img.Set(r.X, rm.Height-r.Y-1, c)
				il.Unlock()
			}
			wg.Done()
		}()
	}
	ar := float64(rm.Width) / float64(rm.Height)
	wf := (2 * rm.FOVRatio) / float64(rm.Width)
	for x := 0; x < rm.Width; x++ {
		for y := 0; y < rm.Height; y++ {
			r := Ray{
				Dir: vector.Vec3{
					float64(x)*wf - rm.FOVRatio,
					float64(y)*wf - rm.FOVRatio/ar,
					1,
				}.Normalize(),
				X:   x,
				Y:   y,
				Pos: rm.Pos,
			}
			rs <- r
		}
	}
	close(rs)
	wg.Wait()
	return img
}

type reflect struct {
	ColorDister,
	toReflect ColorDister
	refVal float64
}

func (r reflect) Dist(v vector.Vec3) float64 {
	return r.ColorDister.Dist(v)
}

func blend(a, b color.Color, bf float64) color.RGBA {
	bs := func(a, b, bf uint32) uint8 {
		return uint8((a*bf + b*(65535-bf)) / (65536*265 - 1))
	}

	ar, ag, ab, aa := a.RGBA()
	br, bg, bb, ba := b.RGBA()
	ibf := uint32(bf * 65535)
	return color.RGBA{
		R: bs(ar, br, ibf),
		G: bs(ag, bg, ibf),
		B: bs(ab, bb, ibf),
		A: bs(aa, ba, ibf),
	}
}

func (ref reflect) Color(r Ray) color.Color {
	n := Normal(ref.ColorDister, r.Pos, 0.001)
	return blend(Ray{
		Pos: r.Pos.Sub(r.Dir.Mult(0.002)),
		Dir: n.Mult(2 * n.Dot(r.Dir)).Sub(r.Dir).Mult(-1),
	}.March(ref.toReflect, 1000, 0.001, 100), ref.ColorDister.Color(r), ref.refVal)
}

func Reflect(refObj ColorDister, refVal float64, toReflect ColorDister) ColorDister {
	return reflect{
		ColorDister: refObj,
		refVal:      math.Max(0, math.Min(refVal, 1)),
		toReflect:   toReflect,
	}
}
