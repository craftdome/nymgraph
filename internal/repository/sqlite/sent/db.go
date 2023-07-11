package sent

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

func (r *Repo) Create(dto CreateDTO) (*entity.Sent, error) {
	q := `INSERT INTO sent (contact_id, text, reply_surbs) 
		VALUES ($1, $2, $3) 
		RETURNING id, contact_id, create_at, text, reply_surbs;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ContactID,
		dto.Text,
		dto.ReplySurbs,
	).Scan(
		&result.ID,
		&result.ContactID,
		&result.CreateAt,
		&result.Text,
		&result.ReplySurbs,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) Delete(dto DeleteDTO) (*entity.Sent, error) {
	q := `DELETE FROM sent 
		WHERE id = $1 
		RETURNING id, contact_id, create_at, text, reply_surbs;`

	result := new(Entity)
	err := r.client.QueryRow(q,
		dto.ID,
	).Scan(
		&result.ID,
		&result.ContactID,
		&result.CreateAt,
		&result.Text,
		&result.ReplySurbs,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "scan")
	}

	return result.ToDomain(), nil
}

func (r *Repo) GetAll(dto GetAllDTO) ([]*entity.Sent, error) {
	q := `SELECT id, contact_id, create_at, text, reply_surbs 
		FROM sent 
		WHERE contact_id = $1 
		ORDER BY id;`

	rows, err := r.client.Query(q, dto.ContactID)
	if err != nil {
		return nil, errors.Wrap(err, "query")
	}
	defer rows.Close()

	var results []*entity.Sent
	for rows.Next() {
		result := new(Entity)
		err := rows.Scan(
			&result.ID,
			&result.ContactID,
			&result.CreateAt,
			&result.Text,
			&result.ReplySurbs,
		)

		if err != nil {
			return nil, errors.Wrapf(err, "scan")
		}

		results = append(results, result.ToDomain())
	}

	return results, nil
}

func (r *Repo) Truncate() error {
	q := `DELETE FROM sent;delete from sqlite_sequence where name='sent'`

	if _, err := r.client.Exec(q); err != nil {
		return errors.Wrapf(err, "Exec")
	}
	return nil
}
