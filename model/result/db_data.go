package result

type DbData struct {
	data interface{}
}

func CreateDbData(data interface{}) DbData {
	return DbData{
		data,
	}
}

func (dat DbData) GetData() interface{} {
	return dat.data
}
