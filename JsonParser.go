package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
)

func OpenDB() (*sql.DB, error) {
	dsn := "admin@tcp(localhost:3306)/go_test?parseTime=true&charset=utf8mb4,utf8"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func ParseUserJson(user User) {

	file, err := os.ReadFile("./user_data.json")
	if err != nil {
		fmt.Println("err reading file: ", err)
	}

	var data User

	er := json.Unmarshal(file, &data)
	if er != nil {
		fmt.Println("err parsing json :", er)
	}

	// Open DB connection to check user existence
	db, err := OpenDB()
	if err != nil {
		fmt.Println("failed to connect to database: ", err)
		return
	}
	defer db.Close()

	// Check if the user exists by ID
	exists, err := GetUserById(data.Id, db)
	if err != nil {
		fmt.Println("Error checking user existence: ", err)
		return
	}

	fmt.Println("exisits here : ", exists)

	if exists {
		// If the user exists, update the user data
		err = updateUserData(data)
		if err != nil {
			fmt.Println("Error updating user data: ", err)
		} else {
			fmt.Println("User data updated successfully")
		}
	} else {
		// If the user doesn't exist, save new user data
		err = saveUserData(data)
		if err != nil {
			fmt.Println("Error saving user data: ", err)
		} else {
			fmt.Println("User data saved successfully")
		}
	}
}

func saveUserData(user User) error {
	// Open DB connection
	db, err := OpenDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Insert user data
	err = insertUser(user, db)
	if err != nil {
		return fmt.Errorf("failed to insert user: %v", err)
	}

	// Insert experiences
	for _, experience := range user.Experiences {
		err := insertExperience(user.Id, experience, db)
		if err != nil {
			return fmt.Errorf("failed to insert experience: %v", err)
		}
	}

	// Insert earnings and associated contributions
	for _, earning := range user.Earnings {
		err := insertEarning(user.Id, earning, db)
		if err != nil {
			return fmt.Errorf("failed to insert earning: %v", err)
		}
	}

	// Insert notes
	for _, note := range user.Notes {
		err := insertNote(user.Id, note, db)
		if err != nil {
			return fmt.Errorf("failed to insert note: %v", err)
		}
	}

	return nil
}

func insertUser(user User, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO users (id, name, email, contact, income, is_present, joined_at) VALUES (?, ?, ?, ?, ?, ?, ?)",
		user.Id, user.Name, user.Email, user.Contact, user.Income, user.IsPresent, user.JoinedAt,
	)
	return err
}

func insertExperience(userId string, experience Experience, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO experiences (user_id, title, duration) VALUES (?, ?, ?)",
		userId, experience.Title, experience.Duration,
	)
	return err
}

func insertEarning(userId string, earning Earning, db *sql.DB) error {
	// Insert the earning
	result, err := db.Exec(
		"INSERT INTO earnings (user_id, date, amount, deductible, bonus, skill_incentive, net_payable) VALUES (?, ?, ?, ?, ?, ?, ?)",
		userId, earning.Date, earning.Amount, earning.Deductible, earning.Bonus, earning.SkillIncentive, earning.NetPayable,
	)
	if err != nil {
		return err
	}

	earningID, _ := result.LastInsertId()

	// Insert contributions for this earning
	err = insertContribution(int(earningID), earning.Contribution, db)
	return err
}

func insertContribution(earningID int, contribution Contribution, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO contributions (earning_id, provident_fund, sst, rt, cit, attendance_deduction, welfare_fund) VALUES (?, ?, ?, ?, ?, ?, ?)",
		earningID, contribution.ProvidentFund, contribution.SST, contribution.RT, contribution.CIT, contribution.AttendanceDeduction, contribution.WelfareFund,
	)
	return err
}

func insertNote(userId string, note Note, db *sql.DB) error {
	_, err := db.Exec(
		"INSERT INTO notes (user_id, created_at, created_by, note, priority_level) VALUES (?, ?, ?, ?, ?)",
		userId, note.CreatedAt, note.CreatedBy, note.Note, note.Priority.Level,
	)
	return err
}

func updateUserData(user User) error {
	db, err := OpenDB()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %v", err)
	}
	defer db.Close()

	// Update user data
	err = updateUser(user, db)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	// Delete existing experiences, earnings, and notes (assuming user is updating them)
	// err = deleteExperiences(user.Id, db)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete experiences: %v", err)
	// }
	//
	// err = deleteEarnings(user.Id, db)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete earnings: %v", err)
	// }
	//
	// err = deleteNotes(user.Id, db)
	// if err != nil {
	// 	return fmt.Errorf("failed to delete notes: %v", err)
	// }
	//
	// // Insert new experiences, earnings, and notes
	// for _, experience := range user.Experiences {
	// 	err := insertExperience(user.Id, experience, db)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert experience: %v", err)
	// 	}
	// }
	//
	// for _, earning := range user.Earnings {
	// 	err := insertEarning(user.Id, earning, db)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert earning: %v", err)
	// 	}
	// }
	//
	// for _, note := range user.Notes {
	// 	err := insertNote(user.Id, note, db)
	// 	if err != nil {
	// 		return fmt.Errorf("failed to insert note: %v", err)
	// 	}
	// }
	//
	return nil
}

func updateUser(user User, db *sql.DB) error {
	_, err := db.Exec(
		"UPDATE users SET name = ?, email = ?, contact = ?, income = ?, is_present = ?, joined_at = ? WHERE id = ?",
		user.Name, user.Email, user.Contact, user.Income, user.IsPresent, user.JoinedAt, user.Id,
	)
	return err
}

func deleteExperiences(userId string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM experiences WHERE user_id = ?", userId)
	return err
}

func deleteEarnings(userId string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM earnings WHERE user_id = ?", userId)
	return err
}

func deleteNotes(userId string, db *sql.DB) error {
	_, err := db.Exec("DELETE FROM notes WHERE user_id = ?", userId)
	return err
}

func GetUserById(userId string, db *sql.DB) (bool, error) {
	var user User

	// Query to check if user exists by VARCHAR user ID
	row := db.QueryRow("SELECT id FROM users WHERE id = ?", userId)
	if err := row.Scan(&user.Id); err != nil {
		if err == sql.ErrNoRows {
			// No user found
			return false, nil
		}
		// Return any other errors encountered
		return false, fmt.Errorf("error querying user by id: %v", err)
	}

	// User found
	return true, nil
}
