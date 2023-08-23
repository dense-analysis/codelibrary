package database

import (
	"context"
	"strconv"

	"github.com/dense-analysis/codelibrary/internal/api/models"
	"github.com/google/uuid"
)

func (db *databaseAPIImpl) FindCodeSamples(
	ctx context.Context,
	search models.CodeSampleSearch,
) (models.CodeSamplePage, error) {
	var page models.CodeSamplePage

	filters := ` WHERE search_index @@ websearch_to_tsquery('english', $1)`
	params := make([]any, 1, 4)
	params[0] = search.Query

	if len(search.Languages) > 0 {
		filters += ` AND language_id = ANY($2)`
		params = append(params, search.Languages)
	}

	// Count the results first.
	countRow := db.pool.QueryRow(
		ctx,
		`SELECT COUNT(*) FROM codesample`+filters,
		params...,
	)
	err := countRow.Scan(&page.Count)

	if err != nil {
		return page, err
	}

	// If there are no results, don't run the second query to fetch a page.
	if page.Count == 0 {
		// Ensure we have an empty slice to avoid serialization issues.
		page.Results = []models.CodeSample{}
		return page, nil
	}

	offset := (search.Page - 1) * search.PageSize

	// Fetch a page of results.
	orderBy := ` ORDER BY (search_index @@ websearch_to_tsquery('english', $1))`
	pagination := ` LIMIT $` +
		strconv.Itoa(len(params)+1) +
		` OFFSET $` +
		strconv.Itoa(len(params)+2)
	params = append(params, search.PageSize)
	params = append(params, offset)
	pageRows, err := db.pool.Query(
		ctx,
		`
			SELECT
				codesample.id,
				submitted_by_id,
				username,
				language_id,
				language.name AS language_name,
				title,
				description,
				body,
				created,
				modified
			FROM codesample
			INNER JOIN "user"
			ON "user".id = codesample.submitted_by_id
			INNER JOIN language
			ON language.id = codesample.language_id
		`+filters+orderBy+pagination,
		params...,
	)

	if err != nil {
		return page, err
	}

	// Determine the capacity for the page.
	// Either it's the page size, or the remaining
	// items on the last page.
	pageCapacity := search.PageSize

	if page.Count < offset+search.PageSize {
		pageCapacity = page.Count - offset
	}

	page.Results = make([]models.CodeSample, 0, pageCapacity)

	for pageRows.Next() {
		sample := models.CodeSample{}
		err = pageRows.Scan(
			&sample.ID,
			&sample.SubmittedBy.ID,
			&sample.SubmittedBy.Username,
			&sample.Language.ID,
			&sample.Language.Name,
			&sample.Title,
			&sample.Description,
			&sample.Body,
			&sample.Created,
			&sample.Modified,
		)

		if err != nil {
			return page, err
		}

		page.Results = append(page.Results, sample)
	}

	return page, nil
}

func (db *databaseAPIImpl) GetCodeSample(ctx context.Context, id uuid.UUID) (models.CodeSample, error) {
	row := db.pool.QueryRow(
		ctx,
		`
			SELECT
				submitted_by_id,
				username,
				language_id,
				language.name AS language_name,
				title,
				description,
				body,
				created,
				modified
			FROM codesample
			INNER JOIN "user"
			ON "user".id = codesample.submitted_by_id
			INNER JOIN language
			ON language.id = codesample.language_id
			WHERE codesample.id = $1
		`,
		id,
	)

	sample := models.CodeSample{ID: id}
	err := row.Scan(
		&sample.SubmittedBy.ID,
		&sample.SubmittedBy.Username,
		&sample.Language.ID,
		&sample.Language.Name,
		&sample.Title,
		&sample.Description,
		&sample.Body,
		&sample.Created,
		&sample.Modified,
	)

	return sample, err
}

func (db *databaseAPIImpl) CreateCodeSample(ctx context.Context, sample models.CodeSample) error {
	_, err := db.pool.Exec(
		ctx,
		`
			INSERT INTO codesample (
				id, submitted_by_id, language_id,
				title, description, body,
				created, modified,
				search_index
			)
			VALUES (
				$1, $2, $3,
				$4, $5, $6,
				$7, $8,
				setweight(to_tsvector($4), 'A') ||
					setweight(to_tsvector($5), 'B') ||
					setweight(to_tsvector($6), 'C')
			)
		`,
		sample.ID, sample.SubmittedBy.ID, sample.Language.ID,
		sample.Title, sample.Description, sample.Body,
		sample.Created, sample.Modified,
	)

	return err
}

func (db *databaseAPIImpl) UpdateCodeSample(ctx context.Context, sample models.CodeSample) error {
	_, err := db.pool.Exec(
		ctx,
		`
			UPDATE codesample
			SET
				submitted_by_id = $2, language_id = $3,
				title = $4, description = $5, body = $6,
				created = $7, modified = $8,
				search_index = setweight(to_tsvector($4), 'A') ||
					setweight(to_tsvector($5), 'B') ||
					setweight(to_tsvector($6), 'C')
			WHERE codesample.id = $1
		`,
		sample.ID, sample.SubmittedBy.ID, sample.Language.ID,
		sample.Title, sample.Description, sample.Body,
		sample.Created, sample.Modified,
	)

	return err
}

func (db *databaseAPIImpl) DeleteCodeSample(ctx context.Context, id uuid.UUID) error {
	_, err := db.pool.Exec(
		ctx,
		`DELETE FROM codesample WHERE id = $1`,
		id,
	)

	return err
}
