package replies

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

func (r *Repo) Create(dto CreateDTO) (*entity.Reply, error) {
	q := `INSERT INTO replies (received_id, text) 
		VALUES ($1, $2) 
		RETURNING id, received_id, create_at, text;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ReceivedID,
		dto.Text,
	).Scan(
		&result.ID,
		&result.ReceivedID,
		&result.CreateAt,
		&result.Text,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.Reply, error) {
	q := `DELETE FROM replies 
		WHERE id = $1 
		RETURNING id, received_id, create_at, text;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.ReceivedID,
		&result.CreateAt,
		&result.Text,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll(dto GetAllDTO) ([]*entity.Reply, error) {
	q := `SELECT id, received_id, create_at, text 
		FROM replies 
		WHERE received_id = $1 
		ORDER BY id;`

	rows, err := r.client.Query(q, dto.ReceivedID)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.Reply
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.ReceivedID,
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

func (r *Repo) Truncate() error {
	q := `DELETE FROM replies;delete from sqlite_sequence where name='replies'`

	if _, err := r.client.Exec(q); err != nil {
		return errors.Wrapf(err, "Exec")
	}
	return nil
}
