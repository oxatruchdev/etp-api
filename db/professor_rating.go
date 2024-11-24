package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
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

func (prs *ProfessorRatingService) CreateProfessorRating(ctx context.Context, professorRating *etp.ProfessorRating, tagIds []int) error {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err, id := createProfessorRating(ctx, tx, professorRating)
	if err != nil {
		return err
	}

	if len(tagIds) > 0 {
		err = associateProfessorRatingTags(ctx, tx, id, tagIds)
		if err != nil {
			return err
		}
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

func (prs *ProfessorRatingService) GetProfessorRatingsWithStats(ctx context.Context, filter etp.ProfessorRatingFilter) (*etp.ProfessorRatingsStats, error) {
	tx, err := prs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	stats, err := getProfessorRatingStats(ctx, tx, filter)
	if err != nil {
		return nil, err
	}

	return &stats, tx.Commit(ctx)
}

func getProfessorRatingStats(ctx context.Context, tx *Tx, filter etp.ProfessorRatingFilter) (etp.ProfessorRatingsStats, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if v := filter.ProfessorId; v != nil {
		where = append(where, "p_r.professor_id = @professorId")
		args["professorId"] = *v
	}
	if v := filter.CourseId; v != nil {
		where = append(where, "p_r.course_id = @courseId")
		args["courseId"] = *v
	}

	if v := filter.IsApproved; v {
		where = append(where, "p_r.is_approved = @isApproved")
		args["isApproved"] = v
	}

	// First, get the professor ratings with window functions for avg and total count
	query := `
		WITH 
		    rating_tags AS (
			SELECT 
			    prt.professor_rating_id,
			    jsonb_agg(
				jsonb_build_object(
				    'id', t.id,
				    'name', t.name
				)
			    ) as tags
			FROM professor_rating_tag prt
			JOIN tag t ON t.id = prt.tag_id
			GROUP BY prt.professor_rating_id
		    ),
		    rating_distribution AS (
			SELECT 
			    professor_id,
			    jsonb_agg(
				jsonb_build_object(
				    'rating', rating,
				    'count', count
				) ORDER BY rating
			    ) as distribution
			FROM (
			    SELECT 
				professor_id,
				rating, 
				COUNT(*) as count
				FROM professor_rating p_r
				WHERE ` + strings.Join(where, " AND ") + ` 
				GROUP BY professor_id, rating
			) rd
			GROUP BY professor_id
		    )
		SELECT 
		    p_r.id,
		    rating,
		    comment,
		    would_take_again,
		    mandatory_attendance,
		    grade,
		    textbook_required,
		    approvals_count,
		    is_approved,
		    p_r.created_at,
		    course_id,
		    p_r.professor_id,
		    p_r.updated_at,
		    updated_count,
		    COALESCE(difficulty, 0),
		    AVG(rating) OVER () AS avg_rating,
		    COUNT(p_r.id) OVER () AS total_ratings,
		    AVG(CAST(would_take_again = true AS int)) OVER () as would_take_again_rating,
		    COALESCE(AVG(difficulty) OVER (), 0) as avg_difficulty,
		    course.id,
		    course.name,
		    course.code,
		    COALESCE(rt.tags, '[]'::jsonb) as tags,
		    COALESCE(rd.distribution, '[]'::jsonb) as rating_distribution
		FROM professor_rating p_r
		LEFT JOIN course on course.id = p_r.course_id
		LEFT JOIN rating_tags rt ON rt.professor_rating_id = p_r.id
		LEFT JOIN rating_distribution rd ON rd.professor_id = p_r.professor_id
		WHERE ` + strings.Join(where, " AND ") +
		` ORDER BY created_at DESC ` + `
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return etp.ProfessorRatingsStats{}, err
	}
	defer rows.Close()

	var stats etp.ProfessorRatingsStats
	var avgRating float64
	var totalRatings int
	var avgWouldTakeAgainRating float64
	var avgDifficulty float64

	for rows.Next() {
		var professorRating etp.ProfessorRating
		professorRating.Course = &etp.Course{}
		var tagsJSON []byte
		var distributionJSON []byte
		err = rows.Scan(
			&professorRating.ID,
			&professorRating.Rating,
			&professorRating.Comment,
			&professorRating.WouldTakeAgain,
			&professorRating.MandatoryAttendance,
			&professorRating.Grade,
			&professorRating.TextbookRequired,
			&professorRating.ApprovalsCount,
			&professorRating.IsApproved,
			&professorRating.CreatedAt,
			&professorRating.CourseId,
			&professorRating.ProfessorId,
			&professorRating.UpdatedAt,
			&professorRating.UpdatedCount,
			&professorRating.Difficulty,
			&avgRating,
			&totalRatings,
			&avgWouldTakeAgainRating,
			&avgDifficulty,
			&professorRating.Course.ID,
			&professorRating.Course.Name,
			&professorRating.Course.Code,
			&tagsJSON,
			&distributionJSON,
		)
		if err != nil {
			return etp.ProfessorRatingsStats{}, err
		}

		// Unmarshal tags
		var tags []*etp.Tag
		if err := json.Unmarshal(tagsJSON, &tags); err != nil {
			return etp.ProfessorRatingsStats{}, fmt.Errorf("failed to unmarshal tags: %w", err)
		}
		professorRating.Tags = tags

		// Parse distribution
		var distributionItems []struct {
			Rating int
			Count  int
		}
		if err := json.Unmarshal(distributionJSON, &distributionItems); err != nil {
			return etp.ProfessorRatingsStats{}, fmt.Errorf("failed to unmarshal distribution: %w", err)
		}

		var ratingDistribution []*etp.RatingDistribution
		// Convert to your RatingDistribution format
		for _, item := range distributionItems {
			ratingDistribution = append(ratingDistribution, &etp.RatingDistribution{
				Rating: item.Rating,
				Count:  item.Count,
			})
		}
		stats.RatingsDistribution = ratingDistribution
		stats.Ratings = append(stats.Ratings, &professorRating)
		slog.Info("rating", "profrating", professorRating)
	}

	if rows.Err() != nil {
		return etp.ProfessorRatingsStats{}, err
	}

	stats.TotalCount = totalRatings
	stats.RatingAvg = avgRating
	stats.WouldTakeAgainAvg = avgWouldTakeAgainRating
	stats.DifficultyAvg = avgDifficulty

	stats.EnsureFullDistribution()

	return stats, nil
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

	if professorRating.UpdatedCount == 3 {
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

	professorRating.UpdatedCount++

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
			is_approved = @isApproved,
			updated_count = @updatedCount,
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
		"approvalsCount":      0,
		"isApproved":          false,
		"updatedCount":        professorRating.UpdatedCount,
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
	professorRatings, n, err := getProfessorRatings(ctx, tx, &etp.ProfessorRatingFilter{ProfessorRatingId: &id})
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

	if v := filter.ProfessorRatingId; v != nil {
		where = append(where, "id = @id")
		args["id"] = *v
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

	slog.Info("total", "n", n)
	if n == 0 {
		return nil, 0, nil
	}

	query = `
		SELECT 
			id,
			rating,
			comment,
			would_take_again,
			mandatory_attendance,
			grade,
			textbook_required,
			approvals_count,
			is_approved,
			created_at,
			course_id,
			professor_id,
			updated_at,
			updated_count
		FROM professor_rating
		WHERE ` + strings.Join(where, " AND ") + `
		ORDER BY updated_at DESC ` + ` 
		` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	professorRatings, err := pgx.CollectRows(rows, pgx.RowToStructByName[etp.ProfessorRating])
	if err != nil {
		return nil, 0, err
	}

	professorRatingsPtr := make([]*etp.ProfessorRating, len(professorRatings))
	for i := range professorRatings {
		professorRatingsPtr[i] = &professorRatings[i]
	}

	return professorRatingsPtr, n, nil
}

func createProfessorRating(ctx context.Context, tx *Tx, professorRating *etp.ProfessorRating) (error, int) {
	query := `
		INSERT INTO professor_rating (
			rating, 
			comment, 
			would_take_again, 
			mandatory_attendance, 
			grade, 
			textbook_required, 
			is_approved, 
			approvals_count, 
			professor_id, 
			course_id,
			difficulty,
			school_id
		)
		VALUES (
			@rating, 
			@comment, 
			@wouldTakeAgain, 
			@mandatoryAttendance, 
			@grade, 
			@textbookRequired, 
			@isApproved, 
			@approvalsCount, 
			@professorId, 
			@courseId,
			@difficulty,
			@schoolId
		) 
		returning id
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
		"difficulty":          professorRating.Difficulty,
		"schoolId":            professorRating.SchoolId,
	}

	var id int
	err := tx.QueryRow(ctx, query, args).Scan(&id)
	if err != nil {
		return err, 0
	}

	return nil, id
}

func associateProfessorRatingTags(ctx context.Context, tx *Tx, id int, tags []int) error {
	for _, tagId := range tags {
		_, err := tx.Exec(ctx, `
			INSERT INTO professor_rating_tag (professor_rating_id, tag_id)
			VALUES ($1, $2)
		`, id, tagId)
		if err != nil {
			return err
		}
	}

	return nil
}
