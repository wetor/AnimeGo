// SPDX-License-Identifier: AGPL-3.0-only
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.
// See the GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program. If not, see <https://www.gnu.org/licenses/>

package res

import (
	"github.com/wetor/AnimeGo/third_party/bangumi/model"
)

type Episode struct {
	Airdate     string              `json:"airdate"`
	Name        string              `json:"name"`
	NameCN      string              `json:"name_cn"`
	Duration    string              `json:"duration"`
	Description string              `json:"desc"`
	Ep          float32             `json:"ep"`
	Sort        float32             `json:"sort"`
	ID          model.EpisodeIDType `json:"id"`
	SubjectID   model.SubjectIDType `json:"subject_id"`
	Comment     uint32              `json:"comment"`
	Type        model.EpTypeType    `json:"type"`
	Disc        uint8               `json:"disc"`
}

type Paged struct {
	Data   []*Episode `json:"data"`
	Total  int64      `json:"total"`
	Limit  int        `json:"limit"`
	Offset int        `json:"offset"`
}
