package service

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
	"vinyl-party/internal/config"
	"vinyl-party/internal/entity"
)

type SpotifyService interface {
	GetAlbumFromLink(link string) (*entity.SpotifyAlbum, error)
}

type spotifyService struct {
	clientID     string
	clientSecret string
	client       *http.Client
	token        *entity.SpotifyAccessToken
}

func NewSpotifyService(proxyServer config.ProxyServer, spotifyCred config.SpotifyCredentials) (SpotifyService, error) {
	client, err := getProxyClient(proxyServer)
	if err != nil {
		slog.Error("failed to init proxy client", "error", err)
		return nil, err
	}

	service := &spotifyService{
		clientID:     spotifyCred.ClientID,
		clientSecret: spotifyCred.ClientSecret,
		client:       client,
	}

	if err := service.getAuthToken(); err != nil {
		slog.Error("failed to get an auth token", "error", err)
		return nil, err
	}

	return service, nil
}

func (s *spotifyService) GetAlbumFromLink(link string) (*entity.SpotifyAlbum, error) {
	albumID, err := s.parseAlbumIDFromLink(link)
	if err != nil {
		slog.Error("failed to parse id from spotify link", "link", link, "error", err)
		return nil, err
	}

	album, err := s.getAlbumByID(albumID)
	if err != nil {
		slog.Error("failed to get album", "albumID", album, "error", err)
		return nil, err
	}
	album.Url = link

	for _, image := range album.Images {
		if image.Width == 300 {
			album.CoverUrl = image.Url
		}
	}
	if album.CoverUrl == "" {
		album.CoverUrl = album.Images[0].Url
	}

	artistsStr := album.Artists[0].Name
	for i := 1; i < len(album.Artists); i++ {
		artistsStr += " ," + album.Artists[i].Name
	}
	album.ArtistsString = artistsStr

	return album, nil
}

func getProxyClient(proxyServer config.ProxyServer) (*http.Client, error) {
	proxy, err := url.Parse(proxyServer.Address)
	if err != nil {
		slog.Error("failed to parse proxy url", "error", err)
		return nil, err
	}

	proxy.User = url.UserPassword(proxyServer.Username, proxyServer.Password)

	transport := &http.Transport{
		Proxy: http.ProxyURL(proxy),
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   proxyServer.Timeout,
	}

	return client, nil
}

func (s *spotifyService) getAuthToken() error {
	authHeader := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", s.clientID, s.clientSecret)))

	data := url.Values{}
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(data.Encode()))
	if err != nil {
		slog.Error("failed to create SotifyAPI token request", "error", err)
		return err
	}

	req.Header.Set("Authorization", "Basic "+authHeader)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.client.Do(req)
	if err != nil {
		slog.Error("failed to do SotifyAPI token request", "error", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		slog.Error("failed to recieve SotifyAPI token request", "error", errors.New("failed to recieve token: "+string(body)))
		return errors.New("failed to recieve token: " + string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&s.token); err != nil {
		slog.Error("failed to decode Spotify API token", "error", err)
		return err
	}

	go s.autoRefreshToken()

	return nil
}

func (s *spotifyService) autoRefreshToken() {
	time.Sleep(time.Duration(s.token.ExpiresIn-60) * time.Second)
	s.getAuthToken()
}

func (s *spotifyService) parseAlbumIDFromLink(link string) (string, error) {
	u, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	pathParts := strings.SplitAfter(u.Path, "/")

	if len(pathParts) < 3 {
		return "", errors.New("invalid URL")
	}

	return pathParts[2], nil
}

func (s *spotifyService) getAlbumByID(id string) (*entity.SpotifyAlbum, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.spotify.com/v1/albums/%s", id), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+s.token.AccessToken)
	resp, err := s.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed Spotify API request: %s", string(body))
	}

	var album entity.SpotifyAlbum
	if err := json.NewDecoder(resp.Body).Decode(&album); err != nil {
		return nil, err
	}

	return &album, nil
}
