package contacts

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
		SET address = $1, alias = $2 
		WHERE id = $3 
		RETURNING id, pseudonym_id, address, alias;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.Address,
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

func (r *Repo) GetAll(dto GetAllDTO) ([]*entity.Contact, error) {
	q := `SELECT id, pseudonym_id, address, alias 
		FROM contacts 
		WHERE pseudonym_id = $1 
		ORDER BY id;`

	rows, err := r.client.Query(q, dto.PseudonymID)
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
