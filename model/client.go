package model

import (
	"time"

	"github.com/starudream/go-lib/sqlite/v2"

	"github.com/starudream/secret-tunnel/util"
)

type Client struct {
	Id           uint      `json:"id" gorm:"primaryKey,autoIncrement"`
	Name         string    `json:"name"`
	Key          string    `json:"key" gorm:"uniqueIndex"`
	Ver          string    `json:"ver"`
	Active       bool      `json:"active" gorm:"default:true"`
	Online       bool      `json:"online"`
	Addr         string    `json:"addr"`
	GO           string    `json:"go"`
	OS           string    `json:"os"`
	ARCH         string    `json:"arch"`
	Hostname     string    `json:"hostname"`
	LastOnlineAt time.Time `json:"last_online_at"`

	CreateAt time.Time `json:"create_at" gorm:"autoCreateTime"`
	UpdateAt time.Time `json:"update_at" gorm:"autoUpdateTime"`
}

func CreateClient(client *Client) (*Client, error) {
	client.Key = util.UUIDShort()
	client.Active = true
	return client, sqlite.DB().Select("name", "key").Create(client).Error
}

func DeleteClient(id uint) error {
	client := &Client{Id: id}
	return sqlite.DB().Delete(client).Error
}

func UpdateClient(client *Client) (*Client, error) {
	return client, sqlite.DB().Select("name").Updates(client).Error
}

func UpdateClientActive(id uint, active bool) error {
	return sqlite.DB().Model(&Client{}).Where("id=?", id).Update("active", active).Error
}

func UpdateClientOnline(client *Client) error {
	client.Online = true
	return sqlite.DB().Select("ver", "online", "addr", "go", "os", "arch", "hostname", "last_online_at").Updates(client).Error
}

func UpdateClientOffline(id uint) error {
	return sqlite.DB().Model(&Client{}).Where("id=?", id).Update("online", false).Error
}

func UpdateAllClientOffline() error {
	return sqlite.DB().Model(&Client{}).UpdateColumn("online", false).Error
}

func GetClientById(id uint) (*Client, error) {
	v := &Client{}
	return v, sqlite.DB().First(v, "id=?", id).Error
}

func GetClientByKey(key string) (*Client, error) {
	v := &Client{}
	return v, sqlite.DB().First(v, "key=?", key).Error
}

func ListClient() (clients []*Client, err error) {
	err = sqlite.DB().Order("id").Find(&clients).Error
	return truncClients(clients), err
}

func truncClients(clients []*Client) []*Client {
	for i := range clients {
		clients[i].LastOnlineAt = clients[i].LastOnlineAt.Local().Truncate(time.Second)
		clients[i].CreateAt = clients[i].CreateAt.Local().Truncate(time.Second)
		clients[i].UpdateAt = clients[i].UpdateAt.Local().Truncate(time.Second)
	}
	return clients
}
