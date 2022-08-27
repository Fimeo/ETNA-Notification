package database

import "log"

func (s *Service) CreateNotification(notification Notification) error {
	_, err := s.DB.Exec(
		"INSERT INTO notification (time, external_id, \"user\") VALUES(NOW(),$1, $2);",
		notification.ExternalID, notification.User)
	if err != nil {
		log.Print("[ERROR] Insert into database failed: ", err)
		return err
	}

	return nil
}

func (s *Service) IsAlreadyNotified(notification Notification) (bool, error) {
	row, err := s.DB.Query(
		"SELECT * FROM notification WHERE external_id=$1 and user=$2",
		notification.ExternalID, notification.User)
	if err != nil {
		log.Print("[ERROR] Cannot read database: ", err)
		return false, nil
	}
	count := 0
	for row.Next() {
		count++
	}
	return count != 0, nil
}

func (s *Service) GetEtnaUsers() ([]EtnaUser, error) {
	rows, err := s.DB.Query("SELECT * FROM users;")
	if err != nil {
		log.Print("[ERROR] Retrieve users failed in database: ", err)
		return nil, err
	}

	var got []EtnaUser
	for rows.Next() {
		var r EtnaUser
		err = rows.Scan(&r.ID, &r.Time, &r.UserID, &r.ChannelID, &r.Login, &r.Password)
		if err != nil {
			return nil, err
		}
		got = append(got, r)
	}

	return got, nil
}
