package received

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

func (r *Repo) Create(dto CreateDTO) (*entity.Received, error) {
	q := `INSERT INTO received (pseudonym_id, text, sender_tag) 
		VALUES ($1, $2, $3) 
		RETURNING id, pseudonym_id, create_at, text, sender_tag;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.PseudonymID,
		dto.Text,
		dto.SenderTag,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.Received, error) {
	q := `DELETE FROM received 
		WHERE id = $1 
		RETURNING id, pseudonym_id, create_at, text, sender_tag;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Get(dto GetDTO) (*entity.Received, error) {
	q := `SELECT id, pseudonym_id, create_at, text, sender_tag 
		FROM received 
		WHERE id = $1;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.PseudonymID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll(dto GetAllDTO) ([]*entity.Received, error) {
	q := `SELECT id, pseudonym_id, create_at, text, sender_tag 
		FROM received
		WHERE pseudonym_id = $1 
		ORDER BY id;`

	rows, err := r.client.Query(q, dto.PseudonymID)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.Received
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.PseudonymID,
			&result.CreateAt,
			&result.Text,
			&result.SenderTag,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "scan")
		}

		results = append(results, result.ToDomain())
	}

	return results, nil
}

func (r *Repo) Truncate() error {
	q := `DELETE FROM received;delete from sqlite_sequence where name='received'`

	if _, err := r.client.Exec(q); err != nil {
		return errors.Wrapf(err, "Exec")
	}
	return nil
}
