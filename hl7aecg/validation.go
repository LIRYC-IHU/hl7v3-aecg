package hl7aecg

import (
	"context"

	"github.com/LIRYC-IHU/hl7v3-aecg/hl7aecg/types"
)

type Validator interface {
	Validate(ctx context.Context, vctx *types.ValidationContext) error
}

func (e *Hl7xml) Validate() error {
	e.validateAll(&e.HL7AEcg)
	return e.vctx.GetError()
}

func (e *Hl7xml) validateAll(objs ...Validator) {
	for _, obj := range objs {
		obj.Validate(e.ctx, e.vctx)
	}
}
