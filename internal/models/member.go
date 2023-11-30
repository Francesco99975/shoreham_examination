package models

type Member struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateMember(memeber *Member) error {
	statement := `INSERT INTO users(email, password) VALUES($1, $2);`

	_, err := db.Exec(statement, memeber.Email, memeber.Password)

	return err
}

func GetMemeber(email string) (Member, error) {
	var memeber Member

	statement := `SELECT * FROM members WHERE email=$1;`

	rows, err := db.Query(statement, email)

	if err != nil {
		return Member{}, err
	}

	for rows.Next() {
		err = rows.Scan(&memeber.ID, &memeber.Email, &memeber.Password)
		if err != nil {
			return Member{}, err
		}

	}

	return memeber, nil
}
