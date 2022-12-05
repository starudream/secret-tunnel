package model

type Client struct {
	Id     uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	Name   string `json:"name"`
	Key    string `json:"key" gorm:"uniqueIndex"`
	Active bool   `json:"active"`
	Online bool   `json:"online"`

	Addr     string `json:"addr"`
	GO       string `json:"go"`
	OS       string `json:"os"`
	ARCH     string `json:"arch"`
	Hostname string `json:"hostname"`

	Meta
}

func CreateClient(client *Client) (*Client, error) {
	return client, _db.Select("name", "key").Create(client).Error
}

func DeleteClient(id uint) error {
	return _db.Delete(&Client{}, "id = ?", id).Error
}

func UpdateClient(client *Client) (*Client, error) {
	return client, _db.Select("name").Where("id = ?", client.Id).Updates(client).Error
}

func UpdateClientActive(id uint, active bool) error {
	return _db.Model(&Client{}).Where("id = ?", id).Update("active", active).Error
}

func UpdateClientOnline(client *Client) error {
	client.Online = true
	return _db.Select("online", "addr", "go", "os", "arch", "hostname").Where("id = ?", client.Id).Updates(client).Error
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
	return clients, _db.Find(&clients).Error
}
