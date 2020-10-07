package vector

import (
	"fmt"
	"math"
)

// Vec3 vector with 3 elements
type Vec3 [3]float64

// VecX -> Vec3 with only x component
func VecX(x float64) Vec3 {
	return Vec3{
		x, 0, 0,
	}
}

// VecY -> Vec3 with only y component
func VecY(y float64) Vec3 {
	return Vec3{
		0, y, 0,
	}
}

// VecZ -> Vec3 with only z component
func VecZ(z float64) Vec3 {
	return Vec3{
		0, 0, z,
	}
}

// Add adds 2 Vec3s
func (v Vec3) Add(w Vec3) Vec3 {
	return Vec3{
		v[0] + w[0],
		v[1] + w[1],
		v[2] + w[2],
	}
}

// Mult multiplies Vec3 by scalar
func (v Vec3) Mult(r float64) Vec3 {
	return Vec3{
		v[0] * r,
		v[1] * r,
		v[2] * r,
	}
}

// Sub subtracts 2 Vec3s
func (v Vec3) Sub(w Vec3) Vec3 {
	return v.Add(w.Mult(-1))
}

// Dot calculates the dot product of 2 Vec3s
func (v Vec3) Dot(w Vec3) float64 {
	return v[0]*w[0] + v[1]*w[1] + v[2]*w[2]
}

// Norm calculates the euclidean Norm
func (v Vec3) Norm() float64 {
	return math.Sqrt(v.Dot(v))
}

// Dist calculates the euclidean distance of 2 Vec3s
func (v Vec3) Dist(w Vec3) float64 {
	return v.Sub(w).Norm()
}

// Cross calculates the cross-product of 2 Vec3s
func (v Vec3) Cross(w Vec3) Vec3 {
	return Vec3{
		v[1]*w[2] - v[2]*w[1],
		v[2]*w[0] - v[0]*w[2],
		v[0]*w[1] - v[1]*w[0],
	}
}

// Normalize normalizes the Vec3 (sets v.Norm() = 1)
func (v Vec3) Normalize() Vec3 {
	n := v.Norm()
	if n == 0 {
		return v
	}
	return v.Mult(1 / n)
}

func (v Vec3) String() string {
	return fmt.Sprintf("(%f, %f, %f)", v[0], v[1], v[2])
}
