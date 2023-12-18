package pseudonyms

import (
	"github.com/craftdome/nymgraph/internal/entity"
	"github.com/craftdome/nymgraph/pkg/client"
	"github.com/pkg/errors"
)

type Repo struct {
	client client.Client
}

func NewRepo(client client.Client) *Repo {
	return &Repo{client: client}
}

func (r *Repo) Create(dto CreateDTO) (*entity.Pseudonym, error) {
	q := `INSERT INTO pseudonyms (name, server) 
		VALUES ($1, $2) 
		RETURNING id, name, server;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.Name,
		dto.Server,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Server,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Update(dto UpdateDTO) (*entity.Pseudonym, error) {
	q := `UPDATE pseudonyms 
		SET name = $1, server = $2 
		WHERE id = $3 
		RETURNING id, name, server;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.Name,
		dto.Server,
		dto.ID,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Server,
	)
	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.Pseudonym, error) {
	q := `DELETE FROM pseudonyms 
		WHERE id = $1 
		RETURNING id, name, server;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Server,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Get(dto GetDTO) (*entity.Pseudonym, error) {
	q := `SELECT id, name, server 
		FROM pseudonyms 
		WHERE id = $1 LIMIT 1;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Server,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll() ([]*entity.Pseudonym, error) {
	q := `SELECT id, name, server 
		FROM pseudonyms 
		ORDER BY id;`

	rows, err := r.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.Pseudonym
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Server,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "scan")
		}

		results = append(results, result.ToDomain())
	}

	return results, nil
}
