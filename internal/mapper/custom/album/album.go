package album

import (
	"vinyl-party/internal/dto"
	"vinyl-party/internal/entity"
)

func CreateDTOToEntity(dto dto.AlbumCreateDTO) entity.Album {
	return entity.Album{
		Title:      dto.Title,
		Artist:     dto.Artist,
		CoverUrl:   dto.CoverUrl,
		SpotifyUrl: dto.SpotifyUrl,
	}
}

func EntityToShortInfoDTO(album entity.Album) dto.AlbumShortInfoDTO {
	return dto.AlbumShortInfoDTO{
		ID:         album.ID,
		Title:      album.Title,
		Artist:     album.Artist,
		CoverUrl:   album.CoverUrl,
		SpotifyUrl: album.SpotifyUrl,
	}
}

func EntityToInfoDTO(album *entity.Album, ratingDTOs []dto.RatingInfoDTO, averageRating int) dto.AlbumInfoDTO {
	return dto.AlbumInfoDTO{
		ID:            album.ID,
		Title:         album.Title,
		Artist:        album.Artist,
		CoverUrl:      album.CoverUrl,
		SpotifyUrl:    album.SpotifyUrl,
		Ratings:       ratingDTOs,
		AverageRating: averageRating,
	}
}

func SpotifyDTOToEntity(album *entity.SpotifyAlbum) entity.Album {
	return entity.Album{
		Title:      album.Name,
		Artist:     album.ArtistsString,
		CoverUrl:   album.CoverUrl,
		SpotifyUrl: album.Url,
	}
}
