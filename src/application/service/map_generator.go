package service

import (
	"github.com/sirupsen/logrus"
)

const (
	top    = 0
	right  = 1
	bottom = 2
	left   = 3
	asc    = 0
	desc   = 1
)

var planePosition = map[int]map[int][]int{
	0: {
		top:    {3, bottom, asc},
		right:  {5, left, asc},
		bottom: {1, top, asc},
		left:   {4, right, asc},
	},
	1: {
		top:    {0, bottom, asc},
		right:  {5, bottom, asc},
		bottom: {2, top, asc},
		left:   {4, bottom, asc},
	},
	2: {
		top:    {1, bottom, asc},
		right:  {5, right, desc},
		bottom: {3, top, asc},
		left:   {4, left, desc},
	},
	3: {
		top:    {2, bottom, asc},
		right:  {5, top, desc},
		bottom: {0, top, asc},
		left:   {4, top, asc},
	},
	4: {
		top:    {3, left, asc},
		right:  {0, left, asc},
		bottom: {1, left, desc},
		left:   {2, left, desc},
	},
	5: {
		top:    {3, right, desc},
		right:  {2, right, desc},
		bottom: {1, right, asc},
		left:   {0, right, asc},
	},
}

func fetchLineFromPlane(plane *Plane, direction int, order int) []int {
	lineLen := len(plane.Cells[0])
	cellIDs := []int{}
	rawIDs := []int{}
	for i := 0; i < lineLen; i++ {
		cellIDs = append(cellIDs, -1)
		rawIDs = append(rawIDs, -1)
	}
	if direction == top {
		for k, v := range plane.Cells[0] {
			rawIDs[k] = v.ID
		}
	} else if direction == bottom {
		for k, v := range plane.Cells[lineLen-1] {
			rawIDs[k] = v.ID
		}
	} else if direction == right {
		for k, v := range plane.Cells {
			rawIDs[k] = v[lineLen-1].ID
		}
	} else if direction == left {
		for k, v := range plane.Cells {
			rawIDs[k] = v[0].ID
		}
	}

	if order == 0 {
		for k, v := range rawIDs {
			cellIDs[k] = v
		}
	} else {
		for k, v := range rawIDs {
			cellIDs[(lineLen-1)-k] = v
		}
	}
	return cellIDs
}

// MapGenerator マップ生成サービス
type MapGenerator struct{}

// Cell マップの１つの領域
type Cell struct {
	ID        int      `json:"id"`
	Lat       int      `json:"lat"`       // 緯度 -90 ~ 90
	Height    int      `json:"height"`    // 高度 -1000 1000
	Neighbors []*int   `json:"neighbors"` // 隣接マスの番号
	Attrs     []string `json:"attrs"`     // 付属物
	Owner     int      `json:"owner"`     // 所有者
}

// Plane １平面
type Plane struct {
	ID    int `json:"id"` // 平面の位置. 0底面,1正面,2右側面,3裏面,4左側面,6天面
	Cells map[int]map[int]*Cell
}

func (m *MapGenerator) Generate(n int) []*Plane { // 分割数. 2n+1
	lineCells := (2*n + 1)              // 一辺あたりの分割数
	planeCells := lineCells * lineCells // 一面あたりの分割数

	// マスと面を生成する
	planes := make([]*Plane, 6, 6)
	cells := map[int]*Cell{}
	for i := 0; i < 6*planeCells; i++ {
		neighbors := []*int{}
		for i := 0; i < 4; i++ {
			neighbors = append(neighbors, nil)
		}
		cells[i] = &Cell{ID: i, Neighbors: neighbors}
	}

	// 面にマスを所属させる
	for i := 0; i < 6; i++ {
		cellsOnPlane := map[int]map[int]*Cell{}
		for j := 0; j < lineCells; j++ {
			cellsOnPlane[j] = map[int]*Cell{}
			bias := i*planeCells + j*lineCells
			for k := 0; k < lineCells; k++ {
				id := bias + k
				target, _ := cells[id]
				cellsOnPlane[j][k] = target
			}
		}
		planes[i] = &Plane{ID: i, Cells: cellsOnPlane}
	}

	// マスに位置関係を設定する
	for i := 0; i < 6; i++ {
		plane := planes[i]
		for j := 0; j < lineCells; j++ {
			for k := 0; k < lineCells; k++ {
				target := plane.Cells[j][k]
				// 上端は他の面と接するので後で設定する
				if j > 0 {
					id := target.ID - lineCells // 1行上を隣接に設定
					target.Neighbors[top] = &id
				} else {
					position := planePosition[i][top]
					neighborLine := fetchLineFromPlane(planes[position[0]], position[1], position[2])
					target.Neighbors[top] = &neighborLine[k]
				}

				if j < lineCells-1 {
					id := target.ID + lineCells // 1行下を隣接に設定
					target.Neighbors[bottom] = &id
				} else {
					position := planePosition[i][bottom]
					neighborLine := fetchLineFromPlane(planes[position[0]], position[1], position[2])
					target.Neighbors[bottom] = &neighborLine[k]
				}

				if k > 0 {
					id := target.ID - 1 // 1つ左を隣接に設定
					target.Neighbors[left] = &id
				} else {
					position := planePosition[i][left]
					neighborLine := fetchLineFromPlane(planes[position[0]], position[1], position[2])
					target.Neighbors[left] = &neighborLine[j]
				}

				if k < lineCells-1 {
					id := target.ID + 1 // 1つ右を隣接に設定
					target.Neighbors[right] = &id
				} else {
					position := planePosition[i][right]
					neighborLine := fetchLineFromPlane(planes[position[0]], position[1], position[2])
					target.Neighbors[right] = &neighborLine[j]
				}
			}
		}

	}
	logrus.Infof("hoge")
	return planes
}
