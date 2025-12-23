package types

import (
	"context"
)

func (c *Code[T, U]) ValidateCode(ctx context.Context, vctx *ValidationContext, key string) error {
	switch any(c).(type) {
	case *Code[ConfidentialityCode, string]:
		return validateConfidentialityCode(vctx, any(c).(*Code[ConfidentialityCode, string]))
	case *Code[ReasonCode, string]:
		return validateReasonCode(vctx, any(c).(*Code[ReasonCode, string]))
	case *Code[CodeRole, CodeSystemOID]:
		return validateTrialSubjectCode(vctx, any(c).(*Code[CodeRole, CodeSystemOID]))
	case *Code[GenderCode, CodeSystemOID]:
		return validateGenderCode(vctx, any(c).(*Code[GenderCode, CodeSystemOID]))
	default:
		return nil
	}
}

func validateConfidentialityCode(vctx *ValidationContext, code *Code[ConfidentialityCode, string]) error {
	if code.Code != CONFIDENTIALITY_CUSTOM && code.Code != CONFIDENTIALITY_INVESTIGATOR_BLINDED && code.Code != CONFIDENTIALITY_SPONSOR_BLINDED && code.Code != CONFIDENTIALITY_BOTH {
		vctx.AddError(ErrConfidentialityCode)
	}

	return nil
}

func validateReasonCode(vctx *ValidationContext, code *Code[ReasonCode, string]) error {
	if code.Code != REASON_PER_PROTOCOL && code.Code != REASON_NOT_IN_PROTOCOL && code.Code != REASON_WRONG_EVENT {
		vctx.AddError(ErrReasonCode)
	}
	return nil
}

func validateTrialSubjectCode(vctx *ValidationContext, code *Code[CodeRole, CodeSystemOID]) error {
	if code.Code != SUBJECT_ROLE_SCREENING && code.Code != SUBJECT_ROLE_ENROLLED {
		vctx.AddError(ErrTrialSubjectCode)
	}
	return nil
}

func validateGenderCode(vctx *ValidationContext, code *Code[GenderCode, CodeSystemOID]) error {
	if code.Code != GENDER_FEMALE && code.Code != GENDER_MALE && code.Code != GENDER_UNDIFFERENTIATED {
		vctx.AddError(ErrRaceCode)
	}
	return nil
}
