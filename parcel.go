package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	res, err := s.db.Exec("insert into parcel (client,status, address, created_at) values (:client, :status, :address, :created_at)", sql.Named("client", p.Client),
		sql.Named("status", p.Status), sql.Named("address", p.Address), sql.Named("created_at", p.CreatedAt))
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	number, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return int(number), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	res := s.db.QueryRow("select * from parcel where number =:number", sql.Named("number", number))
	// здесь из таблицы должна вернуться только одна строка

	// заполните объект Parcel данными из таблицы
	p := Parcel{}
	err := res.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {

		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	date, err := s.db.Query("select * from parcel where client =:client", sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	// здесь из таблицы может вернуться несколько строк

	// заполните срез Parcel данными из таблицы
	var res []Parcel

	for date.Next() {
		p := Parcel{}
		err := date.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if err != nil {
			return nil, err
		}
		res = append(res, p)
	}
	err = date.Err()
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	_, err := s.db.Exec("update parcel set status =:status where number =:number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	res := s.db.QueryRow("select status from parcel where number =:number", sql.Named("number", number))

	p := Parcel{}
	err := res.Scan(&p.Status)
	if err != nil {
		return err
	}
	if p.Status == "registered" {
		_, err := s.db.Exec("update parcel set address =:address where number =:number", sql.Named("address", address), sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	res := s.db.QueryRow("select status from parcel where number =:number", sql.Named("number", number))

	p := Parcel{}
	err := res.Scan(&p.Status)
	if err != nil {
		return err
	}
	if p.Status == "registered" {
		_, err := s.db.Exec("Delete from parcel where number =:number", sql.Named("number", number))
		if err != nil {
			return err
		}
	}
	return nil
}
