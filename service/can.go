// +build ignore

package service

type (
	CanType uint8
)

const (
	CanLike CanType = iota
)

func (s *service) Can(t CanType) error {
	switch t {
	case CanLike:

	}
}

func (s *service) canLike() {

}
