package author

import (
	"fmt"
	"time"
)

type Author struct {
	ID        uint      `gorm:"primarykey"`
	CreatedAt time.Time `gorm:"<-:create"`
	UpdatedAt time.Time
	Name      string
	Surname   string
}

func NewAuthor(name string, surname string) *Author {
	author := &Author{
		Name:    name,
		Surname: surname,
	}
	return author
}

func (Author) TableName() string {
	return "author"
}

func (c *Author) ToString() string {
	return fmt.Sprintf("ID : %d, Name : %s, Surname : %s, CreatedAt : %s", c.ID, c.Name, c.Surname, c.CreatedAt.Format("2006-01-02 15:04:05"))
}

func (a *Author) GetFullName() string {
	return fmt.Sprintf("%s %s", a.Name, a.Surname)
}
