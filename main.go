package main

import (
	"fmt"
	"log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type CreditCard struct {
	gorm.Model `gorm:"tableName:creditcards1"`
	Number string `gorm:"size:255"`
	UserID uint `gorm:"primaryKey"`
}
type User struct {
	gorm.Model

	Name    string `gorm:"size:255"`
    Email string `gorm:"type:varchar(100); uniqueIndex"`
	Age int `gorm:"default:19"`
	Profile Profile //"Has One"
	Orders []Order  //"Has Many"
}
type Profile struct {
	gorm.Model 
	UserID uint  `gorm:"primaryKey"`
	Bio string 
}
type Order struct {
	gorm.Model
	UserID uint
	Item string 
	User User
}

func main() {
	
    dsn := "root:123456789@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("failed to connect database")
    }
    defer sqlDB.Close()
	//db.AutoMigrate(&User{})
	//db.AutoMigrate(&CreditCard{})
	db.AutoMigrate(&Profile{})
	db.AutoMigrate(&Order{})
	var orders = []Order{{Item: "Item 1"}, {Item: "Item 2"}}
	newUser,err:=NewUser("Keulen","aru@gmail.com", 20,"sdu",orders, db)
	if err!=nil {
		log.Fatal("failed to create:", err)
	} 
	order,err:=NewOrder("kitap1",*newUser,db)
	if err != nil {
		// Обработка ошибки
	}
	fmt.Print(order)
	fmt.Print("Created user with ID:", newUser.Name)
	
	allUsers :=GetAllUsers(db)
	fmt.Println("all user:", allUsers)

	getLast := GetLastUser(db)
	fmt.Println("last user witd ID", getLast.ID, getLast.Name)
	if err != nil {
        log.Fatal("failed to connect database")
    }
	getById, err := GetByID(db,3)
	if err!=nil {
      log.Fatal("failed to get user by ID", err)
	}
	fmt.Println("User by id", getById)
	updateUs,err := UpdateUser(db,8,"Haruka", "haru@gmail.com",300)
	if err!=nil {
		log.Fatal("error on update",err)
	}
	fmt.Print("updated users:",updateUs )
	deletedUser,err := DeleteUser(db, 24)
	if err!=nil {
		log.Fatal("error on delete", err)
	}
	fmt.Print("deleted user", deletedUser)
}
func NewUser(name, email string, age int,bio string,orders []Order, db*gorm.DB) (*User, error) {
	tx := db.Begin()
	defer rollbackTransaction(tx)
	if tx.Error!=nil {
		return nil, tx.Error
	}
	profile :=&Profile {
		
		Bio:bio,
	}
	if err:=tx.Create(profile).Error; err!=nil {
		return nil,err
	}
	
     NewUser := &User {
		Name: name,
		Email:email,
		Age:age,
		Profile:*profile,
		Orders:orders,
	 }
	 if NewUser.Profile.Bio == "panic" {
        panic("something went wrong")
    }
	 if err:=tx.Create(NewUser).Error; err!=nil {
		return nil, err
	 }
	 if err:=tx.Commit().Error; err!=nil {
		return nil, err
	 }
	 return NewUser, nil
}
func GetAllUsers(db *gorm.DB) []User{
	var allUsers []User
	result:=db.Find(&allUsers)
	if result.Error != nil {
		log.Fatal("failed to return all users:", result.Error)
	}
	return allUsers
}
func GetFirstUser (db *gorm.DB) User {
	var getFirst User
	result:=db.First(&getFirst)
	if result.Error!=nil {
		log.Fatal("failed to get first user:", result.Error)
	}
	return getFirst
}
func GetLastUser (db *gorm.DB) User {
	var getLast User
	result:=db.Order("id DESC").Take(&getLast)
	if result.Error!=nil {
		log.Fatal("failed to get last user", result.Error)
	}
	return getLast
}
func GetByID (db *gorm.DB, id uint) ([]User, error){
	var getByID []User
	result:=db.Find(&getByID, id)
	if result.Error!=nil {
		return nil, result.Error
	}
	return getByID, nil
}   
func UpdateUser (db *gorm.DB,id uint, name, email string, age int) ([]User, error) {
	var updateUser []User
	result:=db.Model(&updateUser).Where("id=?", id).Updates(User{
		Name:name,
		Email:email,
		Age:age,
	})
	if result.Error!=nil {
		log.Fatal("failed to update the user", result.Error)
	}
	fmt.Printf("success. %d записей\n", result.RowsAffected)
	return updateUser, nil
}
func DeleteUser(db *gorm.DB, id uint) ([]User,error) {
	var deletedUser []User
	result:=db.Delete(&deletedUser,id)
    if result.Error != nil {
		return nil, result.Error
	}
	return deletedUser, nil
}
func rollbackTransaction(tx *gorm.DB) {
	if r:=recover(); r!=nil {
		tx.Rollback()
	}
}
func NewOrder (item string, user User, db*gorm.DB)(*Order,error){
  tx:=db.Begin()
  defer rollbackTransaction(tx)
  if tx.Error!=nil {
     return nil, tx.Error
  }
  orders:=&Order{
	Item:item,
	User: user,
  }
  if orders.Item=="panic" {
	panic("something went wrong")
  }
  if err:=tx.Create(orders).Error; err!=nil {
	return nil, err
 }
 if err:=tx.Commit().Error; err!=nil {
	return nil, err
 }
 return orders, nil
}



