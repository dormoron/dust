package dust

import (
	"github.com/nothingZero/dust/internal/valuer"
	"github.com/nothingZero/dust/model"
)

type core struct {
	model   *model.Model
	dialect Dialect
	creator valuer.Creator
	r       model.Registry
	mils    []Middleware
}
