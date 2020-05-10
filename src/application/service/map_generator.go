package service

const (
	top    = 0
	right  = 1
	bottom = 2
	left   = 3
	tr     = 4
	br     = 5
	bl     = 6
	tl     = 7
	asc    = 0
	desc   = 1
)

// SetNeighbors 隣接マスを読み込む
func (c *Cell) SetNeighbors(n int, planes []*Plane) {
	maxIdx := 2 * n
	i := c.Coordinate[0] // 面
	j := c.Coordinate[1] // 行
	k := c.Coordinate[2] // 列

	// 端っこ以外のマスの場合の上下処理
	if j-1 >= 0 && j+1 <= maxIdx {
		c.Neighbors[top] = &planes[i].Cells[j-1][k].ID
		c.Neighbors[bottom] = &planes[i].Cells[j+1][k].ID
	} else {
		// 上端の場合
		if j-1 < 0 {
			switch i {
			case 0:
				c.Neighbors[top] = &planes[3].Cells[maxIdx][k].ID
			case 1:
				c.Neighbors[top] = &planes[0].Cells[maxIdx][k].ID
			case 2:
				c.Neighbors[top] = &planes[1].Cells[maxIdx][k].ID
			case 3:
				c.Neighbors[top] = &planes[2].Cells[maxIdx][k].ID
			case 4:
				c.Neighbors[top] = &planes[3].Cells[k][0].ID
			case 5:
				c.Neighbors[top] = &planes[3].Cells[maxIdx-k][maxIdx].ID
			}
		} else {
			c.Neighbors[top] = &planes[i].Cells[j-1][k].ID
		}

		// 下端の場合
		if j+1 > 2*n {
			switch i {
			case 0:
				c.Neighbors[bottom] = &planes[1].Cells[0][k].ID
			case 1:
				c.Neighbors[bottom] = &planes[2].Cells[0][k].ID
			case 2:
				c.Neighbors[bottom] = &planes[3].Cells[0][k].ID
			case 3:
				c.Neighbors[bottom] = &planes[0].Cells[0][k].ID
			case 4:
				c.Neighbors[bottom] = &planes[1].Cells[maxIdx-k][0].ID
			case 5:
				c.Neighbors[bottom] = &planes[1].Cells[k][maxIdx].ID
			}
		} else {
			c.Neighbors[bottom] = &planes[i].Cells[j+1][k].ID
		}
	}

	// 端っこ以外のマスの場合の左右処理
	if k-1 >= 0 && k+1 <= 2*n {
		c.Neighbors[right] = &planes[i].Cells[j][k+1].ID
		c.Neighbors[left] = &planes[i].Cells[j][k-1].ID
	} else {
		// 左端の場合
		if k-1 < 0 {
			switch i {
			case 0:
				c.Neighbors[left] = &planes[4].Cells[j][maxIdx].ID
			case 1:
				c.Neighbors[left] = &planes[4].Cells[maxIdx][maxIdx-j].ID
			case 2:
				c.Neighbors[left] = &planes[5].Cells[maxIdx-j][0].ID
			case 3:
				c.Neighbors[left] = &planes[4].Cells[0][j].ID
			case 4:
				c.Neighbors[left] = &planes[2].Cells[maxIdx-j][0].ID
			case 5:
				c.Neighbors[left] = &planes[0].Cells[j][maxIdx].ID
			}
		} else {
			c.Neighbors[left] = &planes[i].Cells[j][k-1].ID
		}

		// 右端の場合
		if k+1 > 2*n {
			switch i {
			case 0:
				c.Neighbors[right] = &planes[5].Cells[j][0].ID
			case 1:
				c.Neighbors[right] = &planes[5].Cells[maxIdx][j].ID
			case 2:
				c.Neighbors[right] = &planes[4].Cells[maxIdx-j][maxIdx].ID
			case 3:
				c.Neighbors[right] = &planes[5].Cells[0][maxIdx-j].ID
			case 4:
				c.Neighbors[right] = &planes[0].Cells[j][0].ID
			case 5:
				c.Neighbors[right] = &planes[2].Cells[maxIdx-j][maxIdx].ID
			}
		} else {
			c.Neighbors[right] = &planes[i].Cells[j][k+1].ID
		}
	}

}

// MapGenerator マップ生成サービス
type MapGenerator struct {
	generated map[int]*Cell // TODO use DB
}

// Cell マップの１つの領域
type Cell struct {
	ID         int      `json:"id"`
	Coordinate []int    `json:"coordinate"` // i,j,k
	Lat        int      `json:"lat"`        // 緯度 -90 ~ 90
	Height     int      `json:"height"`     // 高度 -1000 1000
	Neighbors  []*int   `json:"neighbors"`  // 隣接マスの番号
	Attrs      []string `json:"attrs"`      // 付属物
	Owner      int      `json:"owner"`      // 所有者
}

type Current struct {
	Center    *Cell   `json:"center"`
	Neighbors []*Cell `json:"neighbors"` // 直接接している4マス
	Indirects []*Cell `json:"indirects"` // 間接的に接している8マス
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
				target.Coordinate = []int{i, j, k}
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

				target.SetNeighbors(n, planes)
			}
		}
	}

	generated := map[int]*Cell{}
	for _, plane := range planes {
		for _, line := range plane.Cells {
			for _, cell := range line {
				generated[cell.ID] = cell
			}
		}
	}

	m.generated = generated

	return planes
}

func (m *MapGenerator) Get(id int) *Current {
	res := &Current{
		Center:    m.generated[id],
		Neighbors: []*Cell{top: nil, bottom: nil, left: nil, right: nil},
		Indirects: []*Cell{top: nil, bottom: nil, left: nil, right: nil, tr: nil, br: nil, bl: nil, tl: nil},
	}
	res.Neighbors[top] = m.generated[*res.Center.Neighbors[top]]
	res.Neighbors[bottom] = m.generated[*res.Center.Neighbors[bottom]]
	res.Neighbors[left] = m.generated[*res.Center.Neighbors[left]]
	res.Neighbors[right] = m.generated[*res.Center.Neighbors[right]]
	res.Indirects[top] = m.generated[*res.Neighbors[top].Neighbors[top]]          // うえのうえ
	res.Indirects[tr] = m.generated[*res.Neighbors[top].Neighbors[right]]         // うえのみぎ
	res.Indirects[bottom] = m.generated[*res.Neighbors[bottom].Neighbors[bottom]] // したのした
	res.Indirects[bl] = m.generated[*res.Neighbors[bottom].Neighbors[left]]       // したのひだり
	res.Indirects[right] = m.generated[*res.Neighbors[right].Neighbors[right]]    // 右の右
	res.Indirects[br] = m.generated[*res.Neighbors[right].Neighbors[bottom]]      // 右の下
	res.Indirects[left] = m.generated[*res.Neighbors[left].Neighbors[left]]       // 左の左
	res.Indirects[tl] = m.generated[*res.Neighbors[left].Neighbors[top]]          // 左の↑

	return res

}
