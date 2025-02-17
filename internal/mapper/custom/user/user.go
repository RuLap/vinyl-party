package user

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func RegisterDTOToEntity(dto *dto.UserRegisterDTO) entity.User {
	return entity.User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		Password:  dto.Password,
		AvatarUrl: dto.AvatarUrl,
	}
}

func LoginDTOToEntity(dto *dto.UserLoginDTO) entity.User {
	return entity.User{
		Email:    dto.Email,
		Password: dto.Password,
	}
}

func EntityToShortInfoDTO(user *entity.User) dto.UserShortInfoDTO {
	return dto.UserShortInfoDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		AvatarUrl: user.AvatarUrl,
	}
}

func EntityToInfoDTO(user *entity.User) dto.UserInfoDTO {
	return dto.UserInfoDTO{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		AvatarUrl: user.AvatarUrl,
	}
}

func EntityToLoginRepsponseDTO(user *entity.User, token string) dto.UserLoginResponseDTO {
	return dto.UserLoginResponseDTO{
		User:  EntityToInfoDTO(user),
		Token: token,
	}
}
