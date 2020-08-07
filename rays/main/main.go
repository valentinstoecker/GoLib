package main

import (
	"image/color"
	"image/png"
	"os"

	"github.com/valentinstoecker/GoLib/math/vector"
	"github.com/valentinstoecker/GoLib/rays"
)

func main() {
	trm := rays.RayMarcher{
		Epsilon:  0.001,
		FOVRatio: 01,
		Width:    1920,
		Height:   1080,
		MaxDist:  10,
		MaxIt:    1000,
	}

	col := color.RGBA{
		R: 255,
		G: 255,
		B: 255,
		A: 255,
	}

	cube := rays.Cap{
		rays.Plane{
			Norm: vector.VecX(1),
			Pos:  vector.VecX(1),
			Col:  col,
		},
		rays.Plane{
			Norm: vector.VecX(-1),
			Pos:  vector.VecX(-1),
			Col:  col,
		},
		rays.Plane{
			Norm: vector.VecY(1),
			Pos:  vector.VecY(-1),
			Col:  col,
		},
		rays.Plane{
			Norm: vector.VecY(-1),
			Pos:  vector.VecY(-2),
			Col:  col,
		},
		rays.Plane{
			Norm: vector.VecZ(1),
			Pos:  vector.VecZ(6),
			Col:  col,
		},
		rays.Plane{
			Norm: vector.VecZ(-1),
			Pos:  vector.VecZ(4),
			Col:  col,
		},
	}

	cds := make(rays.ColorDisters, 3)
	pl := rays.Reflect(
		rays.Plane{
			Col: color.White,
			Pos: vector.VecY(-2),
			Norm: vector.Vec3{
				0, 1, 0,
			}.Normalize(),
		}, 0.2, cds)
	sp1 := rays.Reflect(
		rays.Cap{
			rays.Sphere{
				Pos: vector.Vec3{0.5, 0, 6},
				Rad: 1,
				Col: color.NRGBA{
					R: 128,
					G: 128,
					B: 128,
					A: 255,
				},
			}, rays.Sphere{
				Pos: vector.Vec3{-0.5, 0, 5},
				Rad: 1,
				Col: color.NRGBA{
					R: 128,
					G: 128,
					B: 128,
					A: 255,
				},
			},
		}, 0.0, cds)
	// sp2 := rays.Reflect(rays.Sphere{
	// 	Pos: vector.Vec3{0.75, 0, 5},
	// 	Rad: 0.5,
	// 	Col: color.NRGBA{
	// 		R: 128,
	// 		G: 128,
	// 		B: 128,
	// 		A: 255,
	// 	},
	// }, 0.4, cds)
	cds[0] = pl
	cds[1] = sp1
	//cds[2] = sp2
	cds[2] = cube

	lights := []rays.Light{{
		Pos: vector.Vec3{
			-3, 5, 2,
		},
		Lum: 30,
	}, {
		Pos: vector.Vec3{
			3, 5, 2,
		},
		Lum: 20,
	}}

	l := rays.Lighter{
		ColorDister: cds,
		Lights:      lights,
	}

	f, err := os.Create("rm.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, trm.March(l))
	if err != nil {
		panic(err)
	}
	f.Close()
}
