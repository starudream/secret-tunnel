package model

type Task struct {
	Id       uint   `json:"id" gorm:"primaryKey,autoIncrement"`
	ClientId uint   `json:"client_id"`
	Name     string `json:"name"`
	Secret   string `json:"secret" gorm:"uniqueIndex"`
	Addr     string `json:"addr"`
	Active   bool   `json:"active" gorm:"default:true"`

	Meta
}

func CreateTask(task *Task) (*Task, error) {
	task.Active = true
	return task, _db.Select("client_id", "name", "secret", "addr").Create(task).Error
}

func DeleteTask(id uint) error {
	task := &Task{Id: id}
	return _db.Delete(task).Error
}

func UpdateTask(task *Task) (*Task, error) {
	return task, _db.Select("name", "addr").Updates(task).Error
}

func UpdateTaskActive(id uint, active bool) error {
	return _db.Model(&Task{}).Where("id = ?", id).Update("active", active).Error
}

func GetTaskById(id uint) (*Task, error) {
	v := &Task{}
	return v, _db.First(v, "id = ?", id).Error
}

func GetTaskBySecret(clientId uint, secret string) (*Task, error) {
	v := &Task{}
	if clientId == 0 {
		return v, _db.First(v, "secret = ?", secret).Error
	}
	return v, _db.First(v, "client_id = ? AND secret = ?", clientId, secret).Error
}

func ListTaskByClientId(clientId uint) (tasks []*Task, err error) {
	return tasks, _db.Order("id").Find(&tasks, "client_id = ?", clientId).Error
}
