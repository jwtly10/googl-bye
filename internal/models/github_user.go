package models

type GithubUser struct {
	Id        int64  `json:"id"`
	AvatarUrl string `json:"avatar_url"`
	Login     string `json:"login"`
	Url       string `json:"url"`
	Name      string `json:"name"`
}
