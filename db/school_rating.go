package db

import (
	"context"
	"log/slog"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type SchoolRatingService struct {
	db *DB
}

func NewSchoolRatingService(db *DB) *SchoolRatingService {
	return &SchoolRatingService{
		db: db,
	}
}

func (srs *SchoolRatingService) CreateSchoolRating(ctx context.Context, schoolRating *etp.SchoolRating) error {
	tx, err := srs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createSchoolRating(ctx, tx, schoolRating)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (srs *SchoolRatingService) GetSchoolRatings(ctx context.Context, filter *etp.SchoolRatingFilter) ([]*etp.SchoolRating, int, error) {
	tx, err := srs.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	ratings, n, err := getSchoolRatings(ctx, tx, filter)
	if err != nil {
		return nil, 0, err
	}

	return ratings, n, tx.Commit(ctx)
}

func (srs *SchoolRatingService) ApproveSchoolRating(ctx context.Context, id int) error {
	tx, err := srs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = approveSchoolRating(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (srs *SchoolRatingService) UpdateSchoolRating(ctx context.Context, id int, upd *etp.SchoolRatingUpdate) (*etp.SchoolRating, error) {
	tx, err := srs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	schoolRating, err := updateSchoolRating(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return schoolRating, tx.Commit(ctx)
}

func (srs *SchoolRatingService) DeleteSchoolRating(ctx context.Context, id int) error {
	tx, err := srs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = deleteSchoolRating(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteSchoolRating(ctx context.Context, tx *Tx, id int) error {
	query := `
		delete from school_rating
		where id = @id
	`
	args := pgx.NamedArgs{
		"id": id,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

// Updates a school rating by ID in the database.
//
// This function increments the approval count and updates the school rating's
// state to unapproved. If a school rating update is provided, it also updates
// the school rating's rating and comment accordingly.
func updateSchoolRating(ctx context.Context, tx *Tx, id int, upd *etp.SchoolRatingUpdate) (*etp.SchoolRating, error) {
	// Retrieve school ratings by ID from the database.
	schoolRatings, _, err := getSchoolRatings(ctx, tx, &etp.SchoolRatingFilter{RatingID: &id})
	if err != nil {
		return nil, err
	}

	// Check if a matching school rating was found.
	if len(schoolRatings) == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "school ratings not found"}
	}
	schoolRating := schoolRatings[0]

	// We have to reset the review state of approved, also we're going to check if the user has updated this review before
	// If the updated count is greater than 2, we won't allow the user to update the review
	if schoolRating.UpdatedCount == 2 {
		return nil, &etp.Error{Code: etp.ECONFLICT, Message: "school rating has been updated too many times"}
	}

	schoolRating.ApprovalCount = 0
	schoolRating.IsApproved = false

	// Update the school rating.
	if upd.Rating != nil {
		schoolRating.Rating = *upd.Rating
	}

	if upd.Comment != nil {
		schoolRating.Comment = *upd.Comment
	}

	query := `
		update 
			school_rating
		set 
			approval_count = @approvalCount,
			is_approved = @isApproved,
			rating = @rating,
			comment = @comment
		where 
			id = @id
	`

	args := pgx.NamedArgs{
		"id":            id,
		"approvalCount": schoolRating.ApprovalCount,
		"isApproved":    schoolRating.IsApproved,
		"rating":        schoolRating.Rating,
		"comment":       schoolRating.Comment,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return schoolRating, nil
}

func approveSchoolRating(ctx context.Context, tx *Tx, id int) error {
	// Retrieve school rating by ID from the database.
	schoolRating, _, err := getSchoolRatings(ctx, tx, &etp.SchoolRatingFilter{RatingID: &id})
	if err != nil {
		return err
	}

	// Check if a matching school rating was found.
	if len(schoolRating) == 0 {
		return &etp.Error{Code: etp.ENOTFOUND, Message: "school rating not found"}
	}

	// Verify that the school rating has not been already approved.
	if schoolRating[0].IsApproved || schoolRating[0].ApprovalCount >= 3 {
		return &etp.Error{Code: etp.ECONFLICT, Message: "school rating is already approved"}
	}

	// Increment the approval count and update the 'is_approved' flag accordingly.
	schoolRating[0].ApprovalCount++
	if schoolRating[0].ApprovalCount == 3 {
		schoolRating[0].IsApproved = true
	}

	// Construct a named query for updating the school rating in the database.
	query := `
		update
			school_rating
		set
			approval_count = @approvalCount,
			is_approved = @isApproved
		where
			id = @id
	`

	// Prepare the arguments for the named query.
	args := pgx.NamedArgs{
		"approvalCount": schoolRating[0].ApprovalCount,
		"isApproved":    schoolRating[0].IsApproved,
		"id":            id,
	}

	// Execute the update query against the database.
	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func getSchoolRatings(ctx context.Context, tx *Tx, filter *etp.SchoolRatingFilter) ([]*etp.SchoolRating, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if filter.RatingID != nil {
		where = append(where, "id = @id")
		args["id"] = *filter.RatingID
	}

	if filter.IsApproved != nil {
		where = append(where, "is_approved = true")
	}

	if filter.SchoolID != nil {
		where = append(where, "school_id = @schoolId")
		args["schoolId"] = *filter.SchoolID
	}

	query := `
		select 
			id,
			rating,
			comment,
			school_id,
			is_approved,
			approval_count,
			created_at,
			updated_at,
			count(*) over()
		from school_rating
	` + strings.Join(where, " and ") + `
	` + FormatLimitOffset(filter.Offset, filter.Limit) +
		`;`

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var n int

	schoolRatings := make([]*etp.SchoolRating, 0)
	for rows.Next() {
		var schoolRating *etp.SchoolRating

		if err := rows.Scan(
			&schoolRating.ID,
			&schoolRating.Rating,
			&schoolRating.Comment,
			&schoolRating.SchoolID,
			&schoolRating.IsApproved,
			&schoolRating.ApprovalCount,
			&schoolRating.CreatedAt,
			&schoolRating.UpdatedAt,
			&n,
		); err != nil {
			return nil, 0, err
		}

		schoolRatings = append(schoolRatings, schoolRating)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return schoolRatings, n, nil
}

func createSchoolRating(ctx context.Context, tx *Tx, schoolRating *etp.SchoolRating) error {
	query := `
		insert into school_rating
			(
				rating,
				comment,
				school_id,
				is_approved
				approval_count
			)
		values
			(
				@rating,
				@comment,
				@schoolId,
				@isApproved
				@approvalCount
			)
	`

	args := pgx.NamedArgs{
		"rating":        schoolRating.Rating,
		"comment":       schoolRating.Comment,
		"schoolId":      schoolRating.SchoolID,
		"approvalCount": 0,
		"isApproved":    false,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		slog.Error("error while creating school rating", schoolRating, err)
		return err
	}
	return nil
}
