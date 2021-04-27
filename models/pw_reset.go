package models

import (
	"gorm.io/gorm"
	"picapp/hash"
	"picapp/rand"
)

type pwReset struct {
	gorm.Model
	UserId    uint   `gorm:"not null"`
	Token     string `gorm:"-"`
	TokenHash string `gorm:"not null;uniqueIndex"`
}

type pwResetDB interface {
	ByToken(token string) (*pwReset, error)
	Create(pwr *pwReset) error
	Delete(id uint) error
}

func newPwResetValidtory(pwrdb pwResetDB, hmac hash.HMAC) *pwResetValidator {
	return &pwResetValidator{
		pwResetDB: pwrdb,
		hmac:      hmac,
	}
}

type pwResetValidator struct {
	pwResetDB
	hmac hash.HMAC
}

func (pwrv *pwResetValidator) ByToken(token string) (*pwReset, error) {
	pwr := pwReset{Token: token}
	err := runPwResetValFns(&pwr, pwrv.hmacToken)
	if err != nil {
		return nil, err
	}
	return pwrv.pwResetDB.ByToken(pwr.TokenHash)
}

func (pwrv *pwResetValidator) Create(pwr *pwReset) error {
	err := runPwResetValFns(pwr,
		pwrv.requireUserID,
		pwrv.setDefaultToken,
		pwrv.hmacToken,
	)
	if err != nil {
		return err
	}
	return pwrv.pwResetDB.Create(pwr)
}

func (pwrv *pwResetValidator) Delete(id uint) error {
	if id <= 0 {
		return ErrIDInvalid
	}
	return pwrv.pwResetDB.Delete(id)
}

type pwResetGorm struct {
	db *gorm.DB
}

func (pwrg *pwResetGorm) ByToken(tokenHash string) (*pwReset, error) {
	var pwr pwReset
	err := first(pwrg.db.Where("token_hash = ?", tokenHash), &pwr)
	if err != nil {
		return nil, err
	}
	return &pwr, nil
}

func (pwrg *pwResetGorm) Create(pwr *pwReset) error {
	return pwrg.db.Create(pwr).Error
}

func (pwrg *pwResetGorm) Delete(id uint) error {
	pwr := pwReset{Model: gorm.Model{ID: id}}
	return pwrg.db.Delete(&pwr).Error
}

// VALIDATION

func (pwrv *pwResetValidator) requireUserID(pwr *pwReset) error {
	if pwr.UserId <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (pwrv *pwResetValidator) setDefaultToken(pwr *pwReset) error {
	if pwr.Token != "" {
		return nil
	}
	token, err := rand.RememberToken()
	if err != nil {
		return err
	}
	pwr.Token = token
	return nil
}

func (pwrv *pwResetValidator) hmacToken(pwr *pwReset) error {
	if pwr.Token == "" {
		return nil
	}
	pwr.TokenHash = pwrv.hmac.Hash(pwr.Token)
	return nil
}

type pwResetValFn func(*pwReset) error

func runPwResetValFns(pwr *pwReset, fns ...pwResetValFn) error {
	for _, fn := range fns {
		if err := fn(pwr); err != nil {
			return err
		}
	}
	return nil
}