package bbox

import (
)

type Point interface {
     Latitude() float64
     Longitude() float64
}

type BoundingBox interface {
     SW() Pt
     NE() Pt
     South() float64
     North() float64
     East() float64
     West() float64
}

