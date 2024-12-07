package db

import (
	"github.com/TOMMy-Net/air-sell/models"
	"github.com/TOMMy-Net/air-sell/tools"
)

func (s *Storage) CreateUser(u models.User) (int64, error) {
	var userID int
	u.Password = tools.Sum256([]byte(u.Password))
	tx, err := s.DB.Beginx()
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = tx.NamedExec(`INSERT INTO users(email, password) VALUES (:email, :password)`, u)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Get(&userID, `SELECT id FROM users WHERE email = $1 AND password = $2`, u.Email, u.Password)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	tx.Commit()
	return int64(userID), nil
}

func (s *Storage) GetUser(u models.User) (*models.User, error) {
	user := models.User{}
	u.Password = tools.Sum256([]byte(u.Password))
	err := s.DB.Get(&user, `SELECT * FROM users WHERE email = $1 AND password = $2`, u.Email, u.Password)
	return &user, err
}

func (s *Storage) SetBuyHistory(h models.BuyHistory) error {
	_, err := s.DB.Exec(`INSERT INTO buy_history(ticket_id, user_id, buy_time, count) VALUES ($1, $2, $3, $4)`, h.Ticket.ID, h.UserId, h.BuyTime, h.Count)
	return err
}

func (s *Storage) GetBuyHistory(userID int) ([]models.BuyHistory, error) {
	var history []models.BuyHistory
	rows, err := s.DB.Queryx("SELECT h.id, t.id, t.airline, t.departure_time, t.arrival_time, t.luggage, t.hand_baggage, t.price, h.buy_time, h.count FROM buy_history as h JOIN tickets as t ON t.id = h.ticket_id WHERE h.user_id = $1", userID)
	if err != nil {
		return []models.BuyHistory{}, err
	}
	for rows.Next() {
		var m models.BuyHistory
		err := rows.Scan(&m.ID, &m.Ticket.ID, &m.Ticket.Airline, &m.Ticket.DepartureTime, &m.Ticket.ArrivalTime, &m.Ticket.Luggage, &m.Ticket.HandBaggage, &m.Ticket.Price, &m.BuyTime, &m.Count)
		if err == nil {
			history = append(history, m)
		}
	}
	return history, nil
	
}
