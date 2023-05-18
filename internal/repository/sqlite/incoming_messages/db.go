package incoming_messages

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

func (r *Repo) Create(dto CreateDTO) (*entity.IncomingMessage, error) {
	q := `INSERT INTO incoming_messages (text, sender_tag) 
		VALUES ($1, $2) 
		RETURNING id, create_at, text, sender_tag;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.Text,
		dto.SenderTag,
	).Scan(
		&result.ID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.IncomingMessage, error) {
	q := `DELETE FROM incoming_messages 
		WHERE id = $1 
		RETURNING id, create_at, text, sender_tag;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Get(dto GetDTO) (*entity.IncomingMessage, error) {
	q := `SELECT id, create_at, text, sender_tag 
		FROM incoming_messages 
		WHERE id = $1 LIMIT 1;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.CreateAt,
		&result.Text,
		&result.SenderTag,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll() ([]*entity.IncomingMessage, error) {
	q := `SELECT id, create_at, text, sender_tag 
		FROM incoming_messages 
		ORDER BY id;`

	rows, err := r.client.Query(q)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.IncomingMessage
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
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
