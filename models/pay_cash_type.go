package models

func (d_o *DbOrm) PayCashType(cash_code string) PayCashType {
	var c_list PayCashType
	err := d_o.GDb.Model(&PayCashType{}).Where("code = ?", cash_code).First(&c_list)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayCashType{}).Where("code = ?", cash_code).First(&c_list)
	}

	return c_list
}
