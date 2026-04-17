package handler

import (
	_ "github.com/parvez3019/go-swagger3/model"
)

// @Title Get balance
// @Description Get credit balance
// @Accept json
// @Produce json
// @Success 200 {object} model.ResponseData[model.CreditResponse]
// @Failure 400 {object} model.ResponseData[string]
// @Router /v1/profile/balance [get]
func GetBalance() {
}
