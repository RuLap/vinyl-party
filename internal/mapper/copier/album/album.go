package album

import (
	"github.com/jinzhu/copier"
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.AlbumCreateDTO) (entity.Album, error) {
	var album entity.Album
	err := copier.Copy(&album, &dto)
	if err != nil {
		return entity.Album{}, err
	}

	return album, nil
}

func EntityToShortInfoDTO(album entity.Album) (dto.AlbumShortInfoDTO, error) {
	var shortInfoDTO dto.AlbumShortInfoDTO
	err := copier.Copy(&shortInfoDTO, &album)
	if err != nil {
		return dto.AlbumShortInfoDTO{}, err
	}

	return shortInfoDTO, nil
}

func EntityToInfoDTO(album entity.Album) (dto.AlbumInfoDTO, error) {
	var infoDTO dto.AlbumInfoDTO
	err := copier.Copy(&infoDTO, &album)
	if err != nil {
		return dto.AlbumInfoDTO{}, err
	}

	return infoDTO, nil
}
