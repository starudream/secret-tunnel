package model

import (
	"time"

	"github.com/starudream/go-lib/seq"
)

type Client struct {
	Id     uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	Name   string `json:"name"`
	Key    string `json:"key" gorm:"uniqueIndex"`
	Ver    string `json:"ver"`
	Active bool   `json:"active" gorm:"default:true"`
	Online bool   `json:"online"`

	Addr     string `json:"addr"`
	GO       string `json:"go"`
	OS       string `json:"os"`
	ARCH     string `json:"arch"`
	Hostname string `json:"hostname"`

	LastOnlineAt time.Time `json:"last_online_at"`

	CreateAt time.Time `json:"create_at" gorm:"autoCreateTime:milli"`
	UpdateAt time.Time `json:"update_at" gorm:"autoUpdateTime:milli"`
}

func CreateClient(client *Client) (*Client, error) {
	client.Key = seq.UUIDShort()
	client.Active = true
	return client, _db.Select("name", "key").Create(client).Error
}

func DeleteClient(id uint) error {
	client := &Client{Id: id}
	return _db.Delete(client).Error
}

func UpdateClient(client *Client) (*Client, error) {
	return client, _db.Select("name").Updates(client).Error
}

func UpdateClientActive(id uint, active bool) error {
	return _db.Model(&Client{}).Where("id = ?", id).Update("active", active).Error
}

func UpdateClientOnline(client *Client) error {
	client.Online = true
	client.LastOnlineAt = time.Now().Truncate(time.Millisecond)
	return _db.Select("ver", "online", "addr", "go", "os", "arch", "hostname", "last_online_at").Updates(client).Error
}

func UpdateClientOffline(id uint) error {
	return _db.Model(&Client{}).Where("id = ?", id).Update("online", false).Error
}

func GetClientById(id uint) (*Client, error) {
	v := &Client{}
	return v, _db.First(v, "id = ?", id).Error
}

func GetClientByKey(key string) (*Client, error) {
	v := &Client{}
	return v, _db.First(v, "key = ?", key).Error
}

func ListClient() (clients []*Client, err error) {
	return clients, _db.Order("id").Find(&clients).Error
}
