package internal
import (
 
)
func (r *User) GetByID(id domain.UserID) (*domain.User, error) {
	query := `SELECT * FROM "user" WHERE id=?;`
	uRow:=userRow{}
}