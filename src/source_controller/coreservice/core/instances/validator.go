/*
 * Tencent is pleased to support the open source community by making 蓝鲸 available.,
 * Copyright (C) 2017-2018 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the ",License",); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an ",AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 */

package instances

import (
	"configcenter/src/common"
	"configcenter/src/common/errors"
	"configcenter/src/common/http/rest"
	"configcenter/src/common/language"
	"configcenter/src/common/metadata"
)

type validator struct {
	errIf         errors.DefaultCCErrorIf
	properties    map[string]metadata.Attribute
	idToProperty  map[int64]metadata.Attribute
	propertySlice []metadata.Attribute
	require       map[string]bool
	requireFields []string
	dependent     OperationDependences
	objID         string
	language      language.CCLanguageIf
}

// Init init
func NewValidator(kit *rest.Kit, dependent OperationDependences, objID string, bizID int64, language language.CCLanguageIf) (*validator, error) {
	valid := &validator{}
	valid.properties = make(map[string]metadata.Attribute)
	valid.idToProperty = make(map[int64]metadata.Attribute)
	valid.propertySlice = make([]metadata.Attribute, 0)
	valid.require = make(map[string]bool)
	valid.requireFields = make([]string, 0)
	valid.errIf = kit.CCError
	result, err := dependent.SelectObjectAttWithParams(kit, objID, bizID)
	if nil != err {
		return valid, err
	}
	for _, attr := range result {
		if attr.PropertyID == common.BKChildStr || attr.PropertyID == common.BKParentStr {
			continue
		}
		valid.properties[attr.PropertyID] = attr
		valid.idToProperty[attr.ID] = attr
		valid.propertySlice = append(valid.propertySlice, attr)
		if attr.IsRequired {
			valid.require[attr.PropertyID] = true
			valid.requireFields = append(valid.requireFields, attr.PropertyID)
		}
	}
	valid.objID = objID
	valid.dependent = dependent
	valid.language = language
	return valid, nil
}
