package models

type Member struct {
	ID       uint64 `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

func CreateMember(member *Member) error {
	statement := `INSERT INTO users(email, password) VALUES($1, $2);`

	_, err := db.Exec(statement, member.Email, member.Password)

	return err
}

func GetMember(email string) (Member, error) {
	var member Member

	statement := `SELECT * FROM members WHERE email=$1;`

	rows, err := db.Query(statement, email)

	if err != nil {
		return Member{}, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&member.ID, &member.Email, &member.Password)
		if err != nil {
			return Member{}, err
		}

	}

	return member, nil
}
