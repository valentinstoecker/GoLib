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

	cds := rays.ColorDisters{}

	floor := rays.Plane{
		Norm: vector.VecY(1),
		Pos:  vector.VecY(-2),
		Col:  col,
	}

	light := rays.Lighter{
		ColorDister: &cds,
		Lights: []rays.Light{{
			Pos: vector.Vec3{
				2, 4, 8,
			},
			Lum: 30,
		}},
	}
	floorReflect := rays.Reflect(
		floor,
		0.2,
		&light,
	)

	sphere := rays.Sphere{
		Pos: vector.Vec3{
			0, -1, 10,
		},
		Rad: 1,
		Col: color.RGBA{
			R: 255,
			A: 255,
		},
	}

	cds = append(cds, floorReflect, sphere)

	f, err := os.Create("rm.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(f, trm.March(light))
	if err != nil {
		panic(err)
	}
	f.Close()
}
