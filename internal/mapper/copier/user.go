package copier

import (
	"github.com/jinzhu/copier"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func RegisterDTOToEntity(dto dto.UserRegisterDTO) (entity.User, error) {
	var user entity.User
	err := copier.Copy(&user, &dto)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func LoginDTOToEntity(dto dto.UserLoginDTO) (entity.User, error) {
	var user entity.User
	err := copier.Copy(&user, &dto)
	if err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func EntityToShortDTO(user entity.User) (dto.UserShortInfoDTO, error) {
	var shortInfoDTO dto.UserShortInfoDTO
	err := copier.Copy(&shortInfoDTO, &user)
	if err != nil {
		return dto.UserShortInfoDTO{}, err
	}

	return shortInfoDTO, nil
}

func EntityToInfoDTO(user entity.User) (dto.UserInfoDTO, error) {
	var infoDto dto.UserInfoDTO
	err := copier.Copy(&infoDto, &user)
	if err != nil {
		return dto.UserInfoDTO{}, err
	}

	return infoDto, nil
}
