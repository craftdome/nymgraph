package outcoming_messages

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

func (r *Repo) Create(dto CreateDTO) (*entity.OutcomingMessage, error) {
	q := `INSERT INTO outcoming_messages (contact_id, text) 
		VALUES ($1, $2) 
		RETURNING id, contact_id, create_at, text;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ContactID,
		dto.Text,
	).Scan(
		&result.ID,
		&result.ContactID,
		&result.CreateAt,
		&result.Text,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.OutcomingMessage, error) {
	q := `DELETE FROM outcoming_messages 
		WHERE id = $1 
		RETURNING id, contact_id, create_at, text;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.ContactID,
		&result.CreateAt,
		&result.Text,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Get(dto GetDTO) (*entity.OutcomingMessage, error) {
	q := `SELECT id, contact_id, create_at, text 
		FROM outcoming_messages 
		WHERE id = $1 LIMIT 1;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.ContactID,
		&result.CreateAt,
		&result.Text,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll() ([]*entity.OutcomingMessage, error) {
	q := `SELECT id, contact_id, create_at, text 
		FROM outcoming_messages 
		ORDER BY id;`

	rows, err := r.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.OutcomingMessage
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.ContactID,
			&result.CreateAt,
			&result.Text,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "scan")
		}

		results = append(results, result.ToDomain())
	}

	return results, nil
}
