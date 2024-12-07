package db

import "github.com/TOMMy-Net/air-sell/models"

func (s *Storage) GetPassportData(id int) ([]models.Passport, error) {
	var passports []models.Passport
	err := s.DB.Select(&passports, `SELECT * FROM passport_data WHERE user_id = $1`, id)
	return passports, err
}

func (s *Storage) AddPassport(p models.Passport) (int64, error) {
	res, err := s.DB.NamedExec(`INSERT INTO passport_data(name, surname, patronymic, passport_series_and_number, gender, validity_period, date_of_birth, passport_type, citizenship,user_id)
	VALUES(:name, :surname, :patronymic, :passport_series_and_number, :gender, :validity_period, :date_of_birth, :passport_type, :citizenship, :user_id)`, p)
	if err != nil {
		return 0, err
	}
	c, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}

	return c, nil
}

func (s *Storage) DeletePassport(id int) error  {
	_, err := s.DB.Exec(`DELETE FROM passport_data WHERE id = $1`, id)
	return err
}