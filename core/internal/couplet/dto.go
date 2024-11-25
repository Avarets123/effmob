package couplet

type CoupletCreateDto struct {
	Couplets []string `json:"couplets" validate:"required,dive"`
}

func (c *CoupletCreateDto) MapToCreate(currentCoupletNum int, songId string) []CoupletModel {

	var newCouplets []CoupletModel

	for _, v := range c.Couplets {
		currentCoupletNum++
		newCouplets = append(newCouplets, *NewCouplet(songId, v, currentCoupletNum))
	}

	return newCouplets

}

type CoupletDeleteDto struct {
	CoupletsIds []string `json:"coupletsIds" validate:"required,dive"`
}
