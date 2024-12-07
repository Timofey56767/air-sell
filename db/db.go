package db

import (
	"fmt"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/TOMMy-Net/air-sell/models"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattes/migrate/source/file"
)

type Storage struct {
	DB *sqlx.DB
}


type PostgresConnector struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSL      string
}

func ConnectPostgres(p PostgresConnector) (*Storage, error) {
	data, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		p.Host, p.Port, p.User, p.Password, p.Database, p.SSL))
	if err != nil {
		return &Storage{}, err
	}

	err = migratePostgres(data)
	if err != nil {
		return &Storage{}, err
	}
	return &Storage{data}, nil
}

func migratePostgres(db *sqlx.DB) error {
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/postgres",
		"postgres", driver)
	if err != nil {
		return err
	}
	m.Up()

	return nil
}

func ConnectSqlite(path string) (*Storage, error) {

	conn, err := sqlx.Open("sqlite3", path)
	if err != nil {
		return &Storage{}, err
	}

	if err := migrateSqlite(conn); err != nil {
		return &Storage{}, err
	}

	return &Storage{conn}, nil
}

func migrateSqlite(d *sqlx.DB) error {
	driver, err := sqlite3.WithInstance(d.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations/sqlite",
		"sqlite3", driver)
	if err != nil {
		return err
	}

	fmt.Println(m.Up())
	return nil

}

func (s *Storage) AllTickets() ([]models.Ticket, error) {
	m := []models.Ticket{}
	row, err := s.DB.Queryx(`SELECT t.id, t.airline, a1.city, a1.id, a1.iata_code, a1.country, a2.city, a2.id, a2.iata_code, a2.country, t.departure_time, t.arrival_time, t.quantity, t.luggage, t.hand_baggage, t.price
				FROM tickets as t 
				JOIN airports as a1 ON a1.id = t.departure_from
				JOIN airports as a2 ON a2.id = t.arrival_at
				WHERE t.quantity > 0`)

	if err != nil {
		return m, err
	}
	defer row.Close()

	for row.Next() {
		var ticket models.Ticket
		err = row.Scan(&ticket.ID, &ticket.Airline,
			&ticket.DepartureFrom.City, &ticket.DepartureFrom.ID, &ticket.DepartureFrom.Iata, &ticket.DepartureFrom.Country,
			&ticket.ArrivalAt.City, &ticket.ArrivalAt.ID, &ticket.ArrivalAt.Iata, &ticket.ArrivalAt.Country,
			&ticket.DepartureTime, &ticket.ArrivalTime, &ticket.Quantity, &ticket.Luggage, &ticket.HandBaggage, &ticket.Price)
		if err != nil {
			continue
		}
		m = append(m, ticket)
	}
	return m, nil
}

func (s *Storage) FindTickets(t *models.TicketsSearch) ([]models.Ticket, error) {
	tickets := []models.Ticket{}
	row, err := s.DB.Queryx(`SELECT t.id, t.airline, a1.city, a1.id, a1.iata_code, a1.country, a2.city, a2.id, a2.iata_code, a2.country, TO_CHAR(t.departure_time, 'YYYY-MM-DD HH24:MI:SS'), TO_CHAR(t.arrival_time, 'YYYY-MM-DD HH24:MI:SS'), t.quantity, t.luggage, t.hand_baggage, t.price
				FROM tickets as t 
				JOIN airports as a1 ON a1.id = t.departure_from
				JOIN airports as a2 ON a2.id = t.arrival_at
				WHERE t.quantity > 0 AND (a1.city = $1 AND a2.city = $2) AND (DATE(t.departure_time) = $3 )`, t.From, t.To, t.Date_from)

	if err != nil {
		return []models.Ticket{}, err
	}
	defer row.Close()
	for row.Next() {
		var ticket models.Ticket
		err = row.Scan(&ticket.ID, &ticket.Airline,
			&ticket.DepartureFrom.City, &ticket.DepartureFrom.ID, &ticket.DepartureFrom.Iata, &ticket.DepartureFrom.Country,
			&ticket.ArrivalAt.City, &ticket.ArrivalAt.ID, &ticket.ArrivalAt.Iata, &ticket.ArrivalAt.Country,
			&ticket.DepartureTime, &ticket.ArrivalTime, &ticket.Quantity, &ticket.Luggage, &ticket.HandBaggage, &ticket.Price)
		if err != nil {
			continue
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

func (s *Storage) AllAirports() ([]models.AirPorts, error) {
	airPorts := []models.AirPorts{}

	err := s.DB.Select(&airPorts, `SELECT * FROM airports`)
	return airPorts, err
}


func (s *Storage) MinusTicketCount(id string, count int) error {
	var ticket models.Ticket
	tx, err := s.DB.Beginx()
	if err != nil {
		return err
	}
	err = tx.Get(&ticket, `SELECT quantity FROM tickets WHERE id = $1`, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`UPDATE tickets SET quantity = $1 WHERE id = $2`, ticket.Quantity - count, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	
	return tx.Commit()
}