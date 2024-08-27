package repository

import "busha/models"

func (r *Repository) GetTotalComments() (map[int]int, error) {
	rows, err := r.DB.Table("comments").Select("COUNT(comment) as total_comment, film_id").Group("film_id").Rows()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var (
		TotalComment int
		FilmID       int
	)
	m := make(map[int]int)
	for rows.Next() {
		err = rows.Scan(&TotalComment, &FilmID)
		if err != nil {
			return nil, err
		}
		m[FilmID] = TotalComment
	}
	return m, nil
}

func (r *Repository) CreateComment(comment string, filmId int, ip string) (*models.Comments, error) {
	c := models.Comments{
		Comment:   comment,
		FilmID:    filmId,
		Ipaddress: ip,
	}
	if err := r.DB.Create(&c).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) GetComment(movieId string) (*[]models.Comments, error) {
	var comment []models.Comments

	if err := r.DB.Where("film_id = ?", movieId).Select("ipaddress, created_at, comment").Order("id desc").Find(&comment).Error; err != nil {
		return nil, err
	}
	return &comment, nil
}
