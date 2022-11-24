// Filename:  internal/data/permissions.go

package data

import (
	"context"
	"database/sql"
	"time"
)

//Define the slice to hold the permission codes

type Permissions []string

func (p Permissions) Includes(code string) bool {
	for i := range p {
		if code == p[i] {
			return true
		}
	}

	return false
}

type PermissionsModel struct {
	DB *sql.DB
}

func (m PermissionsModel) GetAllForUser(userID int64) (Permissions, error) {
	query := `
		SELECT p.code FROM 
		permissions p 
		INNER JOIN
		users_permissions up
		ON up.id = p.id
		INNER JOIN users u
		ON up.user_id = u.id
		WHERE u.id = $1
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var permissions Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)

		if err != nil {
			return nil, err
		}

		permissions = append(permissions, permission)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}
