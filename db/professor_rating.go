package db

import (
	"context"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type ProfessorRatingService struct {
	db *DB
}

func NewProfessorRatingService(db *DB) *ProfessorRatingService {
	return &ProfessorRatingService{
		db: db,
	}
}

func (prs *ProfessorRatingService) CreateProfessorRating(ctx context.Context, professorRating *etp.ProfessorRating) error {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = createProfessorRating(ctx, tx, professorRating)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (prs *ProfessorRatingService) GetProfessorRatings(ctx context.Context, filter etp.ProfessorRatingFilter) ([]*etp.ProfessorRating, int, error) {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	ratings, n, err := getProfessorRatings(ctx, tx, &filter)
	if err != nil {
		return nil, 0, err
	}

	return ratings, n, tx.Commit(ctx)
}

func (prs *ProfessorRatingService) ApproveProfessorRating(ctx context.Context, id int) error {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = approveProfessorRating(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (prs *ProfessorRatingService) UpdateProfessorRating(ctx context.Context, id int, upd *etp.ProfessorRatingUpdate) (*etp.ProfessorRating, error) {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	professorRating, err := updateProfessorRating(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return professorRating, tx.Commit(ctx)
}

func (prs *ProfessorRatingService) DeleteProfessorRating(ctx context.Context, id int) error {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = deleteProfessorRating(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

// Functions for handling Professor Ratings
func deleteProfessorRating(ctx context.Context, tx *Tx, id int) error {
	query := `
		DELETE FROM professor_rating
		WHERE id = @id
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

func updateProfessorRating(ctx context.Context, tx *Tx, id int, upd *etp.ProfessorRatingUpdate) (*etp.ProfessorRating, error) {
	professorRating, err := getProfessorRatingById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if professorRating.UpdatedCount == 2 {
		return nil, &etp.Error{Code: etp.ECONFLICT, Message: "professor rating has been updated too many times"}
	}

	if v := upd.WouldTakeAgain; v != nil {
		professorRating.WouldTakeAgain = *v
	}

	if v := upd.MandatoryAttendance; v != nil {
		professorRating.MandatoryAttendance = *v
	}

	if v := upd.Grade; v != nil {
		professorRating.Grade = *v
	}

	if v := upd.TextbookRequired; v != nil {
		professorRating.TextbookRequired = *v
	}

	if v := upd.Rating; v != nil {
		professorRating.Rating = *v
	}

	if v := upd.Comment; v != nil {
		professorRating.Comment = *v
	}

	query := `
		UPDATE professor_rating
		SET
			rating = @rating,
			comment = @comment,
			would_take_again = @wouldTakeAgain,
			mandatory_attendance = @mandatoryAttendance,
			grade = @grade,
			textbook_required = @textbookRequired,
			approvals_count = @approvalsCount,
			updated_at = NOW()
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"id":                  id,
		"rating":              professorRating.Rating,
		"comment":             professorRating.Comment,
		"wouldTakeAgain":      professorRating.WouldTakeAgain,
		"mandatoryAttendance": professorRating.MandatoryAttendance,
		"grade":               professorRating.Grade,
		"textbookRequired":    professorRating.TextbookRequired,
		"approvalsCount":      professorRating.ApprovalsCount,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return professorRating, nil
}

func approveProfessorRating(ctx context.Context, tx *Tx, id int) error {
	professorRating, err := getProfessorRatingById(ctx, tx, id)
	if err != nil {
		return err
	}

	if professorRating.IsApproved || professorRating.ApprovalsCount >= 3 {
		return &etp.Error{Code: etp.ECONFLICT, Message: "professor rating is already approved"}
	}

	professorRating.ApprovalsCount++
	if professorRating.ApprovalsCount == 3 {
		professorRating.IsApproved = true
	}

	query := `
		UPDATE professor_rating
		SET
			approvals_count = @approvalsCount,
			is_approved = @isApproved,
			updated_at = NOW()
		WHERE id = @id
	`

	args := pgx.NamedArgs{
		"approvalsCount": professorRating.ApprovalsCount,
		"isApproved":     professorRating.IsApproved,
		"id":             id,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func getProfessorRatingById(ctx context.Context, tx *Tx, id int) (*etp.ProfessorRating, error) {
	professorRatings, n, err := getProfessorRatings(ctx, tx, &etp.ProfessorRatingFilter{ProfessorId: &id})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "professor rating not found"}
	}

	return professorRatings[0], nil
}

func getProfessorRatings(ctx context.Context, tx *Tx, filter *etp.ProfessorRatingFilter) ([]*etp.ProfessorRating, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if v := filter.ProfessorId; v != nil {
		where = append(where, "professor_id = @professorId")
		args["professorId"] = *v
	}

	if v := filter.CourseId; v != nil {
		where = append(where, "course_id = @courseId")
		args["courseId"] = *v
	}

	query := `
		select count(*) 
		from professor_rating
		WHERE ` + strings.Join(where, " AND ")

	var n int
	err := tx.QueryRow(ctx, query, args).Scan(&n)
	if err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return nil, 0, nil
	}

	query = `
		SELECT * FROM professor_rating
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY created_at DESC ` + ` 
		` + FormatLimitOffset(filter.Offset, filter.Limit)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	professorRatings, err := pgx.CollectRows(rows, pgx.RowToStructByName[*etp.ProfessorRating])
	if err != nil {
		return nil, 0, err
	}

	return professorRatings, n, nil
}

func createProfessorRating(ctx context.Context, tx *Tx, professorRating *etp.ProfessorRating) error {
	query := `
		INSERT INTO professor_rating (rating, comment, would_take_again, mandatory_attendance, grade, textbook_required, 
			is_approved, approvals_count, created_at, updated_at, professor_id, course_id)
		VALUES (@rating, @comment, @wouldTakeAgain, @mandatoryAttendance, @grade, @textbookRequired, @isApproved, @approvalsCount, NOW(), NOW(), @professorId, @courseId)
	`

	args := pgx.NamedArgs{
		"rating":              professorRating.Rating,
		"comment":             professorRating.Comment,
		"wouldTakeAgain":      professorRating.WouldTakeAgain,
		"mandatoryAttendance": professorRating.MandatoryAttendance,
		"grade":               professorRating.Grade,
		"textbookRequired":    professorRating.TextbookRequired,
		"isApproved":          professorRating.IsApproved,
		"approvalsCount":      professorRating.ApprovalsCount,
		"professorId":         professorRating.ProfessorId,
		"courseId":            professorRating.CourseId,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}
