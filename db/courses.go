package db

import (
	"context"
	"strings"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/jackc/pgx/v5"
)

type CourseService struct {
	db *DB
}

func NewCourseService(db *DB) *CourseService {
	return &CourseService{
		db: db,
	}
}

func (cs *CourseService) CreateCourse(ctx context.Context, course *etp.Course) error {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return err
	}

	err = createCourse(ctx, tx, course)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (cs *CourseService) UpdateCourse(ctx context.Context, id int, upd *etp.CourseUpdate) (*etp.Course, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	course, err := updateCourse(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}

	return course, tx.Commit(ctx)
}

func (cs *CourseService) GetCourses(ctx context.Context, filter etp.CourseFilter) ([]*etp.Course, int, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback(ctx)

	courses, n, err := getCourses(ctx, tx, filter)
	if err != nil {
		return nil, n, err
	}

	return courses, n, nil
}

func (cs *CourseService) GetCourseById(ctx context.Context, id int) (*etp.Course, error) {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	course, err := getCourseById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func (cs *CourseService) DeleteCourse(ctx context.Context, id int) error {
	tx, err := cs.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	err = deleteCourse(ctx, tx, id)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func deleteCourse(ctx context.Context, tx *Tx, id int) error {
	_, err := tx.Exec(ctx, `delete from course where id = @id`, pgx.NamedArgs{"id": id})
	if err != nil {
		return err
	}

	return nil
}

func updateCourse(ctx context.Context, tx *Tx, id int, upd *etp.CourseUpdate) (*etp.Course, error) {
	course, err := getCourseById(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if v := upd.Code; v != nil {
		course.Code = *v
	}
	if v := upd.Name; v != nil {
		course.Name = *v
	}
	if v := upd.Credits; v != nil {
		course.Credits = *v
	}

	query := `
		update
			course
		set
			name = @name
			code = @code
			credits = @credits
		where id = @id
	`
	args := pgx.NamedArgs{
		"name":    course.Name,
		"code":    course.Code,
		"credits": course.Credits,
	}

	_, err = tx.Exec(ctx, query, args)
	if err != nil {
		return nil, err
	}

	return course, nil
}

func createCourse(ctx context.Context, tx *Tx, course *etp.Course) error {
	query := `
		insert into course
		(
			name,
			code,
			credits,
			department_id,
			school_id,
			created_at
		)
		values
		(
			@name,
			@code,
			@credits,
			@departmentId,
			@schoolId,
			now()
		)
	`

	args := pgx.NamedArgs{
		"name":         course.Name,
		"code":         course.Code,
		"credits":      course.Credits,
		"schoolId":     course.SchoolID,
		"departmentId": course.DepartmentID,
	}

	_, err := tx.Exec(ctx, query, args)
	if err != nil {
		return err
	}

	return nil
}

func getCourseById(ctx context.Context, tx *Tx, id int) (*etp.Course, error) {
	courses, n, err := getCourses(ctx, tx, etp.CourseFilter{ID: &id})
	if err != nil {
		return nil, err
	}

	if n == 0 {
		return nil, &etp.Error{Code: etp.ENOTFOUND, Message: "course not found"}
	}

	return courses[0], nil
}

func getCourses(ctx context.Context, tx *Tx, filter etp.CourseFilter) ([]*etp.Course, int, error) {
	where, args := []string{"1 = 1"}, pgx.NamedArgs{}

	if v := filter.SchoolID; v != nil {
		where = append(where, "school_id = @schoolID")
		args["schoolID"] = *v
	}

	if v := filter.DepartmentId; v != nil {
		where = append(where, "department_id = @departmentId")
		args["departmentId"] = *v
	}

	countQuery := `
		select 
			count(*)
		from 
			course
		where ` + strings.Join(where, " and ")

	var n int
	err := tx.QueryRow(ctx, countQuery, args).Scan(&n)
	if err != nil {
		return nil, 0, err
	}

	if n == 0 {
		return nil, 0, nil
	}

	query := `
		select 
			id,
			code,
			credits,
			department_id,
			school_id
		from
			course
		where` + strings.Join(where, " and ") + `
	` + FormatLimitOffset(filter.Limit, filter.Offset)

	rows, err := tx.Query(ctx, query, args)
	if err != nil {
		return nil, 0, err
	}

	courses, err := pgx.CollectRows(rows, pgx.RowToStructByName[*etp.Course])
	if err != nil {
		return nil, 0, err
	}

	return courses, n, nil
}
