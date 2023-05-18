package contacts

import (
	"github.com/Tyz3/nymgraph/internal/entity"
	"github.com/Tyz3/nymgraph/pkg/client"
	"github.com/pkg/errors"
)

type Repo struct {
	client client.Client
}

func NewRepo(client client.Client) *Repo {
	return &Repo{client: client}
}

func (r *Repo) Create(dto CreateDTO) (*entity.Contact, error) {
	q := `INSERT INTO contacts (pseudonym_id, address, alias) 
		VALUES ($1, $2, $3) 
		RETURNING id, pseudonym_id, address, alias;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.PseudonymID,
		dto.Address,
		dto.Alias,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.Address,
		&result.Alias,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Update(dto UpdateDTO) (*entity.Contact, error) {
	q := `UPDATE contacts 
		SET alias = $1 
		WHERE id = $2 
		RETURNING id, pseudonym_id, address, alias;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.Alias,
		dto.ID,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.Address,
		&result.Alias,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.Contact, error) {
	q := `DELETE FROM contacts 
		WHERE id = $1 
		RETURNING id, pseudonym_id, address, alias;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.Address,
		&result.Alias,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Get(dto GetDTO) (*entity.Contact, error) {
	q := `SELECT id, pseudonym_id, address, alias 
		FROM contacts 
		WHERE id = $1 LIMIT 1;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.Address,
		&result.Alias,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll() ([]*entity.Contact, error) {
	q := `SELECT id, pseudonym_id, address, alias 
		FROM contacts 
		ORDER BY id;`

	rows, err := r.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.Contact
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.PseudonymID,
			&result.Address,
			&result.Alias,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "scan")
		}

		results = append(results, result.ToDomain())
	}

	return results, nil
}
