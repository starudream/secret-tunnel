package model

import (
	"time"

	"github.com/starudream/go-lib/sqlite/v2"

	"github.com/starudream/secret-tunnel/util"
)

type Task struct {
	Id         uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	ClientId   uint   `json:"client_id"`
	Name       string `json:"name"`
	Secret     string `json:"secret" gorm:"uniqueIndex"`
	Addr       string `json:"addr"`
	Active     bool   `json:"active" gorm:"default:true"`
	Compress   bool   `json:"compress"`
	TrafficIn  Size   `json:"traffic_in" gorm:"default:0"`
	TrafficOut Size   `json:"traffic_out" gorm:"default:0"`

	CreateAt time.Time `json:"create_at" gorm:"autoCreateTime:milli"`
	UpdateAt time.Time `json:"update_at" gorm:"autoUpdateTime:milli"`
}

func CreateTask(task *Task) (*Task, error) {
	task.Secret = util.UUIDShort()
	task.Active = true
	return task, sqlite.DB().Select("client_id", "name", "secret", "addr").Create(task).Error
}

func DeleteTask(id uint) error {
	task := &Task{Id: id}
	return sqlite.DB().Delete(task).Error
}

func DeleteTaskByClientId(clientId uint) error {
	return sqlite.DB().Delete(&Task{}, "client_id=?", clientId).Error
}

func UpdateTask(task *Task) (*Task, error) {
	return task, sqlite.DB().Select("name", "addr").Updates(task).Error
}

func UpdateTaskActive(id uint, active bool) error {
	return sqlite.DB().Model(&Task{}).Where("id=?", id).Update("active", active).Error
}

func UpdateTaskCompress(id uint, compress bool) error {
	return sqlite.DB().Model(&Task{}).Where("id=?", id).Update("compress", compress).Error
}

func UpdateTaskTraffic(id uint, in, out uint) error {
	if in == 0 && out == 0 {
		return nil
	}
	return sqlite.DB().Model(&Task{}).Where("id=?", id).Updates(map[string]any{"traffic_in": Expr("traffic_in+?", in), "traffic_out": Expr("traffic_out+?", out)}).Error
}

func GetTaskById(id uint) (*Task, error) {
	v := &Task{}
	return v, sqlite.DB().First(v, "id=?", id).Error
}

func GetTaskBySecret(clientId uint, secret string) (*Task, error) {
	v := &Task{}
	if clientId == 0 {
		return v, sqlite.DB().First(v, "secret=?", secret).Error
	}
	return v, sqlite.DB().First(v, "client_id=? AND secret=?", clientId, secret).Error
}

func ListTaskByClientId(clientId uint) (tasks []*Task, err error) {
	if clientId == 0 {
		return tasks, sqlite.DB().Order("client_id").Order("id").Find(&tasks).Error
	}
	return tasks, sqlite.DB().Order("id").Find(&tasks, "client_id=?", clientId).Error
}
