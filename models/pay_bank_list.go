package models

func (d_o *DbOrm) PayBankList(is_mobile, cash_code, class_id string) PayBankList {
	var pay_bank PayBankList
	err := d_o.GDb.Model(&PayBankList{}).Where("website_type=? and cash_type_code=? and class_id=?", is_mobile, cash_code, class_id).First(&pay_bank)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayBankList{}).Where("website_type=? and cash_type_code=? and class_id=?", is_mobile, cash_code, class_id).First(&pay_bank)
	}

	return pay_bank
}

func (d_o *DbOrm) BankList(is_mobile, cash_code string) []PayBankList {
	var b_list []PayBankList
	err := d_o.GDb.Model(&PayBankList{}).Where("website_type=? and cash_type_code=? and class_id=1", is_mobile, cash_code).Find(&b_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayBankList{}).Where("website_type=? and cash_type_code=? and class_id=1", is_mobile, cash_code).Find(&b_list)
	}

	return b_list
}
