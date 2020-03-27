package models

func (d_o *DbOrm) PayClass() []PayClass {
	var p_class []PayClass
	err := d_o.GDb.Model(&PayClass{}).Find(&p_class)

	if err != nil {
		d_o.getSqlDb()
		d_o.GDb.Model(&PayClass{}).Find(&p_class)
	}

	return p_class
}
