package db

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"

	"github.com/jamolh/notice-board/helpers"
	"github.com/jamolh/notice-board/models"
)

// CreateNotice - add notice to table notices
func CreateNotice(ctx context.Context, model *models.Notice) error {
	var (
		query = `INSERT INTO notices(id, title, description, price, created_at, images) 
			VALUES($1, $2, $3, $4, $5, $6)`
		err error
	)
	model.ID = uuid.NewString()
	model.CreatedAt = time.Now()
	_, err = pool.Exec(ctx, query,
		model.ID,
		model.Title,
		model.Description,
		model.Price,
		model.CreatedAt,
		model.Image)
	if err != nil {
		log.Println("db:CreateNotice failed:", err)
		return err
	}
	return nil
}

// GetNoticeByID - get notice by id
// getAllImages=true - get path all of images
// else get fisrt path from array
func GetNoticeByID(ctx context.Context, request *models.GetNoticeRequestDto) (*models.Notice, error) {
	var (
		model models.Notice
		query string
		err   error
	)

	if request.GetAllImages {
		query = `SELECT id, title, description, price, created_at, images
		FROM notices WHERE id = $1`
	} else {
		query = `SELECT id, title, description, price, created_at, json_build_array(images->0)
				FROM notices WHERE id = $1`
	}

	err = pool.QueryRow(ctx, query, request.ID).Scan(
		&model.ID,
		&model.Title,
		&model.Description,
		&model.Price,
		&model.CreatedAt,
		&model.Image)

	if err != nil {
		log.Println("db:GetNoticeByID failed:", err)
	}
	return &model, err
}

// GetNotices - get from table notices all notices
// sortField - by what field sort data
// by default it is by field created_at
// sortType - how to sort, by ascending or descending
func GetNotices(ctx context.Context, request models.GetNoticesRequestDto) ([]*models.Notice, error) {

	query := `SELECT id, 
			title, 
			description, 
			price, 
			created_at, 
			json_build_array(images->0)
				FROM notices ORDER BY ` +
		validateSortFields(request)

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Println("db:GetNotices failed:", err)
		return nil, err
	}
	defer rows.Close()

	var notices []*models.Notice
	for rows.Next() {
		notice := models.Notice{}
		err = rows.Scan(
			&notice.ID,
			&notice.Title,
			&notice.Description,
			&notice.Price,
			&notice.CreatedAt,
			&notice.Image)
		if err != nil {
			log.Println("db:GetNotices scanning failed:", err)
			return nil, err
		}
		notices = append(notices, &notice)
	}
	return notices, nil
}

func CheckNoticeExistsByTitle(ctx context.Context, title string) (bool, error) {
	var (
		query  = `SELECT EXISTS(SELECT 1 FROM notices WHERE title = $1)`
		err    error
		exists bool
	)
	err = pool.QueryRow(ctx, query, title).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		log.Println("db:CheckNoticeExistsByTitle error:", err)
		return false, err
	}

	return exists, nil
}

const (
	fieldCreatedAt = "created_at"
	fieldPrice     = "price"

	sortByDesc = "desc"
	sortByAsc  = "asc"
)

func validateSortFields(request models.GetNoticesRequestDto) (sql string) {
	// remove needless symbols
	request.Field = helpers.RemoveNonLetter(request.Field)
	request.Order = helpers.RemoveNonLetter(request.Order)

	// if we got field price
	// then sort by price
	// else sort field by created_at
	switch request.Field {
	case fieldPrice:
		sql = fieldPrice
	default:
		sql = fieldCreatedAt
	}

	switch request.Order {
	case sortByAsc:
		sql += " " + sortByAsc
	default:
		sql += " " + sortByDesc
	}

	return sql
}
