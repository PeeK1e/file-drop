package db

func RunTransaction(q []byte, args ...any) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(string(q), args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
