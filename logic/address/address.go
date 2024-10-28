package address

import (
	"errors"
	"gomall/global"
	receive "gomall/interaction/receive/address"
	"gomall/models/address"
)

func GetAddressById(id int64, memberId int64) (result address.UmsMemberReceiveAddress, err error) {
	if err = global.Db.Model(&address.UmsMemberReceiveAddress{}).
		Where("member_id", memberId).Where("id", id).First(&result).Error; err != nil {
		return address.UmsMemberReceiveAddress{}, errors.New("根据id查询地址failed:" + err.Error())
	}
	return result, nil
}

func AddAddress(data *receive.AddAddressReqStruct) (err error) {
	newAddress := &address.UmsMemberReceiveAddress{
		MemberID:      data.MemberId,
		Name:          data.Name,
		PhoneNumber:   data.PhoneNumber,
		PostCode:      data.PostCode,
		Province:      data.Province,
		City:          data.City,
		Region:        data.Region,
		DetailAddress: data.DetailAddress,
		DefaultStatus: data.DefaultStatus,
	}
	return global.Db.Model(&address.UmsMemberReceiveAddress{}).Create(&newAddress).Error
}

// DeleteAddress 这里是真的删除地址项，不是软删除
func DeleteAddress(id int64, memberId int64) (err error) {
	return global.Db.Model(&address.UmsMemberReceiveAddress{}).Where("id = ? and member_id = ?", id, memberId).
		Delete(&address.UmsMemberReceiveAddress{}, id, memberId).Error
}

func List() (result []address.UmsMemberReceiveAddress, err error) {
	if err = global.Db.Model(address.UmsMemberReceiveAddress{}).Find(&result).Error; err != nil {
		return nil, errors.New("查询地址列表失败：" + err.Error())
	}
	return
}

func UpdateAddress(data *receive.UpdateAddressReqStruct, memberId int64) (err error) {
	newAddress := &address.UmsMemberReceiveAddress{
		ID:            data.Id,
		MemberID:      data.MemberId,
		Name:          data.Name,
		PhoneNumber:   data.PhoneNumber,
		DefaultStatus: data.DefaultStatus,
		PostCode:      data.PostCode,
		Province:      data.Province,
		City:          data.City,
		Region:        data.Region,
		DetailAddress: data.DetailAddress,
	}

	return global.Db.Model(&address.UmsMemberReceiveAddress{}).
		Where("id = ? and member_id = ?", data.Id, memberId).Updates(&newAddress).Error
}
